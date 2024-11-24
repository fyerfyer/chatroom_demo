package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var serverURL = "ws://localhost:8000/ws?"
var clientname string

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	flag.StringVar(&clientname, "name", "Alice", "the chatroom login name")
	flag.Parse()
	conn, _, err := websocket.Dial(ctx, serverURL+"name="+clientname, nil)
	defer conn.Close(websocket.StatusAbnormalClosure, "close client connection")
	if err != nil {
		log.Fatal("Failed to set up client")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter message (or 'exit' to quit): ")
		input, _ := reader.ReadString('\n')

		input = input[:len(input)-1]
		if input == "exit" {
			fmt.Println("Exit client side")
			break
		}

		msg := map[string]interface{}{
			"id":         1,
			"name":       clientname,
			"created_at": time.Now(),
			"IsNew":      true,
			"content":    input,
		}

		err := wsjson.Write(ctx, conn, msg)
		if err != nil {
			log.Println("Failed to send message: ", err)
			break
		}

		err = wsjson.Read(ctx, conn, &msg)
		if err != nil {
			log.Println("Failed to receive message: ", err)
		}

		log.Println("Received message: ", msg)
	}
}
