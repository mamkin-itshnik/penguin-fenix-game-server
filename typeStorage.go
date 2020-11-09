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
	TASK_UPDATESCORE   = 4
)

// message type from server
const (
	MSG_YOUR_ID       = 0
	MSG_STATE         = 1
	MSG_HISCORE       = 2
	MSG_DELETE_PLAYER = 3
	MSG_ADD_PLAYER    = 4
	MSG_PLAYER_COUNT  = 5
)

// message type from client
const (
	MSG_CLIENT_WANT_PLAY             = 0
	MSG_CLIENT_WANT_MOVE             = 1
	MSG_CLIENT_WANT_GET_PLAYER_COUNT = 2
)

// engine shit
const (
	SCORE_ENTYTY_COUNT         = 5
	MOVESPEED                  = 0.6 // 0.3 standart
	SHOOTDISTANCE              = 15.0
	VISABILITY_ZONE            = 15.5
	STARTHEALTHPOINT   int64   = 50
	WEAPONBASEDAMAGE   int64   = 1    //1
	OBJECTRADIUS       float64 = 1.1  //0.4
	MAX_XPOS           float64 = 35.5 // 35.5
	MAX_YPOS           float64 = 35.5 // 35.5
	HPHEALLERP         float64 = 0.5
	MAXSCORELINE       int64   = 5
)

// some constants
const (
	LOG_FILE_NAME string        = "penguin_royale_logs.txt"
	TICKPERIOD    time.Duration = 100
)

type Task struct {
	clientId string
	taskType int
	taskArgs []string
}

type ScoreEntyty struct {
	id         string
	scorePoint int64
	nikName    string
}

type Player struct {

	// by engine
	id          string
	pos         Position
	healthPoint int64
	scorePoint  int64

	// multiplayer data
	net.Conn
	isPlay               bool
	nikName              string
	skinID               int64
	visiblePlayersId_old map[string]bool
	visiblePlayersId_new map[string]bool


	// from network
	wannaPos Position
}

type Position struct {
	x, y     float64
	angle    int64
	isAttack bool
}
