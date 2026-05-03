package ui

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/ingodinho/isiedocker/docker"
)

type Model struct {
    containerInfos  []docker.ContainerInfo
    cursor   int
}

func InitialModel() Model {
	return Model{
		containerInfos: make([]docker.ContainerInfo, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return getContainers
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyPressMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.containerInfos)-1 {
                m.cursor++
            }
        }

    case containersFetchedMsg:
    	m.containerInfos = msg
    }

    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m Model) View() tea.View {
    // The header
    s := "ID NAME STATUS\n\n"

    // Iterate over our choices
    for i, containerInfo := range m.containerInfos {
        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Render the row
        s += fmt.Sprintf("%s %s %s %s\n", cursor, containerInfo.Id, containerInfo.Status, containerInfo.Name)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return tea.NewView(s)
}


func getContainers() tea.Msg {
	containers, err := docker.ListContainers(true)

	if err != nil {
		return nil
	}

	return containersFetchedMsg(containers)
}


type containersFetchedMsg []docker.ContainerInfo
