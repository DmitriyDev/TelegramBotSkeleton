package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
)

type UserListFileStorage struct {
	userList *UserList
}

func (uls UserListFileStorage) New() *UserListFileStorage {
	uls.userList = UserList{}.New()
	return &uls
}

func (uls UserListFileStorage) GetUserList() *UserList {
	return uls.userList
}

func (uls UserListFileStorage) AddUser(u *tb.User) {
	uls.GetUserList().append(*u)
}

func (uls UserListFileStorage) RemoveUser(u *tb.User) {
	uls.GetUserList().remove(*u)
}

func (uls UserListFileStorage) GetById(id int) (tb.User, error) {
	return uls.GetUserList().byId(id)
}
