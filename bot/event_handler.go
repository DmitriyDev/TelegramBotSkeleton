package main

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

type EventHandler struct {
	events map[int]*UserEventBox
}

type UserEventBox struct {
	user   *tb.User
	events map[int64]InternalEvent
}

type InternalEvent struct {
	name      string
	message   string
	createdAt time.Time
}

func (ub *UserEventBox) AddEvent(e InternalEvent) {
	id := time.Now().UnixNano()
	ub.events[id] = e
}

func (ub *UserEventBox) DeleteEvent(id int64) {
	delete(ub.events, id)
}

func (ub *UserEventBox) GetEvents() map[int64]InternalEvent {
	return ub.events
}

func (ub *UserEventBox) HastEvents() bool {
	return ub.EventsCount() > 0
}

func (ub *UserEventBox) EventsCount() int {
	return len(ub.events)
}

func (eh EventHandler) New() *EventHandler {
	eh.events = map[int]*UserEventBox{}
	return &eh
}

func (eh *EventHandler) RegisterUser(u *tb.User) {
	eh.events[u.ID] = &UserEventBox{u, map[int64]InternalEvent{}}
}

func (eh *EventHandler) UnregisterUser(u tb.User) {
	delete(eh.events, u.ID)
}

func (eh *EventHandler) AddEventToUser(e InternalEvent, u *tb.User) {
	if _, ok := eh.events[u.ID]; !ok {
		eh.RegisterUser(u)
	}
	eh.events[u.ID].AddEvent(e)
}

func (eh *EventHandler) AddEventToAll(e InternalEvent, ul UserList) {
	for _, u := range ul.users {
		eh.AddEventToUser(e, &u)
	}
}

func (eh *EventHandler) GetUserEvents(u *tb.User) *UserEventBox {
	if _, ok := eh.events[u.ID]; !ok {
		eh.RegisterUser(u)
	}
	return eh.events[u.ID]
}

func (eh *EventHandler) GetAllEvents() map[int]*UserEventBox {
	return eh.events
}

func (eh *EventHandler) UserHasEvents(u tb.User) bool {
	if _, ok := eh.events[u.ID]; !ok {
		return false
	}
	return eh.events[u.ID].HastEvents()
}

func (eh *EventHandler) HasUnfinishedEvents() bool {
	for _, ub := range eh.events {
		if ub.HastEvents() {
			return true
		}
	}

	return false
}
