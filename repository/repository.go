package repository

import (
	"container/list"
	"math/rand"
	"time"
)

type CmdType int

const (
	GetViews = iota
	View
	Update
)

// Command is a representation of commands sent to repository manager
type Command struct {
	Typ        CmdType
	Newbie     string
	ReqCh      chan string
	AnalyticCh chan map[string]int
}

// StartRepositoryManager initializes underlying storage containers and handle Commands received from returned channel
func StartRepositoryManager(clientsCapacity int) chan Command {
	seen := make(map[string]int)
	pending := list.New()
	declined := list.New()
	cmds := make(chan Command, clientsCapacity)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 50; i++ {
		pending.PushBack(newReq())
	}

	SpawnReqs(cmds)

	go func() {
		for cmd := range cmds {
			switch cmd.Typ {
			case GetViews:
				all := list.New()
				all.PushBackList(pending)
				all.PushBackList(declined)
				resMap := make(map[string]int)
				for e := all.Front(); e != nil; e = e.Next() {
					resMap[e.Value.(string)] = seen[e.Value.(string)]
				}
				cmd.AnalyticCh <- resMap
			case View:
				r := pending.Front().Value.(string)
				seen[r] += 1
				cmd.ReqCh <- r
			case Update:
				victim := pending.Front()
				if v, ok := victim.Value.(string); ok {
					declined.PushBack(v)
				}
				pending.Remove(victim)
				pending.PushBack(cmd.Newbie)
			}

		}
	}()

	return cmds
}

// SpawnReqs creates a new random req and send it to cmds channel when Ticker ticks
func SpawnReqs(cmds chan Command) {
	ticker := time.NewTicker(200 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				cmds <- Command{
					Typ:    Update,
					Newbie: newReq(),
				}
			}
		}
	}()
}

func newReq() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	r := make([]rune, 2)
	for i := range r {
		r[i] = letters[rand.Intn(len(letters))]
	}
	return string(r)
}
