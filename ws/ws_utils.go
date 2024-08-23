package ws

import (
	"encoding/json"
	"github.com/dmitriibb/go-common/logging"
	"github.com/gorilla/websocket"
)

var logger = logging.NewLogger("commonWsUtils")

func HandleWsConnectionMessagesWrapper(clientId string, conn *websocket.Conn, handlerFunc func(message Message)) {
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
		handlerFunc(messageAsStruct)
	}
}
