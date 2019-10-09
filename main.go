package main

import "math"

func main() {
	
}

type Rectangle struct {
	Width float64
	Height float64
}

func (r Rectangle) Area() float64  {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64  {
	return math.Pi * c.Radius * c.Radius
}

func Perimeter(width, height float64) float64 {
	return 2 * (width + height)
}


type Shape interface {
	Area() float64
}