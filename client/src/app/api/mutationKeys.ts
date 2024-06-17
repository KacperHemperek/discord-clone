export class MutationKeys {
  static markFriendRequestNotificationAsSeen() {
    return ["mark-friend-request-notifications-as-seen"];
  }

  static markNewMessageNotificationsAsSeen(chatId: number) {
    return ["mark-new-message-notifications-as-seen", chatId];
  }
}
