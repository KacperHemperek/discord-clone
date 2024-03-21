import React from "react";
import FriendListItem from "../../components/friends/FriendListItem";
import { Container } from "../../components/friends/FriendPageContainer";
import DCSearchBar from "../../components/SearchBar";
import { useAllFriends } from "@app/api";

export default function AllFriendsPage() {
  const { data } = useAllFriends();

  const [search, setSearch] = React.useState("");

  const filteredFriends = data?.friends.filter((friend) =>
    friend.username.toLowerCase().includes(search.toLowerCase()),
  );

  return (
    <>
      <Container className="pt-4">
        <DCSearchBar value={search} setValue={setSearch} />
      </Container>
      <Container className="pt-4">
        <h1 className="uppercase text-xs font-semibold tracking-[0.02em] text-dc-neutral-300">
          Friends - {filteredFriends?.length ?? 0}
        </h1>
      </Container>
      <Container className="py-4 overflow-auto">
        {filteredFriends?.map((friend) => (
          <FriendListItem
            id={friend.id}
            username={friend.username}
            key={friend.id}
          />
        ))}
      </Container>
    </>
  );
}
