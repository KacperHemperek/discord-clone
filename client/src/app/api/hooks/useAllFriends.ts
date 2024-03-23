import { useQuery } from "@tanstack/react-query";
import { AllFriendResponse, api, QueryKeys } from "../index.ts";
import { useDeferredFlag } from "@app/hooks/useDeferredFlag.ts";
import { ClientError } from "@app/utils/clientError.ts";

export function useAllFriends() {
  const showLoading = useDeferredFlag(100);
  const query = useQuery<AllFriendResponse, ClientError>({
    queryKey: QueryKeys.getAllFriends(),
    queryFn: async () => api.get<AllFriendResponse>(`/friends`),
  });
  return { ...query, showLoading };
}
