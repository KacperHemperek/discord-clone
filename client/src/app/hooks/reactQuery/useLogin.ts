import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "../../api";

import type { AuthTypes } from "@discord-clone-v2/types";
import { MutationHookOptions } from "../../types/utils";
import { useNavigate } from "react-router-dom";

type AuthLoginMutationOptions = MutationHookOptions<
  AuthTypes.LoginUserSuccessfullyResponseType["user"],
  Error,
  AuthTypes.LoginUserBodyType
>;

export function useLogin(options?: AuthLoginMutationOptions) {
  const navigate = useNavigate();

  const queryClient = useQueryClient();

  return useMutation({
    ...options,
    mutationFn: async (data) => {
      const json = await api.post<AuthTypes.LoginUserSuccessfullyResponseType>(
        "/auth/login",
        {
          body: JSON.stringify(data),
        },
      );

      return json.user;
    },
    onSuccess: (data, variables, context) => {
      queryClient.setQueryData(["user"], data);
      navigate("/home/friends/");
      options?.onSuccess?.(data, variables, context);
    },
  });
}
