package setting

import (
	"log"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	HTTPPort           string
	MessageQueueLength int
	OfflineMsgNum      int
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/chatroom.ini")
	if err != nil {
		log.Fatalf("Failed to parse 'conf/chatroom.ini': %v", err)
	}

	var chatroom = Cfg.Section("chatroom")
	HTTPPort = Cfg.Section("server").
		Key("HTTP_PORT").String()

	MessageQueueLength = chatroom.
		Key("Message_Queue_Length").
		MustInt(1024)
	OfflineMsgNum = chatroom.
		Key("Offline_Message_Num").
		MustInt(10)
}
