package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	NUM_VEHICLES          = 1000
	ACCIDENT_PROBABILITY  = 0.001
	ROAD_CAPACITY         = 10
	INTERSECTION_CAPACITY = 5
	GREEN_DURATION        = 60
	YELLOW_DURATION       = 30
	RED_DURATION          = 60
)

func main() {
	startTime := time.Now()
	rand.New(rand.NewSource(0))

	nodes, err := readNodes("nodes.csv")
	if err != nil {
		fmt.Println("Error reading nodes:", err)
		return
	}

	links, err := readLinks("links.csv")
	if err != nil {
		fmt.Println("Error reading links:", err)
		return
	}

	graph := &Graph{
		Nodes: nodes,
		Links: links,
	}

	// Initialize traffic signals for specific intersections
	for _, node := range graph.Nodes {
		initialStates := []string{"red", "green", "yellow"}
		initialState := initialStates[rand.Intn(len(initialStates))]
		initialDuration := rand.Intn(20) + 10

		node.Signal = &TrafficSignal{
			State:    initialState,
			Duration: initialDuration,
		}
	}

	// Associate links with nodes
	for _, link := range graph.Links {
		fromNode, ok := graph.Nodes[link.FromNodeID]
		if !ok {
			log.Fatalf("FromNodeID %s not found in nodes", link.FromNodeID)
		}
		toNode, ok := graph.Nodes[link.ToNodeID]
		if !ok {
			log.Fatalf("ToNodeID %s not found in nodes", link.ToNodeID)
		}
		link.FromNode = fromNode
		link.ToNode = toNode
	}

	vehicles := []*Vehicle{}

	for i := 0; i < NUM_VEHICLES; i++ {
		startNodeID := randomNodeID(graph.Nodes)
		endNodeID := randomNodeID(graph.Nodes)

		// Ensure start and end nodes are different
		for startNodeID == endNodeID {
			endNodeID = randomNodeID(graph.Nodes)
		}

		path, err := findPath(graph, startNodeID, endNodeID)
		if err != nil {
			log.Printf("Error finding path for Vehicle %d: %v", i, err)
			continue
		}

		vehicles = append(vehicles, &Vehicle{
			ID:       fmt.Sprintf("V%d", i+1),
			Path:     path,
			Position: 0,
			Status:   "waiting",
		})
	}

	// game loop
	game := &Game{
		Graph:    graph,
		Vehicles: vehicles,
		Step:     0,
		Running:  true,
	}
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Traffic Simulation")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	densities := calculateTrafficDensity(graph, game.Step)
	densityData, err := json.MarshalIndent(densities, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("go_traffic_density_data.json", densityData, 0644)
	if err != nil {
		fmt.Println("Error writing JSON file:", err)
		return
	}
	fmt.Println("done writing json data")

	endTime := time.Now()
	executionTime := endTime.Sub(startTime).Seconds()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	currentMemory := m.Alloc / 1024 / 1024
	peakMemory := m.TotalAlloc / 1024 / 1024

	fileExists := true
	if _, err := os.Stat("go_simulation_stats.csv"); os.IsNotExist(err) {
		fileExists = false
	}

	file, err := os.OpenFile("go_simulation_stats.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if !fileExists {
		writer.Write([]string{"Language", "NUM_VEHICLES", "ACCIDENT_PROBABILITY", "ROAD_CAPACITY", "INTERSECTION_CAPACITY", "GREEN_DURATION", "YELLOW_DURATION", "RED_DURATION", "Execution Time (s)", "Current Memory (MB)", "Peak Memory (MB)"})
	}
	writer.Write([]string{
		"Go",
		fmt.Sprintf("%d", NUM_VEHICLES),
		fmt.Sprintf("%f", ACCIDENT_PROBABILITY),
		fmt.Sprintf("%d", ROAD_CAPACITY),
		fmt.Sprintf("%d", INTERSECTION_CAPACITY),
		fmt.Sprintf("%d", GREEN_DURATION),
		fmt.Sprintf("%d", YELLOW_DURATION),
		fmt.Sprintf("%d", RED_DURATION),
		fmt.Sprintf("%f", executionTime),
		fmt.Sprintf("%d", currentMemory),
		fmt.Sprintf("%d", peakMemory),
	})
}
