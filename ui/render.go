package ui

import (
	"ddz/model"
	"fmt"
)

func Rend(cards []model.Card) [][]rune {
	if len(cards) == 0 {
		return nil
	}
	arr := make([][]rune, 4)
	l := len(cards)*3 + 1
	arr[0] = make([]rune, 0, l)
	arr[0] = append(arr[0], '┌')
	arr[1] = make([]rune, 0, l)
	arr[1] = append(arr[1], '│')
	arr[2] = make([]rune, 0, l)
	arr[2] = append(arr[2], '│')
	arr[3] = make([]rune, 0, l)
	arr[3] = append(arr[3], '└')
	for _, card := range cards {
		arr = rendOneCard(card, arr)
	}
	return arr
}

func RendBack(num int) [][]rune {
	if num <= 0 {
		return nil
	}
	arr := make([][]rune, 4)
	l := num*3 + 1
	arr[0] = make([]rune, 0, l)
	arr[0] = append(arr[0], '┌')
	arr[1] = make([]rune, 0, l)
	arr[1] = append(arr[1], '│')
	arr[2] = make([]rune, 0, l)
	arr[2] = append(arr[2], '│')
	arr[3] = make([]rune, 0, l)
	arr[3] = append(arr[3], '└')
	for i := 0; i < num; i++ {
		arr = rendOneCard(model.Card{Type: -1, Number: 'M'}, arr)
	}
	return arr
}

func rendOneCard(c model.Card, cards [][]rune) [][]rune {
	cards[0] = append(cards[0], '─', '─', '┐')
	cards[1] = append(cards[1], rune(c.Number), ' ', '|')
	switch c.Type {
	case -1:
		cards[2] = append(cards[2], ' ')
	case 0:
		cards[2] = append(cards[2], '♠')
	case 1:
		cards[2] = append(cards[2], '♥')
	case 2:
		cards[2] = append(cards[2], '♣')
	case 3:
		cards[2] = append(cards[2], '♦')
	}
	cards[2] = append(cards[2], ' ', '|')
	cards[3] = append(cards[3], '─', '─', '┘')
	return cards
}

func Println(cards [][]rune) {
	for _, bts := range cards {
		fmt.Printf("%s\n", string(bts))
	}
}
