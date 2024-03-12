import { useQuery } from "@tanstack/react-query";
import { api } from "../../api";
import { ChatsTypes } from "@discord-clone-v2/types";

export function useChats() {
  return useQuery({
    queryKey: ["chats"],
    queryFn: async () =>
      api.get<ChatsTypes.GetChatsSuccessResponseType>("/chats"),
  });
}
