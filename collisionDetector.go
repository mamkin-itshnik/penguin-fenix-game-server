package main

import (
	"fmt"
	"log"
	"math"

	"./collision2d"
	//"github.com/mamkin-itshnik/collision2d"
)

//---------------------------------OBJECT IN LEVEL
//circles
var circleArray []collision2d.Circle

//boxes
var boxsArray []collision2d.Box

func init() {

	fmt.Println("Create collision detector ")

	//---------------------------------------- Create objects
	// circles
	circleArray = append(circleArray, collision2d.Circle{collision2d.Vector{0, 0}, 2})

	//boxes
	boxsArray = append(boxsArray, collision2d.Box{collision2d.Vector{0, 36}, 100, 1})

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

func c_checkCollisionInBoxes(point collision2d.Vector) (bool, collision2d.Vector) {

	for _, box := range boxsArray {
		if collision2d.PointInPolygon(point, box.ToPolygon()) {
			// TODO - make nearest point

			for _, qwe := range box.ToPolygon().Edges {

			}

			newX := 0.0
			newY := 0.0
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
