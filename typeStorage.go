package main

import (
	"net"
	"time"
)

// Tasks types
const (
	TASK_DELCLIENT     = 1
	TASK_CLIENTMOVE    = 2
	TASK_RESPAWNCLIENT = 3
)

// message type from server
const (
	MSG_YOURID        = 0
	MSG_STATE         = 1
	MSG_HISCORE       = 2
	MSG_KILLPLAYER    = 3
	MSG_RESPAWNPLAYER = 4
)

// engine shit
const (
	MOVESPEED                = 0.3
	SHOOTDISTANCE            = 40.0
	STARTHEALTHPOINT int64   = 50
	OBJECTRADIUS     float64 = 1.1 //0.4
	MINPOS           float64 = -20.5
	MAXPOS           float64 = 20.5
	HPHEALLERP       float64 = 0.5
)

// some constants
const (
	TICKPERIOD time.Duration = 10
)

type Task struct {
	clientId string
	taskType int
	taskArgs []string
}

type Player struct {

	// by engine
	id          string
	skin        int64
	nickname    string
	pos         Position
	healthPoint int64
	scorePoint  int64

	net.Conn

	// from network
	wannaPos Position
}

type Position struct {
	x, y     float64
	angle    int64
	isAttack bool
}
