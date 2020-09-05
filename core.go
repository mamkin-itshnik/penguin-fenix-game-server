package main

import (
	"fmt"
)

func StartServer_core(adress string) {
	go taskAcceptor_core(ConnectionChan)
	StartServer_CN(adress)
}

func init() {
	fmt.Println("Create core ")
}

func taskAcceptor_core(c chan Task) {
	for {
		newTask := <-c
		switch newTask.TaskType {
		case ADDCLIENT:
			AddPlayer(newTask.ClientID)
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

func AddPlayer(playerID string) {
	fmt.Println("func AddPlayer(playerID string)")
}
func AddTask(playerID string, newTask Task) {
	fmt.Println("Hello TASK!")
}

var players map[string]*Player
