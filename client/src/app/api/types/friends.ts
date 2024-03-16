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

export type PendingFriendRequestsResponse = {
  requests: Array<FriendRequest>;
};
