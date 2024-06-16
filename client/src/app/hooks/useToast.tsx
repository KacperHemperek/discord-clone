import toast, { type Toast } from "react-hot-toast";
import DCToast from "../components/Toast";
import React from "react";

export enum ToastDuration {
  short = 1500,
  medium = 2000,
  long = 4000,
}

type ToastOptions = Partial<
  Pick<
    Toast,
    | "id"
    | "style"
    | "className"
    | "icon"
    | "position"
    | "ariaProps"
    | "iconTheme"
  > & { duration: ToastDuration }
>;

export function useToast() {
  const error = React.useCallback((message: string, options?: ToastOptions) => {
    toast.custom(<DCToast message={message} variant="error" />, options);
  }, []);

  const info = React.useCallback((message: string, options?: ToastOptions) => {
    toast.custom(<DCToast message={message} variant="info" />, options);
  }, []);

  const success = React.useCallback(
    (message: string, options?: ToastOptions) => {
      toast.custom(<DCToast message={message} variant="success" />, options);
    },
    [],
  );

  return React.useMemo(
    () => ({
      error,
      info,
      success,
    }),
    [error, info, success],
  );
}
