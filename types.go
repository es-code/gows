package gows

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Connection struct {
	lock *sync.Mutex
	socket *websocket.Conn
	Channel string
	ConnId string
}

type connections struct {
	connections map[string]Connection
}


type Hub struct {
	channels map[string]connections
	channelsLock *sync.Mutex `json:"channels_lock"`
	Upgrader *websocket.Upgrader
	WriteWait time.Duration
	PongWait time.Duration
	PingPeriod time.Duration
	MaxMessageSize int64
}

