import { UserResponse } from "@app/api";

export enum FriendRequestStatus {
  PENDING = "pending",
  ACCEPTED = "accepted",
  REJECTED = "rejected",
}

export type FriendRequest = {
  id: number;
  status: FriendRequestStatus;
  requestedAt: string;
  statusChangedAt: string;
  user: UserResponse;
};

export type PendingFriendsResponse = {
  requests: Array<FriendRequest>;
};

export type AllFriendResponse = {
  friends: Array<UserResponse>;
};
