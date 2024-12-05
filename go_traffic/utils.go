package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
)

var xCoords = map[string]int{
	"N Wacker": 0,
	"Franklin": 1,
	"Wells":    2,
	"Lasalle":  3,
	"Clark":    4,
	"Dearborn": 5,
	"State":    6,
	"Wabash":   7,
	"Michigan": 8,
}

var yCoords = map[string]int{
	"Jackson":    0,
	"Adams":      1,
	"Monroe":     2,
	"Madison":    3,
	"Washington": 4,
	"Randolph":   5,
	"Lake":       6,
	"W Wacker":   7,
}

func readNodes(filename string) (map[string]*Intersection, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	nodes := make(map[string]*Intersection)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		nodeID := record[0]
		name := record[1]

		// place street names on the grid
		streets := strings.Split(name, " / ")
		xStreet := streets[0]
		yStreet := streets[1]

		x, xOk := xCoords[xStreet]
		y, yOk := yCoords[yStreet]

		if !xOk || !yOk {
			return nil, fmt.Errorf("unknown street names in node: %s", name)
		}

		nodes[nodeID] = &Intersection{
			ID:   nodeID,
			Name: name,
			X:    x,
			Y:    y,
		}
	}
	return nodes, nil
}

func readEdges(filename string) ([]*Road, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var edges []*Road

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		edge := &Road{
			ID:           record[0],
			FromNodeName: record[1],
			ToNodeName:   record[2],
			FromNodeID:   record[3],
			ToNodeID:     record[4],
		}
		edges = append(edges, edge)
	}
	return edges, nil
}

func randomNodeID(nodes map[string]*Intersection) string {
	keys := make([]string, 0, len(nodes))
	for k := range nodes {
		keys = append(keys, k)
	}
	return keys[rand.Intn(len(keys))]
}

func checkForAccidents(graph *Graph, vehicles []*Vehicle) {
	for _, link := range graph.Links {
		link.VehiclesOnRoad = nil
	}

	for _, vehicle := range vehicles {
		if vehicle.Status != "arrived" && vehicle.Position < len(vehicle.Path)-1 {
			currentNode := vehicle.Path[vehicle.Position]
			nextNode := vehicle.Path[vehicle.Position+1]

			for _, link := range graph.Links {
				if link.FromNode == currentNode && link.ToNode == nextNode {
					link.VehiclesOnRoad = append(link.VehiclesOnRoad, vehicle)
					break
				}
			}
		}
	}

	for _, link := range graph.Links {
		if link.Accident == nil {
			if len(link.VehiclesOnRoad) >= 2 {
				if rand.Float64() < ACCIDENT_PROBABILITY {
					accidentPosition := rand.Float64()
					accidentDuration := rand.Intn(300) + 300
					link.Accident = &Accident{Road: link, Position: accidentPosition, Duration: accidentDuration}
				}
			}
		}
	}
}

func updateAccidents(graph *Graph) {
	for _, link := range graph.Links {
		if link.Accident != nil {
			link.Accident.ElapsedTime++
			if link.Accident.ElapsedTime >= link.Accident.Duration {
				link.Accident = nil
			}
		}
	}
}
