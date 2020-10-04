package main

import (
	"fmt"
	"log"
	"math"

	//"./collision2d"
	"github.com/mamkin-itshnik/collision2d"
)

//---------------------------------OBJECT IN LEVEL
//lines
var topLine float64
var downLine float64
var leftLine float64
var rightLine float64

//circles
var circleArray []collision2d.Circle

func init() {

	fmt.Println("Create collision detector ")

	//---------------------------------------- Create objects
	//lines
	topLine = 34.5
	downLine = -34.0
	leftLine = -39.0
	rightLine = 33.0
	// circles
	circleArray = append(circleArray, collision2d.Circle{collision2d.Vector{0, 0}, 2})

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

func c_checkCollisionInLines(point collision2d.Vector) (bool, collision2d.Vector) {

	// TOP
	fmt.Println("ckeck collision on TOP ", point.Y)
	topCollision := false
	downCollision := false
	leftCollision := false
	rightCollision := false
	newVector := collision2d.NewVector(point.X, point.Y)
	if point.Y > topLine {
		newVector.Y = topLine
		topCollision = true
	}
	// DOWN
	if point.Y < downLine {
		newVector.Y = downLine
		downCollision = true
	}
	// LEFT
	if point.X < leftLine {
		newVector.X = leftLine
		leftCollision = true
	}
	// RIGHT
	if point.X > rightLine {
		newVector.X = rightLine
		rightCollision = true
	}
	if topCollision || downCollision || leftCollision || rightCollision {
		return true, newVector
	}
	return false, collision2d.NewVector(0, 0)
}

// Return  nearest point and bool isCollision
func checkCollision(x, y float64) (bool, collision2d.Vector) {

	// check all circles
	//fmt.Println("ckeck collision")
	//fmt.Println("ckeck collision")
	isCollisionCircle, pointCircle := c_checkCollisionInCircles(collision2d.NewVector(x, y))
	if isCollisionCircle {
		return isCollisionCircle, pointCircle
	}
	//	fmt.Println("ckeck collision on line")
	isCollisionLine, pointOnLine := c_checkCollisionInLines(collision2d.NewVector(x, y))
	if isCollisionLine {
		return isCollisionLine, pointOnLine
	}

	return false, collision2d.NewVector(0, 0)
}
