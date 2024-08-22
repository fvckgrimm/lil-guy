package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

const maxOutputLines = 60 // Adjust this value based on your screen size

type model struct {
	character    Character
	message      string
	face         string
	faceIndex    int
	frame        string
	frames       []string
	output       []string
	outputCursor int
	quitting     bool
	inputChan    chan string
}

func initialModel(character Character, message string, inputChan chan string) model {
	return model{
		character: character,
		message:   message,
		face:      character.Faces[0],
		frames:    []string{"<", "-", ">", " "},
		output:    []string{},
		inputChan: inputChan,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), m.readInput)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "up":
			if m.outputCursor > 0 {
				m.outputCursor--
			}
		case "down":
			if m.outputCursor < len(m.output)-maxOutputLines {
				m.outputCursor++
			}
		}
	case tickMsg:
		m.frame = m.frames[time.Now().Unix()%int64(len(m.frames))]
		m.face, m.faceIndex = m.character.NextFace(m.faceIndex)
		return m, tick()
	case inputMsg:
		m.output = append(m.output, string(msg))
		if len(m.output) > maxOutputLines {
			m.outputCursor = len(m.output) - maxOutputLines
		}
		return m, m.readInput
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	s := strings.Builder{}
	width, _, _ := term.GetSize(int(os.Stdout.Fd()))
	if width == 0 {
		width = 80 // fallback width if unable to get terminal size
	}

	leftPadding := 3
	contentWidth := width - leftPadding

	// Message
	wrappedMessage := wrapText(m.message, contentWidth)
	for _, line := range wrappedMessage {
		s.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", leftPadding), line))
	}
	s.WriteString("\n")

	// Character
	characterLines := strings.Split(m.face, "\n")
	if len(characterLines) == 1 {
		// Single-line character, add arms
		leftArm, rightArm := getArms(m.frame)
		s.WriteString(fmt.Sprintf("%s%s %s %s\n", strings.Repeat(" ", leftPadding), leftArm, m.face, rightArm))
	} else {
		// Multi-line character, don't add arms
		for _, line := range characterLines {
			s.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", leftPadding), line))
		}
	}

	// Output
	pipedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	endIndex := m.outputCursor + maxOutputLines
	if endIndex > len(m.output) {
		endIndex = len(m.output)
	}

	for i, line := range m.output[m.outputCursor:endIndex] {
		if i == 0 {
			// First line of output, place it next to the character
			s.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", leftPadding+20), pipedStyle.Render(line)))
		} else {
			// Subsequent lines, align with the first output line
			s.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(" ", leftPadding+20), pipedStyle.Render(line)))
		}
	}

	// Scrollbar indicator
	if len(m.output) > maxOutputLines {
		s.WriteString(fmt.Sprintf("\n%s(%d/%d)", strings.Repeat(" ", leftPadding), m.outputCursor+1, len(m.output)-maxOutputLines+1))
	}

	return s.String()
}

func wrapPreserveNewlines(text string, width int) string {
	lines := strings.Split(text, "\n")
	var result []string
	for _, line := range lines {
		if len(line) <= width {
			result = append(result, line)
		} else {
			wrapped := ""
			for len(line) > width {
				wrapped += line[:width] + "\n"
				line = line[width:]
			}
			if len(line) > 0 {
				wrapped += line
			}
			result = append(result, wrapped)
		}
	}
	return strings.Join(result, "\n")
}

func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	line := ""

	for _, word := range words {
		if len(line)+len(word)+1 > width {
			lines = append(lines, strings.TrimSpace(line))
			line = word
		} else {
			if line != "" {
				line += " "
			}
			line += word
		}
	}

	if line != "" {
		lines = append(lines, strings.TrimSpace(line))
	}

	return lines
}

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*250, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type inputMsg string

func (m model) readInput() tea.Msg {
	input, ok := <-m.inputChan
	if !ok {
		return nil
	}
	return inputMsg(input)
}

func isMultiline(s string) bool {
	return strings.Contains(s, "\n")
}
