export { api } from "./api.ts";
export type {
  LoginUserBodyType,
  LoginUserResponse,
  RegisterUserResponseType,
  RegisterUserBodyType,
  GetLoggedInUserResponse,
} from "./types/auth.ts";
export type { ErrorResponse, SuccessMessageResponse } from "./types/default.ts";
export type { UserResponse } from "./types/user.ts";
export type { PendingFriendRequestsResponse } from "./types/friends.ts";

export { QueryKeys } from "./queryKeys.ts";
