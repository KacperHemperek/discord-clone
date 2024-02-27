export function getWebsocketConnection(path: string) {
  try {
    const wsUrl = `${import.meta.env.VITE_WS_URL}${path}`;

    return new WebSocket(wsUrl);
  } catch (err) {
    console.error(err);
    return null;
  }
}
