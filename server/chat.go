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

// ConnMap holds a map of the connected users
type ConnMap struct {
	connMap map[string]*websocket.Conn
	sync.RWMutex
}

// Message holds information about messages
type Message struct {
	Sender string `json:"sender"`
	ConvoID uint64 `json:"convoid"`
	Msg    string `json:"msg"`
	TimeStamp int64 `json:"timestamp"`
}

// Conversation holds the messages and users in a conversation
type Conversation struct {
	ID uint64
	Users []*User
	Messages []Message
	sync.RWMutex
}

type ConvoMap struct {
	sync.RWMutex
}

// UsageTracker keeps track of the total time spent sending messages and
// the number of messages sent
type UsageTracker struct {
	time int64 // Total time in milliseconds
	num  int64 // Number of messages sent
	sync.Mutex
}

// AddTime adds the input time to the total UsageTracker and
// increments the number of messages
func (ut *UsageTracker) AddTime(t int64) {
	ut.Lock()
	defer ut.Unlock()
	ut.time += t
	ut.num++
}

// Reset resets the usage tracker values to 0 and returns the values
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
	lowerBound float32       = 100    // If the usage time is below this threshold, a hub is removed; in ms
)

var (
	chatLogger *log.Logger
	hubList    SLList
	connMap      ConnMap
	convoMap ConvoMAp
	usage      UsageTracker // Tracks the number of messages send and time to send them
)

func init() {
	chatLogger = log.New(os.Stdout, "Chat Server: ", log.LstdFlags)

	connMap.conns = make(map[string]*websocket.Conn)

	c := make(chan *Message, chanBuffer)
	hubList.Append(c)
	go startChatHub(c, hubList.Length())
	go monitorUsage()
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

// Adds a user to the map of connMap
func addUser(username string, ws *websocket.Conn) {
	// Do a full lock since the user will always be added,
	// even if they are already logged in
	connMap.Lock()
	defer connMap.Unlock()
	// Check to see if the user is logged in elsewhere
	// If so, let that socket know that they are being logged out
	user := connMap.conns[username]
	if user != nil {
		user.Write([]byte("Signed in elsewhere"))
	}
	connMap.conns[username] = ws
}

// Removes a user from the map of connMap
func removeUser(username string, ws *websocket.Conn) {
	// Use an RLock to check if the user to be removed is the same as the one
	// in the map by comparing websockets
	connMap.RLock()
	if connMap.conns[username] == ws {
		// If the connMap are the same, do a full lock and check one more time
		connMap.Lock()
		if connMap.conns[username] == ws {
			// If they are still the same, delete the user from the map
			delete(connMap.conns, username)
		}
		connMap.Unlock()
	}
	connMap.RUnlock()
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
		// Make sure the connMap map is safe for reading
		connMap.RLock()
		for _, user := range msg.to {
			// Check to see if the user is connected
			// If so, send the message
			ws := connMap.conns[user]
			if ws != nil {
				websocket.JSON.Send(ws, *msg)
			}
			// Database message
		}
		connMap.RUnlock()
	}
}

func monitorUsage() {
	useChan := make(chan float32)
	timer := time.AfterFunc(time.Minute*mins, func() {
		ms, n := usage.Reset()
<<<<<<< HEAD
		if n > 0 {
			useChan <- float32(ms) / float32(n)
		} else {
			useChan <- 0
		}
=======
    if n > 0 {
		  useChan <- float32(ms) / float32(n)
    } else {
      useChan <- 0
    }
>>>>>>> 57400ee8095253ca17f362f4866a4070c362bf27
	})
	for {
		select {
		case avg := <-useChan:
      log.Printf("Average time/msg: %f ms\tNum hubs: %d\n", avg, hubList.Length())
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
