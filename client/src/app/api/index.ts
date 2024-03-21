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
export type {
  PendingFriendsResponse,
  AllFriendResponse,
} from "./types/friends.ts";

export { QueryKeys } from "./queryKeys.ts";

export { useLogin } from "./hooks/useLogin";
export { useRegister } from "./hooks/useRegister";
export { usePendingFriendRequests } from "./hooks/usePendingFriendRequests";
export { useLogout } from "./hooks/useLogout";
export { useChats } from "./hooks/useChats";
export { useAllFriends } from "./hooks/useAllFriends.ts";
