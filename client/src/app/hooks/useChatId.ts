import { useParams } from "react-router-dom";

export function useChatId() {
  const params = useParams();
  const chatId = Number(params.chatId);

  if (!chatId || Number.isNaN(chatId)) {
    throw new Error("Chat id is invalid or not present");
  }
  return chatId;
}
