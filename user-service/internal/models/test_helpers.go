package models

import "github.com/stretchr/testify/mock"

func AnyUser() interface{} {
	return mock.MatchedBy(func(u *User) bool {
		return u != nil
	})
}