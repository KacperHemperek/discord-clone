import { useMutation, useQueryClient } from "@tanstack/react-query";
import { api, SuccessMessageResponse } from "@app/api";
import { MutationHookOptions } from "@app/types/utils";
import { useNavigate } from "react-router-dom";

type LogoutMutationOptions = MutationHookOptions;

export function useLogout(options?: LogoutMutationOptions) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    ...options,
    mutationFn: () => api.post<SuccessMessageResponse>("/auth/logout"),
    onSuccess: () => {
      queryClient.setQueryData(["user"], null);
      navigate("/login");
    },
  });
}
