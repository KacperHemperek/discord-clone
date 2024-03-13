import { useNavigate } from "react-router-dom";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { MutationHookOptions } from "@app/types/utils";
import { api, RegisterUserBodyType, RegisterUserResponseType } from "@app/api";

type RegisterMutationOptions = MutationHookOptions<
  RegisterUserResponseType["user"],
  Error,
  RegisterUserBodyType
>;

export function useRegister(options?: RegisterMutationOptions) {
  const navigate = useNavigate();

  const queryClient = useQueryClient();

  return useMutation({
    ...options,
    mutationFn: async (data) => {
      const json = await api.post<RegisterUserResponseType>("/auth/register", {
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
