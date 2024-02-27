import { cn } from '../../utils/cn';

export function Container({
  children,
  className,
}: {
  children?: React.ReactNode;
  className?: string;
}) {
  return <div className={cn('px-8', className)}>{children}</div>;
}
