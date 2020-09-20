package main

import (
	"math"
	"strconv"
)

func makePlayerPos(currentPlayer *Player) {

	//-----------------------------------------------POSITION

	currentPlayer.pos.x += (currentPlayer.wannaPos.x * MOVESPEED)
	currentPlayer.pos.y += (currentPlayer.wannaPos.y * MOVESPEED)
	currentPlayer.pos.angle = currentPlayer.wannaPos.angle
	currentPlayer.pos.isAttack = currentPlayer.wannaPos.isAttack

	//-----------------------------------------------SHOOT
	if currentPlayer.pos.isAttack {
		for id_other, otherPlayer := range players {
			if id_other != currentPlayer.id {
				dist := currentPlayer.pos.distance(otherPlayer.pos)
				if dist < 0 {
					otherPlayer.healthPoint--
					if otherPlayer.healthPoint < 0 {

						//--- Score++ HP++
						currentPlayer.scorePoint += (otherPlayer.scorePoint / 2) + 10
						currentPlayer.healthPoint += int64(HPHEALLERP *
							float64(STARTHEALTHPOINT-currentPlayer.healthPoint))

						var scoreTask Task
						scoreTask.clientId = currentPlayer.id
						scoreTask.taskType = TASK_UPDATESCORE
						taskChan <- scoreTask

						//--- Looser's respawn
						var looseTask Task
						looseTask.clientId = otherPlayer.id
						looseTask.taskType = TASK_RESPAWNCLIENT
						taskChan <- looseTask
					}
				}
			}
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
	// NOTICE - no more score point in player pos
	//message = append(message, strconv.FormatInt(int64(currentPlayer.scorePoint), 10))
	return message
}
func getPlayerScore(currentPlayer *Player) []string {
	var message []string
	message = append(message, strconv.FormatInt(MSG_HISCORE, 10))
	message = append(message, currentPlayer.id)
	message = append(message, strconv.FormatInt(int64(currentPlayer.scorePoint), 10))
	return message
}

func (startPoint Position) distance(target Position) float64 {
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
		distance = ((startPoint.y-endPointY)*target.x + (endPointX-startPoint.x)*target.y + (endPointY*startPoint.x - endPointX*startPoint.y)) /
			math.Sqrt(math.Pow((endPointX-startPoint.x), 2)+math.Pow((endPointY-startPoint.y), 2))
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
