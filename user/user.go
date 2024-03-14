package user

import (
	"ddz/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"sync/atomic"
	"time"
)

const (
	Master = 1
	Slave  = 0

	MessageTypeNotice            = 1
	MessageTypeGetCards          = 2
	MessageTypeOutCards          = 3
	MessageTypeOutScore          = 4
	MessageTypeNoticeScore       = 5
	MessageTypeSelfCards         = 6
	MessageTypeNormalInfo        = 7
	MessageTypeNoticeGroundCards = 8
	MessageTypeOtherLeft         = 9

	StatWaitingCards = 1
	StatIdle         = 0
	StatWaitingScore = 2
)

var (
	ErrHandIsNotGreater = errors.New("手牌无法压过当前牌面")
	ErrHandHasNoCards   = errors.New("没有足够的手牌")
	ErrHandNeedCards    = errors.New("您必须最少出一张牌")
)

type User struct {
	Name         string
	Cards        []model.Card
	Role         uint8 // 1 地主，0 农民
	Conn         *websocket.Conn
	GettingCards atomic.Int32
	outCh        chan []model.Card
	score        chan int
}

type Message struct {
	MessageType int
	Hand        string
	Desc        string
}

type OutMessage struct {
	MessageType int
	Hand        model.Hand
	Desc        string
}

func (u *User) SortCards() {
	model.SortCards(u.Cards)
}

func (u *User) AddCards(cards ...model.Card) {
	u.Cards = append(u.Cards, cards...)
}

func (u *User) Reset() {
	u.Cards = make([]model.Card, 0, 21)
	u.Role = Slave
}

func (u *User) GetScore(score int, timer *time.Timer) int {
	u.GettingCards.Store(StatWaitingScore)
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: MessageTypeOutScore,
		Hand:        model.Hand{},
		Desc:        fmt.Sprintf("当前分数%d", score),
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
		return 0
	}
	select {
	case s := <-u.score:
		return s
	case <-timer.C:
		fmt.Printf("玩家%s超时放弃叫分", u.Name)
		return 0
	}
}

func (u *User) OutCards(hand model.Hand, timer *time.Timer, e error) model.Hand {
	// 进入等待出牌消息状态
	u.GettingCards.Store(StatWaitingCards)

	v := model.Hand{}
	msg := ""
	if e != nil {
		msg = e.Error()
	}
	preHand, err := u.preGetCards(hand, timer, msg)
	if err != nil {
		fmt.Println("Error:", err.Error())
		u.NoticeAny(MessageTypeNormalInfo, model.Hand{}, err.Error())
		return u.OutCards(hand, timer, err)
	}
	// 放弃出牌
	if preHand.Type == model.HandTypeNone {
		u.GettingCards.Store(StatIdle)
		return v
	}
	// 预出牌比大小
	isGreater, err := model.IsGreaterHand(preHand, hand)
	if err != nil {
		fmt.Println("Error:", err.Error())
		u.NoticeAny(MessageTypeNormalInfo, model.Hand{}, err.Error())
		return u.OutCards(hand, timer, err)
	}
	if !isGreater {
		fmt.Println("Error:", ErrHandIsNotGreater.Error())
		u.NoticeAny(MessageTypeNormalInfo, model.Hand{}, ErrHandIsNotGreater.Error())
		return u.OutCards(hand, timer, ErrHandIsNotGreater)
	}
	// 出牌
	outHand, err := u.outCards(preHand)
	if err != nil {
		fmt.Println("Error:", err.Error())
		u.NoticeAny(MessageTypeNormalInfo, model.Hand{}, err.Error())
		return u.OutCards(hand, timer, err)
	}

	u.GettingCards.Store(StatIdle)
	return outHand
}

func (u *User) HandleMessage() {
	go func() {
		for {
			println("等待新的消息...")
			_, p, err := u.Conn.ReadMessage()
			if err != nil {
				fmt.Println("Error:", err.Error())
				u.Conn.Close()
				break
			}
			println("收到消息：", string(p))
			msg := Message{}
			err = json.Unmarshal(p, &msg)
			if err != nil {
				fmt.Println("Error:", err.Error())
			} else {
				go u.handle(msg)
			}
		}
	}()
}

func (u *User) handle(msg Message) {
	fmt.Println(msg, u.GettingCards.Load())
	if u.GettingCards.Load() == StatWaitingCards {
		u.GettingCards.Store(StatIdle)
		// 正在等待出牌阶段，收到消息后退出等待消息状态
		cards := make([]model.Card, 0, len(msg.Hand))
		for _, b := range []byte(msg.Hand) {
			if b >= 'a' && b <= 'z' {
				b -= 32
			}
			_, ok := model.Score[b]
			if ok {
				cards = append(cards, model.Card{
					Number: b,
				})
			} else {
				fmt.Printf("Info: 丢弃不能识别的卡牌%c\n", b)
			}
		}
		u.outCh <- cards
	} else if u.GettingCards.Load() == StatWaitingScore {
		u.GettingCards.Store(StatIdle)
		if len(msg.Hand) == 0 {
			u.score <- 0
		} else {
			if msg.Hand[0] >= '0' && msg.Hand[0] <= '3' {
				println("推送score", int(msg.Hand[0]-'0'))
				u.score <- int(msg.Hand[0] - '0')
			} else {
				u.score <- 0
			}
		}
		println("玩家", u.Name, "出分完成...")
	} else {
		// 非出牌阶段，丢弃出牌消息
		fmt.Println("当前状态", u.GettingCards.Load(), "已丢弃消息", msg)
		u.NoticeAny(MessageTypeNormalInfo, model.Hand{}, "当前不是你的出牌轮次，无效操作")
	}
}

func (u *User) NoticeCurrentHand(hand model.Hand, msg string) {
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: MessageTypeNotice,
		Hand:        hand,
		Desc:        msg,
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func (u *User) NoticeSelfCards(cards []model.Card) {
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: MessageTypeSelfCards,
		Hand:        model.Hand{Cards: cards},
		Desc:        "",
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func (u *User) NoticeOtherLeft(num int, name string) {
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: MessageTypeOtherLeft,
		Hand:        model.Hand{Cards: make([]model.Card, num)},
		Desc:        fmt.Sprintf("玩家%s剩余", name),
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func (u *User) NoticeAny(tp int, hand model.Hand, msg string) {
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: tp,
		Hand:        hand,
		Desc:        msg,
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func (u *User) Close() error {
	return u.Conn.Close()
}

func (u *User) NoticeScore(score int, name string) {
	err := u.Conn.WriteJSON(OutMessage{
		MessageType: MessageTypeNoticeScore,
		Hand:        model.Hand{},
		Desc:        fmt.Sprintf("玩家%s 叫分%d", name, score),
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func (u *User) outCards(hand model.Hand) (model.Hand, error) {
	leftCards := u.Cards
	outCards := make([]model.Card, 0)
	outHand := model.Hand{}
	for _, card := range hand.Cards {
		tempCards := make([]model.Card, 0)
		find := false
		for _, leftCard := range leftCards {
			if !find && leftCard.Number == card.Number {
				find = true
				outCards = append(outCards, leftCard)
			} else {
				tempCards = append(tempCards, leftCard)
			}
		}
		if !find {
			return outHand, ErrHandHasNoCards
		}
		leftCards = tempCards
	}
	u.Cards = leftCards
	return model.ParseHand(outCards)
}

func (u *User) preGetCards(h model.Hand, timer *time.Timer, msg string) (model.Hand, error) {
	// 通知出牌
	u.sendMsg(OutMessage{
		MessageType: MessageTypeGetCards,
		Hand:        h,
		Desc:        msg,
	})
	// 等待出牌
	for {
		select {
		case <-timer.C:
			println("等待出牌超时")
			// 超时，如果是banker则随便出一张最小的，否则放弃出牌
			if h.Type == model.HandTypeNone {
				nh, _ := model.ParseHand([]model.Card{u.Cards[0]})
				return nh, nil
			} else {
				return model.Hand{}, nil
			}
		case ph := <-u.outCh:
			println("上轮手牌类型", model.GetHandTypeName(h.Type), "此轮手牌数量", len(ph))
			if len(ph) == 0 && h.Type == model.HandTypeNone {
				println("进入错误")
				return model.Hand{}, ErrHandNeedCards
			}
			return model.ParseHand(ph)
		}
	}
}

func (u *User) sendMsg(msg OutMessage) {
	err := u.Conn.WriteJSON(msg)
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}

func NewUser(name string) *User {
	return &User{
		Name:         name,
		Cards:        make([]model.Card, 0, 21),
		Role:         Slave,
		GettingCards: atomic.Int32{},
		outCh:        make(chan []model.Card),
		score:        make(chan int),
	}
}
