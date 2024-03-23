import React from "react";
import FriendListItem from "../../components/friends/FriendListItem";
import { Container } from "../../components/friends/FriendPageContainer";
import DCSearchBar from "../../components/SearchBar";
import { AllFriendResponse, useAllFriends } from "@app/api";
import { LoadingSpinner } from "@app/components/LoadingSpinner.tsx";
import { UserPlus, UsersIcon } from "lucide-react";
import { Link } from "react-router-dom";
import Button from "@app/components/Button.tsx";
import { ClientError } from "@app/utils/clientError.ts";
import { ErrorPageWithRetry } from "@app/components/ErrorPageWithRetry.tsx";

export default function AllFriendsPage() {
  const { data, isLoading, error, showLoading, refetch } = useAllFriends();

  const [search, setSearch] = React.useState("");

  return (
    <>
      <Container className="pt-4">
        <DCSearchBar value={search} setValue={setSearch} />
      </Container>
      <Container className="pt-4">
        <h1 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Friends - {data?.friends.length ?? 0}
        </h1>
      </Container>
      <FriendsList
        search={search}
        friends={data?.friends}
        error={error}
        isLoading={isLoading}
        showLoading={showLoading}
        onRetry={refetch}
      />
    </>
  );
}

function FriendsList({
  search,
  isLoading,
  error,
  showLoading,
  friends,
  onRetry,
}: {
  search: string;
  isLoading: boolean;
  error: ClientError | null;
  showLoading: boolean;
  friends: AllFriendResponse["friends"] | undefined;
  onRetry: () => void;
}) {
  const filteredFriends = friends?.filter((friend) =>
    friend.username.toLowerCase().includes(search.toLowerCase()),
  );

  if (error) {
    return (
      <Container className="flex flex-col items-center pt-20">
        <ErrorPageWithRetry
          error={error}
          retry={onRetry}
          defaultErrorMessage="Could not retrieve your friend list, you can try again later"
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

  if (!isLoading && !friends?.length) {
    return (
      <Container className="flex items-center p-20 h-full flex-col max-w-md mx-auto text-center">
        <UsersIcon size={48} className="text-dc-neutral-300" />
        <h1 className="text-dc-neutral-300 text-xl font-medium">
          No friends here, yet!
        </h1>
        <p className="text-dc-neutral-300 pb-4">
          Looking for new friends, or want to talk to someone you know? If you
          know their email just add them!
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
    <Container className="py-4 overflow-auto">
      {filteredFriends?.map((friend) => (
        <FriendListItem
          id={friend.id}
          username={friend.username}
          key={friend.id}
        />
      ))}
    </Container>
  );
}
