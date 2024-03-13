export type UserResponse = {
  id: number;
  username: string;
  email: string;
  active: boolean;
  createdAt: string;
  updatedAt: string;
};

export type LoginUserBodyType = {
  email: string;
  password: string;
};

export type LoginUserResponse = {
  message: string;
  user: UserResponse;
};

export type RegisterUserBodyType = {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
};

export type RegisterUserResponseType = {
  message: string;
  user: UserResponse;
};
