package router

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/tanjl855/blog/internal/blog/controller"
	response "github.com/tanjl855/blog/internal/blog/http"
)

func SetupRouter() *gin.Engine {
	r := gin.New()

	blogCors := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	})

	if viper.GetString("env") != "production" {
		pprof.Register(r)
	}

	r.Use()

	r.Use(gin.Recovery(), blogCors)

	r.NoRoute(func(c *gin.Context) {
		rsp := gin.H{"errCode": response.GatewayNotFound, "message": response.ResponseMessages[response.GatewayNotFound]}
		c.JSON(http.StatusOK, rsp)
	})

	r.GET("/blog/v1/configs", controller.Configs)
	InitWs("192.168.136.129:2023", "111", 100, 10, 1024*1024*1000)
	r.GET("/ws", WS.WsHandler)

	return r
}

type UserConn struct {
	*websocket.Conn
	w   *sync.Mutex
	UId string
}

type WServer struct {
	rwLock       *sync.RWMutex
	wsOutAddr    string
	grpcAddr     string
	wsMaxConnNum int
	wsUpGrader   *websocket.Upgrader
	wsUserToConn map[string]*UserConn
}

var WS *WServer

func InitWs(outAddr string, grpcAddr string, connMaxNum int, timeout int, maxMsgLen int) {
	WS = &WServer{
		rwLock:       &sync.RWMutex{},
		wsOutAddr:    outAddr,
		grpcAddr:     grpcAddr,
		wsMaxConnNum: connMaxNum,
		wsUpGrader: &websocket.Upgrader{
			HandshakeTimeout: time.Duration(timeout) * time.Second,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsUserToConn: make(map[string]*UserConn),
	}
	// 启协程保活ws的外链地址
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("更新服务器心跳失败，维持在线状态失败：%+v", err)
				return
			}
		}()
		// 放到redis的zset里面，score存用户map的长度，member放外链ip地址，做一层负载均衡

		// 定期更新心跳信息（外链地址）
		minTick := time.NewTicker(60 * time.Second)
		for {
			// redis存key外链地址，value随便(eg:1)
			fmt.Println("定期更新到redis，维持在线状态", outAddr)
			<-minTick.C
		}
	}()
}

type wsHandlerReq struct {
	SendId string `json:"send_id"`
	UserId string `json:"user_id"`
}

func (ws *WServer) WsHandler(ctx *gin.Context) {
	req := &wsHandlerReq{}
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		// fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errCode": -1, "message": "参数错误"})
		return
	}
	// fmt.Println(req)
	operationID := ws.grpcAddr
	fmt.Printf("ws pool服务，机器号；%v, sendId: %s \n", operationID, req.SendId)
	conn, err := ws.wsUpGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusOK, gin.H{"errCode": -1, "message": err.Error()})
		return
	}
	newConn := &UserConn{
		Conn: conn,
		w:    &sync.Mutex{},
		UId:  req.SendId,
	}
	fmt.Println("lalala", newConn.UId)

	// read msg
	go ws.readMsg(newConn)
}

func (ws *WServer) readMsg(conn *UserConn) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("读取消息失败：%+v", err)
			return
		}
	}()
	for {
		messageType, msgReader, err := conn.NextReader()
		if messageType == websocket.PingMessage {
			fmt.Printf("uid:%s pingMessage", conn.UId)
			continue
		}
		if err != nil {
			fmt.Printf("uid:%s, remoteIp:%s read msg err:%+v\n", conn.UId, conn.RemoteAddr(), err)
			return
		}

		msg, err := io.ReadAll(msgReader)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.w.Lock()
		defer conn.w.Unlock()
		conn.SetWriteDeadline(time.Now().Add(time.Duration(300) * time.Second))
		// 丢到kafka消息队列

		// 写到ws里面
		fmt.Println(string(msg))
		conn.WriteMessage(websocket.BinaryMessage, msg)
	}
}
