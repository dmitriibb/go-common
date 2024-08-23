package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ClientHandler
// TODO
// - add channel
type ClientHandler struct {
	ClientId string
	conn     *websocket.Conn
}

// NewClientHandler must contain client id in request header
// TODO add on close handler
func NewClientHandler(
	w http.ResponseWriter,
	r *http.Request,
	handleMessageFunc func(client *ClientHandler, message Message),
) (client *ClientHandler, err error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, logger.NewError(err.Error())
	}

	clientId := r.Header.Get("ClientId")
	if len(clientId) == 0 {
		return nil, logger.NewError("clientId is empty. Reject ws connection")
	}

	logger.Debug("established new ws connection for %v", clientId)
	clientHandler := &ClientHandler{
		ClientId: clientId,
		conn:     conn,
	}
	go func() {
		defer func() {
			logger.Debug("close ws for %s", clientId)
			conn.Close()
		}()
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				logger.Error("%s can't read message from ws because %s", clientId, err.Error())
				break
			}
			logger.Debug("%s received ws message - type %v`", clientId, messageType)
			messageAsStruct := Message{}
			err = json.Unmarshal(p, &messageAsStruct)
			if err != nil {
				logger.Error("can't unmarshal ws message %s", string(p))
				continue
			}
			logger.Info("%s received message - %+v", clientId, messageAsStruct)
			// TODO refactor this
			handleMessageFunc(clientHandler, messageAsStruct)
		}
	}()
	return clientHandler, nil
}

func (client *ClientHandler) SendMessage(message Message) {
	messageBytes, _ := json.Marshal(message)
	err := client.conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		logger.Error("%v can't send ws message to %v because %v", client.ClientId, client.conn.RemoteAddr().String(), err.Error())
	}
}
