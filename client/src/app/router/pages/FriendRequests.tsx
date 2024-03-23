import React from "react";
import FriendRequestItem from "@app/components/friends/FriendRequestItem";
import { Container } from "@app/components/friends/FriendPageContainer";
import DCSearchBar from "@app/components/SearchBar";
import { usePendingFriendRequests } from "@app/api";
import { LoadingSpinner } from "@app/components/LoadingSpinner.tsx";
import { AlertTriangle, UserPlus, UsersIcon } from "lucide-react";
import Button from "@app/components/Button.tsx";
import { Link } from "react-router-dom";

export default function FriendRequestsPage() {
  const [search, setSearch] = React.useState("");

  return (
    <>
      <Container className="pt-4">
        <DCSearchBar value={search} setValue={setSearch} />
      </Container>
      <FriendRequestList searchQuery={search} />
    </>
  );
}

function FriendRequestList({ searchQuery }: { searchQuery: string }) {
  const {
    data: requests,
    isLoading,
    error,
    showLoading,
  } = usePendingFriendRequests();

  if (error) {
    return (
      <Container className="flex flex-col items-center pt-20">
        <AlertTriangle size={48} className="text-dc-red-500 mb-2" />
        <h1 className="text-xl text-dc-red-500">
          Error getting your friend requests
        </h1>
        <p className="text-dc-neutral-300">{error.message}</p>
      </Container>
    );
  }

  if (isLoading && showLoading) {
    return (
      <Container className="flex items-center p-20 flex-col">
        <LoadingSpinner size="lg" className="mb-2" />
        <h1 className="text-dc-neutral-300 text-xl">
          Loading your friend requests...
        </h1>
      </Container>
    );
  }

  const filteredRequests =
    requests?.filter((request) =>
      request.user.username.toLowerCase().includes(searchQuery.toLowerCase()),
    ) ?? [];

  if (!filteredRequests.length) {
    return (
      <Container className="flex items-center p-20 h-full flex-col max-w-md mx-auto text-center">
        <UsersIcon size={48} className="text-dc-neutral-300" />
        <h1 className="text-dc-neutral-300 text-xl font-medium">
          No friend requests
        </h1>
        <p className="text-dc-neutral-300 pb-4">
          Waiting for a friend request that is not arriving? You can just send
          it yourself here!
        </p>
        <Link to="/home/friends/invite">
          <Button variant="success" className="flex gap-2 items-center">
            <span>Add friend</span> <UserPlus size={20} />
          </Button>
        </Link>
      </Container>
    );
  }

  return (
    <>
      <Container className="py-4">
        <h1 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Waiting - {filteredRequests.length}
        </h1>
      </Container>
      <Container className="pb-4 overflow-auto">
        {filteredRequests.map((request) => (
          <FriendRequestItem
            id={request.id}
            username={request.user.username}
            key={request.id}
          />
        ))}
      </Container>
    </>
  );
}
