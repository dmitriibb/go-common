package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientHandlerStatus string

const (
	ClientHandlerStatusCreated    ClientHandlerStatus = "created"
	ClientHandlerStatusAuthorized ClientHandlerStatus = "authorized"
)

// ClientHandler
// TODO
// - add channel
type ClientHandler struct {
	ClientId string
	conn     *websocket.Conn
	status   ClientHandlerStatus
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

	clientHandler := &ClientHandler{
		ClientId: "",
		conn:     conn,
		status:   ClientHandlerStatusCreated,
	}

	authorizedClientChan := make(chan *ClientHandler)

	clientId := r.Header.Get("ClientId")
	if len(clientId) != 0 {
		clientHandler.ClientId = clientId
		clientHandler.status = ClientHandlerStatusAuthorized
		authorizedClientChan <- clientHandler
	}
	logger.Debug("established new ws connection for '%v'", clientHandler.ClientId)

	go func() {
		defer func() {
			logger.Debug("close ws for '%v'", clientId)
			conn.Close()
		}()
		for {
			_, p, errMsg := conn.ReadMessage()
			if errMsg != nil {
				logger.Error("%s can't read message from ws because %s", clientId, errMsg.Error())
				break
			}
			//logger.Debug("%s received ws message - type %v`", clientId, messageType)
			messageAsStruct := Message{}
			errMsg = json.Unmarshal(p, &messageAsStruct)
			if errMsg != nil {
				logger.Error("can't unmarshal ws message %s", string(p))
				continue
			}
			logger.Debug("%s received message - %+v", clientId, messageAsStruct)
			if clientHandler.status != ClientHandlerStatusAuthorized {
				if messageAsStruct.Type != MessageTypeLogin {
					msg := "ws client must login via LoginMessage first"
					err = logger.NewError(msg)
					clientHandler.conn.WriteMessage(websocket.TextMessage, []byte(msg))
					clientHandler.conn.Close()
					authorizedClientChan <- clientHandler
					break
				} else {
					msgLogin := MessageLogin{}
					errMsg = json.Unmarshal([]byte(messageAsStruct.Payload), &msgLogin)
					if errMsg != nil {
						logger.Error("can't unmarshal ws login message %v", errMsg.Error())
						continue
					} else {
						clientHandler.ClientId = msgLogin.Username
						clientHandler.status = ClientHandlerStatusAuthorized
						authorizedClientChan <- clientHandler
						continue
					}
				}
			}
			// TODO refactor this
			handleMessageFunc(clientHandler, messageAsStruct)
		}
	}()
	// workaround so client will be always authorized and has clientId
	clientHandler = <-authorizedClientChan
	return clientHandler, err
}

func (client *ClientHandler) SendMessage(message Message) {
	messageBytes, _ := json.Marshal(message)
	err := client.conn.WriteMessage(websocket.TextMessage, messageBytes)
	if err != nil {
		logger.Error("%v can't send ws message to %v because %v", client.ClientId, client.conn.RemoteAddr().String(), err.Error())
	}
}
