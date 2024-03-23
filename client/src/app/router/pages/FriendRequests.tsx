import React from "react";
import FriendRequestItem from "@app/components/friends/FriendRequestItem";
import { Container } from "@app/components/friends/FriendPageContainer";
import DCSearchBar from "@app/components/SearchBar";
import { PendingFriendsResponse, usePendingFriendRequests } from "@app/api";
import { LoadingSpinner } from "@app/components/LoadingSpinner.tsx";
import { UserPlus, UsersIcon } from "lucide-react";
import Button from "@app/components/Button.tsx";
import { Link } from "react-router-dom";
import { ErrorPageWithRetry } from "@app/components/ErrorPageWithRetry.tsx";
import { ClientError } from "@app/utils/clientError.ts";

export default function FriendRequestsPage() {
  const {
    data: requests,
    isLoading,
    error,
    showLoading,
    refetch,
  } = usePendingFriendRequests();
  const [search, setSearch] = React.useState("");

  return (
    <>
      <Container className="pt-4">
        <DCSearchBar value={search} setValue={setSearch} />
      </Container>
      <Container className="pt-4">
        <h1 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Waiting - {requests?.length ?? 0}
        </h1>
      </Container>
      <FriendRequestList
        search={search}
        requests={requests}
        error={error}
        onRetry={refetch}
        isLoading={isLoading}
        showLoading={showLoading}
      />
    </>
  );
}

function FriendRequestList({
  search,
  onRetry,
  isLoading,
  showLoading,
  error,
  requests,
}: {
  search: string;
  isLoading: boolean;
  error: ClientError | null;
  showLoading: boolean;
  requests: PendingFriendsResponse["requests"] | undefined;
  onRetry: () => void;
}) {
  const filteredRequests =
    requests?.filter((request) =>
      request.user.username.toLowerCase().includes(search.toLowerCase()),
    ) ?? [];

  if (error) {
    return (
      <Container className="flex flex-col items-center pt-20">
        <ErrorPageWithRetry
          error={error}
          retry={onRetry}
          defaultErrorMessage="Could not retrieve your friend requests, you can try again later"
        />
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

  if (!requests?.length) {
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
    <Container className="pb-4 overflow-auto">
      {filteredRequests.map((request) => (
        <FriendRequestItem
          id={request.id}
          username={request.user.username}
          key={request.id}
        />
      ))}
    </Container>
  );
}
