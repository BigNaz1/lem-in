package antfarm

import (
	"fmt"
	"strings"
	"time"
)

type Visualizer struct {
	field      *Field
	movements  []string
	antSymbols map[int]rune
	layout     []string
}

func NewVisualizer(field *Field, movements []string) *Visualizer {
	v := &Visualizer{
		field:      field,
		movements:  movements,
		antSymbols: make(map[int]rune),
	}

	// Assign a unique symbol to each ant
	symbols := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9'}
	for i, ant := range field.ants {
		v.antSymbols[ant.id] = symbols[i%len(symbols)]
	}

	v.generateLayout()

	return v
}

func (v *Visualizer) generateLayout() {
	// This is a simplified layout generation.
	// For a more accurate representation, you'd need a more complex algorithm.
	v.layout = []string{
		"        _________________",
		"       /                 \\",
		"  ____[%s]----[%s]--[%s]     |",
		" /            |    /      |",
		"[%s]---[%s]----[%s]  /       |",
		" \\   ________/|  /        |",
		"  \\ /        [%s]/________/",
		"  [%s]_________/",
	}
}

func (v *Visualizer) Run() {
	// Clear console and hide cursor
	fmt.Print("\033[2J")
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor when done

	// Initial state
	v.drawFarm()
	time.Sleep(2 * time.Second)

	// Process each move
	for _, move := range v.movements {
		v.processMove(move)
		v.drawFarm()
		time.Sleep(500 * time.Millisecond)
	}
}

func (v *Visualizer) processMove(move string) {
	parts := strings.Split(move, " ")
	for _, part := range parts {
		antMove := strings.Split(part, "-")
		if len(antMove) == 2 {
			antName := antMove[0]
			roomName := antMove[1]
			for i, ant := range v.field.ants {
				if fmt.Sprintf("L%d", ant.id) == antName {
					v.field.ants[i].currentRoom = roomName
					break
				}
			}
		}
	}
}

func (v *Visualizer) drawFarm() {
	// Clear screen
	fmt.Print("\033[H")

	roomContents := make(map[string]string)
	for _, room := range v.field.rooms {
		content := room.name
		for _, ant := range v.field.ants {
			if ant.currentRoom == room.name {
				content += string(v.antSymbols[ant.id])
			}
		}
		roomContents[room.name] = content
	}

	// Print the layout
	for _, line := range v.layout {
		formattedLine := line
		for _, room := range v.field.rooms {
			formattedLine = strings.Replace(formattedLine, "[%s]", "["+roomContents[room.name]+"]", 1)
		}
		fmt.Println(formattedLine)
	}

	// Print legend
	fmt.Println("\nLegend:")
	for _, ant := range v.field.ants {
		fmt.Printf("%c: Ant L%d\n", v.antSymbols[ant.id], ant.id)
	}
}

func RunVisualizer(field *Field, movements []string) {
	v := NewVisualizer(field, movements)
	v.Run()
}
