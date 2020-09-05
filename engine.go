package main

import "strconv"

func engine_SolveTask(currentPlayer *Player) {

	for _, task := range currentPlayer.TaskMap {
		switch {
		case task.TaskType == CLIENTSHOOT:
			//--------------------------------------------player moves
			engine_makePlayerShoot(currentPlayer)
			//------------------------------------------------------END
		case task.TaskType == CLIENTMOVE:
			//--------------------------------------------player shoot
			engine_makePlayerPos(currentPlayer)
			//------------------------------------------------------END
		default:
			//WRONG
		}
	}

}

func engine_makePlayerPos(currentPlayer *Player) {

	tryX, errX := strconv.ParseFloat(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[0], 32)
	tryY, errY := strconv.ParseFloat(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[1], 32)
	tryA, errA := strconv.ParseFloat(currentPlayer.TaskMap[CLIENTMOVE].TaskArgs[2], 32)

	if (errX == nil) && (errY == nil) && (errA == nil) {
		currentPlayer.Pos.X += tryX
		currentPlayer.Pos.Y += tryY
		currentPlayer.Pos.Angle += tryA
	} else {
		//WRONG
	}
}

func engine_makePlayerShoot(currentPlayer *Player) {

	//check bullets collisions
	/*for id_2, otherPlayer := range players {
			if currentPlayer.isAttack && g.Clients[id].isAttack && (id_2 != id) {
			//fmt.Printf(" call distance \n")
			var dist float64 = g.Clients[id].Pos.Distance(g.Clients[id_2].Pos)
			if dist < 0 {
				println("HEALF --")
				g.Clients[id_2].HealPoint--
				println("HEALF  = %d", g.Clients[id_2].HealPoint)
			} else {
				// vse ok
			}
		}
	}*/
}
