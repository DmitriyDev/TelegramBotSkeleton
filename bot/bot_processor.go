package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type BotProcessor struct {
	Token           string
	bot             *tb.Bot
	handlers        map[string]BotHandler
	eventHandler    *EventHandler
	UserListStorage UserListStorage
}

type BotHandler struct {
	name        string
	command     string
	handler     interface{}
	description string
	isAdmin     bool
	isVisible   bool
}

func (bp BotProcessor) New(token string, eh *EventHandler, uls UserListStorage) *BotProcessor {
	bp.Token = token
	bp.eventHandler = eh
	bp.UserListStorage = uls
	bp.handlers = map[string]BotHandler{}

	return &bp
}

func (bp *BotProcessor) Start() {
	var err error
	bp.bot, err = tb.NewBot(tb.Settings{
		Token:  bp.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	go bp.internalEvenHandlerListener()
	bp.createHandlers()

	bp.bot.Start()
}

func (bp *BotProcessor) createHandlers() {
	c := GetConfig()
	for _, cc := range c.Commands {
		toLogF("%s - %s : handler created", cc.Name, cc.HandlerMethod)
		bp.createHandlerFromConfig(cc)
	}
}

func (bp *BotProcessor) createHandlerFromConfig(cc HandlerConfig) {
	handlerFunc := func(m *tb.Message) {

		if cc.Admin == true && m.Sender.ID != MASTER_USER {
			bp.Send(m.Sender, "Данная команда запрещена")
			return
		}
		toLogF("%s : Called", cc.Name)

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
	u, err := bp.UserListStorage.GetById(MASTER_USER)
	if err != nil {
		bp.Send(&u, "Admin id not set")
		toLogFatalF("User not found. %v", err)
		return
	}
	bp.Send(&u, message)
}

func (bp *BotProcessor) commandsForUser(u *tb.User) map[string]BotHandler {
	uHandlers := map[string]BotHandler{}

	for k, bh := range bp.handlers {
		if !isAdmin(*u) && (bh.isAdmin == true || bh.isVisible == false) {
			continue
		}
		uHandlers[k] = bh
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

func (bp *BotProcessor) registerAdminEvent(name string, message string) {
	ie := InternalEvent{
		name,
		message,
		time.Now(),
	}
	u := bp.getAdminUser()
	bp.eventHandler.AddEventToUser(ie, &u)
	toLog("Register event " + name + " for user admin")
}

func (bp *BotProcessor) getAdminUser() tb.User {
	u, err := bp.UserListStorage.GetById(MASTER_USER)
	if err != nil {
		toLogFatal("Admin user not set")
	}
	return u
}

func (bp *BotProcessor) registerEventForAll(name string, message string) {
	ie := InternalEvent{
		name,
		message,
		time.Now(),
	}

	bp.eventHandler.AddEventToAll(ie, *bp.UserListStorage.GetUserList())
	toLog("Register event " + name + " for all users")
}

func (bp *BotProcessor) internalEvenHandlerListener() {
	for true {
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
				toLogF("event (#%d) %s finished", id, ub.name)
				ue.DeleteEvent(id)
				toLogF("event (#%d) removed from queue", id)
			}
		}

	}
}
