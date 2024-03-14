package main

import (
	"bufio"
	"context"
	"ddz/ui"
	"ddz/user"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"sync"
	"time"
)

var (
	username = flag.String("u", "mack", "用户名")
	host     = flag.String("h", "192.168.1.16", "服务器IP")
)

func main() {
	flag.Parse()
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:32123/?username=%s", *host, *username), nil)
	if err != nil {
		panic(err)
	}
	println("username", *username)
	println("server", *host)
	defer conn.Close()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				conn.Close()
				wg.Done()
				return
			}
			handle(conn, p)
		}
	}()
	wg.Wait()
}

func handle(conn *websocket.Conn, p []byte) {
	msg := user.OutMessage{}
	err := json.Unmarshal(p, &msg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch msg.MessageType {
	case user.MessageTypeSelfCards:
		println("你的剩余手牌为：")
		ui.Println(ui.Rend(msg.Hand.Cards))
	case user.MessageTypeOutScore:
		println(msg.Desc)
		println("请输入您的出分(0-3,回车结束):")
		input(conn, user.MessageTypeOutScore)
	case user.MessageTypeNormalInfo:
		println(msg.Desc)
	case user.MessageTypeNoticeGroundCards:
		println(msg.Desc)
		ui.Println(ui.Rend(msg.Hand.Cards))
	case user.MessageTypeGetCards:
		println("请输入您要出的出牌(30秒出牌时间，回车结束):")
		//input(conn, user.MessageTypeGetCards)
		withTimeoutInput(conn, user.MessageTypeGetCards, 30*time.Second)
	case user.MessageTypeNotice:
		println(msg.Desc)
		ui.Println(ui.Rend(msg.Hand.Cards))
	case user.MessageTypeOtherLeft:
		println(msg.Desc)
		ui.Println(ui.RendBack(len(msg.Hand.Cards)))
	case user.MessageTypeNoticeScore:
		println(msg.Desc)
	default:
		println("消息：", string(p))
	}
}

func withTimeoutInput(conn *websocket.Conn, t int, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultCh := make(chan int, 1)

	go func() {
		input(conn, t)
		resultCh <- 1
	}()

	select {
	case <-ctx.Done():
		println("已超时，自动出牌或弃牌")
	case <-resultCh:
		return
	}
}

func input(conn *websocket.Conn, t int) {
	reader := bufio.NewReader(os.Stdin)
	line, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = conn.WriteJSON(user.Message{
		MessageType: t,
		Hand:        string(line),
		Desc:        "",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
