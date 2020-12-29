package main

import (
	"encoding/json"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type UserListFileStorage struct {
	userList      *UserList
	storageFolder string
}

func (uls UserListFileStorage) New(storageFolder string) *UserListFileStorage {
	uls.userList = UserList{}.New()
	uls.storageFolder = storageFolder
	uls.loadUsers()
	return &uls
}

func (uls *UserListFileStorage) GetUsersIterator() map[int]tb.User {
	return uls.userList.users
}

func (uls *UserListFileStorage) GetUsersAmount() int {
	return len(uls.userList.users)
}

func (uls *UserListFileStorage) GetUserList() *UserList {
	return uls.userList
}

func (uls *UserListFileStorage) AddUser(u *tb.User) {
	uls.GetUserList().append(*u)
	uls.dumpUserFile(*u)
}

func (uls *UserListFileStorage) HasUser(u *tb.User) bool {
	_, err := uls.GetUserList().byId(u.ID)
	return err == nil
}

func (uls *UserListFileStorage) RemoveUser(u *tb.User) {
	uls.GetUserList().remove(*u)
	uls.removeUserFile(*u)
}

func (uls *UserListFileStorage) GetById(id int) (tb.User, error) {
	return uls.GetUserList().byId(id)
}

func (uls *UserListFileStorage) getUserPath(userId int) string {
	return uls.storageFolder + strconv.Itoa(userId) + ".json"
}

func (uls *UserListFileStorage) extractUserByPath(path string) tb.User {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Get user by file error: ", err)
	}

	user := tb.User{}
	json.Unmarshal(content, &user)
	return user
}

func (uls *UserListFileStorage) extractUserById(id int) tb.User {
	userPath := uls.getUserPath(id)
	return uls.extractUserByPath(userPath)
}

func (uls *UserListFileStorage) loadUsers() {
	dir, _ := os.Open(uls.storageFolder)
	files, _ := dir.Readdir(0)

	for i := range files {
		file := files[i]
		path := uls.storageFolder + file.Name()
		if !file.IsDir() && filepath.Ext(path) == ".json" {
			user := uls.extractUserByPath(path)
			uls.userList.append(user)
		}
	}
}

func (uls *UserListFileStorage) dumpUserFile(u tb.User) {
	userPath := uls.getUserPath(u.ID)
	if fileExists(userPath) {
		return
	}

	uJson, _ := json.Marshal(u)
	f, err := os.Create(userPath)
	check(err)
	f.Write(uJson)
	f.Close()
}

func (uls *UserListFileStorage) removeUserFile(u tb.User) {
	userPath := uls.getUserPath(u.ID)
	if !fileExists(userPath) {
		return
	}
	err := os.Remove(userPath)
	check(err)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func check(e error) {
	if e != nil {
		toLogFatalF("%v", e)
	}
}
