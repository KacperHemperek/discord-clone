import { WsMessages } from "@app/api/wstypes/messages.ts";
import { z } from "zod";

export const UpdateAuthTokenSchema = z.object({
  type: z.literal(WsMessages.updateAccessToken),
  accessToken: z.string(),
  refreshToken: z.string(),
});

export type UpdateAuthToken = z.infer<typeof UpdateAuthTokenSchema>;
