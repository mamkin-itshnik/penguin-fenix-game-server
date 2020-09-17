package main

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var players map[string]*Player
var taskChan chan Task

func main() {

	wg.Add(1)

	players = make(map[string]*Player)
	taskChan = make(chan Task)

	arg := os.Args[1] //0.0.0.0:8080

	go startServer(arg)
	go readPlayersInput()
	go taskWorker()
	go tickTockWorker()

	wg.Wait()
}

func startServer(arg string) {
	i := 0
	log.Printf("starting server on %v\n", arg)
	listener, err := net.Listen("tcp", arg)
	if err != nil {
		log.Println("server error ", err)
		return
	}
	defer listener.Close()
	log.Println("server START on ", arg)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
		} else {
			createPlayer(conn, "ID_"+strconv.Itoa(i))
			log.Println("accepted connection from ", conn.RemoteAddr())
			i++
		}
	}
	return
}

func createPlayer(conn net.Conn, id string) {
	if _, ok := players[id]; !ok {

		println("createPlayer %s", id)

		var newPlayer *Player = new(Player)
		newPlayer.id = id
		newPlayer.healthPoint = STARTHEALTHPOINT
		newPlayer.Conn = conn
		newMessage := strconv.Itoa(YOURID) + ";"
		newMessage += id + ";\n"
		players[id] = newPlayer
		newPlayer.Conn.Write([]byte(newMessage))
	} else {
		println("client %s exist.\nWFT?????????", id)
	}
}

func readPlayersInput() {
	for {
		for _, pl := range players {
			message, err := bufio.NewReader(pl.Conn).ReadString('\n')
			if err == nil {
				parsePlayersInput(message, pl)
			} else {
				if err == io.EOF {
					println("NewReader io.EOF", err)
					println("NewReader io.EOF", err.Error())
					pl.Conn.Close()
				} else {
					println("NewReader error____________ hz___", err, err.Error)
				}

				// TODO: check this
				//make task
				var newTask Task
				newTask.taskType = DELCLIENT
				newTask.clientId = pl.id
				taskChan <- newTask
			}
		}
	}
}

func parsePlayersInput(str string, currentPlayer *Player) {
	strArr := strings.Split(str, ";")
	if len(strArr) < 2 {
		println("player str input len = ", len(strArr))
		println("player str =", str)
		return
	}

	switch {
	// case strArr[0] == "0":
	case strArr[0] == "2": // player moves
		if len(strArr) < 4 {
			println("player str input len = ", len(strArr))
			println("player str =", str)
			return
		}
		x, err_x := strconv.ParseFloat(strArr[1], 64)
		y, err_y := strconv.ParseFloat(strArr[2], 64)
		angle, err_a := strconv.ParseInt(strArr[3], 10, 64)
		isAttack, err_attack := strconv.ParseBool(strArr[4])
		if (err_x != nil) || (err_y != nil) || (err_a != nil) || (err_attack != nil) {
			return
		}

		currentPlayer.wannaPos.x = x
		currentPlayer.wannaPos.y = y
		currentPlayer.wannaPos.angle = angle
		currentPlayer.wannaPos.isAttack = isAttack
	// case strArr[0] == "3":
	default:
		log.Println("WTF? There shouldn't be default value")
		return
	}
}

func taskWorker() {
	for {
		newTask := <-taskChan
		switch newTask.taskType {
		case DELCLIENT:
			{
				_, ok := players[newTask.clientId]
				if !ok {
					//log.Println("WTF? Deleting player that doesn't exist in map",
					//	newTask.clientId)
					break
				}
				delete(players, newTask.clientId)
				log.Println("func core_DelPlayer(playerID string)")
				log.Println("NOW PLAYER COUNT = ", len(players))

				var newMsg string
				newMsg += strconv.FormatInt(DELCLIENT, 10) + ";"
				newMsg += newTask.clientId + ";\n"

				sendToPlayers(&newMsg)
			}

		case RESPAWNCLIENT:
			{
				player, ok := players[newTask.clientId]
				if !ok {
					log.Println("WTF? Respawn player that doesn't exist in map",
						newTask.clientId)
				}

				// make random state
				player.healthPoint = STARTHEALTHPOINT
				player.pos.x = MINPOS + rand.Float64()*(MAXPOS-MINPOS)
				player.pos.y = MINPOS + rand.Float64()*(MAXPOS-MINPOS)
				player.scorePoint = 0

				// make message
				var newMsg string
				newMsg += strconv.FormatInt(RESPAWNCLIENT, 10) + ";"
				newMsg += newTask.clientId + ";"
				newMsg += "\n"

				sendToPlayers(&newMsg)
			}
		}
	}
}

func sendToPlayers(msg *string) {
	for _, pl := range players {
		pl.Conn.Write([]byte(*msg))
	}
}

func tickTockWorker() {
	var newmessage string
	for {
		time.Sleep(time.Millisecond * 100)
		newmessage = ""
		//make some physics works
		for _, player := range players {
			makePlayerPos(player)
			newmessage += getPlayerPosMsg(player)
		}
		sendToPlayers(&newmessage)
	}
}
