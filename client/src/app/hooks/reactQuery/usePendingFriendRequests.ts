import { useQuery } from "@tanstack/react-query";
import { api, PendingFriendRequestsResponse, QueryKeys } from "@app/api";

export function usePendingFriendRequests() {
  return useQuery({
    queryKey: QueryKeys.getPendingFriendRequests(),
    queryFn: async () => {
      const response =
        await api.get<PendingFriendRequestsResponse>(`/friends/requests`);

      console.log({ response });

      return response.requests;
    },
  });
}
