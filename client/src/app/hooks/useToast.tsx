import toast, { type Toast } from 'react-hot-toast';
import DCToast from '../components/Toast';

export enum ToastDuration {
  short = 1500,
  medium = 2000,
  long = 4000,
}

type ToastOptions = Partial<
  Pick<
    Toast,
    | 'id'
    | 'style'
    | 'className'
    | 'icon'
    | 'position'
    | 'ariaProps'
    | 'iconTheme'
  > & { duration: ToastDuration }
>;

export function useToast() {
  function error(message: string, options?: ToastOptions) {
    toast.custom(<DCToast message={message} variant='error' />, options);
  }

  function info(message: string, options?: ToastOptions) {
    toast.custom(<DCToast message={message} variant='info' />, options);
  }

  function success(message: string, options?: ToastOptions) {
    toast.custom(<DCToast message={message} variant='success' />, options);
  }

  return {
    error,
    info,
    success,
  };
}
