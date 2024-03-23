import { useQuery } from "@tanstack/react-query";
import { api, QueryKeys, GetAllChats } from "@app/api/index.ts";
import { useDeferredFlag } from "@app/hooks/useDeferredFlag.ts";
import { ClientError } from "@app/utils/clientError.ts";

export function useChats() {
  const showLoading = useDeferredFlag(100);
  const query = useQuery<GetAllChats, ClientError>({
    queryKey: QueryKeys.getAllChats(),
    queryFn: async () => api.get<GetAllChats>("/chats"),
  });

  return { ...query, showLoading };
}
