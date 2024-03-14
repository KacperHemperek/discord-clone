package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WS struct {
	OnMessage func(conn *websocket.Conn, msg []byte)
	OnClose   func(conn *websocket.Conn)
}

func NewWS() *WS {
	return &WS{}
}

func (ws *WS) Serve(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			fmt.Println("error when closing connection ", err)
		}
	}()

	for {
		mType, msg, err := conn.ReadMessage()

		if err != nil {
			return
		}

		switch mType {
		case websocket.TextMessage:
			ws.OnMessage(conn, msg)
		case websocket.CloseMessage:
			ws.OnClose(conn)
			return
		default:
			fmt.Printf("unknown message type %d", mType)
		}
	}
}

type HandlerContext struct {
	Conn *websocket.Conn
	Msg  []byte
}

type Handler = func(w http.ResponseWriter, r *http.Request, ctx HandlerContext)
