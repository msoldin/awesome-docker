// Package cmd /*
package cmd

import (
	"awesome-docker/docker"
	"encoding/json"
	"fmt"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
	"time"
)

// termCmd represents the term command
var termCmd = &cobra.Command{
	Use:   "term",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		docker.Read(func(stats []docker.DockerStat, err error) {
			if err != nil {
				log.Fatalf("docker.Read() failed with %s\n", err)
			}
			//fmt.Println(stats)
		})
	},
}

func init() {
	rootCmd.AddCommand(termCmd)
}

func term() {
	app := tview.NewApplication()

	// Create a new grid and set its dimensions
	grid := tview.NewGrid().
		SetRows(6, 0, 6).
		SetColumns(60, 0, 60).
		SetBorders(true)

	// Create and configure the panels
	panel1 := tview.NewTextView().SetText("Containers").SetScrollable(true)

	// Add the panels to the grid
	grid.AddItem(panel1, 0, 0, 10, 10, 0, 0, false)

	// Create a goroutine that continuously fetches the Docker stats and updates the TextView
	go func() {
		for {
			// Fetch the Docker stats
			stats := dockerStream()

			headers := fmt.Sprintf("%-12s %-12s %-12s %-12s %-12s %-12s %-12s\n", "ID", "Name", "CPU", "Mem", "NetIO", "BlockIO", "PIDs")
			app.QueueUpdateDraw(func() {
				panel1.SetText(headers)
			})

			// Create a variable to hold the formatted string
			var statsText string

			// Format the Docker stats and append to the statsText string
			for _, stat := range stats {
				statsText += fmt.Sprintf("%-12s %-12s %-12s %-12s %-12s %-12s %-12s\n",
					stat.ContainerID,
					stat.Name,
					stat.CPUPerc,
					stat.MemUsage,
					stat.NetIO,
					stat.BlockIO,
					stat.PIDs)
			}

			// Use QueueUpdateDraw to ensure that the update is thread-safe
			app.QueueUpdateDraw(func() {
				// Set the text of the panel to the statsText string
				panel1.SetText(panel1.GetText(true) + statsText)
			})

			// Sleep for a while before fetching the stats again
			time.Sleep(time.Second * 1)
		}
	}()

	// Set the grid as the application root
	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}

func dockerStream() []docker.DockerStat {
	dockerCommand := "docker"
	dockerArgs := []string{"stats", "--no-stream", "--format", `{{json .}}`}
	cmd := exec.Command(dockerCommand, dockerArgs...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	var containers []docker.DockerStat
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var container docker.DockerStat
		err = json.Unmarshal([]byte(line), &container)
		if err != nil {
			log.Fatalf("json.Unmarshal() failed with %s\n", err)
		}
		containers = append(containers, container)
	}

	return containers
}
