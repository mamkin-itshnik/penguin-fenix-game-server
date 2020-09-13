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
	players = make(map[string]*Player)
}

func core_TicTack() {
	var newmessage string
	for {
		time.Sleep(time.Millisecond * 100)
		//fmt.Println("watafuck", len(players))
		newmessage = ""

		//make tasks
		for _, player := range players {
			//fmt.Println("player pos ", player.Pos)
			engine_SolveTask(player)
			CN_addPlayerStringTask(player, &newmessage)
		}
		// send tasks
		for _, player := range players {
			CN_writeClientData(player, &newmessage)
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
			core_DelPlayer(newTask)

		case CLIENTMOVE:
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
func core_DelPlayer(newTask Task) {
	_, ok := players[newTask.ClientID]
	if ok {
		delete(players, newTask.ClientID)
		fmt.Println("func core_DelPlayer(playerID string)")
	}
}
func core_AddPlayer(ClientID string) {
	_, ok := players[ClientID]
	if !ok {
		var newPlayer *Player = new(Player)
		newPlayer.TaskMap = make(map[int]Task)
		newPlayer.Id = ClientID
		newPlayer.HealfPoint = startHealfPoint
		players[ClientID] = newPlayer
		fmt.Println("func AddPlayer(playerID string)")
	}
}

/*func core_AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}*/

var players map[string]*Player
