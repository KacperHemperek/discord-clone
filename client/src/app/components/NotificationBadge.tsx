import { cn } from "@app/utils/cn.ts";

type NotificationBadgeProps = {
  className?: string;
};

export function NotificationBadge({ className }: NotificationBadgeProps) {
  return (
    <span
      className={cn("p-1.5 rounded-full bg-dc-red-500 absolute", className)}
    />
  );
}
