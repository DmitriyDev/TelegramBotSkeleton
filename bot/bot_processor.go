package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"reflect"
	"time"
)

type BotProcessor struct {
	bot      *tb.Bot
	handlers map[string]BotHandler
	UserList UserList
}

type BotHandler struct {
	name        string
	command     string
	handler     interface{}
	description string
	isAdmin     bool
	isVisible   bool
}

func GetBotProcessor(token string) BotProcessor {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	bp := BotProcessor{
		b,
		map[string]BotHandler{},
		UserList{},
	}
	bp.createHandlers()

	return bp
}

func (bp *BotProcessor) createHandlers() {
	c := GetConfig()
	for _, cc := range c.Commands {
		toLog(fmt.Sprintf("%s - %s : handler created", cc.Name, cc.HandlerMethod))
		bp.createHandlerFromConfig(cc)
	}
}

func (bp *BotProcessor) createHandlerFromConfig(cc HandlerConfig) {
	handlerFunc := func(m *tb.Message) {

		if cc.Admin == true && m.Sender.ID != MASTER_USER {
			bp.Send(m.Sender, "Данная команда запрещена")
			return
		}
		toLog(fmt.Sprintf("%s : Called", cc.Name))

		in := []reflect.Value{reflect.ValueOf(m)}
		reflect.ValueOf(bp).MethodByName(cc.HandlerMethod).Call(in)
	}
	bh := BotHandler{cc.Name,
		cc.Command,
		handlerFunc,
		cc.Description,
		cc.Admin,
		cc.Visible,
	}
	bp.handlers[bh.name] = bh
	bp.bot.Handle(bh.command, bp.handlers[bh.name].handler)
}

func (bp *BotProcessor) Send(u *tb.User, message string) {
	bp.bot.Send(u, message)
	toLog("Messages sent to " + u.FirstName)
}

func (bp *BotProcessor) SendMany(ul *UserList, message string) {
	for _, u := range ul.users {
		bp.Send(&u, message)
		toLog("Messages sent to " + u.FirstName)
	}
}

func (bp *BotProcessor) SendToAdmin(message string) {
	u, err := bp.UserList.byId(MASTER_USER)
	if err != nil {
		bp.Send(&u, "Admin id not set")
		toLogFatal(fmt.Sprintf("User not found. %v", err))
		return
	}
	bp.Send(&u, message)
}

func (bp *BotProcessor) commandsForUser(u *tb.User) []BotHandler {
	uHandlers := []BotHandler{}
	for _, bh := range bp.handlers {
		if !isAdmin(*u) && (bh.isAdmin == true || bh.isVisible == false) {
			continue
		}
		uHandlers = append(uHandlers, bh)
	}

	return uHandlers
}

func (bp *BotProcessor) InfoHandler(m *tb.Message) {
	iMessage := "Список доступных комманд.\n"
	for _, bh := range bp.commandsForUser(m.Sender) {
		iMessage += bh.command + " " + bh.description + "\n"
	}
	bp.Send(m.Sender, iMessage)
}
