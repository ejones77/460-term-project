package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

const ACCIDENT_PROBABILITY = 0.0001
const NUM_VEHICLES = 1000

func main() {
	rand.New(rand.NewSource(0))

	nodes, err := readNodes("nodes.csv")
	if err != nil {
		fmt.Println("Error reading nodes:", err)
		return
	}

	links, err := readEdges("links.csv")
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
	}
	ebiten.SetWindowSize(600, 600)
	ebiten.SetWindowTitle("Traffic Simulation")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
