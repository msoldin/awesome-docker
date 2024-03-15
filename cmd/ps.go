package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os/exec"
	"strings"
)

type DockerContainer struct {
	ID         string `json:"ID"`
	Image      string `json:"Image"`
	Command    string `json:"Command"`
	CreatedAt  string `json:"CreatedAt"`
	RunningFor string `json:"RunningFor"`
	Ports      string `json:"Ports"`
	Status     string `json:"Status"`
	Size       string `json:"Size"`
	Names      string `json:"Names"`
	Labels     string `json:"Labels"`
	Networks   string `json:"Networks"`
	Mounts     string `json:"Mounts"`
}

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dockerCommand := "docker"
		dockerArgs := []string{"ps", "--format", `{{json .}}`}
		executeCommand(dockerCommand, dockerArgs)
	},
}

func executeCommand(command string, args []string) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	var containers []DockerContainer
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		var container DockerContainer
		err = json.Unmarshal([]byte(line), &container)
		if err != nil {
			log.Fatalf("json.Unmarshal() failed with %s\n", err)
		}
		containers = append(containers, container)
	}

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()

	for _, container := range containers {
		fmt.Printf("%s: %s\n%s: %s\n%s: %s\n\n\n",
			red("ID"), container.ID,
			green("Image"), container.Image,
			blue("Names"), container.Names)
	}
}

func init() {
	rootCmd.AddCommand(psCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
