package websocket

import (
	"github.com/yuyenews/Beerus/application/websocket/params"
	"github.com/yuyenews/Beerus/application/websocket/wroute"
	"github.com/yuyenews/Beerus/commons/util"
	"log"
	"net"
)

// SessionMap Save each linked session
var sessionMap = make(map[string]*params.WebSocketSession)

var snowflake *util.SnowFlake

// ExecuteConnection Triggers the onConnection function inside the wroute
func ExecuteConnection(routePath string, conn net.Conn) {

	var snowflakeId uint64
	var err error
	if snowflake == nil {
		snowflake, err = util.New(2)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	snowflakeId, err = snowflake.Generate()
	if err != nil {
		log.Println(err.Error())
		return
	}

	session := new(params.WebSocketSession)
	session.Connection = conn
	session.Id = snowflakeId
	sessionMap[routePath] = session

	wroute.GetWebSocketRoute(routePath, wroute.OnConnection)(session, "The client is already connected")
}

// ExecuteMessage Triggers the onMessage function inside the wroute
func ExecuteMessage(routePath string, message string) {
	wroute.GetWebSocketRoute(routePath, wroute.OnMessage)(sessionMap[routePath], message)
}

// ExecuteClose Triggers the onClose function inside the wroute
func ExecuteClose(routePath string) {
	if sessionMap[routePath] == nil {
		return
	}

	delete(sessionMap, routePath)
	wroute.GetWebSocketRoute(routePath, wroute.OnClose)(sessionMap[routePath], "The client is disconnected")
}
