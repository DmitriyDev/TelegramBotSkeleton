package main

import (
	"fmt"
	"os"
	"strconv"
)

var TOKEN = os.Getenv("BOT_TOKEN")
var MASTER_USER, _ = strconv.Atoi(os.Getenv("MASTER_USER"))

const ConfigFilePath = "./config.yaml"
const UserStorageFolder = "./u/"
const TimeLayout = "2006-01-02 15:04"

func main() {
	eh := EventHandler{}.New()
	uls := UserListFileStorage{}.New(UserStorageFolder)
	fmt.Printf("%v", uls)
	bp := BotProcessor{}.New(TOKEN, eh, uls)
	bp.Start()
}
