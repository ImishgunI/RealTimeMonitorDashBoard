package main

import (
	"fmt"
	"math"
	"os"
	metrics "real_time_monitor_dashboard/backend/metrics/cpu_metrics"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().
			Margin(1, 2).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Width(50)
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffffff"))

	warningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF0000"))

	okStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF00"))

	lowUsageColor    = lipgloss.Color("#00FF00")
	mediumUsageColor = lipgloss.Color("#FFFF00")
	highUsageColor   = lipgloss.Color("#FF0000")
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
		"Name: %s\nCores: %d\nThreads: %d\nFrequency: %.1fMHz\nTemreture: %d℃\n",
		m.cpu.Name, m.cpu.Cores, m.cpu.Threads, m.cpu.Frequency/1000,
		int(float64(m.cpu.Temreture)),
	)
	pbar := progressBar(30, int(math.Ceil(float64(m.cpu.Workload))))
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		cpuinfo,
		labelStyle.Render("CPU Usage:"),
		fmt.Sprintf("%s	[%s%d%%]", "CPU", pbar, int(math.Ceil(float64(m.cpu.Workload)))),
	)
	return appStyle.Render(content)
}

func progressBar(width int, percent int) string {
	filledWidth := int(width * percent / 100)
	bar := strings.Builder{}

	barColor := lowUsageColor
	if percent > 50 {
		barColor = mediumUsageColor
	}
	if percent > 80 {
		barColor = highUsageColor
	}

	filledStyle := lipgloss.NewStyle().Foreground(barColor)
	for i := range width {
		if i < filledWidth {
			bar.WriteString(filledStyle.Render("|"))
		} else {
			bar.WriteString(" ")
		}
	}
	return bar.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
