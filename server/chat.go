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
	"time"

	"golang.org/x/net/websocket"
)

type ConnMap struct {
	conns map[string]*websocket.Conn
	sync.RWMutex
}

type Message struct {
	Sender string `json:"sender"`
	to     []string
	Msg    string `json:"msg"`
}

type UsageTracker struct {
	time int64 // Total time in milliseconds
	num  int64 // Number of messages sent
	sync.Mutex
}

func (ut *UsageTracker) AddTime(t int64) {
	ut.Lock()
	defer ut.Unlock()
	ut.time += t
	ut.num++
}

func (ut *UsageTracker) Reset() (t, n int64) {
	ut.Lock()
	defer ut.Unlock()
	t, n = ut.time, ut.num
	ut.time = 0
	ut.num = 0
	return
}

const (
	chanBuffer int           = 100  // The number of messages a hub can have buffered
	mins       time.Duration = 1    // Number of minutes to wait before checking usage
	upperBound float32       = 1000 // If the usage time is above this threshold, another hub is added; in ms
	lowerBound float32       = 0    // If the usage time is below this threshold, a hub is removed; in ms
)

var (
	chatLogger *log.Logger
	hubList    SLList
	conns      ConnMap
	usage      UsageTracker // Tracks the number of messages send and time to send them
)

func init() {
	chatLogger = log.New(os.Stdout, "Chat Server: ", log.LstdFlags)
	conns.conns = make(map[string]*websocket.Conn)
	c := make(chan *Message, chanBuffer)
	hubList.Append(c)
	go startChatHub(c, hubList.Length())
}

func check(err error) bool {
	if err != nil {
		chatLogger.Println(err)
		// ws.Write([]byte("ERROR"))
	}
	return err != nil
}

func chatSocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	// Get username
	var username string
	err := websocket.Message.Receive(ws, &username)
	if err != nil {
		if err.Error() != "EOF" {
			chatLogger.Println(err)
		}
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
			sendMessage(msg)
		}
	}
}

func sendMessage(msg *Message) {
	hubList.RLock()
	defer hubList.RUnlock()
	start := time.Now()
	for {
		for node := hubList.head; node != nil; node = node.next {
			c := node.data.(chan *Message)
			select {
			case c <- msg:
				d := time.Since(start)
				usage.AddTime(d.Milliseconds())
				return
			default:
			}
		}
	}
}

// Adds a user to the map of conns
func addUser(username string, ws *websocket.Conn) {
	// Do a full lock since the user will always be added,
	// even if they are already logged in
	conns.Lock()
	defer conns.Unlock()
	// Check to see if the user is logged in elsewhere
	// If so, let that socket know that they are being logged out
	user := conns.conns[username]
	if user != nil {
		user.Write([]byte("Signed in elsewhere"))
	}
	conns.conns[username] = ws
}

// Removes a user from the map of conns
func removeUser(username string, ws *websocket.Conn) {
	// Use an RLock to check if the user to be removed is the same as the one
	// in the map by comparing websockets
	conns.RLock()
	if conns.conns[username] == ws {
		// If the conns are the same, do a full lock and check one more time
		conns.Lock()
		if conns.conns[username] == ws {
			// If they are still the same, delete the user from the map
			delete(conns.conns, username)
		}
		conns.Unlock()
	}
	conns.RUnlock()
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

func startChatHub(hubChan chan *Message, num int) {
	chatLogger.Printf("Starting hub %d\n", num)
	defer chatLogger.Printf("Hub %d stopped\n", num)
	// Loop through the channel
	for msg := range hubChan {
		// Make sure the conns map is safe for reading
		conns.RLock()
		for _, user := range msg.to {
			// Check to see if the user is connected
			// If so, send the message
			ws := conns.conns[user]
			if ws != nil {
				websocket.JSON.Send(ws, *msg)
			}
			// Database message
		}
		conns.RUnlock()
	}
}

func monitorUsage() {
	useChan := make(chan float32)
	timer := time.AfterFunc(time.Minute*mins, func() {
		ms, n := usage.Reset()
		if n > 0 {
			useChan <- float32(ms) / float32(n)
		} else {
			useChan <- 0
		}
	})
	for {
		select {
		case avg := <-useChan:
			if avg < lowerBound {
				if hubList.Length() != 1 {
					node := hubList.PopLast()
					c := node.data.(chan *Message)
					close(c)
				}
			} else if avg > upperBound {
				c := make(chan *Message, chanBuffer)
				hubList.Append(c)
				go startChatHub(c, hubList.Length())
			}
			timer.Reset(time.Minute * mins)
		}
	}
}
