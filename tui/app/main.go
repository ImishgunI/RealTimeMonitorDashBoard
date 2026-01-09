package main

import (
	"fmt"
	"math"
	"os"
	"real_time_monitor_dashboard/backend/metrics"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

/*
	Main idea is to use some msg types to update behavour of program, fetchData, Call ticker to fetchdata in period,
	View() only render ui and nothing else, Update() react to some msg in runtime and do some tasks.
*/

type tickMsg time.Time

type model struct {
	cpu metrics.CPUMetrics
}

type cpuUpdateMsg struct {
	cpu metrics.CPUMetrics
}

func fetchData() tea.Msg {
	cpu := metrics.New()
	cpu = metrics.Init(cpu)

	return cpuUpdateMsg{
		cpu: *cpu,
	}
}

func tickCmd(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Constuctor for model struct
func initialModel() model {
	return model{
		cpu: *metrics.New(),
	}
}

// Default function, needs to some i/o operations like saving/loading to/from disk
func (m model) Init() tea.Cmd {
	return tea.Batch(
		fetchData,
		tickCmd(time.Nanosecond),
	)
}

// Update method needs to check what happends and return update model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// which type our msg
	switch msg := msg.(type) {
	// is key pressed?
	case tea.KeyMsg:
		// what key was pressed?
		switch msg.String() {
		// This key should exit program
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case tickMsg:
		return m, tea.Batch(
			fetchData,
			tickCmd(500*time.Millisecond),
		)
	case cpuUpdateMsg:
		m.cpu = msg.cpu
		return m, nil
	}
	return m, nil
}

// Render UI
func (m model) View() string {
	cpuinfo := fmt.Sprintf(
		"Name: %s\nCores: %d\nThreads: %d\nFrequency: %.1fMHz\nTemreture: %d℃\nWorkload: %d%%\n",
		m.cpu.Name, m.cpu.Cores, m.cpu.Threads, m.cpu.Frequency/1000,
		int(math.Ceil(float64(m.cpu.Temreture))),
		int(math.Ceil(float64(m.cpu.Workload))),
	)
	s := cpuinfo
	return s
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
