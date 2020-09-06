package main

import (
	"strconv"
	"strings"
)

func newWrongTask(newTask Task) Task {
	newTask.TaskType = WRONGTASK
	return newTask
}
func TP_makeClientInputsTask(str string, playerID string) Task {
	var newTask Task
	newTask.ClientID = playerID

	strArr := strings.Split(str, ";")
	//println("strArr[0] =", strArr[0])
	//println("len(strArr[0] =", len(strArr[0]))
	if len(strArr) < 2 {
		println("str len = ", len(strArr))
		return newWrongTask(newTask)
	}

	switch {
	case strArr[0] == "0":
		//--------------------------------------------manage clients
		//core_AddPlayer(newTask.ClientID) may be reconnect
		//------------------------------------------------------END
	case strArr[0] == "1":
		//--------------------------------------------player moves
		newTask.TaskType = CLIENTMOVE
		newTask.TaskArgs = make([]string, 3)
		_, err_x := strconv.ParseFloat(strArr[1], 32)
		_, err_y := strconv.ParseFloat(strArr[1], 32)
		_, err_a := strconv.ParseFloat(strArr[1], 32)

		if (err_x == nil) && (err_y == nil) && (err_a == nil) {

			//println("read client data for MOVE: ", strArr)
			newTask.TaskArgs[0] = strArr[1] //X
			newTask.TaskArgs[1] = strArr[2] //Y
			newTask.TaskArgs[2] = strArr[3] //A
			return newTask
		} else {
			return newWrongTask(newTask)
		}
		//------------------------------------------------------END
	case strArr[0] == "2":
		//--------------------------------------------player shoot
		newTask.TaskType = CLIENTSHOOT
		newTask.TaskArgs = make([]string, 1)
		_, err := strconv.ParseFloat(strArr[1], 32)
		if err == nil {

			//println("read client data for SHOOT: ", strArr)
			newTask.TaskArgs[0] = strArr[1]
			return newTask
		} else {
			return newWrongTask(newTask)
		}
		//------------------------------------------------------END
	default:
		//WRONG
		println("______________WRONG")
		println("read client data char value =", strArr[0][0])
		println("client data =", str)
		println("______________WRONG")
		return newWrongTask(newTask)
	}

	return newWrongTask(newTask)
	//AddPlayer(playerID)
}

func TP_makeStringTask(currentPlayer *Player, taskNumber int) string {
	var message string

	switch {
	case taskNumber == ADDCLIENT:
		//--------------------------------------------player moves
		//message += strconv.FormatInt(ADDCLIENT, 10) + ";"
		//message += "add;"
		//message += currentPlayer.Id + ";"
		//message += "\n"
		//_, ok := currentPlayer.TaskMap[ADDCLIENT]
		//if ok {
		//	delete(currentPlayer.TaskMap, ADDCLIENT)
		//}
		//------------------------------------------------------END
	case taskNumber == CLIENTMOVE:
		//--------------------------------------------player shoot
		message += strconv.FormatInt(CLIENTMOVE, 10) + ";"
		message += currentPlayer.Id + ";"
		message += strconv.FormatFloat(currentPlayer.Pos.X, 'f', 1, 64) + ";"
		message += strconv.FormatFloat(currentPlayer.Pos.Y, 'f', 1, 64) + ";"
		message += strconv.FormatFloat(currentPlayer.Pos.Angle, 'f', 1, 64) + ";"
		message += strconv.FormatInt(currentPlayer.HealfPoint, 10) + ";"
		message += "\n"
		//------------------------------------------------------END
	default:
		//WRONG
	}

	return message
}

//
/*func TP_makeNewClientsTask(playerID string) Task {
	var newTask Task
	newTask.TaskType = ADDCLIENT
	newTask.ClientID = playerID
	return newTask
}*/

/*//var dX, dY, nAngl float64
	if strings.Contains(strArr[0], "XD") {
		if strArr[1] == "X" {
			//println("red client data: ", message)
			//Clients[id].isAttack = false
			continue
		}
		//Clients[id].isAttack = true
		//nAngl, _ = strconv.ParseFloat(strArr[1], 32)
		//Clients[id].Pos.ShootAngle = nAngl
		continue
	}

	if len(strArr) > 3 {
		//dX, _ = strconv.ParseFloat(strArr[1], 32)
		//dY, _ = strconv.ParseFloat(strArr[2], 32)
		//nAngl, _ = strconv.ParseFloat(strArr[3], 32)

			g.Clients[id].Pos.TryDeltaX = dX
			g.Clients[id].Pos.TryDeltaY = dY
			g.Clients[id].Pos.Angle = nAngl
	}
	return newTask
}*/
