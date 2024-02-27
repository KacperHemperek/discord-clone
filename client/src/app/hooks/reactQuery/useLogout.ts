import { useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '../../utils/api';
import { MutationHookOptions } from '../../types/utils';
import { useNavigate } from 'react-router-dom';
import { CommonResponsesTypes } from '@discord-clone-v2/types';

type LogoutMutationOptions = MutationHookOptions;

export function useLogout(options?: LogoutMutationOptions) {
  const queryClient = useQueryClient();
  const navigate = useNavigate();

  return useMutation({
    ...options,
    mutationFn: () =>
      api.post<CommonResponsesTypes.MessageSuccessResponseType>('/auth/logout'),
    onSuccess: () => {
      queryClient.setQueryData(['user'], null);
      navigate('/login');
    },
  });
}
