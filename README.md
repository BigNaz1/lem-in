# Ant Farm Simulator

## Overview
This project implements an Ant Farm Simulator in Go. The program efficiently moves a group of ants through a network of interconnected rooms, from a designated start room to an end room, optimizing for the fewest number of turns. It also includes a visualization feature to display the ant movements graphically.

## Features
- Simulates ant movement through a complex room network
- Finds optimal paths for multiple ants simultaneously
- Avoids collisions and bottlenecks in ant movement
- Outputs a turn-by-turn movement log
- Provides a visualization of ant movements (optional)

## Code Structure
The main components of the code are:

1. `Field` struct: Represents the entire ant farm, including ants, rooms, and start/end points.
2. `Ant` struct: Represents individual ants with their current position and status.
3. `Room` struct: Represents rooms in the farm with connections to other rooms.
4. `StartTurnBasedSolving` function: The main algorithm for moving ants through the farm.
5. `getNextRoom` function: Determines the next best move for each ant.
6. Helper functions: `isRoomEmpty`, `getNumOfFinishedAnts`, etc.
7. Visualization functions (details may vary based on implementation)

## Algorithm
The algorithm uses a turn-based approach to move ants:
1. For each turn, it attempts to move as many ants as possible.
2. It prioritizes shorter paths and avoids creating bottlenecks.
3. The simulation continues until all ants reach the end room.

## Usage
To use the simulator:

1. Create a text file (e.g., `example00.txt`) with the ant farm configuration. The file should follow this format:
   - First line: number of ants
   - Subsequent lines: room definitions and connections
   - Use `##start` and `##end` to denote the start and end rooms

   Example:
    4
    ##start
    0 0 3
    2 2 5
    3 4 0
    ##end
    1 8 3
    0-2
    2-3
    3-1

2. Run the program with the input file:
    $ go run . example00.txt

3. To run with visualization (if implemented):
    $ go run . example00.txt visualize

## Example Output
Text output:


    L1-2
    L1-3 L2-2
    L1-1 L2-3 L3-2
    L2-1 L3-3 L4-2
    L3-1 L4-3
    L4-1

This output shows the movement of 4 ants through the farm, completing in 6 turns.

Visualization output: 
        _________________  
       /                 \ 
  ____[0]----[2]--[3]     | 
 /            |    /      |
[0]---[2]----[3]  /       |
 \   ________/|  /        |
  \ /        [0]/________/
  [0]_________/

## Performance
The current implementation achieves optimal or near-optimal solutions for complex ant farm configurations, balancing path utilization and minimizing the number of turns.

## Visualization
The visualization feature provides a graphical representation of the ant farm and the movement of ants.

## Future Improvements
Potential areas for enhancement include:
- Further optimization of initial path assignments
- Implementation of more advanced look-ahead mechanisms
- Fine-tuning of path selection criteria
- Enhancing the visualization with more interactive features

## Contributors
Nezar Jaberi & Abdulaziz Rajab
