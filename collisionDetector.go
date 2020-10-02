package main

import (
	"fmt"
	"log"
	"math"

	//"./collision2d"
	//"github.com/Tarliton/collision2d""
	"github.com/mamkin-itshnik/collision2d"
)

//---------------------------------OBJECT IN LEVEL
//circles
var circleArray []collision2d.Circle

//boxes
//..........

func init() {

	fmt.Println("Create collision detector ")
	//memory allocate
	//circleArray := make([]collision2d.Circle, 5)
	fmt.Println("circles count = ", len(circleArray))

	//---------------------------------------- Create objects
	// circles
	//circleArray[0] = collision2d.NewCircle(collision2d.NewVector(0, 0), 10)
	circleArray = append(circleArray, collision2d.Circle{collision2d.Vector{0, 0}, 2})
	fmt.Println("circles count = ", len(circleArray))
}

func c_checkCollisionInCircles(point collision2d.Vector) (bool, collision2d.Vector) {

	//fmt.Println("ckeck collision in circles. circles count = ", len(circleArray))
	log.Println("TRY detect  collision ", point.X, " ", point.Y)
	for _, circle := range circleArray {
		if collision2d.PointInCircle(point, circle) {
			// TODO - make nearest point
			dX2 := circle.Pos.X - point.X
			dY2 := circle.Pos.Y - point.Y
			k := circle.R / math.Sqrt(math.Pow(dX2, 2)+math.Pow(dY2, 2))
			newX := circle.Pos.X - dX2*k
			newY := circle.Pos.Y - dY2*k
			return true, collision2d.NewVector(newX, newY)
		}
	}
	return false, collision2d.NewVector(0, 0)
}

// Return  nearest point and bool isCollision
func checkCollision(x, y float64) (bool, collision2d.Vector) {

	// check all circles
	//fmt.Println("ckeck collision")
	isCollision, point := c_checkCollisionInCircles(collision2d.NewVector(x, y))
	if isCollision {
		return isCollision, point
	}

	return false, collision2d.NewVector(0, 0)
}
