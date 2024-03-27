import React from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useGroupedMessages } from "../../hooks/useGroupedMessages";
import { useToast } from "../../hooks/useToast";
import {
  api,
  ChatType,
  GetChat,
  Message,
  QueryKeys,
  SuccessMessageResponse,
} from "@app/api";
import { useAuth } from "../../context/AuthProvider";
import OneDayChatMessageGroup from "../../components/chats/OneDayChatMessageGroup";
import { useChatId } from "@app/hooks/useChatId.ts";
import { useWebsocket } from "@app/api/ws.ts";
import {
  NewMessageWsSchema,
  NewMessageWsType,
} from "@app/api/wstypes/chats.ts";

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
  chatId: number;
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
      const oldData = queryClient.getQueryData(QueryKeys.getAllChats());

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
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });
    },
    onError: async (err) => {
      toast.error(err.message);
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });

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
  const [lastMessageId, setLastMessageId] = React.useState(-1);
  const [sentMessages, setSentMessages] = React.useState<Message[]>([]);
  const chatId = useChatId();

  function updateLastMessageId() {
    setLastMessageId((curr) => curr - 1);
  }

  const { data: chatInfo } = useQuery({
    queryKey: QueryKeys.getChat(chatId),
    queryFn: async () => api.get<GetChat>(`/chats/${chatId}`),
  });

  // TODO: handle error when sending messages (revalidate query would be easiest I think)
  const { mutate: sendMessageMutate } = useMutation({
    mutationFn: async (message: string) =>
      api.post<SuccessMessageResponse>(`/chats/${chatId}/messages`, {
        body: JSON.stringify({
          text: message,
        }),
      }),
  });

  const groupedMessages = useGroupedMessages(chatInfo?.messages ?? []);

  const [newMessage, setNewMessage] = React.useState<string>("");

  useWebsocket({
    path: `/chats/${chatId}`,
    onMessage: (data) => {
      const newMessageResult = NewMessageWsSchema.safeParse(data);
      if (newMessageResult.success) {
        addNewMessageToChat(newMessageResult.data.message);
      }
    },
  });

  function addNewMessageToChat(message: NewMessageWsType["message"]) {
    if (message.user.id === user?.id) {
      const lastMessage = sentMessages[sentMessages.length - 1];

      if (!lastMessage) return;

      queryClient.setQueryData(
        QueryKeys.getChat(chatId),
        (oldData: GetChat): GetChat => {
          const newMessages = oldData.messages.map((m) => {
            if (m.id === lastMessage.id) {
              return message;
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
      QueryKeys.getChat(chatId),
      (oldData: GetChat): GetChat => ({
        ...oldData,
        messages: [message, ...oldData.messages],
      }),
    );
  }

  function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    const message = newMessage.trim();

    if (!message || !user) return;
    sendMessageMutate(message);
    const messageObj: Message = {
      id: lastMessageId,
      text: message,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      user,
    };
    updateLastMessageId();
    queryClient.setQueryData(
      QueryKeys.getChat(chatId),
      (oldData: GetChat): GetChat => ({
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
      QueryKeys.getChat(chatId),
      (oldData: GetChat): GetChat => ({
        ...oldData,
        name,
      }),
    );
  }

  if (!chatInfo) return null;

  const placeholder =
    chatInfo.type === ChatType.PRIVATE
      ? `Message @${chatInfo.name}`
      : `Message ${chatInfo.name}`;

  return (
    <div className="max-h-screen h-full flex-grow flex">
      <div className="flex-grow flex flex-col">
        <nav className="border-b flex border-dc-neutral-1000 w-full p-3 gap-4">
          {!!chatInfo.name && !!chatId && (
            <NameChangeElement
              name={chatInfo.name}
              disabled={chatInfo.type === ChatType.PRIVATE}
              chatId={chatId}
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
