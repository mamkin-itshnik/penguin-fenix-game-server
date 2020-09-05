package main

import "net"

type TaskID int

// Tasks constats
const (
	ADDCLIENT   = 0
	DELCLIENT   = 1
	CLIENTMOVE  = 2
	CLIENTSHOOT = 3
	WRONGTASK   = 4
)
const (
	TASKCOUNT = 5
)

//for ConnectManager.go
type Client struct {
	net.Conn
	clientID string
}

// ConnectManager.go + core.go  + engine.go
type Task struct {
	ClientID string
	TaskType int
	TaskArgs []string
}

type Player struct {
	ClientState
	TaskArray []Task
}

type Position struct {
	X, Y, Angle float64
}

type ClientState struct {
	Pos        Position
	Id         string
	isAttack   bool
	HealfPoint int64
}
