package main

import (
	// "github.com/ingodinho/isiedocker/docker"
	"github.com/ingodinho/isiedocker/ui"
	tea "charm.land/bubbletea/v2"
	"fmt"
	"os"
)

func main() {
	// docker.ListContainers(true)
	p := tea.NewProgram(ui.InitialModel())
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}
