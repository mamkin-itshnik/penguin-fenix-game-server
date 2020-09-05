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
}

func core_taskAcceptor(c chan Task) {
	for {
		newTask := <-c
		switch newTask.TaskType {
		case ADDCLIENT:
			core_AddPlayer(newTask.ClientID)
		case DELCLIENT:
			core_DelPlayer(newTask.ClientID)
		case CLIENTMOVE:
			// do something for
		case CLIENTSHOOT:
			// do something for
		}
		//AddPlayer(playerID)
	}
}

func core_DelPlayer(playerID string) {
	fmt.Println("func core_DelPlayer(playerID string)")
}
func core_AddPlayer(playerID string) {
	fmt.Println("func AddPlayer(playerID string)")
}
func core_AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}

var players map[string]*Player
