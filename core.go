package main

import (
	"fmt"
)

func core_StartServer(adress string) {
	go core_taskAcceptor(ConnectionChan)
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
			fmt.Println("_________________ hui")
		case CLIENTMOVE:
			// do something for
		case CLIENTSHOOT:
			// do something for
		}
		//AddPlayer(playerID)
	}
}

func core_AddPlayer(playerID string) {
	fmt.Println("func AddPlayer(playerID string)")
}
func core_AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}

var players map[string]*Player
