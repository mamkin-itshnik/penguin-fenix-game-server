package main

import (
	"fmt"
)

func core_StartServer(adress string) {
	go core_taskAcceptor(TaskChan)
	CN_StartServer(adress)
}

func init() {
	fmt.Println("Create core ")
	players = make(map[string]Player)
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
		for i := 0; i < TASKCOUNT; i++ {
			if i == newTask.TaskType {
				player.TaskArray[i] = newTask
			}
		}
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
		newPlayer.TaskArray = make([]Task, TASKCOUNT)
		players[playerID] = newPlayer
		fmt.Println("func AddPlayer(playerID string)")
	}
}

/*func core_AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}*/

var players map[string]Player
