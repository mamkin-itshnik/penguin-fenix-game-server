package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

var Clients map[string]*Client
var ConnectionChan chan Task

func init() {
	fmt.Println("Create connectManager ")
	Clients = make(map[string]*Client)
	ConnectionChan = make(chan Task)
}

func addClient(conn net.Conn, Id string) bool {
	if _, ok := Clients[Id]; !ok {

		//make client
		var newClient Client
		newClient.clientID = Id
		newClient.Conn = conn
		Clients[Id] = &newClient

		//make task
		var newTask Task
		newTask.TaskType = ADDCLIENT
		newTask.ClientID = Id
		ConnectionChan <- newTask

		println("new client add!, now client count = ", len(Clients))
		return true
	} else {
		println("client %s exist", Id)
		return false
	}
}

func StartServer_CN(adress string) {
	go runAcceptor(adress)
	go readClientsData()
	go writeClientData()
}

func runAcceptor(adress string) error {
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
			addClient(conn, "ID_"+strconv.Itoa(i))
			log.Println("accepted connection from ", conn.RemoteAddr())
			i++
		}
	}
	return err
}

func writeClientData() {

}

func readClientsData() {
	for {
		for id, cli := range Clients {
			fmt.Println("___", cli.clientID)
			message, err := bufio.NewReader(Clients[id].Conn).ReadString('\n')
			if err == nil {
				log.Printf("error accepting connection %v", err)
				println("read client data: ", message)

				strArr := strings.Split(message, ";")
				if len(strArr) < 3 {
					continue
				}

				//var dX, dY, nAngl float64
				if strings.Contains(strArr[0], "XD") {
					if strArr[1] == "X" {
						//println("read client data: ", message)
						//Clients[id].isAttack = false
						continue
					}
					//Clients[id].isAttack = true
					//nAngl, _ = strconv.ParseFloat(strArr[1], 32)
					//Clients[id].Pos.ShootAngle = nAngl
					continue
				}

				if len(strArr) > 3 {
					//dX, _ = strconv.ParseFloat(strArr[1], 32)
					//dY, _ = strconv.ParseFloat(strArr[2], 32)
					//nAngl, _ = strconv.ParseFloat(strArr[3], 32)

					/*	g.Clients[id].Pos.TryDeltaX = dX
						g.Clients[id].Pos.TryDeltaY = dY
						g.Clients[id].Pos.Angle = nAngl*/
				}
			} else {
				if err == io.EOF {
					println("NewReader io.EOF", err)
					//sErr := cli.Conn.Close
					//if sErr == nil {
					//	println("Socket close for ", cli.clientID)
					defer cli.Conn.Close()
					//make task
					var newTask Task
					newTask.TaskType = DELCLIENT
					newTask.ClientID = cli.clientID
					ConnectionChan <- newTask
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
