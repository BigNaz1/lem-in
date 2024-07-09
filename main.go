package main

import (
	"antfarm/antfarm"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: Invalid number of arguments")
		return
	}

	filePath := os.Args[1]
	lines, err := antfarm.ReadFile(filePath)
	if err != nil {
		fmt.Println("ERROR: File not found")
		return
	}

	// Check for multiple start or end points
	startCount,endCount := 0, 0
	for _,line := range lines{
		if line == "##start"{
			startCount++
		}
		if line == "##end"{
			endCount++
		}
	}
	
	if startCount > 1 || endCount > 1{
		fmt.Println("Error: You cannot have more than one starting or ending point")
		return
	}

	// Print input
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()

	field := antfarm.Field{}
	err = antfarm.ResolveInput(lines, &field)
	if err != nil {
		fmt.Println("ERROR: invalid data format, " + err.Error())
		return
	}

	var shortestPaths [][]string
	visited := make(map[string]bool)
	antfarm.FindShortestPaths(field, "", []string{}, visited, &shortestPaths)
	shortestPaths = antfarm.RemoveTooLongPaths(shortestPaths)

	if len(shortestPaths) == 0 {
		fmt.Printf("ERROR: invalid data format, farm is unsolvable")
		return
	}

	// Capture the movements
	movements := antfarm.StartTurnBasedSolving(&field, shortestPaths)

	// Print movements
	fmt.Println("Ant movements:")
	for _, move := range movements {
		fmt.Println(move)
	}

	// Run visualizer if "visualize" argument is provided
	if len(os.Args) > 2 && os.Args[2] == "visualize" {
		fmt.Println("\nStarting visualization...")
		antfarm.RunVisualizer(&field, movements)
	}
}
