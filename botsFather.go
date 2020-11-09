package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
	"math/rand"

	//"container/list"
	"net"
)

const (
	NEEDPLAYERS = 4// 4 bot for 70x70 square
	BOT_MOVE_CHANGE_TIME_OUT             = 1000
	BOT_SHOOT_DISTANCE 					 = 15.0
	BOT_MAX_SHOOT_DEGREE_DEFLECTION 	 = 30
	BOT_MIN_SHOOT_DEGREE_DEFLECTION 	 = 10
	SERVER_QUERY_TIME_OUT 				 = 100
)

var serverQuerryTimeOut time.Duration = SERVER_QUERY_TIME_OUT // 1000 ms
var bots map[int]*Bot
var currentBotId int = -1
//var wg sync.WaitGroup
var botsMutex sync.Mutex


type BotMoveState struct {
	beginTimeOut   int64
	currentTimeOut int64
}

type Bot struct {
	//---------------------Bot feature
	MoveState    BotMoveState
	visibleEnemy map[string]*Player
	visibleConteinerMutex sync.Mutex
	//--------------------- update by engine
	pos      Position
	id       string
	isAttack bool
	shootDegreeDeflection int64
	//----------------------write interface
	net.Conn
	//--------------------- update by client
	tryX          float64
	tryY          float64
	tryShootAngle int64
	tryMoveAngle  int64
	tryAttack     bool
}

func init() {

	//wg.Add(1)
	initBotFather()
	go tikTak()
	//wg.Wait()
}

// ?m.b. func (bot * Bot)
func (bot *Bot) makeTryXYforBot() {
	bot.tryX = -math.Cos((float64(bot.tryMoveAngle) * (math.Pi / 180)))
	bot.tryY = math.Sin((float64(bot.tryMoveAngle) * (math.Pi / 180)))
}
func (bot *Bot) makeShootAngle(targetPos Position) {
	//---------------------------------------------------------------------------TODO CHECK THIS SHIT

	vectorX := bot.pos.x - targetPos.x
	vectorY := bot.pos.y - targetPos.y
	tangAlfa := vectorY / vectorX

	angle := 90.0 - ((math.Atan(tangAlfa) * 180.0) / math.Pi)

	if bot.pos.x > targetPos.x {

		angle -= 90.0
	} else {
		angle += 90.0
	}
	bot.tryShootAngle = int64(angle) + makeShootDegreeDeflection(bot)

}

func makeShootDegreeDeflection(bot *Bot) int64{

	degreeValue := int64(rand.Float64() * float64(bot.shootDegreeDeflection))
	if degreeValue%2 == 0{
		return degreeValue
	}else{
		return -degreeValue
	}
	 
}


func checkPlayersCount(playersCount int64) {
	if playersCount > NEEDPLAYERS {
		removeBot()
		return
	}
	if playersCount < NEEDPLAYERS {
		addBot()
	}
}

func botWalk() {
	for _, bot := range bots {

		bot.MoveState.currentTimeOut -= 100 // TIK PERIOD
		if bot.MoveState.currentTimeOut < 0 {
			bot.MoveState.currentTimeOut = bot.MoveState.beginTimeOut

			//change bot direction
			bot.tryMoveAngle += 60 // rotate N degree
			bot.tryMoveAngle %= 360
			bot.tryShootAngle = bot.tryMoveAngle // foreward rotation
			bot.makeTryXYforBot()
		}
	}
}

func botShoot() {
	for _, bot := range bots {

		// find minimum distance to player
		bot.isAttack = false
		//var minDistance float64 = 9999999.9
		bot.visibleConteinerMutex.Lock()
		for _, enemy := range bot.visibleEnemy {
			if bot.id != enemy.id {
				if enemy.pos.pointDistance(bot.pos) < BOT_SHOOT_DISTANCE {

					//fmt.Println("try ATTACK player with ID = ", player.id, "  my id = ", bot.id)
					bot.isAttack = true
					bot.makeShootAngle(enemy.pos)
					break
				}
			}
		}
		bot.visibleConteinerMutex.Unlock()
	}
}

func addBot() {

	//fmt.Println("try addBot")
	tcpAddr, err_adr := net.ResolveTCPAddr("tcp4", "127.0.0.1:55555")
	if err_adr == nil {
		conn, err_dial := net.DialTCP("tcp", nil, tcpAddr)
		if err_dial == nil {

			// make new Player
			var newPlayer *Bot = new(Bot)
			newPlayer.id = "X"
			newPlayer.Conn = conn
			newPlayer.shootDegreeDeflection =  int64(rand.Float64() * float64(BOT_MAX_SHOOT_DEGREE_DEFLECTION)) + BOT_MIN_SHOOT_DEGREE_DEFLECTION
			newPlayer.visibleEnemy = make(map[string]*Player)
			newPlayer.MoveState.beginTimeOut = BOT_MOVE_CHANGE_TIME_OUT
			newmessage := "0;true;" + "nik name;" + "0" + ";\n"

			currentBotId++

			botsMutex.Lock()
			bots[currentBotId] = newPlayer
			botsMutex.Unlock()

			go readServerDataForBot(bots[currentBotId])

			//fmt.Println("try send First Message from bot")
			newPlayer.Conn.Write([]byte(newmessage))
			//fmt.Println("Message is sended")
		}
	}
	//fmt.Println("end addBot")
}
func createBotFatherListner() {
	//fmt.Println("_______try createBotFatherListner____")
	tcpAddr, err_adr := net.ResolveTCPAddr("tcp4", "127.0.0.1:55555")
	if err_adr == nil {
		conn, err_dial := net.DialTCP("tcp", nil, tcpAddr)
		if err_dial == nil {

			go requestLoopToServer(conn)
			go responseLoopFromServer(conn)
		}
	}
}

func requestLoopToServer(conn net.Conn) {

	var message string
	message += strconv.FormatInt(MSG_CLIENT_WANT_GET_PLAYER_COUNT, 10) + ";" // player move msg type
	message += "\n"
	for {
		conn.Write([]byte(message))
		time.Sleep(time.Millisecond * serverQuerryTimeOut)
	}
}

func responseLoopFromServer(conn net.Conn) {

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('$')
		if err == nil {
			//log.Println("readPlayersInput_____", message)
			parseServerResponse(message)
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

func parseServerResponse(message string) {
	strArr := strings.Split(strings.TrimSuffix(message, "$"), ";")
	if len(strArr) < 2 {
		println("server str input len = ", len(strArr))
		println("server str =", message)
		return
	}

	switch {
	case strArr[0] == strconv.FormatInt(MSG_PLAYER_COUNT, 10):
		//message = append(message, strconv.FormatInt(MSG_STATE, 10)){

		if len(strArr) < 2 {
			println("server str input len = ", len(strArr))
			println("server str =", message)
			return
		}
		playerCount, err_count := strconv.ParseInt(strArr[1], 10, 64)

		if err_count != nil {
			return
		}
		checkPlayersCount(playerCount)

	default:
		return
	}
}

func removeBot() {

	botsMutex.Lock()
	bots[currentBotId].Conn.Close()
	defer delete(bots, currentBotId)
	currentBotId--
	botsMutex.Unlock()
}

// Main game loop
func tikTak() {
	time.Sleep(time.Millisecond * 1000)
	fmt.Println("start botsFather ")
	for {
		//checkPlayersCount()

		botsMutex.Lock()
		botWalk()
		botShoot() // set/unset Bot.isAttack true/false
		botMessage()
		botsMutex.Unlock()

		time.Sleep(time.Millisecond * TICKPERIOD)
	}
}

func botMessage() {
	for _, bot := range bots {
		var message string
		message += strconv.FormatInt(MSG_CLIENT_WANT_MOVE, 10) + ";" // player move msg type
		//message += currentPlayer.Id + ";"
		message += strconv.FormatFloat(bot.tryX, 'f', 1, 64) + ";"
		message += strconv.FormatFloat(bot.tryY, 'f', 1, 64) + ";"
		message += strconv.FormatInt(int64(bot.tryShootAngle), 10) + ";"
		message += strconv.FormatBool(bot.isAttack) + ";"
		message += "\n"
		//fmt.Println("Send from bot ", message)
		bot.Conn.Write([]byte(message))
	}
}

func initBotFather() {
	time.Sleep(time.Millisecond * 4000)
	fmt.Println("init botsFather")
	bots = make(map[int]*Bot)
	go createBotFatherListner()
}

func readServerDataForBot(bot *Bot) {
	//defer makeDeletePlayerTask(bot)
	//fmt.Println("readPlayersInput_____")
	reader := bufio.NewReader(bot.Conn)
	for {
		message, err := reader.ReadString('$')
		if err == nil {
			parseServerInput(message, bot)
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

func parseServerInput(message string, bot *Bot) {
	//fmt.Println("__________________read server message ", strings.TrimSuffix(message, "$"))
	//fmt.Println("______+________read server message ",bot.id, message)
	strArr := strings.Split(strings.TrimSuffix(message, "$"), "/")
	if len(strArr) < 1 {
		println("server str input len = ", len(strArr))
		println("server str =", message)
		return
	}

	for _, subMessage := range strArr {
		parseStringTaskForBot(subMessage, bot)
	}

}

func parseStringTaskForBot(message string, bot *Bot) {

	strArr := strings.Split(message, ";")
	if len(strArr) < 2 {
		return
	}

	switch {
	// set id
	case strArr[0] == strconv.FormatInt(MSG_YOUR_ID, 10):
		println("add bot", strArr[1])
		bot.id = strArr[1]

	case strArr[0] == strconv.FormatInt(MSG_STATE, 10):
		//message = append(message, strconv.FormatInt(MSG_STATE, 10))

		if len(strArr) < 5 {
			println("server str input len = ", len(strArr))
			println("server str MSG_STATE=", strArr)
			return
		}

		if strArr[1] == bot.id {
			x, err_x := strconv.ParseFloat(strArr[2], 64)
			y, err_y := strconv.ParseFloat(strArr[3], 64)
			if (err_x != nil) || (err_y != nil) {
				return
			}

			bot.pos.x = x
			bot.pos.y = y
		} else {
			if otherEnemy, ok := bot.visibleEnemy[strArr[1]]; ok {
				x, err_x := strconv.ParseFloat(strArr[2], 64)
				y, err_y := strconv.ParseFloat(strArr[3], 64)
				if (err_x != nil) || (err_y != nil) {
					return
				}

				otherEnemy.pos.x = x
				otherEnemy.pos.y = y
			}
		}

	case strArr[0] == strconv.FormatInt(MSG_DELETE_PLAYER, 10):
		//remove
		if len(strArr) < 2 {
			println("server str input len = ", len(strArr))
			println("server str MSG_DELETE_PLAYER=", strArr)
			return
		}
		id := strArr[1]
		if _, ok := bot.visibleEnemy[id]; ok {
			//fmt.Println("**************   REMOVE ENEMNY in bot *********")

			bot.visibleConteinerMutex.Lock()
			delete(bot.visibleEnemy, id)
			bot.visibleConteinerMutex.Unlock()
		}

	case strArr[0] == strconv.FormatInt(MSG_ADD_PLAYER, 10):
		//remove
		if len(strArr) < 2 {
			println("server str input len = ", len(strArr))
			println("server str MSG_ADD_PLAYER=", strArr)
			return
		}
		id := strArr[1]
		if _, ok := bot.visibleEnemy[id]; !ok {

			//println("add new player to bot =", id)
			var newPlayer *Player = new(Player)
			newPlayer.id = id

			bot.visibleConteinerMutex.Lock()
			bot.visibleEnemy[id] = newPlayer
			bot.visibleConteinerMutex.Unlock()
		}

	default:
		return
	}
}
