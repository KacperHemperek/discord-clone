import React from "react";
import { useWebsocket } from "@app/api/hooks/useWebsocket";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api, GetAllChats, QueryKeys } from "@app/api";
import {
  FriendRequestNotification,
  FriendRequestNotificationSchema,
  NewMessageNotification,
  NewMessageNotificationSchema,
} from "@app/api/wstypes/notifications";
import { z } from "zod";
import { useAuth } from "@app/context/AuthProvider.tsx";

type GetNotificationsResponse<
  T extends FriendRequestNotification | NewMessageNotification,
> = {
  notifications: T[];
};

function useNotificationsContextValue() {
  const { handleMessage, connect } = useWebsocket();
  const wsRef = React.useRef<WebSocket | null>(null);
  const queryClient = useQueryClient();
  const { refreshToken, accessToken } = useAuth();

  const { data: friendRequestNotifications } = useQuery({
    queryKey: QueryKeys.getFriendRequestNotifications(),
    queryFn: async () =>
      api.get<GetNotificationsResponse<FriendRequestNotification>>(
        `/notifications/friend-requests?seen=false&limit=5`,
      ),
  });

  const { data: newMessageNotifications } = useQuery({
    queryKey: QueryKeys.getNewMessageNotifications(),
    queryFn: async () =>
      api.get<GetNotificationsResponse<NewMessageNotification>>(
        `/notifications/new-messages?seen=false`,
      ),
  });

  React.useEffect(() => {
    const onMessage = (data: unknown) => {
      const friendRequestValidation =
        FriendRequestNotificationSchema.safeParse(data);
      if (friendRequestValidation.success) {
        queryClient.setQueryData(
          QueryKeys.getFriendRequestNotifications(),
          (oldData: GetNotificationsResponse<FriendRequestNotification>) => {
            return {
              notifications: [
                friendRequestValidation.data,
                ...oldData.notifications,
              ],
            };
          },
        );
        return;
      }

      const newMessageValidation = NewMessageNotificationSchema.safeParse(data);
      if (newMessageValidation.success) {
        const chatId = newMessageValidation.data.data.chatId;
        const allChats = queryClient.getQueryData<GetAllChats>(
          QueryKeys.getAllChats(),
        );
        const chatExists = Boolean(
          allChats?.chats.some((chat) => chatId === chat.id),
        );
        if (!chatExists) {
          void queryClient.invalidateQueries({
            queryKey: QueryKeys.getAllChats(),
          });
        }

        queryClient.setQueryData(
          QueryKeys.getNewMessageNotifications(),
          (oldData: GetNotificationsResponse<NewMessageNotification>) => {
            return {
              notifications: [
                newMessageValidation.data,
                ...oldData.notifications,
              ],
            };
          },
        );
        return;
      }
      const objWithTypeSchema = z.object({
        type: z.string(),
      });
      const objWithType = objWithTypeSchema.safeParse(data);
      if (objWithType.success) {
        console.error(
          `unsupported payload type from websocket: ${objWithType.data.type}`,
        );
      }
    };

    wsRef.current = connect({
      path: `/notifications`,
      accessToken,
      refreshToken,
      errMessage: "Could not connect to notifications socket",
    });

    if (wsRef.current) {
      wsRef.current.addEventListener("message", handleMessage(onMessage));
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [accessToken, connect, handleMessage, queryClient, refreshToken]);

  const hasUnseenFriendRequestNotifications =
    friendRequestNotifications &&
    friendRequestNotifications.notifications.some(
      (notification) => !notification.seen,
    );

  return {
    friendRequestNotifications: friendRequestNotifications?.notifications,
    newMessageNotifications: newMessageNotifications?.notifications,
    hasUnseenFriendRequestNotifications,
  };
}

type NotificationsContextType = ReturnType<typeof useNotificationsContextValue>;

const NotificationsContext =
  React.createContext<NotificationsContextType | null>(null);

export function useNotifications() {
  const context = React.useContext(NotificationsContext);

  if (!context) {
    throw new Error(
      "useFriendRequests must be used within a FriendRequestsProvider",
    );
  }

  return context;
}

export default function NotificationsProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const value = useNotificationsContextValue();

  return (
    <NotificationsContext.Provider value={value}>
      {children}
    </NotificationsContext.Provider>
  );
}
