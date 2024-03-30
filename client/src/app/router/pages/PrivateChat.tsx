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
  GetAllChats,
  GetChat,
  QueryKeys,
  SuccessMessageResponse,
} from "@app/api";
import { useAuth } from "../../context/AuthProvider";
import OneDayChatMessageGroup from "../../components/chats/OneDayChatMessageGroup";
import { useChatId } from "@app/hooks/useChatId.ts";
import { useWebsocket } from "@app/api/ws.ts";
import {
  ChatNameUpdatedWsSchema,
  NewMessageWsSchema,
  NewMessageWsType,
} from "@app/api/wstypes/chats.ts";
import { ClientError } from "@app/utils/clientError.ts";
import { cn } from "@app/utils/cn.ts";

const nameChangeSchema = z.object({
  newName: z.string(),
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
  const { mutate, isPending } = useMutation<
    SuccessMessageResponse,
    ClientError,
    NameChangeFormData
  >({
    mutationFn: async (data: NameChangeFormData) =>
      api.put<SuccessMessageResponse>(`/chats/${chatId}/update-name`, {
        body: JSON.stringify(data),
      }),
    onMutate: async (inputData) => {
      updateChatName(inputData.newName);
      const oldData = queryClient.getQueryData<GetAllChats>(
        QueryKeys.getAllChats(),
      );

      if (!oldData) return;

      const newChatList = oldData.chats.map((c) => {
        if (c.id === chatId) {
          return {
            ...c,
            name: inputData.newName,
          };
        }

        return c;
      });

      const newData: GetAllChats = {
        chats: newChatList,
      };

      queryClient.setQueryData(QueryKeys.getAllChats(), newData);
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({
        queryKey: QueryKeys.getAllChats(),
      });
    },
    onError: async (err) => {
      if (err.code === 400) {
        toast.error("New chat name must be between 6 and 32 characters");
      } else {
        toast.error("Unknown error when changing chat name");
      }
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
      newName: name,
    },
  });

  React.useEffect(() => {
    form.setValue("newName", name);
  }, [chatId, name, form]);

  return (
    <form
      onSubmit={form.handleSubmit((data) => {
        setOldChatName(name);
        mutate(data);
      })}
    >
      <input
        type="text"
        {...form.register("newName")}
        className="px-2 rounded-md bg-transparent border ring-0 border-transparent font-semibold hover:border-dc-neutral-950 focus:bg-dc-neutral-950 focus:border-dc-neutral-1000 disabled:hover:bg-transparent disabled:hover:border-transparent disabled:cursor-text"
        size={form.watch("newName").length}
        disabled={disabled || isPending}
      />
      <button className="hidden" type="submit" />
    </form>
  );
}

export default function PrivateChat() {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  const [newMessage, setNewMessage] = React.useState<string>("");

  const chatId = useChatId();
  const toast = useToast();
  const wsRef = React.useRef<WebSocket | null>(null);

  const { data: chatInfo } = useQuery({
    queryKey: QueryKeys.getChat(chatId),
    queryFn: async () => api.get<GetChat>(`/chats/${chatId}`),
  });

  const { mutate: sendMessageMutate, isPending: isSendingMessage } =
    useMutation({
      mutationFn: async (params: { message: string }) =>
        api.post<SuccessMessageResponse>(`/chats/${chatId}/messages`, {
          body: JSON.stringify({
            text: params.message,
          }),
        }),
      onError: async () => {
        toast.error("Could not send message");
        await queryClient.invalidateQueries({
          queryKey: QueryKeys.getChat(chatId),
        });
      },
    });

  const groupedMessages = useGroupedMessages(chatInfo?.messages ?? []);
  const { handleMessage, connect } = useWebsocket();

  const updateChatNameFromWs = React.useCallback(
    (newChatName: string) => {
      queryClient.setQueryData(
        QueryKeys.getChat(chatId),
        (oldChatInfo: GetChat): GetChat => {
          return {
            ...oldChatInfo,
            name: newChatName,
          };
        },
      );
      queryClient.setQueryData(
        QueryKeys.getAllChats(),
        (oldGetChatsResponse: GetAllChats): GetAllChats => {
          return {
            chats: oldGetChatsResponse.chats.map((chat) =>
              chat.id === chatId ? { ...chat, name: newChatName } : chat,
            ),
          };
        },
      );
    },
    [chatId, queryClient],
  );

  const addNewMessageToChat = React.useCallback(
    (message: NewMessageWsType["message"]) => {
      queryClient.setQueryData(
        QueryKeys.getChat(chatId),
        (oldData: GetChat): GetChat => {
          oldData.messages.unshift(message);
          return oldData;
        },
      );
    },
    [chatId, queryClient],
  );

  const onMessage = React.useCallback(
    (data: unknown) => {
      const newMessageResult = NewMessageWsSchema.safeParse(data);
      if (newMessageResult.success) {
        addNewMessageToChat(newMessageResult.data.message);
      }
      const chatNameUpdateResult = ChatNameUpdatedWsSchema.safeParse(data);
      if (chatNameUpdateResult.success) {
        updateChatNameFromWs(chatNameUpdateResult.data.newName);
      }
    },
    [addNewMessageToChat, updateChatNameFromWs],
  );

  React.useEffect(() => {
    const messageHandler = handleMessage(onMessage);
    wsRef.current = connect(
      `/chats/${chatId}`,
      "Could not connect to chat, please try again later",
    );

    if (wsRef.current) {
      wsRef.current.addEventListener("message", messageHandler);
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [chatId]);

  function sendMessage(e: React.FormEvent) {
    if (isSendingMessage) return;

    e.preventDefault();
    const message = newMessage.trim();

    if (!message || !user) return;
    sendMessageMutate({ message });
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
        <form
          onSubmit={sendMessage}
          className={cn(
            "w-full flex pt-0 pb-4 px-4",
            isSendingMessage && "opacity-70",
          )}
        >
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
