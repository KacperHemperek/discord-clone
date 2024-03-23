import { useQuery } from "@tanstack/react-query";
import { api, PendingFriendsResponse, QueryKeys } from "@app/api";
import { useDeferredFlag } from "@app/hooks/useDeferredFlag.ts";

export function usePendingFriendRequests() {
  const showLoading = useDeferredFlag(100);
  const query = useQuery({
    queryKey: QueryKeys.getPendingFriendRequests(),
    queryFn: async () => {
      const response =
        await api.get<PendingFriendsResponse>(`/friends/requests`);

      console.log({ response });

      return response.requests;
    },
  });

  return { ...query, showLoading };
}
