package main

/*
 * Possibly have a chat id as the query
 * Have a special value in the queury for the bot
 */

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"golang.org/x/net/websocket"
)

type UserMap struct {
	users map[string]*websocket.Conn
	sync.RWMutex
}

type Message struct {
	Sender string `json:"sender"`
	to     []string
	Msg    string `json:"msg"`
}

const (
	runSecondHub bool = false
)

var (
	chatLogger *log.Logger
	hub1Chan   chan *Message
	hub2Chan   chan *Message
	users      UserMap
)

func init() {
	chatLogger = log.New(os.Stdout, "Chat Server: ", log.LstdFlags)
	users.users = make(map[string]*websocket.Conn)
	hub1Chan = make(chan *Message, 100)
	go startChatHub(hub1Chan)
	if runSecondHub {
		hub2Chan = make(chan *Message, 100)
		go startChatHub(hub2Chan)
	}
}

func chatSocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	// check := func(err error) bool {
	//   if err != nil {
	//     chatLogger.Println(err)
	//     ws.Write([]byte("ERROR"))
	//   }
	//   return err != nil
	// }

	// Get username
	var username string
	err := websocket.Message.Receive(ws, &username)
	if err != nil {
		chatLogger.Println(err)
		return
	}
	defer removeUser(username, ws)

	for {
		// Read bytes to get different types of messages and hvae more control
		var bmsg [2048]byte
		if l, err := ws.Read(bmsg[:]); err != nil {
			if err.Error() == "EOF" {
				return
			}
		} else {
			msg := newMsg(bmsg[:l])
			if runSecondHub {
				select {
				case hub1Chan <- msg:
				case hub2Chan <- msg:
				}
			} else {
				hub1Chan <- msg
			}
		}
	}
}

func startChatHub(hubChan chan *Message) {
	// Loop through the channel
	for msg := range hubChan {
		// Make sure the users map is safe for reading
		users.RLock()
		for _, user := range msg.to {
			// Check to see if the user is connected
			// If so, send the message
			ws := users.users[user]
			if ws != nil {
				websocket.JSON.Send(ws, *msg)
			}
			// Database message
		}
		users.RUnlock()
	}
}

// Adds a user to the map of users
func addUser(username string, ws *websocket.Conn) {
	// Do a full lock since the user will always be added,
	// even if they are already logged in
	users.Lock()
	defer users.Unlock()
	// Check to see if the user is logged in elsewhere
	// If so, let that socket know that they are being logged out
	user := users.users[username]
	if user != nil {
		user.Write([]byte("Signed in elsewhere"))
	}
	users.users[username] = ws
}

// Removes a user from the map of users
func removeUser(username string, ws *websocket.Conn) {
	// Use an RLock to check if the user to be removed is the same as the one
	// in the map by comparing websockets
	users.RLock()
	if users.users[username] == ws {
		// If the users are the same, do a full lock and check one more time
		users.Lock()
		if users.users[username] == ws {
			// If they are still the same, delete the user from the map
			delete(users.users, username)
		}
		users.Unlock()
	}
	users.RUnlock()
}

// Creates a new message
func newMsg(bytes []byte) *Message {
	var msg *Message
	err := json.Unmarshal(bytes, msg)
	if err != nil {
		chatLogger.Println(err)
		return nil
	}
	return msg
}
