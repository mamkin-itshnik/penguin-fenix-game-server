package connectManager

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"../core"
)

type Client struct {
	net.Conn
	clientID string
}

var Clients map[string]*Client

func init() {
	fmt.Println("Create connectManager ")
	Clients = make(map[string]*Client)
}

func addClient(conn net.Conn, Id string) bool {
	if _, ok := Clients[Id]; !ok {
		var newClient Client
		Clients[Id] = &newClient
		core.AddPlayer(Id)
		println("new client add!, now client count = ", len(Clients))
		return true
	} else {
		println("client %s exist", Id)
		return false
	}
}

func StartServer(adress string) error {
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
