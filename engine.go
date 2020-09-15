package main

import (
	"math"
	"strconv"
)

var moveSpeed = 0.3
var shootDistance = 40.0
var startHealfPoint int64 = 50
var objectRadius float64 = 1.1 //0.4
var minPos float64 = -20.5
var maxPos float64 = 20.5

func engine_SolveTask(currentPlayer *Player) {

	for _, task := range currentPlayer.TaskMap {
		switch {
		case task.TaskType == CLIENTMOVE:
			//--------------------------------------------player shoot
			engine_makePlayerPos(currentPlayer)
			//------------------------------------------------------END
		case task.TaskType == DELCLIENT:
			//--------------------------------------------player shoot
			return
			//------------------------------------------------------END
		default:
			//WRONG
		}
	}

}

func engine_makePlayerPos(currentPlayer *Player) {

	//-----------------------------------------------POSITION
	tryX, errX := strconv.ParseFloat(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[0], 32)
	tryY, errY := strconv.ParseFloat(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[1], 32)
	tryA, errA := strconv.Atoi(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[2])
	tryAttack, errAttack := strconv.ParseBool(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[3])

	if (errX == nil) && (errY == nil) && (errA == nil) && (errAttack == nil) {
		currentPlayer.Pos.X += (tryX * moveSpeed)
		currentPlayer.Pos.Y += (tryY * moveSpeed)
		currentPlayer.Pos.Angle = tryA
		currentPlayer.isAttack = tryAttack
	} else {
		//WRONG
	}

	//-----------------------------------------------SHOOT
	if currentPlayer.isAttack {
		for id_other, otherPlayer := range players {
			if id_other != currentPlayer.Id {
				//fmt.Printf(" call distance \n")
				var dist float64 = currentPlayer.Pos.Distance(otherPlayer.Pos)
				if dist < 0 {
					//println("HEALF --")
					otherPlayer.HealfPoint--
					if otherPlayer.HealfPoint < 0 {
						var newTask Task
						newTask.ClientID = otherPlayer.Id
						newTask.TaskType = REBURNCLIENT
						TaskChan <- newTask
						//newTask.TaskArgs = make([]string, 1)
					}

				} else {
					// vse ok
				}
			}
		}
	}

	//----------------------------------------------END
}

func (startPoint Position) Distance(target Position) float64 {
	//return shoot distance between point "target" and "OTREZOK" started in "startPoint"

	var distance float64
	endPointX := startPoint.X - shootDistance*math.Cos((float64(startPoint.Angle)*(math.Pi/180)))
	endPointY := startPoint.Y + shootDistance*math.Sin((float64(startPoint.Angle)*(math.Pi/180)))

	endPointX_back := startPoint.X + shootDistance*math.Cos((float64(startPoint.Angle)*(math.Pi/180)))
	endPointY_back := startPoint.Y - shootDistance*math.Sin((float64(startPoint.Angle)*(math.Pi/180)))

	foreward_distance := math.Sqrt(math.Pow((endPointX-target.X), 2) + math.Pow((endPointY-target.Y), 2))
	backward_distance := math.Sqrt(math.Pow((endPointX_back-target.X), 2) + math.Pow((endPointY_back-target.Y), 2))

	//fmt.Printf("___________________ \n")
	//fmt.Printf("ANGLE = %d \n", startPoint.Angle)
	if foreward_distance > backward_distance {
		distance = math.Sqrt(math.Pow((startPoint.X-target.X), 2) + math.Pow((startPoint.Y-target.Y), 2))
	} else {
		distance = ((startPoint.Y-endPointY)*target.X + (endPointX-startPoint.X)*target.Y + (endPointY*startPoint.X - endPointX*startPoint.Y)) /
			math.Sqrt(math.Pow((endPointX-startPoint.X), 2)+math.Pow((endPointY-startPoint.Y), 2))
		//fmt.Printf("TARGET   %d : %d \n", target.X, target.Y)
		//fmt.Printf("___________________ \n")
		//fmt.Printf("distanse = %d \n", distance)
		//fmt.Printf("___________________ \n")
	}

	//fmt.Printf("start %d : %d \n", startPoint.X, startPoint.Y)
	//fmt.Printf("end   %d : %d \n", endPointX, endPointY)
	//fmt.Printf("distanse = %d \n", distance)
	distance = (math.Abs(distance) - objectRadius)
	//fmt.Printf("distanse = %d \n", distance)
	return distance
}
