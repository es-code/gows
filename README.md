# GOWS
**GOWS is GoLang web-socket module
Provides you with ease of handling web socket connections with a few lines, it supports multi-connection on one channel, Ping/Pong handler, saves multi concurrency writes on hub and connections.**
<hr>

## Installation :
`go get -u github.com/es-code/gows`

## Usage :

```

var ws *gows.Hub

func main()  {
ws = gows.Init()

//endpoint to open ws connection and read messages , write message to connection
http.HandleFunc("/message/listen",messageSocket)
//start web server
http.ListenAndServe(":9000",nil)

}
```

```
func messageSocket(res http.ResponseWriter, req *http.Request)  {

channel:=req.FormValue("channel")

//open connection for current request
conn,err:=ws.Open(channel,res,req)

if err!=nil{
fmt.Println("ws error open conn")
}
//start listen for messages
err=ws.Listen(conn,func(messageType int, message []byte) {
//your code here to handle every message
//example : re-write message on channel
ws.WriteOnChannel(channel,message)
})

if err!=nil{
fmt.Println("gows connection ended")
}
return

}
```









