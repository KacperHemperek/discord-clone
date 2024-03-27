import { z } from "zod";
import { WsMessages } from "@app/api/wstypes/messages.ts";

export const NewMessageWsSchema = z.object({
  type: z.literal(WsMessages.newMessage),
  message: z.object({
    id: z.number(),
    text: z.string(),
    createdAt: z.string(),
    updatedAt: z.string(),
    user: z.object({
      id: z.number(),
      username: z.string(),
      email: z.string(),
      active: z.boolean(),
      createdAt: z.string(),
      updatedAt: z.string(),
    }),
  }),
});

export type NewMessageWsType = z.infer<typeof NewMessageWsSchema>;
