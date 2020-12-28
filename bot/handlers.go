package main

import tb "gopkg.in/tucnak/telebot.v2"

func (bp *BotProcessor) InfoHandler(m *tb.Message) {
	iMessage := "Список доступных комманд.\n"
	for _, bh := range bp.commandsForUser(m.Sender) {
		iMessage += bh.command + " " + bh.description + "\n"
	}
	bp.Send(m.Sender, iMessage)
}

func (bp *BotProcessor) TestEventHandler(m *tb.Message) {
	bp.Send(m.Sender, "Event registration started \n")

	bp.UserList.append(*m.Sender)

	bp.registerEvent(m.Sender, "Test message1", "Привет тест 1")
	bp.registerEvent(m.Sender, "Test message2", "Привет тест 2")
	bp.registerEventForAll("Test message General1", "Привет Общий тест 1")
	bp.registerEvent(m.Sender, "Test message3", "Привет тест 3")
	bp.registerEventForAll("Test message General2", "Привет Общий тест 2")
	bp.registerEvent(m.Sender, "Test message4", "Привет тест 4")
	bp.registerEventForAll("Test message General3", "Привет Общий тест 3")
	bp.registerEventForAll("Test message General4", "Привет Общий тест 4")

}
