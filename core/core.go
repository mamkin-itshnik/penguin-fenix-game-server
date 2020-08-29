package core

import (
	"fmt"

	"../connectManager"
)

func StartServer(adress string) error {
	go clientAccepter(connectManager.ConnectionChan)
	return connectManager.StartServer(adress)
}

func Hello() {
	fmt.Println("Hello, World!")
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

func clientAccepter(c chan string) {
	for {
		playerID := <-c
		AddPlayer(playerID)
	}
}

func AddPlayer(playerID string) {
	fmt.Println("func AddPlayer(playerID string)")
}
func AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}

var players map[string]*Player
