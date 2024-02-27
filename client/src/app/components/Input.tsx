import React from 'react';
import { cn } from '../utils/cn';

type InputProps = Omit<
  React.InputHTMLAttributes<HTMLInputElement>,
  'className'
> & {
  label: string;
  error?: string;
};

const DCInput = React.forwardRef<HTMLInputElement, InputProps>(
  ({ error, label, ...props }, ref) => {
    return (
      <div className='flex flex-col'>
        <label
          htmlFor={props.name}
          className={cn(
            'text-xs font-bold text-dc-neutral-300 pb-3',
            error && 'text-dc-red-400',
          )}
        >
          <span className='uppercase'>{label}</span>
          {error && <span className='font-normal italic'> - {error}</span>}
        </label>
        <input
          ref={ref}
          className='w-full p-2 rounded-sm bg-dc-neutral-1000 disabled:opacity-60 disabled:cursor-not-allowed'
          {...props}
        />
      </div>
    );
  },
);

export default DCInput;
