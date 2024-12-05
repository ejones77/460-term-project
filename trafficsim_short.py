import trafficSimulator as ts

sim = ts.Simulation()

sim.create_quadratic_bezier_curve((0, 0), (50, 0), (50, 50))
sim.create_vehicle(path=[0])


# Add road segments
sim.create_segment((300, 98), (0, 98))
sim.create_segment((0, 102), (300, 102))
sim.create_segment((180, 60), (0, 60))
sim.create_segment((220, 55), (180, 60))
sim.create_segment((300, 30), (220, 55))
sim.create_segment((180, 60), (160, 98))
sim.create_segment((158, 130), (300, 130))
sim.create_segment((0, 178), (300, 178))
sim.create_segment((300, 182), (0, 182))
sim.create_segment((160, 102), (155, 180))


# Add vehicle generator
sim.create_vehicle_generator(
    vehicle_rate=20,
    vehicles=[
        (10, {'path': [1], 'v': 16.6}),
        (1, {'path': [3], 'v': 16.6, 'l': 7}),
        (10, {'path': [2], 'v': 16.6}),
        (1, {'path': [2], 'v': 16.6, 'l': 7}),
        (10, {'path': [3], 'v': 16.6}),
        (1, {'path': [3], 'v': 16.6, 'l': 7}),
        (10, {'path': [4], 'v': 16.6}),
        (1, {'path': [4], 'v': 16.6, 'l': 7}),
        ]
    )

# Show simulation visualization
win = ts.Window(sim)
win.show()
