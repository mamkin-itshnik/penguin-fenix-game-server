package main

import (
	"log"
	"strconv"

	//"./collision2d"
	//"github.com/mamkin-itshnik/collision2d"
)

func makePlayerState(currentPlayer *Player) {

	calculateMoving(currentPlayer)
	calculateVisibleAndShoots(currentPlayer)
}

func calculateMoving(currentPlayer *Player) {

	makeNewPosition(currentPlayer)

	//collision logic
	isCollision, point := checkCollision(currentPlayer.pos.x, currentPlayer.pos.y)
	if isCollision {
		currentPlayer.pos.x = point.X
		currentPlayer.pos.y = point.Y
	}
}

func calculateVisibleAndShoots(currentPlayer *Player) {
	//----------------------------------------------- CHECK VISIBILITY PLAYERS + SHOOT
	for id_other, otherPlayer := range players {

		if !otherPlayer.isPlay {
			continue
		}
		if id_other != currentPlayer.id {
			isVisibleObjects(currentPlayer, otherPlayer)
		}
	}
}

func isVisibleObjects(currentPlayer, otherPlayer *Player) {
	dist := currentPlayer.pos.pointDistance(otherPlayer.pos)

	//check distance to objects
	if dist < VISABILITY_ZONE {

		//set new state
		currentPlayer.visiblePlayersId_new[otherPlayer.id] = true
		//check attack
		checkPlayerShoot(currentPlayer, otherPlayer)

	} else {

		_, ok := currentPlayer.visiblePlayersId_new[otherPlayer.id]
		if ok{
			currentPlayer.visiblePlayersId_new[otherPlayer.id] = false
		}		
	}
}

func makeOtherPosMsg(currentPlayer *Player) []string {

	var message []string

	if !currentPlayer.isPlay {
		return message
	}
	for id, isVisible := range currentPlayer.visiblePlayersId_new {
		// check actual state
		if isVisible {

			_, ok := currentPlayer.visiblePlayersId_old[id]
			if !ok{
				// make message about new player
				message = append(message, spliceSubMessage(strconv.FormatInt(MSG_ADD_PLAYER, 10), id))

			}
			// add other player position
			message = append(message, spliceSubMessage(getPlayerPosMsg(players[id])...))
			currentPlayer.visiblePlayersId_old[id] = true

		} else {
			isPreviousVisible := currentPlayer.visiblePlayersId_old[id]
			if isPreviousVisible {
				// make message about player destruct
				message = append(message, spliceSubMessage(strconv.FormatInt(MSG_DELETE_PLAYER, 10), id))
				currentPlayer.visiblePlayersId_old[id] = false
			}else{
				delete(currentPlayer.visiblePlayersId_old, id)
				delete(currentPlayer.visiblePlayersId_new, id)
			}
		}
	}

	return message
}

func checkPlayerShoot(currentPlayer, otherPlayer *Player) {
	//-----------------------------------------------SHOOT
	if !currentPlayer.pos.isAttack {
		return
	}
	if currentPlayer.id == otherPlayer.id {
		log.Println("************************************************ WATAFUCK?????", currentPlayer.id)
		return
	}
	dist := currentPlayer.pos.shootDistance(otherPlayer.pos)
	if dist < 0 {
		otherPlayer.healthPoint -= WEAPONBASEDAMAGE
		if otherPlayer.healthPoint < 0 {
			playerKnokPlayer(currentPlayer, otherPlayer)
		}
	}
}

func playerKnokPlayer(winner, looser *Player) {
	//--- Score++ HP++
	winner.scorePoint += 100
	winner.healthPoint += int64(HPHEALLERP *
		float64(STARTHEALTHPOINT-winner.healthPoint))

	looser.healthPoint = STARTHEALTHPOINT

	var scoreTask Task
	scoreTask.clientId = winner.id
	scoreTask.taskType = TASK_UPDATESCORE
	taskChan <- scoreTask

	//--- Looser's respawn
	var looseTask Task
	looseTask.clientId = looser.id
	looseTask.taskType = TASK_RESPAWNCLIENT
	taskChan <- looseTask
}
