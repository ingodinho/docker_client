package ui

import (
	"fmt"
	"log"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/ingodinho/isiedocker/docker"
)

type Model struct {
	containerInfos []docker.ContainerInfo
	cursor         int
}
type containersFetchedMsg []docker.ContainerInfo
type tickMsg time.Time

func InitialModel() Model {
	return Model{
		containerInfos: make([]docker.ContainerInfo, 0),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(getContainers, doTick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		{
			return m, tea.Batch(getContainers, doTick())
		}

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

		// The s key starts a container
		case "s":
			containerId := m.containerInfos[m.cursor].Id
			return m, startContainer(containerId)

		case "t":
			containerId := m.containerInfos[m.cursor].Id
			return m, stopContainer(containerId)

		case "r":
			containerId := m.containerInfos[m.cursor].Id
			return m, restartContainer(containerId)
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
	s := "ID STATUS NAME STATE\n\n"

	// Iterate over our choices
	for i, containerInfo := range m.containerInfos {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s %s %s %s\n", cursor, containerInfo.Id, containerInfo.Status, containerInfo.Name, containerInfo.State)
	}

	// The footer
	s += "\n q to quit. "
	s += "s to start container. "
	s += "t to stop container. "
	s += "r to restart container. "

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

func startContainer(id string) func() tea.Msg {
	return func() tea.Msg {
		containerInfos, err := docker.StartContainer(id)
		if err != nil {
			log.Fatal("error starting container", err)
		}

		return containersFetchedMsg(containerInfos)
	}
}

func stopContainer(id string) func() tea.Msg {
	return func() tea.Msg {
		containerInfos, err := docker.StopContainer(id)
		if err != nil {
			log.Fatal("error stopping container", err)
		}

		return containersFetchedMsg(containerInfos)
	}
}

func restartContainer(id string) func() tea.Msg {
	return func() tea.Msg {
		containerInfos, err := docker.RestartContainer(id)
		if err != nil {
			log.Fatal("error restarting container", err)
		}

		return containersFetchedMsg(containerInfos)
	}
}

func doTick() tea.Cmd {
	return tea.Tick(time.Second*3, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
