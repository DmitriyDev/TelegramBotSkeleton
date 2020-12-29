package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
)

type UserListStorage interface {
	GetUsersIterator() map[int]tb.User
	GetUserList() *UserList
	GetUsersAmount() int
	AddUser(u *tb.User)
	HasUser(u *tb.User) bool
	RemoveUser(u *tb.User)
	GetById(id int) (tb.User, error)
}

type UserList struct {
	users map[int]tb.User
}

func (ul UserList) New() *UserList {
	ul.users = map[int]tb.User{}
	return &ul
}
func (ul *UserList) append(u tb.User) {
	ul.users[u.ID] = u
}

func (ul *UserList) remove(u tb.User) {
	delete(ul.users, u.ID)
}

func (ul *UserList) byId(id int) (tb.User, error) {
	if u, ok := ul.users[id]; ok {
		return u, nil
	}
	return tb.User{}, fmt.Errorf("User with ID %d not found", id)
}

func isAdmin(u tb.User) bool {
	return u.ID == MASTER_USER
}
