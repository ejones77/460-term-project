package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
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

	// Example vehicles
	vehicles := []*Vehicle{}

	path1, err := findPath(graph, "7", "71")
	if err != nil {
		log.Fatalf("Error finding path for Vehicle 1: %v", err)
	}
	vehicles = append(vehicles, &Vehicle{
		ID:   "V1",
		Path: path1,
	})

	path2, err := findPath(graph, "9", "52")
	if err != nil {
		log.Fatalf("Error finding path for Vehicle 2: %v", err)
	}
	vehicles = append(vehicles, &Vehicle{
		ID:   "V2",
		Path: path2,
	})

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
