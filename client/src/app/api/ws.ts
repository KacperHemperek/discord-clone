import { useToast } from "@app/hooks/useToast.tsx";
import { useAuth } from "@app/context/AuthProvider.tsx";
import { UpdateAuthTokenSchema } from "@app/api/wstypes/auth.ts";

export function useWebsocket() {
  const toast = useToast();
  const { accessToken, refreshToken, setAccessToken, setRefreshToken } =
    useAuth();

  function getWsUrl(path: string, accessToken: string, refreshToken: string) {
    const url = new URL(`${import.meta.env.VITE_WS_URL}${path}`);
    url.searchParams.set("accessToken", accessToken);
    url.searchParams.set("refreshToken", refreshToken);
    return url.href;
  }

  function getData(data: string) {
    try {
      return JSON.parse(data);
    } catch (err) {
      return null;
    }
  }

  function connect(
    path: string,
    errMessage: string = "Could not connect to websocket",
  ) {
    const websocket = new WebSocket(getWsUrl(path, accessToken, refreshToken));
    if (!websocket) {
      toast.error(errMessage);
      return null;
    }
    return websocket;
  }

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
