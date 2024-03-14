package test

import (
	"ddz/model"
	"ddz/ui"
	"fmt"
	"testing"
	"time"
)

func TestType(t *testing.T) {
	fmt.Printf("%c\n", '♦')
}

func TestShuffle(t *testing.T) {
	ori := model.GetInitCards()
	//rend := ui.Rend(ori)
	//for _, bts := range rend {
	//	fmt.Printf("%s\n", string(bts))
	//}
	cards := model.Shuffle(ori)
	rend := ui.Rend(cards)
	for _, bts := range rend {
		fmt.Printf("%s\n", string(bts))
	}

}

func TestCh(t *testing.T) {
	//ch := make(chan int)
	var ch chan int
	after := time.After(3 * time.Second)
	go func() {
		println("adding")
		ch <- 10
		println("added")
	}()
	select {
	case v := <-ch:
		println(v)
	case <-after:
		println("time out")
	}
}

func TestBack(t *testing.T) {
	ui.Println(ui.RendBack(5))
}

func TestAscii(t *testing.T) {
	println('A', 'Z', 'a', 'z')
	println('a' - 'A')
}

func TestParseHand(t *testing.T) {
	cards := []model.Card{{'3', 0}, {'4', 0}, {'5', 0}, {'6', 0}, {'7', 0}, {'8', 0}}
	hand, err := model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
	println("-------------------------------------------------------")
	cards = []model.Card{
		{'3', 0}, {'4', 0}, {'5', 0}, {'6', 0}, {'7', 0}, {'8', 0},
		{'3', 3}, {'4', 3}, {'5', 3}, {'6', 3}, {'7', 3}, {'8', 3},
	}
	hand, err = model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
	println("-------------------------------------------------------")
	cards = []model.Card{
		{'3', 0}, {'3', 3}, {'3', 1}, {'3', 2},
	}
	hand, err = model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
	println("-------------------------------------------------------")
	cards = []model.Card{
		{'3', 0}, {'3', 3}, {'3', 1}, {'3', 2},
		{'4', 0}, {'4', 3}, {'4', 1}, {'4', 2},
		{'5', 0}, {'5', 3}, {'5', 1}, {'5', 2},
		{'6', 0}, {'6', 3}, {'6', 1}, {'6', 2},
	}
	hand, err = model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
	println("-------------------------------------------------------")
	cards = []model.Card{
		{'3', 0}, {'3', 3}, {'3', 1},
		{'4', 0}, {'4', 3}, {'4', 1},
		{'5', 0}, {'5', 3}, {'5', 1},
		{'6', 0}, {'6', 3}, {'6', 1},
		{'7', 0}, {'7', 3}, {'7', 1}, {'7', 2},
	}
	hand, err = model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
	println("-------------------------------------------------------")
	cards = []model.Card{
		{'3', 0}, {'3', 3}, {'3', 1},
		{'4', 0}, {'4', 3}, {'4', 1},
		{'5', 0}, {'5', 3}, {'5', 1},
		{'6', 0}, {'6', 3}, {'6', 1},
		{'8', 0}, {'8', 3}, {'8', 1}, {'8', 2},
	}
	hand, err = model.ParseHand(cards)
	if err != nil {
		println(err.Error())
	}
	printHand(hand)
}

func printHand(h model.Hand) {
	println("手牌类型", model.GetHandTypeName(h.Type))
	println("手牌最大值")
	ui.Println(ui.Rend([]model.Card{h.Value}))
	println("手牌")
	ui.Println(ui.Rend(h.Cards))
}
