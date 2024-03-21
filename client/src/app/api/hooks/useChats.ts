import { useQuery } from "@tanstack/react-query";
import { api, QueryKeys } from "../index.ts";

export function useChats() {
  return useQuery({
    queryKey: QueryKeys.getAllChats(),
    queryFn: async () =>
      api.get<ChatsTypes.GetChatsSuccessResponseType>("/chats"),
  });
}
