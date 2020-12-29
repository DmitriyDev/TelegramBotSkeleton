package main

import (
	"os"
	"strconv"
)

var TOKEN = os.Getenv("BOT_TOKEN")
var MASTER_USER, _ = strconv.Atoi(os.Getenv("MASTER_USER"))

const ConfigFilePath = "./config.yaml"
const TimeLayout = "2006-01-02 15:04"

func main() {
	eh := EventHandler{}.New()
	uls := UserListFileStorage{}.New()


	//bp := BotProcessor{Token: TOKEN, eventHandler: eh, UserListStorage: UserListFileStorage{}.New()}
	bp := BotProcessor{}.New(TOKEN, eh, uls)


	//Token: TOKEN, eventHandler: eh, UserListStorage: UserListFileStorage{}.New()}
	bp.Start()
}
