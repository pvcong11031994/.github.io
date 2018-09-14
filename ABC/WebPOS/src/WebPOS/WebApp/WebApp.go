package WebApp

import (
	"WebPOS/Models/ModelItems"
	"github.com/goframework/gf"
)

const (
	_CONTEXT_USER_INFO_KEY = "_CONTEXT_USER_INFO_KEY"
)

func GetContextUser(ctx *gf.Context) *ModelItems.UserItem {
	userData, ok := ctx.Get(_CONTEXT_USER_INFO_KEY)
	if !ok {
		return nil
	}

	user := userData.(*ModelItems.UserItem)
	return user
}

func SetContextUser(ctx *gf.Context, user *ModelItems.UserItem) {
	ctx.Set(_CONTEXT_USER_INFO_KEY, user)
}

func GetSessionUser(ctx *gf.Context) *ModelItems.UserItem {
	userData, ok := ctx.Session.Values[_CONTEXT_USER_INFO_KEY]
	if !ok {
		return nil
	}

	user := userData.(ModelItems.UserItem)
	return &user
}

func SetSessionUser(ctx *gf.Context, user *ModelItems.UserItem) {
	ctx.Session.Values[_CONTEXT_USER_INFO_KEY] = user
}
