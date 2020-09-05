package main
func newWrongTask(newTask Task)Task{
	
}
func TP_makeTask(str string,playerID string) Task {
	var newTask Task
	newTask.ClientID = playerID

	strArr := strings.Split(message, ";")
	if len(strArr) < 2 {
		newTask.TaskType = WRONGTASK
		return newTask
	}

	switch strArr[0]{
	case 'A':
		//--------------------------------------------manage clients
		//core_AddPlayer(newTask.ClientID) may be reconnect
		//------------------------------------------------------END
	case 'B':
		//--------------------------------------------player moves

		//------------------------------------------------------END
	case 'C':
		//--------------------------------------------player shoot
		newTask.TaskType = CLIENTSHOOT

		shootAngle, err = strconv.ParseFloat(strArr[1], 32)
		if err != nil {

			println("read client data: ", message)

			var newTask Task = TP_makeTask(message, id)
			return newTask
		}else{
			return newWrongTask(newTask)
		}
		newTask.TaskArgs
		
		Clients[id].Pos.ShootAngle = nAngl
		//------------------------------------------------------END
	default:
		//WRONG
		return newWrongTask(newTask)
	}
	//AddPlayer(playerID)
}

	//var dX, dY, nAngl float64
	if strings.Contains(strArr[0], "XD") {
		if strArr[1] == "X" {
			//println("read client data: ", message)
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
	}*/
	return newTask
}
