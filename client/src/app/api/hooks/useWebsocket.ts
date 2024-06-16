import { useToast } from "@app/hooks/useToast.tsx";
import { useAuth } from "@app/context/AuthProvider.tsx";
import { UpdateAuthTokenSchema } from "@app/api/wstypes/auth.ts";
import React from "react";

export function useWebsocket() {
  const toast = useToast();
  const { setAccessToken, setRefreshToken } = useAuth();

  const getWsUrl = React.useCallback(
    (path: string, accessToken: string, refreshToken: string) => {
      const url = new URL(`${import.meta.env.VITE_WS_URL}${path}`);
      url.searchParams.set("accessToken", accessToken);
      url.searchParams.set("refreshToken", refreshToken);
      return url.href;
    },
    [],
  );

  const getData = React.useCallback((data: string) => {
    try {
      return JSON.parse(data);
    } catch (err) {
      return null;
    }
  }, []);

  const connect = React.useCallback(
    ({
      refreshToken,
      accessToken,
      path,
      errMessage = "Could not connect to websocket",
    }: {
      path: string;
      accessToken: string;
      refreshToken: string;
      errMessage: string;
    }) => {
      if (accessToken.length === 0 || refreshToken.length === 0) {
        return null;
      }
      const websocket = new WebSocket(
        getWsUrl(path, accessToken, refreshToken),
      );
      if (!websocket) {
        toast.error(errMessage);
        return null;
      }
      return websocket;
    },
    [getWsUrl, toast],
  );

  function handleMessage(
    cb: (data: unknown, eventRaw?: MessageEvent<string>) => void,
  ): (event: MessageEvent<string>) => void {
    return function (event) {
      const data = getData(event.data);
      if (data != null) {
        const updateTokenResult = UpdateAuthTokenSchema.safeParse(data);
        if (updateTokenResult.success) {
          const { accessToken, refreshToken } = updateTokenResult.data;
          setAccessToken(accessToken);
          setRefreshToken(refreshToken);
          return;
        }
      }
      cb(data, event);
    };
  }

  return {
    handleMessage,
    connect,
  };
}
