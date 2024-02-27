import { useNavigate } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { AuthTypes } from '@discord-clone-v2/types';
import { MutationHookOptions } from '../../types/utils';
import { api } from '../../utils/api';

type RegisterMutationOptions = MutationHookOptions<
  AuthTypes.RegisterUserCreatedResponseType['user'],
  Error,
  AuthTypes.RegisterUserBodyType
>;

export function useRegister(options?: RegisterMutationOptions) {
  const navigate = useNavigate();

  const queryClient = useQueryClient();

  return useMutation({
    ...options,
    mutationFn: async (data) => {
      const json = await api.post<AuthTypes.RegisterUserCreatedResponseType>(
        '/auth/register',
        {
          body: JSON.stringify(data),
        },
      );

      return json.user;
    },
    onSuccess: (data, variables, context) => {
      queryClient.setQueryData(['user'], data);
      navigate('/home/friends/');
      options?.onSuccess?.(data, variables, context);
    },
  });
}
