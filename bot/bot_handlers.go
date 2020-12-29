package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"sort"
	"strconv"
)

func (bp *BotProcessor) InfoHandler(m *tb.Message) {
	iMessage := "Список доступных комманд.\n"

	cmdMap := bp.commandsForUser(m.Sender)
	keys := make([]string, 0, len(cmdMap))
	for k := range cmdMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		bh := cmdMap[k]
		iMessage += bh.command + " " + bh.description + "\n"
	}
	bp.Send(m.Sender, iMessage)
}

func (bp *BotProcessor) SubscribeHandler(m *tb.Message) {
	if bp.UserListStorage.HasUser(m.Sender) {
		bp.Send(m.Sender, "Вы уже подписаны на рассылку")
		return
	}
	bp.UserListStorage.AddUser(m.Sender)
	eventMessage := "New sub: #" + strconv.Itoa(m.Sender.ID)
	bp.registerAdminEvent("Registration Event", eventMessage)
	bp.Send(m.Sender, "Вы успешно подписаны на рассылку")
}

func (bp *BotProcessor) UnsubscribeHandler(m *tb.Message) {
	if !bp.UserListStorage.HasUser(m.Sender) {
		bp.Send(m.Sender, "Вы не подписаны на рассылку")
		return
	}
	bp.UserListStorage.AddUser(m.Sender)
	eventMessage := "unsubscribe: #" + strconv.Itoa(m.Sender.ID)
	bp.registerAdminEvent("Registration Event", eventMessage)
	bp.Send(m.Sender, "Вы успешно отписаны от рассылки")
}

func (bp *BotProcessor) GetSubscribersListHandler(m *tb.Message) {
	subs := "Подписчики:\n"
	for _, u := range bp.UserListStorage.GetUsersIterator() {
		subs += strconv.Itoa(u.ID) + "(" + u.Username + "):" + u.FirstName + " " + u.LastName + "\n"
	}
	bp.Send(m.Sender, subs)
}

func (bp *BotProcessor) GetSubscribersCountHandler(m *tb.Message) {
	subCount := strconv.Itoa(bp.UserListStorage.GetUsersAmount())
	bp.Send(m.Sender, "Подписчиков :"+subCount)
}

func (bp *BotProcessor) TestEventHandler(m *tb.Message) {
	bp.Send(m.Sender, "Event registration started \n")

	bp.UserListStorage.AddUser(m.Sender)

	bp.registerEvent(m.Sender, "Test message1", "Привет тест 1")
	bp.registerEvent(m.Sender, "Test message2", "Привет тест 2")
	bp.registerEventForAll("Test message General1", "Привет Общий тест 1")
	bp.registerEvent(m.Sender, "Test message3", "Привет тест 3")
	bp.registerEventForAll("Test message General2", "Привет Общий тест 2")
	bp.registerEvent(m.Sender, "Test message4", "Привет тест 4")
	bp.registerEventForAll("Test message General3", "Привет Общий тест 3")
	bp.registerEventForAll("Test message General4", "Привет Общий тест 4")

}
