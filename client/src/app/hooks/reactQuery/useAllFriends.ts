import { FriendsTypes } from "@discord-clone-v2/types";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../api";

export function useAllFriends() {
  return useQuery({
    queryKey: ["all-friends"],
    queryFn: async () => {
      return api.get<FriendsTypes.GetAllFriendsResponseBodyType>(`/friends`);
    },
  });
}
