import { UserResponse } from "@app/api";

export type PendingFriendRequestsResponse = {
  requests: Array<UserResponse>;
};
