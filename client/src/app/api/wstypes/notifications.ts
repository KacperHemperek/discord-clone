import { z } from "zod";

export enum NotificationType {
  friendRequest = "friend_request",
  newMessage = "new_message",
}

export const FriendRequestNotificationSchema = z.object({
  type: z.literal(NotificationType.friendRequest),
  id: z.number(),
  seen: z.boolean(),
  userId: z.number(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
  data: z.any(),
});

export type FriendRequestNotification = z.infer<
  typeof FriendRequestNotificationSchema
>;
