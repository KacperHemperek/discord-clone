import { useQuery } from "@tanstack/react-query";
import { AllFriendResponse, api, QueryKeys } from "../index.ts";

export function useAllFriends() {
  return useQuery({
    queryKey: QueryKeys.getAllFriends(),
    queryFn: async () => {
      return api.get<AllFriendResponse>(`/friends`);
    },
  });
}
