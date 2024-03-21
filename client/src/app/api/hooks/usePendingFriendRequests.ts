import { useQuery } from "@tanstack/react-query";
import { api, PendingFriendsResponse, QueryKeys } from "@app/api";

export function usePendingFriendRequests() {
  return useQuery({
    queryKey: QueryKeys.getPendingFriendRequests(),
    queryFn: async () => {
      const response =
        await api.get<PendingFriendsResponse>(`/friends/requests`);

      console.log({ response });

      return response.requests;
    },
  });
}
