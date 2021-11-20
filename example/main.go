package main

import (
	"fmt"
	"net/http"
	"github.com/es-code/gows"
)

var ws *gows.Hub

func main()  {
	ws = gows.Init()

	//endpoint to open ws connection and read messages , write message to connection
	http.HandleFunc("/message/listen",messageSocket)
	//endpoint to send message on channel
	http.HandleFunc("/message/send",sendMessage)

	//start web server
	http.ListenAndServe(":9000",nil)

}


func messageSocket(res http.ResponseWriter, req *http.Request)  {

	channel:=req.FormValue("channel")

	//open connection for current request
	conn,err:=ws.Open(channel,res,req)

	if err!=nil{
		fmt.Println("ws error open conn")
	}
	//start listen for messages
	err=ws.Listen(conn,func(messageType int, message []byte) {
		//handle every message
		//re-write message on channel
		ws.WriteOnChannel(channel,message)
	})

	if err!=nil{
		fmt.Println("gows connection ended")
	}
	return

}

//send message on channel
func sendMessage(res http.ResponseWriter,req *http.Request)  {
	channel:=req.FormValue("channel")
	if channel == ""{
		 res.Write([]byte("please send your channel name"))
		return
	}
	ws.WriteOnChannel(channel,[]byte(string("hello every one listen on "+req.FormValue("channel"))))

}