import random
import simpy
import pygame
import pandas as pd
import numpy as np
from collections import deque

"""
Baseline Initializations
"""

NUM_VEHICLES = 1000
ACCIDENT_PROBABILITY = 0.0001

class TrafficSignal:
    def __init__(self, state, duration):
        self.state = state
        self.duration = duration
        self.elapsed_time = 0

class Intersection:
    def __init__(self, id, name, x, y, signal=None):
        self.id = id
        self.name = name
        self.x = x
        self.y = y
        self.signal = signal
        self.capacity = 5
        self.queue = []

class Road:
    def __init__(self, id, from_node, to_node):
        self.id = id
        self.from_node = from_node
        self.to_node = to_node
        self.accident = None
        self.vehicles_on_road = []

class Vehicle:
    def __init__(self, id, path, graph):
        self.id = id
        self.path = path
        self.graph = graph
        self.position = 0
        self.status = 'waiting'
        self.progress = 0.0

class Accident:
    def __init__(self, road, position, duration):
        self.road = road
        self.position = position
        self.duration = duration
        self.elapsed_time = 0

class Graph:
    def __init__(self, nodes, links):
        self.nodes = nodes
        self.links = links

x_coords = {
    "N Wacker": 0,
    "Franklin": 1,
    "Wells": 2,
    "Lasalle": 3,
    "Clark": 4,
    "Dearborn": 5,
    "State": 6,
    "Wabash": 7,
    "Michigan": 8,
}

y_coords = {
    "Jackson": 0,
    "Adams": 1,
    "Monroe": 2,
    "Madison": 3,
    "Washington": 4,
    "Randolph": 5,
    "Lake": 6,
    "W Wacker": 7,
}

def read_nodes(filename):
    nodes_df = pd.read_csv(filename)
    nodes = {}
    for _, row in nodes_df.iterrows():
        node_id = row['Node ID']
        name = row['Name']
        # plots each node on the grid
        x_street, y_street = name.split(' / ')
        x = x_coords.get(x_street)
        y = y_coords.get(y_street)
        nodes[node_id] = Intersection(node_id, name, x, y)
    return nodes

def read_edges(filename, nodes):
    links_df = pd.read_csv(filename)
    links = []
    for _, row in links_df.iterrows():
        from_node = nodes.get(row['From_Node_ID'])
        to_node = nodes.get(row['To_Node_ID'])
        links.append(Road(row['Link_ID'], from_node, to_node))
    return links

"""
Graph navigation
"""

def find_path(graph, start_id, end_id):
    """
    vehicles use dijkstra to travel paths
    """
    start_node = graph.nodes.get(start_id)
    end_node = graph.nodes.get(end_id)
    distances = {node: np.inf for node in graph.nodes.values()}
    distances[start_node] = 0
    previous = {node: None for node in graph.nodes.values()}
    
    unvisited = set(graph.nodes.values())
    while unvisited:
        # Find unvisited node with min distance
        current = min(unvisited, key=lambda node: distances[node])
        
        if current == end_node:
            break
            
        unvisited.remove(current)

        # Check all neighbors of current node
        for link in graph.links:
            if link.from_node == current:
                neighbor = link.to_node
                distance = distances[current] + 1
                
                if distance < distances[neighbor]:
                    distances[neighbor] = distance
                    previous[neighbor] = current

    # reconstruct path
    path = []
    current = end_node
    while current is not None:
        path.append(current)
        current = previous[current]
    
    return path[::-1]

def check_for_accidents(graph, vehicles):
    
    for link in graph.links:
        link.vehicles_on_road = []
        
    for vehicle in vehicles:
        if vehicle.status != 'arrived' and vehicle.position < len(vehicle.path) - 1:
            current_node = vehicle.path[vehicle.position]
            next_node = vehicle.path[vehicle.position + 1]
            
            for link in graph.links:
                if link.from_node == current_node and link.to_node == next_node:
                    link.vehicles_on_road.append(vehicle)
                    break
    
    for link in graph.links:
        if link.accident is None:  
            if len(link.vehicles_on_road) >= 2:
                if random.random() < ACCIDENT_PROBABILITY:
                    
                    accident_position = random.random()  
                    accident_duration = random.randint(300, 600)
                    link.accident = Accident(link, accident_position, accident_duration)

"""
Graph updates
"""


def update_traffic_signals(graph):
    for node in graph.nodes.values():
        if node.signal:
            node.signal.elapsed_time += 1
            if node.signal.elapsed_time >= node.signal.duration:
                if node.signal.state == 'red':
                    node.signal.state = 'green'
                    node.signal.duration = 60
                elif node.signal.state == 'green':
                    node.signal.state = 'yellow'
                    node.signal.duration = 30
                elif node.signal.state == 'yellow':
                    node.signal.state = 'red'
                    node.signal.duration = 60
                node.signal.elapsed_time = 0

def update_vehicles(vehicles):
    for vehicle in vehicles:
        if vehicle.status != 'arrived':
            current_node = vehicle.path[vehicle.position]
            if vehicle.position < len(vehicle.path) - 1:
                next_node = vehicle.path[vehicle.position + 1]
                
                # Check if there's an accident
                current_road = None
                for link in vehicle.graph.links:
                    if link.from_node == current_node and link.to_node == next_node:
                        current_road = link
                        break
                
                if current_road and current_road.accident:
                    # wait upon reaching an accident
                    if vehicle.progress < current_road.accident.position:
                        vehicle.progress += 0.01
                        vehicle.status = 'moving'
                        if vehicle.progress >= current_road.accident.position:
                            vehicle.status = 'waiting'
                    else:
                        vehicle.status = 'waiting'
                    continue
                
            if not current_node.signal or current_node.signal.state == 'green':
                vehicle.progress += 0.01
                vehicle.status = 'moving'
                if vehicle.progress >= 1.0:
                    vehicle.progress = 0.0
                    vehicle.position += 1
                    if vehicle.position >= len(vehicle.path) - 1:
                        vehicle.status = 'arrived'
            else:
                vehicle.status = 'waiting'

def update_accidents(graph):
    for link in graph.links:
        if link.accident:
            link.accident.elapsed_time += 1
            if link.accident.elapsed_time >= link.accident.duration:
                link.accident = None

"""
GUI / Game loop
"""

def draw(screen, graph, vehicles):
    screen.fill((255, 255, 255))
    
    # roads
    for link in graph.links:
        x1, y1 = link.from_node.x * 80 + 50, link.from_node.y * 80 + 50
        x2, y2 = link.to_node.x * 80 + 50, link.to_node.y * 80 + 50
        pygame.draw.line(screen, (128, 128, 128), (x1, y1), (x2, y2), 5)
        pygame.draw.line(screen, (192, 192, 192), (x1+5, y1+5), (x2+5, y2+5), 5)

    # intersections
    for node in graph.nodes.values():
        x, y = node.x * 80 + 50, node.y * 80 + 50
        color = (0, 0, 255) if not node.signal else {
            'red': (255, 0, 0),
            'green': (0, 255, 0),
            'yellow': (255, 255, 0)
        }[node.signal.state]
        pygame.draw.rect(screen, color, (x-10, y-10, 20, 20))

    # vehicles
    for vehicle in vehicles:
        if vehicle.status != 'arrived' and vehicle.position < len(vehicle.path) - 1:
            current_node = vehicle.path[vehicle.position]
            next_node = vehicle.path[vehicle.position + 1]
            x = current_node.x * (1 - vehicle.progress) + next_node.x * vehicle.progress
            y = current_node.y * (1 - vehicle.progress) + next_node.y * vehicle.progress
            x, y = x * 80 + 50, y * 80 + 50
            
            # Change vehicle color based on status
            vehicle_color = (0, 0, 255)  # Default blue for moving vehicles
            if vehicle.status == 'waiting':
                vehicle_color = (255, 165, 0)  # Orange for waiting vehicles
            
            pygame.draw.rect(screen, vehicle_color, (x-5, y-5, 10, 10))

    # accidents
    for link in graph.links:
        if link.accident:
            x1, y1 = link.from_node.x * 80 + 50, link.from_node.y * 80 + 50
            x2, y2 = link.to_node.x * 80 + 50, link.to_node.y * 80 + 50
            
            accident_x = x1 + (x2 - x1) * link.accident.position
            accident_y = y1 + (y2 - y1) * link.accident.position
            
            pygame.draw.circle(screen, (255, 0, 0), (int(accident_x), int(accident_y)), 12)
            pygame.draw.circle(screen, (255, 255, 0), (int(accident_x), int(accident_y)), 8)

def main():
    pygame.init()
    screen = pygame.display.set_mode((800, 800))
    clock = pygame.time.Clock()

    nodes = read_nodes("go_traffic/nodes.csv")
    links = read_edges("go_traffic/links.csv", nodes)
    graph = Graph(nodes, links)

    # ranomize traffic signals
    for node in graph.nodes.values():
        initial_state = random.choice(['red', 'green', 'yellow'])
        initial_duration = random.randint(10, 30)
        node.signal = TrafficSignal(initial_state, initial_duration)

    # create vehicles
    num_vehicles = NUM_VEHICLES
    vehicles = []
    for i in range(num_vehicles):
        start_node_id = random.choice(list(graph.nodes.keys()))
        end_node_id = random.choice(list(graph.nodes.keys()))
        while start_node_id == end_node_id:
            end_node_id = random.choice(list(graph.nodes.keys()))

        try:
            path = find_path(graph, start_node_id, end_node_id)
            vehicles.append(Vehicle(f"V{i+1}", path, graph))
        except ValueError as e:
            print(e)

    running = True
    while running:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                running = False

        update_traffic_signals(graph)
        check_for_accidents(graph, vehicles)
        update_accidents(graph)
        update_vehicles(vehicles)
        draw(screen, graph, vehicles)

        pygame.display.flip()
        clock.tick(60)

    pygame.quit()

if __name__ == "__main__":
    main()