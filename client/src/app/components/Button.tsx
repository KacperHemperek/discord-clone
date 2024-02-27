import { VariantProps, cva } from 'class-variance-authority';
import React from 'react';

const button = cva(
  'rounded-sm duration-150 font-medium disabled:opacity-60 disabled:cursor-not-allowed',
  {
    variants: {
      variant: {
        primary:
          'bg-dc-purple-500 text-dc-neutral-50 hover:opacity-80 transition-opacity',
        success:
          'bg-dc-green-500 text-dc-neutral-50 hover:opacity-80 transition-opacity',
        danger:
          'bg-dc-red-500 text-dc-neutral-50 hover:opacity-80 transition-opacity',
        link: 'bg-transparent text-dc-neutral-50 hover:underline',
      },
      size: {
        sm: 'px-2 py-1',
        md: 'px-3 py-2',
        lg: 'px-4 py-3',
      },
      fontSize: {
        sm: 'text-sm',
        md: 'text-md',
      },
      fontWeight: {
        normal: 'font-normal',
        medium: 'font-medium',
        semibold: 'font-semibold',
      },
    },
    defaultVariants: {
      variant: 'primary',
      size: 'md',
      fontSize: 'sm',
      fontWeight: 'medium',
    },
  },
);

type ButtonProps = VariantProps<typeof button> & {
  children: React.ReactNode;
} & React.ButtonHTMLAttributes<HTMLButtonElement>;

const DCButton = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    { children, size, variant, fontSize, fontWeight, className, ...rest },
    ref,
  ) => {
    return (
      <button
        ref={ref}
        className={button({ size, variant, fontSize, fontWeight, className })}
        {...rest}
      >
        {children}
      </button>
    );
  },
);

export default DCButton;
