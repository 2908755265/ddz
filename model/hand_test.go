package model

import (
	"fmt"
	"testing"
)

func TestThreePair(t *testing.T) {
	cards := []Card{{'Q', 0}, {'Q', 1}, {'K', 0}, {'K', 1}, {'K', 2}}
	hand, err := ParseHand(cards)
	if err != nil {
		panic(err)
	}
	if hand.Type != HandTypeThreePair {
		panic(fmt.Sprintf("期望%s类型，实际%s类型", GetHandTypeName(HandTypeThreePair), GetHandTypeName(hand.Type)))
	}
	if hand.Value.Number != 'K' {
		panic(fmt.Sprintf("期望%c值，实际%c值", 'K', hand.Value.Number))
	}
}
