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
var taskChan chan Task

var playersMutex sync.RWMutex

func main() {

	wg.Add(1)

	players = make(map[string]*Player)
	taskChan = make(chan Task)

	arg := os.Args[1] //0.0.0.0:8080
	//---------------------------------------------------------LOG file setup
	/*f, err := os.OpenFile("penguin_royale_logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("This is a test log entry")*/
	//----------------------------------------------------------END setup logfile

	fmt.Println("startServer")
	go startServer(arg)
	go taskWorker()
	go tickTockWorker()
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
			fmt.Println("accepted connection from ", conn.RemoteAddr())
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
	delete(players, newTask.clientId)
	playersMutex.Unlock()

	taskChan <- newTask
}

func taskWorker() {
	for {
		log.Println("taskWorker start.")
		newTask := <-taskChan
		switch newTask.taskType {
		case TASK_DELCLIENT:
			{
				log.Println("func core_DelPlayer(playerID string)")

				playersMutex.RLock()
				log.Println("NOW PLAYER COUNT = ", len(players))
				playersMutex.RUnlock()

				sendToPlayers(prepareMsg(strconv.FormatInt(MSG_KILLPLAYER, 10), newTask.clientId))
			}

		case TASK_RESPAWNCLIENT:
			{
				playersMutex.RLock()
				player, ok := players[newTask.clientId]
				playersMutex.RUnlock()
				if !ok {
					log.Println("WTF? Respawn player that doesn't exist in map",
						newTask.clientId)
				}

				// make random state
				player.healthPoint = STARTHEALTHPOINT
				player.pos = makeRandomPos()
				player.scorePoint /= 2

				sendToPlayers(prepareMsg(strconv.FormatInt(MSG_RESPAWNPLAYER, 10), newTask.clientId))
			}
		case TASK_UPDATESCORE:
			{
				var allScoreMsg string

				playersMutex.RLock()
				log.Println("update player task:", len(players))
				playersMutex.RUnlock()
				//make score array
				//var lineCount int64
				fmt.Println("_______________000000000000000000_________")
				playersMutex.Lock()
				fmt.Println("_______________111111111111111111_________")
				for _, player := range players {
					fmt.Println("_______________ZZZZZZZZZZZ_________", player.id)
					playerMsg := getPlayerScore(player)
					allScoreMsg += prepareMsg(playerMsg...)

					// split players score in one message.
					allScoreMsg += "#"

				}
				playersMutex.Unlock()
				fmt.Println("_______________END PREPARE UPD HISCORE MAG_________")

				if len(allScoreMsg) != 0 {
					sendToPlayers(allScoreMsg)
				}
			}
		}
		log.Println("taskWorker end.")
	}
}

func tickTockWorker() {
	var newmessage []string
	for {
		//log.Println("tickTockWorker start.")
		time.Sleep(time.Millisecond * TICKPERIOD)
		newmessage = newmessage[:0]

		playersMutex.RLock()
		log.Println("tic player count:", len(players))
		playersMutex.RUnlock()

		//make some physics works
		playersMutex.Lock()
		for _, player := range players {
			if !player.isPlay {
				continue
			}
			makePlayerPos(player)
			playerMsg := getPlayerPosMsg(player)
			newmessage = append(newmessage, prepareMsg(playerMsg...))
		}
		playersMutex.Unlock()

		if len(newmessage) != 0 {
			sendToPlayers(newmessage...)
		}
		//log.Println("tickTockWorker end.")
	}
}

func prepareMsg(parts ...string) string {
	return (strings.Join(parts, ";"))
}

func sendToPlayers(parts ...string) {
	msg := strings.Join(parts, "/")
	// ADD stop byte as $ symbol
	msg += "$"
	log.Println("send to all:", msg)
	playersMutex.Lock()
	for _, pl := range players {
		//log.Println("really send")
		pl.Conn.Write([]byte(msg))
	}
	playersMutex.Unlock()
}

func parsePlayersInput(str string, currentPlayer *Player) {

	//println("player input = ", str)
	strArr := strings.Split(str, ";")
	if len(strArr) < 2 {
		println("player str input len = ", len(strArr))
		println("player str =", str)
		return
	}

	switch {
	//---------------------------------------------------------------- player moves
	case strArr[0] == strconv.FormatInt(MSG_CLIENT_WANT_MOVE, 10):
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

		playersMutex.Lock()
		currentPlayer.wannaPos.x = x
		currentPlayer.wannaPos.y = y
		currentPlayer.wannaPos.angle = angle
		currentPlayer.wannaPos.isAttack = isAttack
		playersMutex.Unlock()

		// --------------------------------------------------------------player go in play
	case strArr[0] == strconv.FormatInt(MSG_CLIENT_WANT_PLAY, 10):
		if len(strArr) < 4 {
			println("read less arg onto needed for player starts = ", len(strArr))
			return
		}
		isPlay, err_Play := strconv.ParseBool(strArr[1])
		if err_Play != nil {
			println("isPlay, err_Play := strconv.ParseBool(strArr[1]) = ERROR", err_Play.Error)
			return
		}
		newNikName := (strArr[2])

		newSkinID, err_skinID := strconv.ParseInt(strArr[3], 10, 64)
		if err_skinID != nil {
			println("newSkinID, err_skinID := strconv.ParseInt(strArr[3], 10, 64) = ERROR", err_skinID.Error)
			return
		}

		playersMutex.Lock()
		currentPlayer.skinID = newSkinID
		currentPlayer.nikName = newNikName
		currentPlayer.isPlay = isPlay
		playersMutex.Unlock()

		playersMutex.Lock()
		if currentPlayer.isPlay {
			playersMutex.Unlock()
			newMessage := strconv.Itoa(MSG_YOURID) + ";"
			newMessage += currentPlayer.id + "$"

			playersMutex.Lock()
			currentPlayer.Conn.Write([]byte(newMessage))
			playersMutex.Unlock()

			println("player want play", str)
			//println("player is play =", currentPlayer.isPlay)

			var startTask Task
			startTask.clientId = currentPlayer.id
			startTask.taskType = TASK_RESPAWNCLIENT
			taskChan <- startTask

			//Player's start task
			var scoreTask Task
			scoreTask.clientId = currentPlayer.id
			scoreTask.taskType = TASK_UPDATESCORE
			taskChan <- scoreTask
		} else {
			playersMutex.Unlock()
			//---- delete player from other player in client
			var newTask Task
			newTask.taskType = TASK_DELCLIENT
			newTask.clientId = currentPlayer.id
			taskChan <- newTask
		}
	default:
		log.Println("WTF? There shouldn't be default value")
		return
	}
}

func createPlayer(conn net.Conn, id string) {

	playersMutex.Lock()
	if _, ok := players[id]; !ok {
		playersMutex.Unlock()
		log.Println("createPlayer %s", id)

		var newPlayer *Player = new(Player)
		newPlayer.id = id
		newPlayer.healthPoint = STARTHEALTHPOINT
		newPlayer.Conn = conn
		newPlayer.scorePoint = 0
		newPlayer.isPlay = false
		//newPlayer.pos = makeRandomPos()

		//log.Println("Try lock")
		//playersMutex.Lock()

		players[id] = newPlayer

		//playersMutex.Unlock()
		//log.Println("UNlock")

		go readClientData(newPlayer)
	} else {
		playersMutex.Unlock()
		log.Println("client %s exist.\nWFT?????????", id)
	}

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
