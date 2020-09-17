package main

import "net"

// Tasks types
const (
	DELCLIENT     = 1
	CLIENTMOVE    = 2
	RESPAWNCLIENT = 3
)

// message type
const (
	YOURID  = 0
	STATE   = 1
	HISCORE = 2
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
