package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var players map[string]*Player
var topScorePlayer []ScoreEntyty
var taskChan chan Task
var playersWriteChan chan string
var playersMutex sync.Mutex
var atomicIntPlayerCount int
var needStoreLog bool

func main() {

	wg.Add(1)

	players = make(map[string]*Player)
	topScorePlayer = make([]ScoreEntyty, SCORE_ENTYTY_COUNT)
	taskChan = make(chan Task, 100000)
	playersWriteChan = make(chan string)
	needStoreLog = false

	arg := os.Args[1] //0.0.0.0:55555 //127.0.0.1:55555
	//---------------------------------------------------------LOG file setup
	if needStoreLog{
		setUpLogFile()
	}

	fmt.Println("startServer")
	go startServer(arg)
	go taskWorker()
	go tickTockWorker()
	go writeLoop()
	//initBotFather()

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
			log.Println("accepted connection from ", conn.RemoteAddr())
			createPlayer(conn, "ID_"+strconv.Itoa(i))
			i++
		}
	}
	return
}

func makeDeletePlayerTask(player *Player) {
	var newTask Task
	newTask.taskType = TASK_DELCLIENT
	newTask.clientId = player.id
	player.Conn.Close()

	playersMutex.Lock()

	//need delete old playerID
	deleteOldPlayerIdFromOtherPlayers(player)
	delete(players, newTask.clientId)
	atomicIntPlayerCount = len(players)

	playersMutex.Unlock()

	taskChan <- newTask
}

func writeLoop() {
	for {
		newMessage := <-playersWriteChan

		playersMutex.Lock()
		for _, pl := range players {
			//log.Println("really send________", newMessage)
			//log.Println("_")
			pl.Conn.Write([]byte(newMessage))
		}
		playersMutex.Unlock()

	}
}

func taskWorker() {
	for {
		taskCoint := len(taskChan)
		if taskCoint > 50000 {
			taskOverflow()
		}

		newTask := <-taskChan
		switch newTask.taskType {
		case TASK_DELCLIENT:
			{
				solveTaskDellClient(&newTask)
			}

		case TASK_RESPAWNCLIENT:
			{
				solveTaskRespawnClient(&newTask)
			}
		case TASK_UPDATESCORE:
			{
				solveTaskUpdateScore(&newTask)
			}
		}

	}
}

func tickTockWorker() {
	for {

		time.Sleep(time.Millisecond * TICKPERIOD)

		//make some physics works
		playersMutex.Lock()
		for _, player := range players {
			if !player.isPlay {
				continue
			}
			workOnPlayer(player)
		}
		playersMutex.Unlock()
	}
}

func workOnPlayer(player *Player) {
	//calculate state

	makePlayerState(player)

	message := makePlayerMessage(player)

	player.Conn.Write([]byte(message))
}

func makePlayerMessage(player *Player) string {

	var subMessage []string

	subMessage = makeOtherPosMsg(player)

	playerMessage := spliceSubMessage(getPlayerPosMsg(player)...)

	subMessage = append(subMessage, playerMessage)

	finalMessage := spliceMessages(subMessage...)

	// add stop-byte
	finalMessage += "$"

	return finalMessage
}

func parsePlayersInput(str string, currentPlayer *Player) {

	//println("player input = ", str)
	strArr := strings.Split(str, ";")
	if len(strArr) < 1 {
		println("player str input len = ", len(strArr))
		println("player str =", str)
		return
	}

	switch {
	//---------------------------------------------------------------- player moves
	case strArr[0] == strconv.FormatInt(MSG_CLIENT_WANT_MOVE, 10):

		parseClientWantToMove(strArr, currentPlayer)

	case strArr[0] == strconv.FormatInt(MSG_CLIENT_WANT_PLAY, 10):

		parseClientWantToPlay(strArr, currentPlayer)

	case strArr[0] == strconv.FormatInt(MSG_CLIENT_WANT_GET_PLAYER_COUNT, 10):

		parseClientWantGetPlayerCount(currentPlayer)

	default:
		log.Println("WTF? There shouldn't be default value")
		return
	}
}

func createPlayer(conn net.Conn, id string) {

	log.Println("createPlayer %s", id)

	playersMutex.Lock()
	log.Println("player count =", len(players))
	playersMutex.Unlock()

	var newPlayer *Player = new(Player)
	newPlayer.id = id
	newPlayer.Conn = conn

	setZeroPropertyForPlayer(newPlayer)

	playersMutex.Lock()
	players[id] = newPlayer
	atomicIntPlayerCount = len(players)
	playersMutex.Unlock()

	go readClientData(newPlayer)
}

func readClientData(player *Player) {
	defer makeDeletePlayerTask(player)
	reader := bufio.NewReader(player.Conn)
	for {
		message, err := reader.ReadString('\n')
		if err == nil {
			//log.Println("readPlayersInput_____ ", player.id, message)
			parsePlayersInput(message, player)
		} else {
			if err == io.EOF {
				log.Println("bufio error io.EOF", err)
			} else {
				log.Println("bufio unknow error ", err)
			}
			log.Println("readPlayersInput player loop err +++")
			return
		}
	}

}