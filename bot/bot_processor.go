package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type BotProcessor struct {
	Token        string
	bot          *tb.Bot
	handlers     map[string]BotHandler
	UserList     UserList
	eventHandler EventHandler
}

type BotHandler struct {
	name        string
	command     string
	handler     interface{}
	description string
	isAdmin     bool
	isVisible   bool
}

func (bp *BotProcessor) init(token string, ul UserList) {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}
	eh := EventHandler{}
	eh.init()

	bp.bot = b
	bp.UserList = ul
	bp.handlers = map[string]BotHandler{}
	bp.eventHandler = eh

	go bp.internalEvenHandlerListener()

	bp.createHandlers()
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

func (bp *BotProcessor) registerEvent(u *tb.User, name string, message string) {
	ie := InternalEvent{
		name,
		message,
		time.Now(),
	}

	bp.eventHandler.AddEventToUser(ie, u)
	toLog("Register event " + name + " for user #" + strconv.Itoa(u.ID))
}

func (bp *BotProcessor) registerEventForAll(name string, message string) {
	ie := InternalEvent{
		name,
		message,
		time.Now(),
	}

	bp.eventHandler.AddEventToAll(ie, &bp.UserList)
	toLog("Register event " + name + " for all users")
}

func (bp *BotProcessor) internalEvenHandlerListener() {
	for true {
		time.Sleep(time.Second)

		if !bp.eventHandler.HasUnfinishedEvents() {
			continue
		}

		for _, ue := range bp.eventHandler.GetAllEvents() {

			keys := make([]int64, 0, len(ue.events))

			for k := range ue.events {
				keys = append(keys, k)
			}
			sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

			for i := range keys {
				id := keys[i]
				ub := ue.events[id]
				bp.Send(ue.user, ub.message)
				toLog(fmt.Sprintf("event (#%d) %s finished", id, ub.name))
				ue.DeleteEvent(id)
				toLog(fmt.Sprintf("event (#%d) removed from queue", id))
			}
		}

	}
}
