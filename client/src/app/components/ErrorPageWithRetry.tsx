import { AlertTriangleIcon, RefreshCwIcon } from "lucide-react";
import Button from "@app/components/Button.tsx";
import { ClientError } from "@app/utils/clientError.ts";
import { cn } from "@app/utils/cn.ts";

export function ErrorPageWithRetry({
  retry,
  error,
  defaultErrorMessage,
  className,
}: {
  retry: () => void;
  error: ClientError | Error;
  defaultErrorMessage: string;
  className?: string;
}) {
  const isClientError = error instanceof ClientError;

  return (
    <div
      className={cn(
        "flex flex-col items-center max-w-md text-center",
        className,
      )}
    >
      <AlertTriangleIcon size={48} className="text-dc-red-500" />
      <h3 className="text-dc-red-500 text-xl font-medium">
        Something went wrong!
      </h3>
      <p className="text-dc-neutral-300 pb-2">
        {isClientError
          ? error.code === 500
            ? defaultErrorMessage
            : error.message
          : error.message}
      </p>
      <Button
        variant="primary"
        className="flex gap-2 items-center"
        onClick={retry}
      >
        <span>Try again</span> <RefreshCwIcon size={16} />
      </Button>
    </div>
  );
}
