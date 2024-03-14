package play

import (
	"ddz/model"
	"ddz/user"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	ErrUserNum     = errors.New("玩家人数必须是3人")
	ErrNoUserScore = errors.New("没有用户叫分")
)

type Game struct {
	Players     []*user.User
	timer       time.Timer
	cards       []model.Card
	Score       int // 基础叫分
	Multi       int // 倍数
	banker      *user.User
	Winner      *user.User
	currentHand model.Hand
}

func (g *Game) Start() error {
	if len(g.Players) != 3 {
		return ErrUserNum
	}
	// 洗牌 发牌
	g.init()
	// 叫分
	err := g.initScore()
	if err != nil {
		// 叫分异常（无人叫分），则重新开始游戏
		log.Println("无人叫分,重新开始游戏", err.Error())
		for _, player := range g.Players {
			player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, "无人叫分,重新开始游戏")
		}
		return g.Start()
	}
	// 开始出牌循环
	g.gameLoop()
	// 输出结果信息
	g.writeResult()

	return nil
}

func (g *Game) Close() error {
	for _, player := range g.Players {
		player.Close()
	}
	return nil
}

func (g *Game) gameLoop() {
	for {
		println("开始回合")
		g.round()
		println("结束回合")
		if g.checkWinner() {
			break
		}
	}
}

func (g *Game) writeResult() {
	score := g.Score
	for i := 1; i < g.Multi; i++ {
		score *= 2
	}
	masterScore := 2 * score
	role := "地主"
	if g.Winner.Role == user.Slave {
		role = "农民"
	}
	for _, player := range g.Players {
		player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("用户%s率先出完牌，%s胜利", g.Winner.Name, role))
		if g.Winner.Role == user.Master {
			if player.Role == user.Master {
				player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("你赢得%d分", masterScore))
			} else {
				player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("你输掉%d分", score))
			}
		} else {
			if player.Role == user.Master {
				player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("你输掉%d分", masterScore))
			} else {
				player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("你赢得%d分", score))
			}
		}
	}
}

func (g *Game) checkWinner() bool {
	for _, player := range g.Players {
		if len(player.Cards) == 0 {
			g.Winner = player
			return true
		}
	}
	return false
}

func (g *Game) round() int {
	// 新一轮出牌，置空当前桌面手牌
	g.currentHand, _ = model.ParseHand([]model.Card{})
	idx := 0
	for i, player := range g.Players {
		if player == g.banker {
			idx = i
			break
		}
	}
	first := true

	for {
		u := g.Players[idx]
		// 如果出牌人是上一个banker，表示后续的两人都pass了，此回合结束
		if u == g.banker && !first {
			return 0
		}
		first = false

		timer := time.NewTimer(30 * time.Second)
		// 当前用户根据当前桌面手牌出牌
		outHand := u.OutCards(g.currentHand, timer, nil)
		timer.Stop()
		// 如果有出牌，表示能压过桌面手牌，banker设置为当前出牌人，桌面手牌更新为当前手牌
		if outHand.Type != model.HandTypeNone {
			g.banker = u
			g.currentHand = outHand
		}
		if outHand.Type == model.HandTypeKingBoom || outHand.Type == model.HandTypeNormalBoom {
			g.Multi++
		}
		// 通知所有玩家当前出牌情况
		for _, player := range g.Players {
			for _, op := range g.Players {
				if op.Name != player.Name {
					player.NoticeOtherLeft(len(op.Cards), op.Name)
				}
			}
			if outHand.Type == model.HandTypeNone {
				player.NoticeAny(user.MessageTypeNormalInfo, outHand, fmt.Sprintf("玩家 %s 弃牌", u.Name))
			} else {
				player.NoticeCurrentHand(outHand, fmt.Sprintf("玩家 %s 出牌", u.Name))
			}
			player.NoticeSelfCards(player.Cards)
		}
		// 如果当前出牌人剩余牌数量为0，则不继续游戏
		if len(u.Cards) == 0 {
			g.Winner = u
			return 1
		}

		// 定位下一位出牌人
		idx++
		idx %= 3
	}
}

func (g *Game) initScore() error {
	for _, player := range g.Players {
		timer := time.NewTimer(30 * time.Second)
		score := player.GetScore(g.Score, timer)
		if score > g.Score {
			g.Score = score
			g.banker = player
		}
		for _, p := range g.Players {
			p.NoticeScore(score, player.Name)
		}
	}

	if g.Score == 0 {
		return ErrNoUserScore
	}

	g.banker.Role = user.Master
	for _, player := range g.Players {
		player.NoticeAny(user.MessageTypeNormalInfo, model.Hand{}, fmt.Sprintf("%s叫分%d分，成功当选地主", g.banker.Name, g.Score))
	}

	// 通知底牌给地主
	g.banker.NoticeAny(user.MessageTypeNoticeGroundCards, model.Hand{Cards: g.cards[51:]}, "底牌")
	// 给地主发剩余底牌
	g.banker.AddCards(g.cards[51:]...)
	g.banker.SortCards()
	g.banker.NoticeSelfCards(g.banker.Cards)

	return nil
}

func (g *Game) init() {
	// 重置玩家手牌
	for _, player := range g.Players {
		player.Reset()
	}
	// 洗牌
	g.cards = model.Shuffle(model.GetInitCards())
	// 发牌
	g.initUserCards()
}

func (g *Game) initUserCards() {
	for i := 0; i < 17; i++ {
		for idx, player := range g.Players {
			player.AddCards(g.cards[i*3+idx])
		}
	}
	// 通知当前手牌
	for _, player := range g.Players {
		player.SortCards()
		player.NoticeSelfCards(player.Cards)
	}
}

func NewGame(u1, u2, u3 *user.User) *Game {
	return &Game{
		Players: []*user.User{u1, u2, u3},
		Score:   0,
		Multi:   1,
	}
}
