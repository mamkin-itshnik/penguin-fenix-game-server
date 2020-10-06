package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	//"container/list"
	"net"
)

var bots map[int]*Bot
var currentPlayerCount int
var currentBotId int = 0
var botBeginTimeOut int64 = 1000    // ms
var botShootDistance float64 = 20.0 // ms

const (
	NEEDPLAYERS = 1
)

type BotMoveState struct {
	beginTimeOut   int64
	currentTimeOut int64
	//tryAngle int64
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
	angle := 90.0 + ((math.Atan(tangAlfa) * 180.0) / math.Pi)
	angle += 180.0
	//fmt.Println("set attack angle =", angle)
	bot.tryShootAngle = int64(angle)
}

func (pos1 Position) pointDistance(pos2 Position) float64 {
	return math.Sqrt((pos1.x-pos2.x)*(pos1.x-pos2.x) + (pos1.y-pos2.y)*(pos1.y-pos2.y))
}

type Bot struct {
	//---------------------Bot feature
	MoveState BotMoveState
	//--------------------- update by engine
	pos      Position
	id       string
	isAttack bool
	//HealfPoint  int64
	//Scores      int64
	//----------------------write interface
	net.Conn
	//--------------------- update by client
	tryX          float64
	tryY          float64
	tryShootAngle int64
	tryMoveAngle  int64
	tryAttack     bool
}

func checkPlayersCount() {
	if currentPlayerCount > NEEDPLAYERS {
		//	removeBot()
		return
	}
	if currentPlayerCount < NEEDPLAYERS {
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
		for _, player := range players {
			if bot.id != player.id {
				if player.pos.pointDistance(bot.pos) < botShootDistance {

					//fmt.Println("try ATTACK player with ID = ", player.id, "  my id = ", bot.id)
					bot.isAttack = true
					bot.makeShootAngle(player.pos)
				}
			}
		}
	}
}

func addBot() {
	currentPlayerCount++ // KOSTUL
	fmt.Println("try addBot")
	tcpAddr, err_adr := net.ResolveTCPAddr("tcp4", "127.0.0.1:55555")
	if err_adr == nil {
		conn, err_dial := net.DialTCP("tcp", nil, tcpAddr)
		if err_dial == nil {

			// make new Player
			var newPlayer *Bot = new(Bot)
			newPlayer.id = "X"
			newPlayer.Conn = conn
			newPlayer.MoveState.beginTimeOut = botBeginTimeOut
			newmessage := "0;true;" + "nik name;" + "0" + ";\n"

			bots[currentBotId] = newPlayer
			go readServerData(newPlayer)
			currentBotId++
			fmt.Println("try send First Message from bot")
			newPlayer.Conn.Write([]byte(newmessage))
			fmt.Println("Message is sended")
		}
	}
	fmt.Println("end addBot")
}

func removeBot() {
	currentPlayerCount-- // KOSTUL
	defer bots[currentBotId].Conn.Close()
	//make task
	var newTask Task
	newTask.taskType = TASK_DELCLIENT
	newTask.clientId = bots[currentBotId].id
	taskChan <- newTask
	//delete from map
	delete(bots, currentBotId)
	currentBotId--
}

// Main game loop
func tikTak() {
	//time.Sleep(time.Millisecond * 4000)
	fmt.Println("start botsFather ")
	for {
		time.Sleep(time.Millisecond * TICKPERIOD)
		checkPlayersCount()
		botWalk()
		botShoot() // set/unset Bot.isAttack true/false
		botMessage()
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

func init() {
	//fmt.Println("Create botsFather ")
	bots = make(map[int]*Bot)
	currentPlayerCount = 0
	go tikTak()
}

func readServerData(bot *Bot) {
	//defer makeDeletePlayerTask(bot)
	reader := bufio.NewReader(bot.Conn)
	for {
		message, err := reader.ReadString('$')
		if err == nil {
			//log.Println("readPlayersInput_____", message)
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
	strArr := strings.Split(strings.TrimSuffix(message, "$"), ";")
	if len(strArr) < 2 {
		println("server str input len = ", len(strArr))
		println("server str =", message)
		return
	}

	switch {
	// set id
	case strArr[0] == strconv.FormatInt(MSG_YOURID, 10):

		bot.id = strArr[1]
	case strArr[0] == strconv.FormatInt(MSG_STATE, 10):
		//message = append(message, strconv.FormatInt(MSG_STATE, 10))
		if strArr[1] == bot.id {
			//message = append(message, strconv.FormatFloat(currentPlayer.pos.x, 'f', 1, 64))
			//message = append(message, strconv.FormatFloat(currentPlayer.pos.y, 'f', 1, 64))
			//message = append(message, strconv.FormatInt(int64(currentPlayer.pos.angle), 10))
			//message = append(message, strconv.FormatInt(currentPlayer.healthPoint, 10))
			//message = append(message, strconv.FormatBool(currentPlayer.pos.isAttack))
			if len(strArr) < 5 {
				println("server str input len = ", len(strArr))
				println("server str =", message)
				return
			}
			x, err_x := strconv.ParseFloat(strArr[2], 64)
			y, err_y := strconv.ParseFloat(strArr[3], 64)
			if (err_x != nil) || (err_y != nil) {
				return
			}

			bot.pos.x = x
			bot.pos.y = y
		}
	default:
		return
	}
}
