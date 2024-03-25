import React from "react";
import { z } from "zod";

const inviteItemSchema = z.object({
  id: z.string(),
  username: z.string(),
  email: z.string(),
  seen: z.boolean(),
});

type FriendInvite = z.infer<typeof inviteItemSchema>;

function useFriendRequestsValue() {
  const [requests, setRequests] = React.useState<FriendInvite[]>([]);

  // useWebsocket({
  //   path: "/friends/invites",
  //   onMessage: (event) => {
  //     const jsonData = JSON.parse(event.data);
  //
  //     if (jsonData?.type === FriendRequestType.allFriendInvites) {
  //       const parsedData = allFriendsRequestSchema.safeParse(jsonData);
  //
  //       if (parsedData.success) {
  //         setRequests(parsedData.data.payload);
  //       }
  //     }
  //
  //     if (jsonData?.type === FriendRequestType.newFriendInvite) {
  //       const parsedData = newFriendRequestSchema.safeParse(jsonData);
  //
  //       if (parsedData.success) {
  //         if (match) {
  //           queryClient.refetchQueries({ queryKey: ["seen-all"] });
  //         }
  //
  //         setRequests((notifications) => [
  //           ...notifications,
  //           parsedData.data.payload,
  //         ]);
  //       }
  //     }
  //   },
  // });

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
