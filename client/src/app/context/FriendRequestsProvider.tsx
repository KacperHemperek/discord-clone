import { FriendRequestType } from '../types/friends';
import React from 'react';
import { getWebsocketConnection } from '../utils/websocket';
import { z } from 'zod';
import { useMatch } from 'react-router-dom';
import { useQueryClient } from '@tanstack/react-query';

const inviteItemSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string(),
  seen: z.boolean(),
});

const allFriendsRequestSchema = z.object({
  type: z.literal(FriendRequestType.allFriendInvites),
  payload: z.array(inviteItemSchema),
});

const newFriendRequestSchema = z.object({
  type: z.literal(FriendRequestType.newFriendInvite),
  payload: inviteItemSchema,
});

type FriendInvite = z.infer<typeof inviteItemSchema>;

function useFriendRequestsValue() {
  const queryClient = useQueryClient();
  const match = useMatch('/home/friends/requests');

  const [requests, setRequests] = React.useState<FriendInvite[]>([]);

  const friendInviteWs = React.useRef<WebSocket | null>(null);

  React.useEffect(() => {
    const ws = getWebsocketConnection('/friends/invites');
    if (!ws) return;

    friendInviteWs.current = ws;

    friendInviteWs.current.addEventListener('message', handleWebsocketMessage);

    return () => {
      friendInviteWs.current?.close();
    };
  }, []);

  function handleWebsocketMessage({ data }: { data: string }) {
    const jsonData = JSON.parse(data);

    if (jsonData?.type === FriendRequestType.allFriendInvites) {
      const parsedData = allFriendsRequestSchema.safeParse(jsonData);

      if (parsedData.success) {
        setRequests(parsedData.data.payload);
      }
    }

    if (jsonData?.type === FriendRequestType.newFriendInvite) {
      const parsedData = newFriendRequestSchema.safeParse(jsonData);

      if (parsedData.success) {
        if (match) {
          queryClient.refetchQueries({ queryKey: ['seen-all'] });
        }

        setRequests((notifications) => [
          ...notifications,
          parsedData.data.payload,
        ]);
      }
    }
  }

  function markAllAsSeen() {
    setRequests((requests) => {
      return requests.map((request) => ({
        ...request,
        seen: true,
      }));
    });
  }

  function removeRequest(id: string) {
    setRequests((requests) => {
      return requests.filter((request) => request.id !== id);
    });
  }

  const hasNewRequests = requests.filter((n) => !n.seen).length > 0;

  return {
    requests,
    markAllAsSeen,
    hasNewRequests,
    removeRequest,
  };
}

type FriendRequestContextType = ReturnType<typeof useFriendRequestsValue>;

const FriendRequestsContext =
  React.createContext<FriendRequestContextType | null>(null);

export function useFriendRequests() {
  const context = React.useContext(FriendRequestsContext);

  if (!context) {
    throw new Error(
      'useFriendRequests must be used within a FriendRequestsProvider',
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
