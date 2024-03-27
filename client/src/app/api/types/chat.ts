import { UserResponse } from "@app/api";

export enum ChatType {
  PRIVATE = "private",
  GROUP = "group",
}

export type Chat = {
  id: number;
  name: string;
  type: ChatType;
  createdAt: string;
  updatedAt: string;
  members: Array<UserResponse>;
};

export type GetAllChats = {
  chats: Array<Chat>;
};

export type CreateChatResponse = {
  message: string;
  chatId: string;
};

export type Message = {
  id: number;
  text: string;
  createdAt: string;
  updatedAt: string;
  user: UserResponse;
};

export type GetChat = {
  id: number;
  name: string;
  type: ChatType;
  createdAt: string;
  updatedAt: string;
  messages: Array<Message>;
};
