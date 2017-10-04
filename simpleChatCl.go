package main

import (
  "log"
  "net/http"

  "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel
var upgrader = websocket.Upgrader{}

type Message struct {

  Email string `json:"email"`
  Username string `json:"username"`
  Message string `json:"message"`

}

func handleConnections(w http.ResponseWriter, r *http.Request){
  //Upgrading the Get Request (initial ) to a WebSocket
  ws,err := upgrader.Upgrade(w,r,nil)
  if err != nil {
    log.Fatal(err)
  }
  //closing the connection to WebSocker
  defer ws.Close()
  clients[ws]=true
  for
  {
    var msg Message
    //Read in a new message as JSON  and add it to the Message struct
    err :=ws.ReadJSON(&msg)
    if err != nil {
      log.Printf("error:%v",err)
      delete(clients,ws)
      break
    }
    //sending the message to the broadcast channel created
    broadcast <- msg
  }
}

func handleMessages(){
  for {
    //grab the message from the broadcast channel (next message )
    log.Println("testing cli ")
    msg := <-broadcast
    log.Println("testing cli 2")
    //Send it out to every one that is conencted
    for client := range clients {
      err:=client.WriteJSON(msg)
      if err != nil{
        log.Printf("error:%v",err)
        client.Close()
        delete(clients,client)//deleting a map vlaue

      }
    }
  }
}
func main(){

  fs := http.FileServer(http.Dir("./public"))
  http.Handle("/",fs)
  http.HandleFunc("/ws", handleConnections)
  go handleMessages()
  log.Println("listening to the http server ")
  err := http.ListenAndServe(":8000",nil)
  if err != nil {
               log.Fatal("ListenAndServe: ", err)
       }




}
