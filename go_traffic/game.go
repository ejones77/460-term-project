package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Graph          *Graph
	Vehicles       []*Vehicle
	Step           int
	Running        bool
	TrafficDensity []map[string]interface{}
}

func (g *Game) Update() error {
	if !g.Running {
		return nil // Stop updating if the game is not running
	}

	// Update traffic signals
	for _, node := range g.Graph.Nodes {
		if node.Signal != nil {
			node.Signal.ElapsedTime++
			if node.Signal.ElapsedTime >= node.Signal.Duration {
				switch node.Signal.State {
				case "red":
					node.Signal.State = "green"
					node.Signal.Duration = GREEN_DURATION
				case "green":
					node.Signal.State = "yellow"
					node.Signal.Duration = YELLOW_DURATION
				case "yellow":
					node.Signal.State = "red"
					node.Signal.Duration = RED_DURATION
				}
				node.Signal.ElapsedTime = 0
			}
		}
	}

	// Update vehicles
	allArrived := true
	for _, vehicle := range g.Vehicles {
		if vehicle.Status != "arrived" {
			allArrived = false
			currentNode := vehicle.Path[vehicle.Position]
			if vehicle.Position < len(vehicle.Path)-1 {
				nextNode := vehicle.Path[vehicle.Position+1]

				// Check for accidents
				var currentRoad *Road
				for _, link := range g.Graph.Links {
					if link.FromNode == currentNode && link.ToNode == nextNode {
						currentRoad = link
						break
					}
				}

				if currentRoad != nil && currentRoad.Accident != nil {
					if vehicle.Progress < currentRoad.Accident.Position {
						vehicle.Progress += 0.01
						vehicle.Status = "moving"
						if vehicle.Progress >= currentRoad.Accident.Position {
							vehicle.Status = "waiting"
						}
					} else {
						vehicle.Status = "waiting"
					}
					continue
				}
			}

			if currentNode.Signal == nil || currentNode.Signal.State == "green" {
				vehicle.Progress += 0.01
				vehicle.Status = "moving"
				if vehicle.Progress >= 1.0 {
					vehicle.Progress = 0.0
					vehicle.Position++
					if vehicle.Position >= len(vehicle.Path)-1 {
						vehicle.Status = "arrived"
					}
				}
			} else {
				vehicle.Status = "waiting"
			}
		}
	}

	if allArrived {
		fmt.Println("All vehicles have arrived. Ending simulation.")
		g.Running = false // Set the flag to stop the game loop
	}

	g.Step++

	checkForAccidents(g.Graph, g.Vehicles)
	updateAccidents(g.Graph)
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

		// Draw two lines for lanes in each direction
		vector.StrokeLine(screen, float32(x1), float32(y1), float32(x2), float32(y2), 5.0, color.Gray{Y: 128}, false)         // Lane 1
		vector.StrokeLine(screen, float32(x1+5), float32(y1+5), float32(x2+5), float32(y2+5), 5.0, color.Gray{Y: 192}, false) // Lane 2
	}

	// intersections
	for _, node := range g.Graph.Nodes {
		x := float64(node.X)*scale + offsetX
		y := float64(node.Y)*scale + offsetY
		var signalColor color.RGBA

		if node.Signal != nil {
			switch node.Signal.State {
			case "red":
				signalColor = color.RGBA{255, 0, 0, 255}
			case "green":
				signalColor = color.RGBA{0, 255, 0, 255}
			case "yellow":
				signalColor = color.RGBA{255, 255, 0, 255}
			default:
				signalColor = color.RGBA{0, 0, 255, 255} // Default blue if state is unknown
			}
		} else {
			signalColor = color.RGBA{0, 0, 255, 255} // Default blue for intersections without signals
		}

		vector.DrawFilledRect(screen, float32(x-10), float32(y-10), 20, 20, signalColor, false) // Larger shape for intersections
	}

	// vehicles
	for _, vehicle := range g.Vehicles {
		if vehicle.Status != "arrived" {
			currentNode := vehicle.Path[vehicle.Position]
			nextNode := vehicle.Path[vehicle.Position+1]

			x := float64(currentNode.X)*(1-vehicle.Progress) + float64(nextNode.X)*vehicle.Progress
			y := float64(currentNode.Y)*(1-vehicle.Progress) + float64(nextNode.Y)*vehicle.Progress

			x = x*scale + offsetX
			y = y*scale + offsetY

			// blue when moving, orange when waiting
			vehicleColor := color.RGBA{0, 0, 255, 255}
			if vehicle.Status == "waiting" {
				vehicleColor = color.RGBA{255, 165, 0, 255}
			}

			vector.DrawFilledRect(screen, float32(x-5), float32(y-5), 10, 10, vehicleColor, false) // Larger squares for vehicles
		}
	}
	// accidents
	for _, link := range g.Graph.Links {
		if link.Accident != nil {
			x1 := float64(link.FromNode.X)*scale + offsetX
			y1 := float64(link.FromNode.Y)*scale + offsetY
			x2 := float64(link.ToNode.X)*scale + offsetX
			y2 := float64(link.ToNode.Y)*scale + offsetY

			accidentX := x1 + (x2-x1)*link.Accident.Position
			accidentY := y1 + (y2-y1)*link.Accident.Position

			vector.DrawFilledCircle(screen, float32(accidentX), float32(accidentY), 12, color.RGBA{255, 0, 0, 255}, false)
			vector.DrawFilledCircle(screen, float32(accidentX), float32(accidentY), 8, color.RGBA{255, 255, 0, 255}, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 800
}
