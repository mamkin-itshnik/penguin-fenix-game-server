package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"os"
	"strings"
	"fmt"

	//"./collision2d"
	"github.com/mamkin-itshnik/collision2d"
)


func parseClientWantGetPlayerCount(currentPlayer *Player) {

	newMessage := strconv.Itoa(MSG_PLAYER_COUNT) + ";"
	newMessage += strconv.Itoa(atomicIntPlayerCount)
	newMessage += "$"

	currentPlayer.Conn.Write([]byte(newMessage))
}

func parseClientWantToMove(strArr []string, currentPlayer *Player) {
	if len(strArr) < 4 {
		println("player str input len = ", len(strArr))
		println("player str =", strArr)
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
}

func parseClientWantToPlay(strArr []string, currentPlayer *Player) {

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

	currentPlayer.skinID = newSkinID
	currentPlayer.nikName = newNikName
	currentPlayer.isPlay = isPlay

	if currentPlayer.isPlay {

		newMessage := strconv.Itoa(MSG_YOUR_ID) + ";"
		newMessage += currentPlayer.id + "$"

		currentPlayer.Conn.Write([]byte(newMessage))

		var startTask Task
		startTask.clientId = currentPlayer.id
		startTask.taskType = TASK_RESPAWNCLIENT
		taskChan <- startTask
	} else {
		//---- delete player from other player in client
		var newTask Task
		newTask.taskType = TASK_DELCLIENT
		newTask.clientId = currentPlayer.id
		taskChan <- newTask
	}
}

func solveTaskDellClient(newTask *Task) {
	log.Println("func core_DelPlayer(playerID string)")

	playersMutex.Lock()
	log.Println("NOW PLAYER COUNT = ", len(players))
	playersMutex.Unlock()

	sendToPlayers(prepareMsg(strconv.FormatInt(MSG_DELETE_PLAYER, 10), newTask.clientId))
}

func solveTaskRespawnClient(newTask *Task) {

	playersMutex.Lock()
	player, ok := players[newTask.clientId]
	playersMutex.Unlock()

	if !ok {
		log.Println("WTF? Respawn player that doesn't exist in map",
			newTask.clientId)
		return
	}

	// make random state
	player.healthPoint = STARTHEALTHPOINT
	player.pos = makeRandomPos()
	player.scorePoint /= 2
}

func solveTaskUpdateScore(newTask *Task) {

	ok := needUpdScore(newTask.clientId)

	if ok {
		updateScore()
	}
}

func needUpdScore(playerId string) bool {

	playersMutex.Lock()
	player, ok := players[playerId]
	playersMutex.Unlock()

	if !ok {
		return false
	}else{
		return checkNewScore(player)
	}
}

func checkNewScore(player * Player)bool{
		for i, scoreEntyty := range topScorePlayer {

		if player.scorePoint >= scoreEntyty.scorePoint {
			moveDownInScorePosition(i)
			topScorePlayer[i].nikName = player.nikName
			topScorePlayer[i].scorePoint = player.scorePoint
			topScorePlayer[i].id = player.id
			return true
		}
	}
	return false
}

func taskOverflow() {

	f, err := os.OpenFile("chan_log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("This is a test log entry")
	for i := 0; i < 10000; i++ {
		qwe := <-taskChan
		log.Println(qwe.taskType)
	}
	fmt.Println("___________")
	wg.Add(-1)
}


func setZeroPropertyForPlayer(newPlayer *Player) {

	newPlayer.healthPoint = STARTHEALTHPOINT
	newPlayer.scorePoint = 0
	newPlayer.pos.angle = 0
	newPlayer.isPlay = false
	newPlayer.visiblePlayersId_new = make(map[string]bool)
	newPlayer.visiblePlayersId_old = make(map[string]bool)
}


func prepareMsg(parts ...string) string {
	msg := strings.Join(parts, ";")
	msg += "/"
	return msg
}

func appendMessageParts(original *string, parts ...string) {
	*original += strings.Join(parts, ";")
	*original += "/"
}

func spliceSubMessage(parts ...string) string {
	msg := strings.Join(parts, ";")
	return msg
}

func spliceMessages(parts ...string) string {
	msg := strings.Join(parts, "/")
	return msg
}

func writeToPlayers(str string) {
	playersWriteChan <- str
}

func sendToPlayers(parts ...string) {
	msg := strings.Join(parts, "/")
	// ADD stop byte as $ symbol
	msg += "$"
	//log.Println("send to all:", msg)
	playersWriteChan <- msg
}

func moveDownInScorePosition(beginPos int) {

	//reverse run
	i := len(topScorePlayer) - 1
	for i > beginPos {
		topScorePlayer[i] = topScorePlayer[i-1]
		i--
	}
}

func updateScore(){

	var allScoreMsg string

	playersMutex.Lock()
		for _, scoreEntyty := range topScorePlayer {

			playerMsg := getScoreFromEntyty(&scoreEntyty)
			allScoreMsg += prepareMsg(playerMsg...)
			// split players score in one message.
			allScoreMsg += "#"
		}
		playersMutex.Unlock()

		if len(allScoreMsg) != 0 {
			sendToPlayers(allScoreMsg)
		}
}

func setUpLogFile(){
	f, err := os.OpenFile(LOG_FILE_NAME, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Println("This is a test log entry")
}

func getScoreFromEntyty(currentScore *ScoreEntyty) []string {
	var message []string
	message = append(message, strconv.FormatInt(MSG_HISCORE, 10))
	message = append(message, currentScore.id)
	message = append(message, currentScore.nikName)
	message = append(message, strconv.FormatInt(int64(currentScore.scorePoint), 10))
	// message = append(message, "/") //end of player state
	return message
}
func makeRandomPos() Position {
	var pos Position
	pos.x = -MAX_XPOS + rand.Float64()*(MAX_XPOS*2.0)
	pos.y = -MAX_YPOS + rand.Float64()*(MAX_YPOS*2.0)
	return pos
}

func deleteOldPlayerIdFromOtherPlayers(player *Player) {
	for id, playedPlayer := range players {
		if id != player.id {
			delete(playedPlayer.visiblePlayersId_new, player.id)
			delete(playedPlayer.visiblePlayersId_old, player.id)
		}
	}
}

func getPlayerPosMsg(currentPlayer *Player) []string {
	var message []string
	message = append(message, strconv.FormatInt(MSG_STATE, 10))
	message = append(message, currentPlayer.id)
	message = append(message, strconv.FormatFloat(currentPlayer.pos.x, 'f', 1, 64))
	message = append(message, strconv.FormatFloat(currentPlayer.pos.y, 'f', 1, 64))
	message = append(message, strconv.FormatInt(int64(currentPlayer.pos.angle), 10))
	message = append(message, strconv.FormatInt(currentPlayer.healthPoint, 10))
	message = append(message, strconv.FormatBool(currentPlayer.pos.isAttack))
	return message
}

func makeNewPosition(currentPlayer *Player) {
	currentPlayer.pos.x += (currentPlayer.wannaPos.x * MOVESPEED)
	currentPlayer.pos.y += (currentPlayer.wannaPos.y * MOVESPEED)
	currentPlayer.pos.angle = currentPlayer.wannaPos.angle
	currentPlayer.pos.isAttack = currentPlayer.wannaPos.isAttack
}

func edgeDistance(startPoint, endPoint, targetPoint collision2d.Vector) float64 {

	var distance float64

	distance = ((startPoint.Y-endPoint.Y)*targetPoint.X +
		(endPoint.X-startPoint.X)*targetPoint.Y + (endPoint.Y*startPoint.X - endPoint.X*startPoint.Y)) /
		math.Sqrt(math.Pow((endPoint.X-startPoint.X), 2)+math.Pow((endPoint.Y-startPoint.Y), 2))
	return distance
}

func (pos1 Position) pointDistance(pos2 Position) float64 {
	return math.Sqrt((pos1.x-pos2.x)*(pos1.x-pos2.x) + (pos1.y-pos2.y)*(pos1.y-pos2.y))
}

func (startPoint Position) shootDistance(target Position) float64 {

	//return shoot distance between point "target" and "OTREZOK" started in "startPoint"

	var distance float64
	endPointX := startPoint.x - SHOOTDISTANCE*math.Cos((float64(startPoint.angle)*(math.Pi/180)))
	endPointY := startPoint.y + SHOOTDISTANCE*math.Sin((float64(startPoint.angle)*(math.Pi/180)))

	endPointX_back := startPoint.x + SHOOTDISTANCE*math.Cos((float64(startPoint.angle)*(math.Pi/180)))
	endPointY_back := startPoint.y - SHOOTDISTANCE*math.Sin((float64(startPoint.angle)*(math.Pi/180)))

	foreward_distance := math.Sqrt(math.Pow((endPointX-target.x), 2) + math.Pow((endPointY-target.y), 2))
	backward_distance := math.Sqrt(math.Pow((endPointX_back-target.x), 2) + math.Pow((endPointY_back-target.y), 2))

	//fmt.Printf("___________________ \n")
	//fmt.Printf("ANGLE = %d \n", startPoint.angle)
	if foreward_distance > backward_distance {
		distance = math.Sqrt(math.Pow((startPoint.x-target.x), 2) + math.Pow((startPoint.y-target.y), 2))
	} else {
		distance = edgeDistance(collision2d.Vector{
			startPoint.x, startPoint.y},
			collision2d.Vector{endPointX, endPointY},
			collision2d.Vector{target.x, target.y})

		//fmt.Printf("TARGET   %d : %d \n", target.x, target.y)
		//fmt.Printf("___________________ \n")
		//fmt.Printf("distanse = %d \n", distance)
		//fmt.Printf("___________________ \n")
	}

	//fmt.Printf("start %d : %d \n", startPoint.x, startPoint.y)
	//fmt.Printf("end   %d : %d \n", endPointX, endPointY)
	//fmt.Printf("distanse = %d \n", distance)
	distance = (math.Abs(distance) - OBJECTRADIUS)
	//fmt.Printf("distanse = %d \n", distance)
	return distance
}
