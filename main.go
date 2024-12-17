package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"
)

// Log configuration
var log = logrus.New()

// CLI flags
var outputFormat string

// Structure to hold system check results
type CheckResult struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type model struct {
	results  []CheckResult
	quitting bool
	spinner  int
}

// Tick message type for spinner animation
type tickMsg time.Time

// Visual styles
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF79C6"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	loadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C")).Bold(true)
	footerStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4")).Italic(true)
)

// Spinner animation frames
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func runChecks() []CheckResult {
	var wg sync.WaitGroup
	results := make([]CheckResult, 0)
	mutex := &sync.Mutex{}

	checks := []struct {
		Name    string
		Cmd     string
		ErrHint string
	}{
		{"System Update", "sudo apt update -y 2>/dev/null | grep -v 'WARNING'", "Failed to fetch updates. Ensure apt is installed and configured."},
		{"Kernel Check", "uname -r", "Kernel information not available."},
		{"UFW Firewall Status", "sudo ufw status | grep -q active", "UFW firewall is inactive or not installed."},
		{"SSH Security", "grep -q 'PermitRootLogin no' /etc/ssh/sshd_config", "Root login over SSH is permitted. Update sshd_config."},
		{"Disk Usage", "df -h", "Disk usage information could not be retrieved."},
		{"Memory Usage", "free -m", "Memory usage data is unavailable."},
		{"Service Status (rsyslog)", "systemctl is-active --quiet rsyslog", "rsyslog service is not active."},
		{"Cron Jobs", "crontab -l", "No cron jobs found for the current user."},
		{"TLS Support", "openssl ciphers -v | grep -q 'TLSv1.2\\|TLSv1.3'", "TLSv1.2 or TLSv1.3 support is missing."},
		{"Password Policy", "grep -q 'minlen' /etc/security/pwquality.conf", "Password policy not enforced. Check pwquality.conf."},
	}

	for _, check := range checks {
		wg.Add(1)
		go func(name, cmd, errHint string) {
			defer wg.Done()
			start := time.Now()
			status, msg := runCommand(cmd)
			if status == "Failed" {
				msg = errHint + " (" + msg + ")"
			}
			elapsed := time.Since(start)

			mutex.Lock()
			results = append(results, CheckResult{
				Name:    name,
				Status:  status,
				Message: fmt.Sprintf("%s (%.2fs)", msg, elapsed.Seconds()),
			})
			mutex.Unlock()
		}(check.Name, check.Cmd, check.ErrHint)
	}

	wg.Wait()
	return results
}

func runCommand(cmd string) (string, string) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		return "Failed", strings.TrimSpace(string(out))
	}
	return "Passed", strings.TrimSpace(string(out))
}

type checkResultsMsg []CheckResult

type quitMsg struct{}

func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return checkResultsMsg(runChecks())
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case checkResultsMsg:
		m.results = msg
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, func() tea.Msg {
				return quitMsg{}
			}
		}
	case tickMsg:
		m.spinner = (m.spinner + 1) % len(spinnerFrames)
		return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	case quitMsg:
		m.quitting = true
		return m, tea.Quit
	}
	return m, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func formatMessage(msg string) string {
	parts := strings.Split(msg, " (")
	if len(parts) != 2 {
		return msg
	}

	content := parts[0]
	timing := "(" + parts[1]

	if strings.Contains(content, "\n") {
		lines := strings.Split(content, "\n")
		indentedLines := make([]string, len(lines))
		indentedLines[0] = ""
		for i := 1; i < len(lines); i++ {
			indentedLines[i] = "        " + lines[i]
		}
		return strings.Join(indentedLines, "\n") + " " + timing
	}

	return content + " " + timing
}

func (m model) View() string {
	if m.quitting {
		return "Exiting...\n"
	}

	if len(m.results) == 0 {
		return loadingStyle.Render(fmt.Sprintf("Performing system checks... %s\n", spinnerFrames[m.spinner]))
	}

	var resultView strings.Builder
	w := tabwriter.NewWriter(&resultView, 2, 4, 2, ' ', 0)

	fmt.Fprintln(w, titleStyle.Render("System Check Results:"))
	fmt.Fprintln(w)

	for _, result := range m.results {
		statusSymbol := successStyle.Render("✔")
		messageStyle := successStyle
		if result.Status == "Failed" {
			statusSymbol = errorStyle.Render("✘")
			messageStyle = errorStyle
		}

		formattedMsg := formatMessage(result.Message)
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			statusSymbol,
			result.Name+"\t",
			messageStyle.Render(formattedMsg))
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, footerStyle.Render("Press 'q' to quit"))

	w.Flush()

	if outputFormat == "json" {
		jsonData, _ := json.MarshalIndent(m.results, "", "  ")
		return string(jsonData)
	}

	return resultView.String()
}

func main() {
	log.Out = os.Stdout
	log.SetLevel(logrus.InfoLevel)

	if len(os.Args) > 1 && os.Args[1] == "--json" {
		outputFormat = "json"
	}

	if os.Geteuid() != 0 {
		log.Fatal("This program must be run as root.")
	}

	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		log.Fatalf("Error starting program: %v", err)
	}
}
