package main

import (
	"fmt"
	"time"
)

func core_StartServer(adress string) {
	go core_taskAcceptor(TaskChan)
	go core_TicTack()
	CN_StartServer(adress)
}

func init() {
	fmt.Println("Create core ")
	players = make(map[string]Player)
}

func core_TicTack() {
	for {
		time.Sleep(time.Millisecond * 100)
		for _, player := range players {
			engine_SolveTask(&player)
		}
	}
}

func core_taskAcceptor(c chan Task) {
	for {
		newTask := <-c
		switch newTask.TaskType {
		case ADDCLIENT:
			core_AddPlayer(newTask.ClientID)

		case DELCLIENT:
			core_DelPlayer(newTask.ClientID)

		case CLIENTMOVE, CLIENTSHOOT:
			core_setTask(newTask)
		}
		//AddPlayer(playerID)
	}
}

func core_setTask(newTask Task) {
	player, ok := players[newTask.ClientID]
	if ok {
		player.TaskMap[newTask.TaskType] = newTask
	}
}
func core_DelPlayer(playerID string) {
	_, ok := players[playerID]
	if ok {
		delete(players, playerID)
		fmt.Println("func core_DelPlayer(playerID string)")
	}
}
func core_AddPlayer(playerID string) {
	_, ok := players[playerID]
	if !ok {
		var newPlayer Player
		newPlayer.TaskMap = make(map[int]Task)
		players[playerID] = newPlayer
		fmt.Println("func AddPlayer(playerID string)")
	}
}

/*func core_AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}*/

var players map[string]Player
