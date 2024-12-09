# 460-term-project

## Overview
This repository offers two side-by-side versions of a simple traffic simulation. One built in Python, and the other built in Go. 

The simulation leverages the similarities between Python's Pygame engine, and Go's ebiten package. 

The Go code is split between a few files
`game.go` -- handles the game engine updates, window layout, and overall visualization of entity behaviors
`graph.go` -- defines the network architecture and traversal through said network
`utils.go` -- provides helper functions to read csv's alongside defined constants
`main.go` -- brings everything together to run a configurable simulation with data exports appended to `go_simulation_stats.csv`

The Python code is contained to `py_traffic.py` -- which primarily serves a translation of the initial Go code. This will also write results to `python_simulation_stats.csv`

Results were compiled and analyzed -- available in `Combined results.xlsx`

## Setup -- Python

in the respository directory, create and activate a virtual environment (assuming powershell)

 py -m venv .venv
./.venv/scripts/activate
install dependencies with

py -m pip install --upgrade pip
py -m pip install -r requirements.txt

## Setup -- Go

- Ensure Go is installed on your machine
- in the `go_traffic` subdirectory of this repository, execute the following command
```
go run .
```
- Alternatively you can build an executable program that will run the simulation with
```
go build .
```
