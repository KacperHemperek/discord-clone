import { z } from "zod";

export enum NotificationType {
  friendRequest = "friend_request",
  newMessage = "new_message",
}

const BaseNotificationFields = {
  id: z.number(),
  seen: z.boolean(),
  userId: z.number(),
  createdAt: z.string().datetime(),
  updatedAt: z.string().datetime(),
};

export const FriendRequestNotificationSchema = z.object({
  ...BaseNotificationFields,
  type: z.literal(NotificationType.friendRequest),
  data: z.any(),
});

export const NewMessageNotificationSchema = z.object({
  type: z.literal(NotificationType.newMessage),
  ...BaseNotificationFields,
  data: z.object({
    chatId: z.number(),
  }),
});

export type NewMessageNotification = z.infer<
  typeof NewMessageNotificationSchema
>;

export type FriendRequestNotification = z.infer<
  typeof FriendRequestNotificationSchema
>;
