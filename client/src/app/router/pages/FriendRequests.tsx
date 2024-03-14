import React from "react";
import { useQuery } from "@tanstack/react-query";
import FriendRequestItem from "@app/components/friends/FriendRequestItem";
import { useFriendRequests } from "@app/context/FriendRequestsProvider";
import { api, SuccessMessageResponse } from "@app/api";
import { Container } from "@app/components/friends/FriendPageContainer";
import DCSearchBar from "@app/components/SearchBar";
import { usePendingFriendRequests } from "@app/hooks/reactQuery/usePendingFriendRequests.ts";
import { LoadingSpinner } from "@app/components/LoadingSpinner.tsx";

export default function FriendRequestsPage() {
  const { markAllAsSeen, hasNewRequests } = useFriendRequests();

  const { data: requests, isLoading, error } = usePendingFriendRequests();

  useQuery({
    enabled: hasNewRequests,
    queryKey: ["seen-all"],
    queryFn: async () => {
      const data = await api.put<SuccessMessageResponse>(
        "/friends/invites/seen",
      );

      markAllAsSeen();
      return data;
    },
  });

  const [search, setSearch] = React.useState("");

  const filteredRequests =
    requests?.filter((request) =>
      request.username.toLowerCase().includes(search.toLowerCase()),
    ) ?? [];

  if (error) {
    return (
      <Container className="flex items-center justify-center h-full">
        <h1 className="text-dc-neutral-300">Error fetching friend requests</h1>
      </Container>
    );
  }

  if (isLoading) {
    return (
      <Container className="flex items-center justify-center h-full">
        <LoadingSpinner />
      </Container>
    );
  }

  if (!filteredRequests.length) {
    return (
      <Container className="flex items-center justify-center h-full">
        <h1 className="text-dc-neutral-300">No friend requests</h1>
      </Container>
    );
  }

  return (
    <>
      <Container className="pt-4">
        <DCSearchBar value={search} setValue={setSearch} />
      </Container>
      <Container className="py-4">
        <h1 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Waiting - {filteredRequests.length}
        </h1>
      </Container>
      <Container className="pb-4 overflow-auto">
        {filteredRequests.map((request) => (
          <FriendRequestItem
            id={request.id}
            userId={request.id}
            username={request.username}
            key={request.id}
          />
        ))}
      </Container>
    </>
  );
}
