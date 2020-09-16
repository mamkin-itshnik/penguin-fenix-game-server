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
	case strArr[0] == "2":
		//--------------------------------------------player moves
		if len(strArr) < 4 {
			return newWrongTask(newTask)
		}
		newTask.TaskType = CLIENTMOVE
		newTask.TaskArgs = make([]string, 4)
		_, err_x := strconv.ParseFloat(strArr[1], 32)
		_, err_y := strconv.ParseFloat(strArr[2], 32)
		_, err_a := strconv.ParseFloat(strArr[3], 32)
		_, err_attack := strconv.ParseBool(strArr[4])

		if (err_x == nil) && (err_y == nil) && (err_a == nil) && (err_attack == nil) {

			//println("read client data for MOVE: ", strArr)
			newTask.TaskArgs[0] = strArr[1] //X
			newTask.TaskArgs[1] = strArr[2] //Y
			newTask.TaskArgs[2] = strArr[3] //A
			newTask.TaskArgs[3] = strArr[4] //A
			return newTask
		} else {
			return newWrongTask(newTask)
		}
		//------------------------------------------------------END
	case strArr[0] == "3":
		//--------------------------------------------

		//------------------------------------------------------END
	default:
		//WRONG
		//println("______________WRONG")
		//println("read client data char value =", strArr[0][0])
		//println("client data =", str)
		//println("______________WRONG")
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
		message += strconv.FormatInt(int64(currentPlayer.Pos.Angle), 10) + ";"
		message += strconv.FormatInt(currentPlayer.HealfPoint, 10) + ";"
		message += strconv.FormatBool(currentPlayer.isAttack) + ";"
		message += strconv.FormatInt(int64(currentPlayer.Scores), 10) + ";"
		message += "\n"
		//------------------------------------------------------END
	case taskNumber == DELCLIENT:
		//--------------------------------------------player delete
		return message
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
