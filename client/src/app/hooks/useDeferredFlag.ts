import React from "react";

export function useDeferredFlag(ms: number) {
  const [value, setValue] = React.useState(false);
  const timeout = React.useRef<NodeJS.Timeout | null>(null);

  React.useEffect(() => {
    timeout.current = setTimeout(() => {
      setValue(true);
    }, ms);

    return () => {
      if (timeout.current) {
        clearTimeout(timeout.current);
      }
    };
  }, [ms]);

  return value;
}
