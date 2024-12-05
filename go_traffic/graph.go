package main

import (
	"fmt"
	"math"
)

type TrafficSignal struct {
	State       string
	Duration    int
	ElapsedTime int
}

type Intersection struct {
	ID       string
	Name     string
	X        int
	Y        int
	Signal   *TrafficSignal
	Capacity int
	Queue    []*Vehicle
}

type Road struct {
	ID             string
	FromNodeID     string
	ToNodeID       string
	FromNode       *Intersection
	ToNode         *Intersection
	FromNodeName   string
	ToNodeName     string
	Accident       *Accident
	VehiclesOnRoad []*Vehicle
}

type Accident struct {
	Road        *Road
	Position    float64
	Duration    int
	ElapsedTime int
}

type Vehicle struct {
	ID       string
	Path     []*Intersection
	Position int
	Status   string
	Progress float64
}

type Graph struct {
	Nodes map[string]*Intersection
	Links []*Road
}

func findPath(graph *Graph, startID, endID string) ([]*Intersection, error) {
	startNode := graph.Nodes[startID]
	endNode := graph.Nodes[endID]
	distances := make(map[*Intersection]float64)
	previous := make(map[*Intersection]*Intersection)
	unvisited := make(map[*Intersection]bool)

	for node := range graph.Nodes {
		distances[graph.Nodes[node]] = math.Inf(1)
		unvisited[graph.Nodes[node]] = true
	}
	distances[startNode] = 0

	// Find unvisited node with min distance
	for len(unvisited) > 0 {
		var current *Intersection
		minDist := math.Inf(1)
		for node := range unvisited {
			if distances[node] < minDist {
				minDist = distances[node]
				current = node
			}
		}

		if current == nil {
			break
		}

		if current == endNode {
			break
		}

		delete(unvisited, current)

		// Check all neighbors of current node
		for _, link := range graph.Links {
			if link.FromNode == current {
				neighbor := link.ToNode
				distance := distances[current] + 1

				if distance < distances[neighbor] {
					distances[neighbor] = distance
					previous[neighbor] = current
				}
			}
		}
	}

	var path []*Intersection
	current := endNode
	for current != nil {
		path = append(path, current)
		current = previous[current]
	}

	// Check if path exists
	if len(path) == 0 || path[len(path)-1] != startNode {
		return nil, fmt.Errorf("path not found between %s and %s", startID, endID)
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path, nil
}
