package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Graph    *Graph
	Vehicles []*Vehicle
	Step     int
}

func (g *Game) Update() error {
	if g.Step%30 == 0 {
		for _, vehicle := range g.Vehicles {
			if vehicle.Status != "arrived" {
				if vehicle.Position < len(vehicle.Path)-1 {
					vehicle.Position++
					vehicle.Status = "moving"
				} else {
					vehicle.Status = "arrived"
				}
			}
		}
	}
	g.Step++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	// center the grid
	offsetX := 50.0
	offsetY := 50.0
	scale := 80.0

	// roads
	for _, link := range g.Graph.Links {
		x1 := float64(link.FromNode.X)*scale + offsetX
		y1 := float64(link.FromNode.Y)*scale + offsetY
		x2 := float64(link.ToNode.X)*scale + offsetX
		y2 := float64(link.ToNode.Y)*scale + offsetY
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 2.0, color.Black, false)
	}

	// intersections
	for _, node := range g.Graph.Nodes {
		x := float64(node.X)*scale + offsetX
		y := float64(node.Y)*scale + offsetY
		vector.DrawFilledRect(screen, float32(x-5), float32(y-5), 10, 10, color.RGBA{0, 0, 255, 255}, false)
	}

	// vehicles
	for _, vehicle := range g.Vehicles {
		if vehicle.Status != "arrived" {
			currentNode := vehicle.Path[vehicle.Position]
			x := float64(currentNode.X)*scale + offsetX
			y := float64(currentNode.Y)*scale + offsetY
			// Use red color for vehicles
			vector.DrawFilledRect(screen, float32(x-5), float32(y-5), 10, 10, color.RGBA{255, 0, 0, 255}, false)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}
