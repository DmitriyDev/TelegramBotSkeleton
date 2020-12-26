package main

import (
	"os"
	"strconv"
)

var TOKEN = os.Getenv("BOT_TOKEN")
var MASTER_USER, _ = strconv.Atoi(os.Getenv("MASTER_USER"))
var CONFIG_FILE_PATH = "./config.yaml"
var TIME_LAYOUT = "2006-01-02 15:04"

func main() {
	bp := GetBotProcessor(TOKEN)
	bp.bot.Start()
}
