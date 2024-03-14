package user

import (
	"ddz/model"
	"testing"
)

func TestOutCards(t *testing.T) {
	user := NewUser("zhangsan")
	user.AddCards([]model.Card{{'3', 1}, {'4', 1}, {'5', 2}, {'6', 1}, {'6', 2}, {'6', 3}}...)
	user.outCards(model.Hand{Cards: []model.Card{{'3', 0}, {'6', 0}, {'6', 0}, {'6', 0}}})
}
