export type LoginUserBodyType = {
  email: string;
  password: string;
};

export type LoginUserResponse = {
  message: string;
  user: {
    id: number;
    username: string;
    email: string;
    active: boolean;
    createdAt: string;
    updatedAt: string;
  };
};
