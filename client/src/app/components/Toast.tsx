import { CheckCircle, Info, XCircle } from 'lucide-react';
import React from 'react';

export type ToastVariant = 'error' | 'success' | 'info';

export default function DCToast({
  message,
  variant,
}: {
  variant: ToastVariant;
  message: string;
}) {
  const size = 20;

  const icons: Record<ToastVariant, React.ReactNode> = {
    error: <XCircle size={size} className='text-dc-red-500' />,
    info: <Info size={size} className='text-sky-500' />,
    success: <CheckCircle size={size} className='text-dc-green-500' />,
  };

  return (
    <div className='bg-dc-neutral-1000 rounded-md shadow-md flex flex-col p-4 max-w-xs w-full'>
      <div className='flex gap-2 items-center'>
        {icons[variant]}
        <h1 className='capitalize font-medium text-dc-neutral-50'>{variant}</h1>
      </div>
      <p className='text-dc-neutral-300 pl-7 text-sm'>{message}</p>
    </div>
  );
}
