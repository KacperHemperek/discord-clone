import { useQuery } from "@tanstack/react-query";
import { api, PendingFriendsResponse, QueryKeys } from "@app/api";
import { useDeferredFlag } from "@app/hooks/useDeferredFlag.ts";
import { ClientError } from "@app/utils/clientError.ts";

export function usePendingFriendRequests() {
  const showLoading = useDeferredFlag(100);
  const query = useQuery<PendingFriendsResponse["requests"], ClientError>({
    queryKey: QueryKeys.getPendingFriendRequests(),
    queryFn: async () => {
      const response =
        await api.get<PendingFriendsResponse>(`/friends/requests`);

      return response.requests;
    },
  });

  return { ...query, showLoading };
}
