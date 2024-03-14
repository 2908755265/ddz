package model

import (
	"errors"
	"math/rand"
	"sort"
)

const (
	HandTypeNone           = 0
	HandTypeKingBoom       = 1
	HandTypeNormalBoom     = 2
	HandTypeSingle         = 3
	HandTypePair           = 4
	HandTypeThree          = 5
	HandTypeThreeSingle    = 6
	HandTypeThreePair      = 7
	HandTypeStraight       = 8
	HandTypeStraightPair   = 9
	HandTypePlane          = 10
	HandTypePlaneTwoSingle = 11
	HandTypePlaneTwoPair   = 12
	HandTypeFourTwoSingle  = 13
	HandTypeFourTwoPair    = 14
)

var (
	Score = map[byte]int8{
		'3': 3,
		'4': 4,
		'5': 5,
		'6': 6,
		'7': 7,
		'8': 8,
		'9': 9,
		'0': 10,
		'J': 11,
		'Q': 12,
		'K': 13,
		'A': 14,
		'2': 15,
		'X': 16,
		'D': 17,
		'T': 18,
	}
	ErrHandType = errors.New("错误的手牌")
)

type Hand struct {
	Type  int
	Cards []Card
	Value Card
}

type Card struct {
	Number byte
	Type   int8 // 0 黑桃 1 红心 2 梅花 3 方块
}

func GetHandTypeName(ht int) string {
	switch ht {
	case HandTypeNone:
		return "弃牌"
	case HandTypeKingBoom:
		return "王炸"
	case HandTypeNormalBoom:
		return "炸弹"
	case HandTypeSingle:
		return "单张"
	case HandTypePair:
		return "对子"
	case HandTypeThree:
		return "三不带"
	case HandTypeThreeSingle:
		return "三带一"
	case HandTypeThreePair:
		return "三带一对"
	case HandTypeStraight:
		return "顺子"
	case HandTypeStraightPair:
		return "连对"
	case HandTypePlane:
		return "飞机"
	case HandTypePlaneTwoSingle:
		return "飞机带单"
	case HandTypePlaneTwoPair:
		return "飞机带对"
	case HandTypeFourTwoSingle:
		return "四带二"
	case HandTypeFourTwoPair:
		return "四带二对"
	default:
		return "未知"
	}
}

func GetInitCards() []Card {
	return []Card{
		{'A', 0}, {'2', 0}, {'3', 0}, {'4', 0}, {'5', 0}, {'6', 0}, {'7', 0}, {'8', 0}, {'9', 0}, {'0', 0}, {'J', 0}, {'Q', 0}, {'K', 0},
		{'A', 1}, {'2', 1}, {'3', 1}, {'4', 1}, {'5', 1}, {'6', 1}, {'7', 1}, {'8', 1}, {'9', 1}, {'0', 1}, {'J', 1}, {'Q', 1}, {'K', 1},
		{'A', 2}, {'2', 2}, {'3', 2}, {'4', 2}, {'5', 2}, {'6', 2}, {'7', 2}, {'8', 2}, {'9', 2}, {'0', 2}, {'J', 2}, {'Q', 2}, {'K', 2},
		{'A', 3}, {'2', 3}, {'3', 3}, {'4', 3}, {'5', 3}, {'6', 3}, {'7', 3}, {'8', 3}, {'9', 3}, {'0', 3}, {'J', 3}, {'Q', 3}, {'K', 3},
		{'D', -1}, {'X', -1}, //{'T', -1},
	}
}

func Shuffle(cards []Card) []Card {
	l := len(cards)
	r := make([]Card, 0, l)
	for {
		if l <= 0 {
			break
		}
		n := rand.Intn(l)
		r = append(r, cards[n])
		next := make([]Card, 0, l-1)
		for i := 0; i < l; i++ {
			if i == n {
				continue
			}
			next = append(next, cards[i])
		}
		cards = next
		l--
	}
	return r
}

func Greater(a, b byte) bool {
	return Score[a] > Score[b]
}

func SortCards(cards []Card) {
	sort.Slice(cards, func(i, j int) bool {
		return !Greater(cards[i].Number, cards[j].Number)
	})
}

func IsBoom(cards []Card) bool {
	SortCards(cards)
	return isNormalBoom(cards) || isKingBoom(cards)
}

func IsNormalBoom(cards []Card) bool {
	SortCards(cards)
	return isNormalBoom(cards)
}

func isNormalBoom(cards []Card) bool {
	return len(cards) == 4 && cards[0].Number == cards[1].Number && cards[0].Number == cards[2].Number && cards[0].Number == cards[3].Number
}

func IsKingBoom(cards []Card) bool {
	SortCards(cards)
	return isKingBoom(cards)
}

func isKingBoom(cards []Card) bool {
	return len(cards) == 2 && cards[0].Number == 'X' && cards[1].Number == 'D'
}

func IsPair(cards []Card) bool {
	SortCards(cards)
	return isPair(cards)
}

func isPair(cards []Card) bool {
	return len(cards) == 2 && cards[0].Number == cards[1].Number
}

func IsThree(cards []Card) bool {
	SortCards(cards)
	return isThree(cards)
}

func isThree(cards []Card) bool {
	return len(cards) == 3 && cards[0].Number == cards[1].Number && cards[0].Number == cards[2].Number
}

func IsThreeSingle(cards []Card) bool {
	SortCards(cards)
	return isThreeSingle(cards)
}

func isThreeSingle(cards []Card) bool {
	return len(cards) == 4 && (isThree(cards[:3]) || isThree(cards[1:]))
}

func IsThreePair(cards []Card) bool {
	SortCards(cards)
	return isThreePair(cards)
}

func isThreePair(cards []Card) bool {
	return len(cards) == 5 && (isThree(cards[:3]) && isPair(cards[3:])) || (isThree(cards[2:]) && isPair(cards[:2]))
}

func IsStraight(cards []Card) bool {
	SortCards(cards)
	return isStraight(cards)
}

func isStraight(cards []Card) bool {
	// 只能是5-A的连续 保障手牌范围
	if len(cards) < 5 || len(cards) > 12 || Score[cards[len(cards)-1].Number] > Score['A'] {
		return false
	}

	// 保障手牌连续
	current := Score[cards[0].Number]
	for i := 1; i < len(cards); i++ {
		if Score[cards[i].Number]-current != 1 {
			return false
		}
		current = Score[cards[i].Number]
	}

	return true
}

func IsStraightPair(cards []Card) bool {
	SortCards(cards)
	return isStraightPair(cards)
}

func isStraightPair(cards []Card) bool {
	if len(cards) < 6 || len(cards)%2 != 0 || Score[cards[len(cards)-1].Number] > Score['A'] {
		return false
	}
	// 保障手牌连续
	current := Score[cards[0].Number]
	if !isPair(cards[:2]) {
		return false
	}
	for i := 2; i < len(cards); i += 2 {
		if Score[cards[i].Number]-current != 1 || !isPair(cards[i:i+2]) {
			return false
		}
		current = Score[cards[i].Number]
	}

	return true
}

func IsPlane(cards []Card) bool {
	SortCards(cards)
	return isPlane(cards)
}

func isPlane(cards []Card) bool {
	if len(cards) < 6 || len(cards)%3 != 0 || Score[cards[len(cards)-1].Number] > Score['A'] {
		return false
	}
	// 保障手牌连续
	current := Score[cards[0].Number]
	if !isThree(cards[:3]) {
		return false
	}
	for i := 3; i < len(cards); i += 3 {
		if Score[cards[i].Number]-current != 1 || !isThree(cards[i:i+3]) {
			return false
		}
		current = Score[cards[i].Number]
	}

	return true
}

func IsFourTwoSingle(cards []Card) bool {
	SortCards(cards)
	return isFourTwoSingle(cards)
}

func isFourTwoSingle(cards []Card) bool {
	return len(cards) == 6 && (isNormalBoom(cards[:4]) || isNormalBoom(cards[1:5]) || isNormalBoom(cards[2:]))
}

func IsPlaneTwoSingle(cards []Card) (bool, Card) {
	SortCards(cards)
	return isPlaneTwoSingle(cards)
}

func isPlaneTwoSingle(cards []Card) (bool, Card) {
	val := Card{}
	if len(cards) < 8 || len(cards)%4 != 0 {
		return false, val
	}
	temp := map[byte]int{}
	for _, card := range cards {
		temp[card.Number]++
	}
	p := len(cards) / 4
	singles := make([]Card, 0)
	threes := make([]Card, 0)
	cnt := 0
	for num, c := range temp {
		if c < 3 {
			for _, card := range cards {
				if card.Number == num {
					singles = append(singles, card)
				}
			}
		} else if c == 3 {
			for _, card := range cards {
				if card.Number == num {
					threes = append(threes, card)
				}
			}
			cnt++
		} else {
			idx := true
			for _, card := range cards {
				if card.Number == num {
					if idx {
						singles = append(singles, card)
						idx = !idx
					} else {
						threes = append(threes, card)
					}
				}
			}
			cnt++
		}
	}

	SortCards(threes)
	if cnt < p {
		return false, val
	} else if cnt == p {
		return isPlane(threes), threes[len(threes)-1]
	} else {
		return hasPlane(p, threes)
	}
}

// 查找是否包含指定长度的飞机
func hasPlane(l int, cards []Card) (bool, Card) {
	mc := Card{}
	if len(cards) < 6 || len(cards)%3 != 0 {
		return false, mc
	}

	r := false
	m := 0
	v := int8(0)
	f := true

	for i := 0; i < len(cards); i += 3 {
		if isThree(cards[i : i+3]) {
			if f {
				v = Score[cards[i].Number]
				f = false
				m = 1
			} else {
				if Score[cards[i].Number]-v == 1 {
					m++
				} else {
					m = 1
				}
				v = Score[cards[i].Number]
			}
			if m >= l {
				r = true
				mc = cards[i]
			}
		}
	}

	return r, mc
}

func IsFourTwoPair(cards []Card) (bool, Card) {
	SortCards(cards)
	return isFourTwoPair(cards)
}

func isFourTwoPair(cards []Card) (bool, Card) {
	if isNormalBoom(cards[:4]) && isNormalBoom(cards[4:]) {
		return true, cards[4]
	} else if isNormalBoom(cards[:4]) {
		return isPair(cards[4:6]) && isPair(cards[6:]), cards[0]
	} else if isNormalBoom(cards[2:6]) {
		return isPair(cards[:2]) && isPair(cards[6:]), cards[2]
	} else if isNormalBoom(cards[4:]) {
		return isPair(cards[:2]) && isPair(cards[2:4]), cards[4]
	} else {
		return false, Card{}
	}
}

func IsPlaneTwoPair(cards []Card) (bool, Card) {
	return isPlaneTwoPair(cards)
}

func isPlaneTwoPair(cards []Card) (bool, Card) {
	val := Card{}
	if len(cards) < 10 || len(cards)%5 != 0 {
		return false, val
	}
	temp := map[byte]int{}
	for _, card := range cards {
		temp[card.Number]++
	}
	p := len(cards) / 5
	singles := make([]Card, 0)
	threes := make([]Card, 0)
	cnt := 0
	for num, c := range temp {
		if c == 1 {
			return false, val
		} else if c == 2 || c == 4 {
			for _, card := range cards {
				if card.Number == num {
					singles = append(singles, card)
				}
			}
		} else if c == 3 {
			for _, card := range cards {
				if card.Number == num {
					threes = append(threes, card)
				}
			}
			cnt++
		}
	}

	SortCards(threes)
	if cnt < p {
		return false, val
	} else if cnt == p {
		return isPlane(threes), threes[len(threes)-1]
	} else {
		return hasPlane(p, threes)
	}
}

func ParseHand(cards []Card) (Hand, error) {
	SortCards(cards)
	h := Hand{Cards: cards}
	switch len(cards) {
	case 0:
		return h, nil
	case 1:
		h.Value = cards[0]
		h.Type = HandTypeSingle
	case 2:
		if isKingBoom(cards) {
			h.Value = cards[1]
			h.Type = HandTypeKingBoom
		} else if isPair(cards) {
			h.Value = cards[0]
			h.Type = HandTypePair
		} else {
			return h, ErrHandType
		}
	case 3:
		if isThree(cards) {
			h.Value = cards[0]
			h.Type = HandTypeThree
		} else {
			return h, ErrHandType
		}
	case 4:
		if isNormalBoom(cards) {
			h.Value = cards[0]
			h.Type = HandTypeNormalBoom
		} else if isThreeSingle(cards) {
			h.Value = cards[1]
			h.Type = HandTypeThreeSingle
		} else {
			return h, ErrHandType
		}
	case 5:
		if isThreePair(cards) {
			h.Value = cards[2]
			h.Type = HandTypeThreePair
		} else if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else {
			return h, ErrHandType
		}
	case 6:
		if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else if isStraightPair(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraightPair
		} else if isPlane(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypePlane
		} else if isFourTwoSingle(cards) {
			h.Value = cards[3]
			h.Type = HandTypeFourTwoSingle
		} else {
			return h, ErrHandType
		}
	case 7:
		if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else {
			return h, ErrHandType
		}
	case 8:
		if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else if isStraightPair(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraightPair
		} else if ok, val := isPlaneTwoSingle(cards); ok {
			h.Value = val
			h.Type = HandTypePlaneTwoSingle
		} else if ok, val := isFourTwoPair(cards); ok {
			h.Value = val
			h.Type = HandTypeFourTwoPair
		} else {
			return h, ErrHandType
		}
	case 9:
		if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else if isPlane(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypePlane
		} else {
			return h, ErrHandType
		}
	default:
		if isStraight(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraight
		} else if isStraightPair(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypeStraightPair
		} else if isPlane(cards) {
			h.Value = cards[len(cards)-1]
			h.Type = HandTypePlane
		} else if ok, val := isPlaneTwoSingle(cards); ok {
			h.Value = val
			h.Type = HandTypePlaneTwoSingle
		} else if ok, val := isPlaneTwoPair(cards); ok {
			h.Value = val
			h.Type = HandTypePlaneTwoPair
		} else {
			return h, ErrHandType
		}
	}

	return h, nil
}

func IsGreaterHand(g Hand, l Hand) (bool, error) {
	if l.Type == HandTypeKingBoom {
		return false, nil
	}
	if l.Type == HandTypeNone {
		return true, nil
	}
	// g可以是王炸
	if g.Type == HandTypeKingBoom {
		return true, nil
	}
	// 同类型，比VAL
	if g.Type == l.Type {
		if len(g.Cards) == len(l.Cards) {
			if Score[g.Value.Number] > Score[l.Value.Number] {
				return true, nil
			}
			return false, nil
		}
		return false, ErrHandType
	} else {
		if g.Type == HandTypeNormalBoom {
			return true, nil
		}
		return false, ErrHandType
	}
}
