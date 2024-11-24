package models

import (
	"log"

	"github.com/fyerfyer/chatroom/pkg/setting"
)

type broadcast struct {
	users map[string]*User

	loginChannel   chan *User
	loggoutChannel chan *User
	messageChannel chan *Message

	// user judgement
	checkUserChannel  chan string
	checkUserCanLogin chan bool

	// get user list
	requestUsersChannel chan struct{}
	listUsersChannel    chan []*User
}

var Broadcaster = &broadcast{
	users:               make(map[string]*User),
	loginChannel:        make(chan *User),
	loggoutChannel:      make(chan *User),
	messageChannel:      make(chan *Message, setting.MessageQueueLength),
	checkUserChannel:    make(chan string),
	checkUserCanLogin:   make(chan bool),
	requestUsersChannel: make(chan struct{}),
}

func (b *broadcast) Start() {
	for {
		select {
		// case 1: user enter in
		case user := <-b.loginChannel:
			b.users[user.Name] = user
			// tell the user the recent messages
			UserMessageProcessor.Send(user)
		// case 2: user exit
		case user := <-b.loggoutChannel:
			delete(b.users, user.Name)
			// close the goroutine for the user
			// because its goroutine is continuous open
			user.CloseChannel()
		//case 3: send message
		case msg := <-b.messageChannel:
			// send the message to all users
			for _, user := range b.users {
				if user.ID == msg.User.ID {
					continue
				}

				user.MessageChannel <- msg
			}
			UserMessageProcessor.Save(msg)
		// case 4: check if user has already login
		case name := <-b.checkUserChannel:
			if _, ok := b.users[name]; ok {
				b.checkUserCanLogin <- false
			} else {
				b.checkUserCanLogin <- true
			}
		// case 5: get the user list
		case <-b.requestUsersChannel:
			usersList := make([]*User, 0, len(b.users))
			for _, user := range b.users {
				usersList = append(usersList, user)
			}

			b.listUsersChannel <- usersList
		}
	}
}

func (b *broadcast) UserLogin(user *User) {
	b.loginChannel <- user
}

func (b *broadcast) UserLoggout(user *User) {
	b.loggoutChannel <- user
}

func (b *broadcast) CheckUserCanLogin(name string) bool {
	b.checkUserChannel <- name
	return <-b.checkUserCanLogin
}

func (b *broadcast) GetUserList() []*User {
	b.requestUsersChannel <- struct{}{}
	return <-b.listUsersChannel
}

func (b *broadcast) Broadcast(msg *Message) {
	if len(b.messageChannel) >= setting.MessageQueueLength {
		log.Println("the broadcast queue has been full!")
	}

	b.messageChannel <- msg
}
