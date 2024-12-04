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
