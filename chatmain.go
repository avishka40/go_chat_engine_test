package main

import "github.com/gorilla/websocket"
import(
  "net/http"
  "encode/json"
)
type ClientManager struct {

  clients map[*Client]bool
  broadcast chan []byte
  register chan *Client
  unregister chan *Client

}

type Client struct {

  id string
  socket *websocket.Content
  send chan []byte

}

type Message struct {

  Sender string `json:"sender,omitempty"`
  Recipient string `json:"recipient,omitempty"`
  Content string `json:"content,omitempty"`

}

var manager = ClientManager{

  broadcast:make(chan []byte),
  register:make(chan *Client),
  unregister:make(chan *Client),
  clients: make(map[*Client]bool),

}
/**
 * [should run forever with a go routine]
 * (manager* clientmanger)---> as i understand prohibits other type of data not to acesss the function
 * @param  {[type]} manager [description]
 * @return {[type]}         [description]
 */
func (manager *ClientManager) start(){

  for
  {
    select
    {
    case conn := <-manager.register:

      //addding a new client(registering)
      manager.clients[conn] = true
      jsonMessage,_ := json.Marshal(&Message{Content:"/A new client(socket) is connected"})
      manager.send(jsonMessage,conn)

    case conn := <-manager.unregister:
      //remove client(unregistering) partially clear on the idea
      if _,ok :=manager.clients[conn]
      {
      close(conn.send)
      delete(manager.clients, conn)
      jsonMessage,_:=json.Marshal(&Message{Content:"/A client has disconncted"})
      manager.send(jsonMessage,conn)
      }
      //BROADCASTING A message to all the clients
    case message := <-manager.broadcast:
      for conn :=range manager.clients
      {
        select
        {
        case conn.send <-message:
        default:
          close(conn.send)
          delete(manager.clients,conn)
        }

      }
    }
  }
// all funcs belows not mine copy paseted from the client example in github so try to understand them
  func (manager *ClientManager) send(message []byte, ignore *Client) {
      for conn := range manager.clients {
          if conn != ignore {
              conn.send <- message
          }
      }
  }


  func (c *Client) read() {
    defer func() {
        manager.unregister <- c
        c.socket.Close()
    }()

    for {
        _, message, err := c.socket.ReadMessage()
        if err != nil {
            manager.unregister <- c
            c.socket.Close()
            break
        }
        jsonMessage, _ := json.Marshal(&Message{Sender: c.id, Content: string(message)})
        manager.broadcast <- jsonMessage
    }
}

func (c *Client) write() {
    defer func() {
        c.socket.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.socket.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            c.socket.WriteMessage(websocket.TextMessage, message)
        }
    }
}

func wsPage(res http.ResponseWriter, req *http.Request) {
    conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
    if error != nil {
        http.NotFound(res, req)
        return
    }
    client := &Client{id: uuid.NewV4().String(), socket: conn, send: make(chan []byte)}

    manager.register <- client

    go client.read()
    go client.write()
}

}
