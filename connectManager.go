package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

var Clients map[string]*Client
var TaskChan chan Task

func init() {
	fmt.Println("Create connectManager ")
	Clients = make(map[string]*Client)
	TaskChan = make(chan Task)
}

func CN_addClient(conn net.Conn, Id string) bool {
	if _, ok := Clients[Id]; !ok {

		println("CN_addClient %s", Id)
		//make client
		var newClient *Client
		newClient = new(Client)
		newClient.clientID = Id
		newClient.Conn = conn
		Clients[Id] = newClient

		//make task
		var newTask Task
		newTask.TaskType = ADDCLIENT
		newTask.ClientID = Id
		TaskChan <- newTask

		newmessage := "0;" + Id + ";"
		newClient.Conn.Write([]byte(newmessage))

		println("new client add!, now client count = ", len(Clients))
		return true
	} else {
		println("client %s exist", Id)
		return false
	}
}

func CN_StartServer(adress string) {
	go CN_runAcceptor(adress)
	go CN_readClientsData()
	//go CN_writeClientData()
}

func CN_runAcceptor(adress string) error {
	var i int
	i = 0
	log.Printf("try starting server on %v\n", adress)
	listener, err := net.Listen("tcp", adress)
	if err != nil {
		fmt.Println("server error ", err)
		return err
	}
	fmt.Println("server START on ", adress)
	defer listener.Close()
	for true {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
		} else {
			CN_addClient(conn, "ID_"+strconv.Itoa(i))
			log.Println("accepted connection from ", conn.RemoteAddr())
			i++
		}
	}
	return err
}

func CN_writeClientData(currentPlayer Player) {

	log.Println("____________1", len(currentPlayer.TaskMap))
	for i, _ := range currentPlayer.TaskMap {
		Clients[currentPlayer.Id].Conn.Write([]byte(TP_makeStringTask(&currentPlayer, i)))
	}
}

func CN_readClientsData() {
	for {
		for id, cli := range Clients {
			//fmt.Println("___", cli.clientID)
			message, err := bufio.NewReader(Clients[id].Conn).ReadString('\n')
			if err == nil {

				//println("read client data: ", message)

				var newTask Task = TP_makeClientInputsTask(message, id)
				TaskChan <- newTask
			} else {
				if err == io.EOF {
					println("NewReader io.EOF", err)
					println("NewReader io.EOF", err.Error())
					//sErr := cli.Conn.Close
					//if sErr == nil {
					//	println("Socket close for ", cli.clientID)
					defer cli.Conn.Close()
					//make task
					var newTask Task
					newTask.TaskType = DELCLIENT
					newTask.ClientID = cli.clientID
					TaskChan <- newTask
					delete(Clients, cli.clientID)
					//} else {
					//	println("Socket close error ", sErr)
					//}
				} else {
					println("NewReader error", err)
				}
			}
		}
	}
}
