package core

import (
	"fmt"

	"../connectManager"
)

func StartServer(adress string) {
	go taskAcceptor(connectManager.ConnectionChan)
	connectManager.StartServer(adress)
}

func init() {
	fmt.Println("Create core ")
}

type Position struct {
	X, Y, Angle float64
}

type ClientState struct {
	Pos        Position
	Id         string
	isAttack   bool
	HealfPoint int64
}

type Task struct {
	NewAngle, TryDeltaX, TryDeltaY float64
	tryAttack                      bool
	CoolDown                       int
}

type Player struct {
	ClientState
	Task
}

func taskAcceptor(c chan connectManager.Task) {
	for {
		newTask := <-c
		switch newTask.TaskType {
		case connectManager.ADDCLIENT:
			AddPlayer(newTask.ClientID)
		case connectManager.DELCLIENT:
			fmt.Println("_________________ hui")
		case connectManager.CLIENTMOVE:
			// do something for
		case connectManager.CLIENTSHOOT:
			// do something for
		}
		//AddPlayer(playerID)
	}
}

func AddPlayer(playerID string) {
	fmt.Println("func AddPlayer(playerID string)")
}
func AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}

var players map[string]*Player
