import React from "react";
import { useWebsocket } from "@app/api/ws.ts";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { api, QueryKeys } from "@app/api";
import {
  FriendRequestNotification,
  FriendRequestNotificationSchema,
} from "@app/api/wstypes/notifications.ts";

type GetFriendRequestNotificationsResponse = {
  notifications: FriendRequestNotification[];
};

function useFriendRequestsValue() {
  const { handleMessage, connect } = useWebsocket();
  const wsRef = React.useRef<WebSocket | null>(null);
  const queryClient = useQueryClient();

  const { data: notificationsData } = useQuery({
    queryKey: QueryKeys.getFriendRequestNotifications(),
    queryFn: async () =>
      api.get<GetFriendRequestNotificationsResponse>(
        `/friends/requests/notifications?seen=false&limit=5`,
      ),
  });

  React.useEffect(() => {
    const onMessage = (data: unknown) => {
      const friendRequestValidation =
        FriendRequestNotificationSchema.safeParse(data);
      if (friendRequestValidation.success) {
        queryClient.setQueryData(
          QueryKeys.getFriendRequestNotifications(),
          (oldData: GetFriendRequestNotificationsResponse) => {
            return {
              notifications: [
                friendRequestValidation.data,
                ...oldData.notifications,
              ],
            };
          },
        );
      }
    };

    const messageHandler = handleMessage(onMessage);
    wsRef.current = connect(
      `/notifications`,
      "Could not to notifications socket",
    );

    if (wsRef.current) {
      wsRef.current.addEventListener("message", messageHandler);
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect, handleMessage, queryClient]);

  const hasUnseenFriendRequestNotifications =
    notificationsData &&
    notificationsData.notifications.some((notification) => !notification.seen);

  return {
    friendRequestNotifications: notificationsData,
    hasUnseenFriendRequestNotifications,
  };
}

type FriendRequestContextType = ReturnType<typeof useFriendRequestsValue>;

const FriendRequestsContext =
  React.createContext<FriendRequestContextType | null>(null);

export function useFriendRequests() {
  const context = React.useContext(FriendRequestsContext);

  if (!context) {
    throw new Error(
      "useFriendRequests must be used within a FriendRequestsProvider",
    );
  }

  return context;
}

export default function FriendRequestsProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const value = useFriendRequestsValue();

  return (
    <FriendRequestsContext.Provider value={value}>
      {children}
    </FriendRequestsContext.Provider>
  );
}
