package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

// ClientConnectionWrapper
// TODO
// - add reconnection
// - add channel
type ClientConnectionWrapper struct {
	clientId string
	conn     *websocket.Conn
}

func NewClientConnectionWrapper(
	clientId string,
	url string,
	handleMessageFunc func(client *ClientConnectionWrapper, message Message),
) (*ClientConnectionWrapper, error) {
	header := http.Header{}
	header.Set("ClientId", clientId)
	conn, response, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		return nil, logger.NewError("%s can't establish ws connection because %v", clientId, err.Error())
	}

	logger.Debug("%s ws connection response status = %s", clientId, response.Status)
	wrapper := &ClientConnectionWrapper{
		clientId: clientId,
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
			// TODO fix this :/
			handleMessageFunc(wrapper, messageAsStruct)
		}
	}()
	return wrapper, nil
}

func (client *ClientConnectionWrapper) SendMessage(message Message) {
	messageBytes, _ := json.Marshal(message)
	err := client.conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		logger.Error("%v can't send ws message to %v because %v", client.clientId, client.conn.RemoteAddr().String(), err.Error())
	}
}

func (client *ClientConnectionWrapper) Close() {
	address := client.conn.RemoteAddr().String()
	logger.Debug("close '%v' ws to %v", client.clientId, address)
	err := client.conn.Close()
	if err != nil {
		logger.Error("can't close '%v' ws to %v because '%v'", client.clientId, address, err.Error())
	}
}
