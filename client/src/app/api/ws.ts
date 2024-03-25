import React from "react";
import { useToast } from "@app/hooks/useToast.tsx";
import { useAuth } from "@app/context/AuthProvider.tsx";
import { UpdateAuthTokenSchema } from "@app/api/wstypes/auth.ts";

type UseWebsocketParams = {
  path: string;
  onMessage?: (ev: MessageEvent<string>) => void;
  onClose?: (ev: CloseEvent) => void;
  onOpen?: (ev: Event) => void;
};

function getWsUrl(path: string, accessToken: string, refreshToken: string) {
  const url = new URL(`${import.meta.env.VITE_WS_URL}${path}`);
  url.searchParams.set("accessToken", accessToken);
  url.searchParams.set("refreshToken", refreshToken);
  return url.href;
}

function getData(data: string) {
  try {
    return JSON.stringify(data);
  } catch (err) {
    return null;
  }
}

export function useWebsocket({
  path,
  onClose,
  onMessage,
  onOpen,
}: UseWebsocketParams) {
  const toast = useToast();
  const { accessToken, refreshToken, setAccessToken, setRefreshToken } =
    useAuth();
  const wsRef = React.useRef<WebSocket | null>(null);
  React.useEffect(() => {
    function connect() {
      const websocket = new WebSocket(
        getWsUrl(path, accessToken, refreshToken),
      );
      if (!websocket) {
        toast.error("Could not connect to websocket");
        return null;
      }
      return websocket;
    }

    function handleMessage(event: MessageEvent<string>) {
      const data = getData(event.data);
      if (data != null) {
        const updateTokenResult = UpdateAuthTokenSchema.safeParse(data);
        if (updateTokenResult.success) {
          const { accessToken, refreshToken } = updateTokenResult.data;
          setAccessToken(accessToken);
          setRefreshToken(refreshToken);
        }
      }
      onMessage?.(event);
    }
    const ws = connect();
    if (!ws) return;
    wsRef.current = ws;
    wsRef.current.addEventListener("message", handleMessage);
    if (onClose) {
      wsRef.current.addEventListener("close", onClose);
    }
    if (onOpen) {
      wsRef.current.addEventListener("open", onOpen);
    }
    return () => {
      wsRef.current?.removeEventListener("message", handleMessage);
      if (onClose) {
        wsRef.current?.removeEventListener("close", onClose);
      }
      if (onOpen) {
        wsRef.current?.removeEventListener("open", onOpen);
      }
      wsRef.current?.close();
    };
  }, [accessToken, path, refreshToken]);
}
