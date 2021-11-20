package gows

import (
	"errors"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

//default gows Upgrader

var upgrader = websocket.Upgrader{
ReadBufferSize:    1024,
WriteBufferSize:   1024,
CheckOrigin:       func(r *http.Request) bool { return true },
EnableCompression: true,
}


//Initialization new Hub

func Init() *Hub  {
	//new channels map
	channelsMap:=make(map[string]connections)
	hub:=&Hub{channels:channelsMap,channelsLock: &sync.Mutex{}}
	//set gows default settings
	hub.SetDefault()
	return hub
}

//default settings

func (h *Hub) SetDefault()  {
	h.WriteWait = 10 * time.Second
	h.PongWait= 15 * time.Second
	h.PingPeriod = (h.PongWait * 9) / 10
	h.MaxMessageSize = 1024
	h.Upgrader = &upgrader
}


//add new connection to hub for channel

func (h *Hub) addConnection(channelId string,con *websocket.Conn) *Connection {
	//check if channel exists in hub
	_, ok := h.channels[channelId]
	//generate new connection id
	connId:=Uuid()
	//new connection
	connection:=Connection{
		lock: &sync.Mutex{},
		socket: con,
		Channel: channelId,
		ConnId: connId,
	}
	//lock hub channels
	h.channelsLock.Lock()
	//unlock hub channels
	defer func() {
		h.channelsLock.Unlock()
	}()

	//if channel exists
	if ok{
		//write new connection to map
		h.channels[channelId].connections[connId]=connection
	}else{
		newConnections:=make(map[string]Connection)
		newConnections[connId] = connection
		h.channels[channelId]=connections{
			connections:newConnections,
		}
	}

	return &connection

}

func (h *Hub) RemoveConn(c *Connection)  {
	conn,check :=h.checkConnection(c)
	 if check{
		 //close socket
		 conn.socket.Close()
		 //remove connection
		 delete(h.channels[conn.Channel].connections, conn.ConnId)

		 if len(h.channels[conn.Channel].connections) == 0{
			delete(h.channels,conn.Channel)
		 }
	 }

}

func (h *Hub) checkConnection(c *Connection) (*Connection,bool)  {
	conn,ok:=h.channels[c.Channel].connections[c.ConnId]
	return &conn,ok
}


func (h *Hub) Open(channel string, res http.ResponseWriter, req *http.Request) (*Connection,error) {

	//upgrade current request with hub Upgrader in safe lock
	var conn, err = h.Upgrader.Upgrade(res, req, nil)

	if err!=nil{
		return &Connection{},err
	}
	//add channel connection to hub in safe lock
	connection:=h.addConnection(channel,conn)
	//return connection
	return connection,nil

}

func (h *Hub) Listen(c *Connection,messageHandler func(messageType int, message []byte )) error {
	//check connection still in hub
	conn,check:=h.checkConnection(c)
	if !check{
		return errors.New("error connection not found in hub")
	}

	//open thread for reading messages from connection
	go h.readMessages(conn,messageHandler)
	//send ping messages
	h.pingHandler(conn)
	//listen ended
	return errors.New("socket connection ended")
}


func (h *Hub) WriteOnChannel(channel string,message []byte){
		for _,connection:= range h.channels[channel].connections{
			   h.writeToConn(&connection,1,message)
		}
}

func (h *Hub) WriteOnConn(c *Connection,messageType int, message []byte) error {
	conn,check:=h.checkConnection(c)
	if !check{
		return errors.New("error connection not found in hub")
	}

	return h.writeToConn(conn,messageType,message)

}

func (h *Hub) writeToConn(c *Connection,messageType int, message []byte) error {
	//lock connection
	c.lock.Lock()
	defer func() {
		//unlock connection
		c.lock.Unlock()
	}()

	c.socket.SetWriteDeadline(time.Now().Add(h.WriteWait))
	err:=c.socket.WriteMessage(messageType,message)
	if err != nil {
		h.RemoveConn(c)
		return err
	}
	return err
}



func (h *Hub) readMessages(Conn *Connection,messageHandler func(messageType int, message []byte )) {

	defer func() {
		h.RemoveConn(Conn)
	}()

	socket:=Conn.socket
	socket.SetReadLimit(h.MaxMessageSize)
	socket.SetReadDeadline(time.Now().Add(h.PongWait))
	socket.SetPongHandler(func(string) error { socket.SetReadDeadline(time.Now().Add(h.PongWait)); return nil })


	for {
		messageType, message, err := socket.ReadMessage()
		if err != nil {
			return
		}
		messageHandler(messageType,message)
	}

}


func (h *Hub) pingHandler(conn *Connection)  {
	socket:=conn.socket
	ticker := time.NewTicker(h.PingPeriod)

	defer func() {
		h.RemoveConn(conn)
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			socket.SetWriteDeadline(time.Now().Add(h.WriteWait))
			if err := socket.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}

}







