export class QueryKeys {
  static getPendingFriendRequests() {
    return ["friend-requests"];
  }

  static getLoggedInUser() {
    return ["user"];
  }

  static getAllFriends() {
    return ["all-friends"];
  }

  static getAllChats() {
    return ["chats"];
  }

  static getChat(chatId: number) {
    return ["chat", chatId];
  }

  static getFriendRequestNotifications() {
    return ["friend-request-notifications"];
  }
}
