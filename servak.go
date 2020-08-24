package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"./gameCore"
)

var wg sync.WaitGroup
var serverIsWork bool
var clients map[string]*gameCore.GameClient

func startServer(adress string) error {
	var i int
	i = 0
	log.Printf("try starting server on %v\n", adress)
	listener, err := net.Listen("tcp", adress)
	if err != nil {
		fmt.Println("server error ", err)
		return err
	} else {
		fmt.Println("server START on ", adress)
	}
	defer listener.Close()
	for serverIsWork {
		conn, err := listener.Accept()
		if err != nil {
			//log.Printf("error accepting connection %v", err)
		} else {
			gameCore.GetInstance().AddClient(conn, "ID_"+strconv.Itoa(i))

			log.Println("accepted connection from ", conn.RemoteAddr())
			log.Println("clients count now = ", len(gameCore.GetInstance().Clients))
			i++
		}
	}
	return err
}

func reciveDataFromClient() {
	for {
		gameCore.GetInstance().ReadClientsData()
	}
}

func sendDataaToClient() {
	for {
		gameCore.GetInstance().WriteClientsData()
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	serverIsWork = true
	wg.Add(1)
	arg := os.Args[1] //192.168.0.105:8080
	go startServer(arg)
	//go startServer("192.168.0.105:8080")
	go reciveDataFromClient()
	go sendDataaToClient()
	//------- waiting like system "pause", but if call wg.Done() once anywhere program end // see
	wg.Wait()
}
