package main

import (
	"bytes"
	"log"
	"net/http"
	"time"
  "github.com/gorilla/websocket"
  )

/* structs*/


type Client struct {
  ws * websocket.Conn
  //sever passes content to this channel
  send chan []byte

}


//server broascasts a new message and this fires

func(c *Client) write(){
  //make sure to close the connection incase loops exits ------> check this function
  defer func() {
    c.ws.Close()
  }()

  for{
    select{
    case message,ok:= <-c.send:
      if !ok{
        c.ws.WriteMessage(websocket.CloseMessage,[]byte{})
        return
      }
      c.ws.WriteMessage(websocket.TextMessage,message)
    }

  }
}


//New message recived so pass it to hub

func (c *Client) read(){
  defer func(){

  }
}
