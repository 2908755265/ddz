package wsserver

import (
	"ddz/model"
	"ddz/play"
	"ddz/user"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	waitUsers    = map[string]*user.User{}
	playingUsers = map[string]*user.User{}
	lock         sync.Mutex
)

// WsServer
// @Description:
type WsServer struct {
	srv      *http.Server
	upGrader *websocket.Upgrader
	mux      *http.ServeMux
}

// Start
//
//	@Description:
//	@receiver s
func (s *WsServer) Start() {
	s.setRouter()
	fmt.Println("websocket server started at :32123")
	err := s.srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

// setRouter
//
//	@Description:
//	@receiver s
func (s *WsServer) setRouter() {
	s.mux.HandleFunc("/", s.game)
	s.srv.Handler = s.mux
}

// notice
//
//	@Description:
//	@receiver s
//	@param w
//	@param r
func (s *WsServer) game(w http.ResponseWriter, r *http.Request) {
	fmt.Println("收到WS连接请求")
	header := http.Header{}
	header.Set("Access-Control-Allow-Origin", "*")
	conn, err := s.upGrader.Upgrade(w, r, header)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		_ = conn.Close()
		fmt.Println("Error:", "can not get username from query, conn has been closed.")
		return
	}

	lock.Lock()
	defer lock.Unlock()

	u, ok := playingUsers[username]
	if ok {
		u.Conn = conn
		return
	}

	newUser := user.NewUser(username)
	newUser.Conn = conn
	waitUsers[username] = newUser
	newUser.HandleMessage()
	println(username, "加入房间。", "当前等待人数", len(waitUsers))
	for _, p := range waitUsers {
		p.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("用户%s加入房间,还需等待%d人", newUser.Name, 3-len(waitUsers)))
	}

	if len(waitUsers) == 3 {
		temp := make([]*user.User, 0, 3)
		for name, u := range waitUsers {
			temp = append(temp, u)
			playingUsers[name] = u
		}
		waitUsers = make(map[string]*user.User)
		game := play.NewGame(temp[0], temp[1], temp[2])
		go func() {
			for _, player := range game.Players {
				player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, "游戏开始...")
			}
			println("游戏开始...")
			err := game.Start()
			if err != nil {
				fmt.Println("Error:", err.Error())
			}
			// 关闭
			game.Close()
			lock.Lock()
			for _, player := range game.Players {
				delete(playingUsers, player.Name)
				println(player.Name, "离开房间")
			}
			println("游戏结束")
			lock.Unlock()
		}()
	}
}

// NewWsServer
//
//	@Description:
//	@param ctx
//	@return service.Starter
func NewWsServer() *WsServer {
	ug := new(websocket.Upgrader)
	ug.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &WsServer{
		srv:      &http.Server{Addr: fmt.Sprintf(":%d", 32123)},
		upGrader: ug,
		mux:      http.NewServeMux(),
	}
}
