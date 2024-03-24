package middlewares

import (
	"github.com/gorilla/websocket"
	"github.com/kacperhemperek/discord-go/utils"
	"net/http"
)

type ConnectWsMiddleware = func(h utils.APIHandler) utils.APIHandler

func NewConnectWsMiddleware() ConnectWsMiddleware {
	return func(h utils.APIHandler) utils.APIHandler {
		return func(w http.ResponseWriter, r *http.Request, c *utils.Context) error {
			upgrader := websocket.Upgrader{
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			}

			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return err
			}
			c.Conn = conn
			err = h(w, r, c)
			if err != nil {
				_ = c.Conn.Close()
			}
			return err
		}
	}
}
