import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, LoginUserBodyType, LoginUserResponse } from "@app/api";
import { MutationHookOptions } from "@app/types/utils";
import { useNavigate } from "react-router-dom";

type AuthLoginMutationOptions = MutationHookOptions<
  LoginUserResponse["user"],
  Error,
  LoginUserBodyType
>;

export function useLogin(options?: AuthLoginMutationOptions) {
  const navigate = useNavigate();

  const queryClient = useQueryClient();

  return useMutation({
    ...options,
    mutationFn: async (data) => {
      const json = await api.post<LoginUserResponse>("/auth/login", {
        body: JSON.stringify(data),
      });

      return json.user;
    },
    onSuccess: (data, variables, context) => {
      queryClient.setQueryData(["user"], data);
      navigate("/home/friends/");
      options?.onSuccess?.(data, variables, context);
    },
  });
}
