package resolv

import (
	"fmt"
	"math"
)

// A Circle represents an ordinary circle, and has a radius, in addition to normal shape properties.
type Circle struct {
	BasicShape
	Radius int32
}

// NewCircle returns a pointer to a new Circle object.
func NewCircle(x, y, radius int32) *Circle {
	c := &Circle{Radius: radius}
	c.X = x
	c.Y = y
	return c
}

// IsColliding returns true if the Circle is colliding with the specified other Shape, including the other Shape
// being wholly within the Circle.
func (c *Circle) IsColliding(other Shape) bool {

	switch b := other.(type) {

	case *Circle:
		return Distance(c.X, c.Y, b.X, b.Y) <= c.Radius+b.Radius
	case *Rectangle:
		closestX := c.X
		closestY := c.Y

		if c.X < b.X {
			closestX = b.X
		} else if c.X > b.X+b.W {
			closestX = b.X + b.W
		}

		if c.Y < b.Y {
			closestY = b.Y
		} else if c.Y > b.Y+b.H {
			closestY = b.Y + b.H
		}

		return Distance(c.X, c.Y, closestX, closestY) <= c.Radius
	case *Line:
		//return b.IsColliding(c)
		// 通过该线段与圆心作三角形，判断线段与圆是否有焦点或在圆内
		return c.isCollidingWithLine(b)
	case *Space:
		return b.IsColliding(c)

	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return false

}

// WouldBeColliding returns whether the Circle would be colliding with the specified other Shape if it were to move
// in the specified direction.
func (c *Circle) WouldBeColliding(other Shape, dx, dy int32) bool {
	c.X += dx
	c.Y += dy
	isColliding := c.IsColliding(other)
	c.X -= dx
	c.Y -= dy
	return isColliding
}

// GetBoundingRect returns a Rectangle which has a width and height of 2*Radius.
func (c *Circle) GetBoundingRect() *Rectangle {
	r := &Rectangle{}
	r.W = c.Radius * 2
	r.H = c.Radius * 2
	r.X = c.X - c.Radius
	r.Y = c.Y - c.Radius
	return r
}

func (c *Circle) isCollidingWithLine(l *Line) bool {
	AC := float64(Distance(c.X, c.Y, l.X, l.Y))
	CB := float64(Distance(c.X, c.Y, l.X2, l.Y2))
	BA := float64(l.GetLength())

	// 线段两点到圆心的距离小于圆半径则一定碰撞
	if AC <= float64(c.Radius) || CB <= float64(c.Radius) {
		return true
	}

	// 线段与圆心作三角形，线段为底边，求其高
	p := (AC + CB + BA) / 2
	h := 2 * math.Sqrt(p*(p-AC)*(p-CB)*(p-BA)) / BA

	// 高若大于圆半径则不然不碰撞
	if h <= float64(c.Radius) {
		cosC := (AC*AC + CB*CB - BA*BA) / (2 * AC * CB)
		primaryLine := func() float64 {
			if AC < CB {
				return CB
			}
			return AC
		}()

		newBA := math.Sqrt(primaryLine*primaryLine + float64(c.Radius*c.Radius) - 2*primaryLine*float64(c.Radius)*cosC)
		// 如果原底边小于以半径做临边的底边，则不碰撞
		if BA >= newBA {
			return true
		}

	}

	return false
}
