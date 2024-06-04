import React from "react";
import { useWebsocket } from "@app/api/ws.ts";

function useFriendRequestsValue() {
  const [friendRequestNotifications, setFriendRequestNotifications] =
    React.useState(0);
  const { handleMessage, connect } = useWebsocket();
  const wsRef = React.useRef<WebSocket | null>(null);

  React.useEffect(() => {
    const onMessage = () => {
      console.log("new request");
      setFriendRequestNotifications((prev) => {
        return prev + 1;
      });
    };

    const messageHandler = handleMessage(onMessage);
    wsRef.current = connect(
      `/notifications`,
      "Could not to notifications socket",
    );

    if (wsRef.current) {
      wsRef.current.addEventListener("message", messageHandler);
      wsRef.current.addEventListener("open", () => {
        console.log("notifications socket opened");
      });
    }

    return () => {
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect, handleMessage]);

  function markAllAsSeen() {
    console.log("marking all as seen");
    setFriendRequestNotifications(0);
  }

  return {
    markAllAsSeen,
    friendRequestNotifications,
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
