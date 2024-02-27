/* package created 12.09.2023. for Mobilisis d.o.o. by Elena Kržina */
/* used to connect to specific websocket, listen on the websocket and post to the websocket */

package websocket

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

/* opens websocket by creating a url object with the provided scheme, host and path */
func OpenWebsocket(scheme, host, path string) (*websocket.Conn, error) {

	serverUrl := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket: ", err)
	}

	/* to close connection write defer conn.Close() in calling package */
	return conn, nil
}

/* listens on the provided websocket connection; use after opening the websocket and open in separate goroutine*/
/* implement with an infinite loop */
func ListenOnWS(conn *websocket.Conn) ([]byte, error) {

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message: ", err)
		return nil, err
	}
	/* returns bytes from the websocket */
	return message, nil

}

/* sends a message to the provided websocket connection; use after opening the websocket */
func SendToWS(conn *websocket.Conn, message []byte) error {
	err := conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Println("Error sending message: ", err)
		return err
	}
	/* if the message is sent, the function does not return anything */
	fmt.Println("Message received.")
	return nil
}
