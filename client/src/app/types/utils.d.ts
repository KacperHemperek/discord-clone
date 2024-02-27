import { DefaultError } from '@tanstack/react-query';
import type { UseMutationOptions } from '@tanstack/react-query';

export type MutationHookOptions<
  TData = unknown,
  TError = DefaultError,
  TVariables = void,
  TContext = unknown,
> = Omit<UseMutationOptions<TData, TError, TVariables, TContext>, 'mutationFn'>;

// @eslint-disable-next-line @typescript-eslint/no-explicit-any
type TODO = any;
