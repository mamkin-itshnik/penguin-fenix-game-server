package gameCore

import (
	"bufio"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
	"sync"
)

var instance *gameCore
var once sync.Once
var moveSpeed = 0.3
var shootDistance = 40.0

//var bulletSpeed = 2.5
//var bulletLifeTime = 5
var startHealfPoint int64 = 50
var objectRadius float64 = 1.1 //0.4
var shootRate = 1

type Position struct {
	X, Y, Angle, ShootAngle, TryDeltaX, TryDeltaY float64
}

func (startPoint Position) Distance(target Position) float64 {
	//return shoot distance between point "target" and "OTREZOK" started in "startPoint"

	var distance float64
	endPointX := startPoint.X - shootDistance*math.Cos((startPoint.ShootAngle*(math.Pi/180)))
	endPointY := startPoint.Y + shootDistance*math.Sin((startPoint.ShootAngle*(math.Pi/180)))

	endPointX_back := startPoint.X + shootDistance*math.Cos((startPoint.ShootAngle*(math.Pi/180)))
	endPointY_back := startPoint.Y - shootDistance*math.Sin((startPoint.ShootAngle*(math.Pi/180)))

	foreward_distance := math.Sqrt(math.Pow((endPointX-target.X), 2) + math.Pow((endPointY-target.Y), 2))
	backward_distance := math.Sqrt(math.Pow((endPointX_back-target.X), 2) + math.Pow((endPointY_back-target.Y), 2))

	fmt.Printf("___________________ \n")
	fmt.Printf("ANGLE = %d \n", startPoint.ShootAngle)
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

	fmt.Printf("start %d : %d \n", startPoint.X, startPoint.Y)
	fmt.Printf("end   %d : %d \n", endPointX, endPointY)
	fmt.Printf("distanse = %d \n", distance)
	distance = (math.Abs(distance) - objectRadius)
	fmt.Printf("distanse = %d \n", distance)
	return distance
}

type GameClient struct {
	net.Conn
	Pos       Position
	Id        string
	CoolDown  int
	canShoot  bool
	HealPoint int64
	isAttack  bool
}

type gameCore struct {
	Clients map[string]*GameClient
}

func (g gameCore) AddClient(conn net.Conn, Id string) {
	var newClient GameClient
	newClient.Id = Id
	newClient.Conn = conn
	newClient.CoolDown = 0
	newClient.isAttack = false
	newClient.canShoot = true
	newClient.HealPoint = startHealfPoint
	g.Clients[newClient.Id] = &newClient
	println("new client add!, now count = ", len(g.Clients))
	newmessage := "Init;" + Id + ";"
	newClient.Conn.Write([]byte(newmessage))
}

func (g gameCore) ReadClientsData() {
	for id, _ := range g.Clients {
		// read message split by \n
		message, err := bufio.NewReader(g.Clients[id].Conn).ReadString('\n')
		if err == nil {
			//log.Printf("error accepting connection %v", err)
			//println("read client data: ", message)

			strArr := strings.Split(message, ";")
			if len(strArr) < 3 {
				continue
			}

			var dX, dY, nAngl float64
			if strings.Contains(strArr[0], "XD") {
				if strArr[1] == "X" {
					//println("read client data: ", message)
					g.Clients[id].isAttack = false
					continue
				}
				g.Clients[id].isAttack = true
				nAngl, _ = strconv.ParseFloat(strArr[1], 32)
				g.Clients[id].Pos.ShootAngle = nAngl
				continue
			}

			if len(strArr) > 3 {
				dX, _ = strconv.ParseFloat(strArr[1], 32)
				dY, _ = strconv.ParseFloat(strArr[2], 32)
				nAngl, _ = strconv.ParseFloat(strArr[3], 32)

				g.Clients[id].Pos.TryDeltaX = dX
				g.Clients[id].Pos.TryDeltaY = dY
				g.Clients[id].Pos.Angle = nAngl
			}
		} else {
			//println("error", err)
		}
	}
}
func (g gameCore) WriteClientsData() {

	var newmessage string

	//generate message
	for id, _ := range g.Clients {
		//generate players position
		g.Clients[id].Pos.X = g.Clients[id].Pos.X + g.Clients[id].Pos.TryDeltaX*moveSpeed
		g.Clients[id].Pos.Y = g.Clients[id].Pos.Y + g.Clients[id].Pos.TryDeltaY*moveSpeed

		//generate bullets position
		//--cool dawn timers
		g.Clients[id].CoolDown++
		if g.Clients[id].CoolDown >= shootRate {
			g.Clients[id].CoolDown = 0
			g.Clients[id].canShoot = true
		} else {
			g.Clients[id].canShoot = false
		}

		//check bullets collisions
		for id_2, _ := range g.Clients {
			if g.Clients[id].canShoot && g.Clients[id].isAttack && (id_2 != id) {
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
		}

		// write data to string
		// write ------- id;x_pos;y_pos;angle;HP;
		newmessage += id + ";" + strconv.FormatFloat(g.Clients[id].Pos.X, 'f', 1, 64) +
			";" + strconv.FormatFloat(g.Clients[id].Pos.Y, 'f', 1, 64) + ";" +
			strconv.FormatFloat(g.Clients[id].Pos.Angle, 'f', 0, 64) + ";" +
			strconv.FormatInt(g.Clients[id].HealPoint, 10) + ";"

		// write bullet info x_pos;y_pos; or X;Y; if not fly
		if g.Clients[id].isAttack {
			newmessage += "1;"
		} else {
			newmessage += "0;"
		}
		newmessage += "\n"
	}
	//println("new msg = ", newmessage)

	//send message
	for id, _ := range g.Clients {
		g.Clients[id].Conn.Write([]byte(newmessage))
		//println("new msg = ", newmessage)
	}
}

func GetInstance() *gameCore {
	once.Do(func() {
		instance = &gameCore{}
		instance.Clients = make(map[string]*GameClient)
	})
	return instance
}

func init() {
	fmt.Println("Create gameCore ")
}
