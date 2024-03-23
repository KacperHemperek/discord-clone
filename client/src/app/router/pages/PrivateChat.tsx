import React from "react";
import { useParams } from "react-router-dom";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChatSocketMessageType } from "../../types/chats.ts";
import { useGroupedMessages } from "../../hooks/useGroupedMessages";
import { useToast } from "../../hooks/useToast";
import { getWebsocketConnection } from "../../utils/websocket";
import { api } from "../../api";
import { useAuth } from "../../context/AuthProvider";
import OneDayChatMessageGroup from "../../components/chats/OneDayChatMessageGroup";
import ChatLinkList from "../../components/chats/ChatLinkList";

const nameChangeSchema = z.object({
  name: z.string().min(1),
});

type NameChangeFormData = z.infer<typeof nameChangeSchema>;

function NameChangeElement({
  name,
  chatId,
  disabled,
  updateChatName,
}: {
  name: string;
  chatId: string;
  updateChatName: (name: string) => void;
  disabled?: boolean;
}) {
  const [oldChatName, setOldChatName] = React.useState<string | null>(null);

  const queryClient = useQueryClient();
  const toast = useToast();
  const { mutate, isPending } = useMutation({
    mutationFn: async (data: NameChangeFormData) =>
      api.put<ChatsTypes.UpdateChatNameSuccessResponseType>(
        `/chats/${chatId}/update-name`,
        {
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(data),
        },
      ),

    onMutate: async (inputData) => {
      updateChatName(inputData.name);
      const oldData =
        queryClient.getQueryData<ChatsTypes.GetChatsSuccessResponseType>([
          "chats",
        ]);

      if (!oldData) return;

      const newChatList = oldData.chats.map((c) => {
        if (c.id === chatId) {
          return {
            ...c,
            name: inputData.name,
          };
        }

        return c;
      });

      const newData: ChatsTypes.GetChatsSuccessResponseType = {
        chats: newChatList,
      };

      queryClient.setQueryData(["chats"], newData);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["chats"] });
    },
    onError: (err) => {
      toast.error(err.message);
      queryClient.invalidateQueries({ queryKey: ["chats"] });

      if (!oldChatName) return;

      updateChatName(oldChatName);
    },
  });

  const form = useForm<NameChangeFormData>({
    resolver: zodResolver(nameChangeSchema),
    defaultValues: {
      name,
    },
  });

  React.useEffect(() => {
    if (name && form.getValues("name") !== name) {
      form.setValue("name", name);
    }
  }, [name, form]);

  return (
    <form
      onSubmit={form.handleSubmit((data) => {
        setOldChatName(name);
        mutate(data);
      })}
    >
      <input
        type="text"
        {...form.register("name")}
        className="px-2 rounded-md bg-transparent border ring-0 border-transparent font-semibold hover:border-dc-neutral-950 focus:bg-dc-neutral-950 focus:border-dc-neutral-1000 disabled:hover:bg-transparent disabled:hover:border-transparent disabled:cursor-text"
        size={form.watch("name").length}
        disabled={disabled || isPending}
      />
      <button className="hidden" type="submit" />
    </form>
  );
}

export default function PrivateChat() {
  const { user } = useAuth();
  const queryClient = useQueryClient();

  const { chatId } = useParams<{ chatId: string }>();
  const websocketRef = React.useRef<WebSocket | null>(null);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [_sentMessages, setSentMessages] = React.useState<
    ChatsTypes.ChatMessage[]
  >([]);

  const { data: chatInfo } = useQuery({
    queryKey: ["chat", chatId],
    queryFn: async () =>
      api.get<ChatsTypes.GetChatInfoWithMessagesSuccessResponseType>(
        `/chats/${chatId}`,
      ),
  });

  const groupedMessages = useGroupedMessages(chatInfo?.messages ?? []);

  const [newMessage, setNewMessage] = React.useState<string>("");

  React.useEffect(() => {
    websocketRef.current = getWebsocketConnection(`/chats/${chatId}`);

    if (websocketRef.current) {
      websocketRef.current.addEventListener("message", (event) => {
        const data = JSON.parse(event.data) as ChatsTypes.NewMessageType;

        if (data.type === ChatSocketMessageType.newMessage) {
          addNewMessageToChat(data);
        }
      });
    }

    return () => {
      websocketRef.current?.close();
    };
  }, [chatId]);

  function addNewMessageToChat(payload: ChatsTypes.NewMessageType) {
    if (payload.message.sender.id === user?.id) {
      const lastMessage = _sentMessages.at(-1);

      if (!lastMessage) return;

      queryClient.setQueryData(
        ["chat", chatId],
        (
          oldData: ChatsTypes.GetChatInfoWithMessagesSuccessResponseType,
        ): ChatsTypes.GetChatInfoWithMessagesSuccessResponseType => {
          const newMessages = oldData.messages.map((m) => {
            if (m.id === lastMessage.id) {
              return payload.message;
            }
            return m;
          });

          return {
            ...oldData,
            messages: newMessages,
          };
        },
      );

      setSentMessages((prevSentMessages) => {
        prevSentMessages.pop();
        return prevSentMessages;
      });
    }

    queryClient.setQueryData(
      ["chat", chatId],
      (
        oldData: ChatsTypes.GetChatInfoWithMessagesSuccessResponseType,
      ): ChatsTypes.GetChatInfoWithMessagesSuccessResponseType => ({
        ...oldData,
        messages: [payload.message, ...oldData.messages],
      }),
    );
  }

  function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    const message = newMessage.trim();

    if (!message || !websocketRef.current || !user) return;

    websocketRef.current.send(
      JSON.stringify({
        type: ChatSocketMessageType.newMessage,
        message,
      }),
    );

    const messageObj: ChatsTypes.ChatMessage = {
      id: window.crypto.randomUUID(),
      text: message,
      createdAt: new Date().toISOString(),
      sender: user,
      image: null,
    };

    queryClient.setQueryData(
      ["chat", chatId],
      (
        oldData: ChatsTypes.GetChatInfoWithMessagesSuccessResponseType,
      ): ChatsTypes.GetChatInfoWithMessagesSuccessResponseType => ({
        ...oldData,
        messages: [messageObj, ...oldData.messages],
      }),
    );

    setSentMessages((prev) => {
      return [...prev, messageObj];
    });
    setNewMessage("");
  }

  function updateChatName(name: string) {
    queryClient.setQueryData(
      ["chat", chatId],
      (
        oldData: ChatsTypes.GetChatInfoWithMessagesSuccessResponseType,
      ): ChatsTypes.GetChatInfoWithMessagesSuccessResponseType => ({
        ...oldData,
        name,
      }),
    );
  }

  if (!chatInfo) return null;

  const chatName = chatInfo.name
    ? chatInfo.name
    : !!user && chatInfo.users.find((u) => u.id !== user.id)!.username;

  const placeholder =
    chatInfo.type === ChatTypes.private
      ? `Message @${chatName}`
      : `Message ${chatName}`;

  return (
    <div className="max-h-screen h-full flex-grow flex">
      <ChatLinkList />
      <div className="flex-grow flex flex-col">
        <nav className="border-b flex border-dc-neutral-1000 w-full p-3 gap-4">
          {!!chatName && !!chatId && (
            <NameChangeElement
              name={chatName}
              chatId={chatId}
              disabled={chatInfo.type === ChatTypes.private}
              updateChatName={updateChatName}
            />
          )}
        </nav>

        <div className="flex flex-col-reverse flex-grow overflow-y-scroll px-4 py-6">
          {groupedMessages.map(({ date, messages }) => (
            <OneDayChatMessageGroup
              key={date.getTime()}
              date={date}
              messages={messages}
            />
          ))}
        </div>
        <form onSubmit={sendMessage} className="w-full flex pt-0 pb-4 px-4">
          <input
            className="w-full p-2 rounded-md bg-dc-neutral-1000 outline-none"
            placeholder={placeholder}
            onChange={(e) => setNewMessage(e.target.value)}
            value={newMessage}
          />
          <button type="submit" hidden />
        </form>
      </div>
    </div>
  );
}
