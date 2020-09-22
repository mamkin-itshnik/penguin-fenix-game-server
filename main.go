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
	//---------------------------------------------------------LOG file setup
	f, err := os.OpenFile("penguin_royale_logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("This is a test log entry")
	//----------------------------------------------------------END setup logfile

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

func readPlayersInput() {
	for {
		for _, pl := range players {
			message, err := bufio.NewReader(pl.Conn).ReadString('\n')
			if err == nil {
				//log.Println("readPlayersInput player loop non err")
				parsePlayersInput(message, pl)
			} else {
				if err == io.EOF {
					log.Println("bufio error io.EOF", err)
					pl.Conn.Close()
				} else {
					log.Println("bufio error ", err)
				}

				// TODO: check this
				//make task
				log.Println("readPlayersInput player loop err +++")
				var newTask Task
				newTask.taskType = TASK_DELCLIENT
				newTask.clientId = pl.id
				delete(players, newTask.clientId)
				taskChan <- newTask
			}
		}
	}
}

func taskWorker() {
	for {
		log.Println("taskWorker start.")
		newTask := <-taskChan
		switch newTask.taskType {
		case TASK_DELCLIENT:
			{
				log.Println("func core_DelPlayer(playerID string)")
				log.Println("NOW PLAYER COUNT = ", len(players))

				sendToPlayers(prepareMsg(strconv.FormatInt(MSG_KILLPLAYER, 10), ";", newTask.clientId, ";"))
			}

		case TASK_RESPAWNCLIENT:
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

				sendToPlayers(prepareMsg(strconv.FormatInt(MSG_RESPAWNPLAYER, 10), ";", newTask.clientId), ";")
			}
		case TASK_UPDATESCORE:
			{
				var newmessage []string
				log.Println("update player task:", len(players))
				//make score array
				for _, player := range players {
					playerMsg := getPlayerScore(player)
					newmessage = append(newmessage, prepareMsg(playerMsg...))
				}
				if len(newmessage) != 0 {
					sendToPlayers(newmessage...)
				}
			}
		}
		log.Println("taskWorker end.")
	}
}

func tickTockWorker() {
	var newmessage []string
	for {
		log.Println("tickTockWorker start.")
		time.Sleep(time.Millisecond * TICKPERIOD)
		newmessage = newmessage[:0]

		log.Println("tic player count:", len(players))
		//make some physics works
		for _, player := range players {
			makePlayerPos(player)
			playerMsg := getPlayerPosMsg(player)
			newmessage = append(newmessage, prepareMsg(playerMsg...))
		}
		if len(newmessage) != 0 {
			sendToPlayers(newmessage...)
		}
		log.Println("tickTockWorker end.")
	}
}

func prepareMsg(parts ...string) string {
	return (strings.Join(parts, ""))
}

func sendToPlayers(parts ...string) {
	msg := strings.Join(parts, "")
	msg += "\n"
	log.Println("send to all:", msg)
	for _, pl := range players {
		log.Println("really send")
		pl.Conn.Write([]byte(msg))
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

func createPlayer(conn net.Conn, id string) {
	if _, ok := players[id]; !ok {

		log.Println("createPlayer %s", id)

		var newPlayer *Player = new(Player)
		newPlayer.id = id
		newPlayer.healthPoint = STARTHEALTHPOINT
		newPlayer.Conn = conn
		newMessage := strconv.Itoa(MSG_YOURID) + ";"
		newMessage += id + ";\n"
		players[id] = newPlayer
		newPlayer.Conn.Write([]byte(newMessage))
	} else {
		log.Println("client %s exist.\nWFT?????????", id)
	}
}
