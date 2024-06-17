import React from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useGroupedMessages } from "@app/hooks/useGroupedMessages";
import { useToast } from "@app/hooks/useToast";
import {
  api,
  ChatType,
  GetAllChats,
  GetChat,
  QueryKeys,
  SuccessMessageResponse,
} from "@app/api";
import { useAuth } from "@app/context/AuthProvider";
import OneDayChatMessageGroup from "@app/components/chats/OneDayChatMessageGroup";
import { useChatId } from "@app/hooks/useChatId";
import { useWebsocket } from "@app/api/hooks/useWebsocket";
import {
  ChatNameUpdatedWsSchema,
  NewMessageWsSchema,
  NewMessageWsType,
} from "@app/api/wstypes/chats";
import { cn } from "@app/utils/cn";
import { MessageCircle } from "lucide-react";
import { MutationKeys } from "@app/api/mutationKeys.ts";
import { PrivateChatHeader } from "@app/components/chats/PrivateChatHeader.tsx";

type SendHelloMessageListProps = {
  sendMessage: (message: string) => void;
};

function SendHelloMessageList({ sendMessage }: SendHelloMessageListProps) {
  const helloMessages = [
    "Hi guys how it's going 👋",
    "Hi y'all nice to meet you 👋",
    "Hello people how are you 👋",
  ];

  return (
    <div className="flex flex-col max-w-96 w-full gap-2">
      {helloMessages.map((message) => (
        <button
          className="text-left border border-dc-neutral-700 rounded-sm px-4 py-2 hover:bg-dc-neutral-850 hover:opacity-95 transition"
          onClick={() => sendMessage(message)}
        >
          {message}
        </button>
      ))}
    </div>
  );
}

export default function PrivateChat() {
  const { user, accessToken, refreshToken } = useAuth();
  const chatId = useChatId();
  const queryClient = useQueryClient();
  const [newMessage, setNewMessage] = React.useState<string>("");
  const { mutate: markNotificationsAsSeen } = useMutation({
    mutationKey: MutationKeys.markNewMessageNotificationsAsSeen(chatId),
    mutationFn: (chatId: number) =>
      api.put<SuccessMessageResponse>(
        "/notifications/new-messages/mark-as-seen",
        {
          body: JSON.stringify({
            chatId,
          }),
        },
      ),
    onSuccess: () => {
      void queryClient.invalidateQueries({
        queryKey: QueryKeys.getNewMessageNotifications(),
      });
    },
  });

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
          return { ...oldData, messages: [message, ...oldData.messages] };
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
    const cleanup = () => {
      if (wsRef.current) {
        wsRef.current?.close();
      }
    };
    if (refreshToken.length === 0 || accessToken.length === 0) {
      return cleanup;
    }
    wsRef.current = connect({
      path: `/chats/${chatId}`,
      accessToken,
      refreshToken,
      errMessage: "Could not connect to chat, please try again later",
    });

    if (wsRef.current) {
      wsRef.current.addEventListener("message", handleMessage(onMessage));
    }

    return cleanup;
  }, [chatId, refreshToken, accessToken, onMessage, connect, handleMessage]);

  React.useEffect(() => {
    markNotificationsAsSeen(chatId);
  }, [chatId, markNotificationsAsSeen]);

  function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    if (isSendingMessage) return;

    const message = newMessage.trim();

    if (!message || !user) return;
    sendMessageMutate({ message });
    setNewMessage("");
  }

  if (!chatInfo) return null;

  const placeholder =
    chatInfo.type === ChatType.PRIVATE
      ? `Message @${chatInfo.name}`
      : `Message ${chatInfo.name}`;

  return (
    <div className="max-h-screen h-full flex-grow flex">
      <div className="flex-grow flex flex-col">
        <PrivateChatHeader name={chatInfo.name} type={chatInfo.type} />
        {groupedMessages.length === 0 && (
          <div className="flex flex-col items-center justify-center mx-auto flex-grow">
            <MessageCircle size={48} className="text-dc-neutral-300" />
            <h1 className="text-dc-neutral-300 text-xl font-medium">
              There are no messages on this chat, yet!
            </h1>
            <p className="text-dc-neutral-300 pb-6">
              You can say hi to all members of this chat
            </p>
            <SendHelloMessageList
              sendMessage={(message) => sendMessageMutate({ message })}
            />
          </div>
        )}

        {groupedMessages.length > 0 && (
          <div className="flex flex-col-reverse flex-grow overflow-y-scroll px-4 py-6">
            {groupedMessages.length === 0 &&
              "There are no messages currently on this chat"}
            {groupedMessages.map(({ date, messages }) => (
              <OneDayChatMessageGroup
                key={date.getTime()}
                date={date}
                messages={messages}
              />
            ))}
          </div>
        )}
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
