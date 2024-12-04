package main

import "fmt"

type Intersection struct {
	ID   string
	Name string
	X    int
	Y    int
}

type Road struct {
	ID           string
	FromNodeID   string
	ToNodeID     string
	FromNode     *Intersection
	ToNode       *Intersection
	FromNodeName string
	ToNodeName   string
}

type Vehicle struct {
	ID       string
	Path     []*Intersection
	Position int
	Status   string
}

type Graph struct {
	Nodes map[string]*Intersection
	Links []*Road
}

func findPath(graph *Graph, startID, endID string) ([]*Intersection, error) {
	startNode, ok := graph.Nodes[startID]
	if !ok {
		return nil, fmt.Errorf("start node %s not found", startID)
	}
	endNode, ok := graph.Nodes[endID]
	if !ok {
		return nil, fmt.Errorf("end node %s not found", endID)
	}

	// Since it's a grid, we can move in the X and Y directions along the route
	path := []*Intersection{startNode}

	current := startNode
	for current.ID != endNode.ID {
		nextX := current.X
		nextY := current.Y

		if current.X < endNode.X {
			nextX++
		} else if current.X > endNode.X {
			nextX--
		} else if current.Y < endNode.Y {
			nextY++
		} else if current.Y > endNode.Y {
			nextY--
		}

		// Find the next node at (nextX, nextY)
		nextNode := findNode(graph.Nodes, nextX, nextY)
		if nextNode == nil {
			return nil, fmt.Errorf("path not found between %s and %s", startID, endID)
		}

		path = append(path, nextNode)
		current = nextNode
	}
	return path, nil
}

func findNode(nodes map[string]*Intersection, x, y int) *Intersection {
	for _, node := range nodes {
		if node.X == x && node.Y == y {
			return node
		}
	}
	return nil
}
