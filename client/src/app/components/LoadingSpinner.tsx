import { Loader2 } from "lucide-react";
import { cn } from "@app/utils/cn.ts";

export function LoadingSpinner({
  className,
  size = "md",
}: {
  size?: "sm" | "md" | "lg" | "xl";
  className?: string;
}) {
  const sizes = {
    sm: 24,
    md: 36,
    lg: 48,
    xl: 64,
  };

  return (
    <Loader2
      size={sizes[size]}
      className={cn("animate-spin text-dc-neutral-300", className)}
    />
  );
}
