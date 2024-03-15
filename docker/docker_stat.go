package docker

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os/exec"
	"strings"
)

type DockerStat struct {
	ContainerID string `json:"ID"`
	Name        string `json:"Name"`
	CPUPerc     string `json:"CPUPerc"`
	MemUsage    string `json:"MemUsage"`
	MemPerc     string `json:"MemPerc"`
	NetIO       string `json:"NetIO"`
	BlockIO     string `json:"BlockIO"`
	PIDs        string `json:"PIDs"`
}

func Read(callback func([]DockerStat, error)) {
	command := "docker"
	args := []string{"stats", "--format", `{{json .}}`}
	cmd := exec.Command(command, args...)

	// Get the stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("cmd.StdoutPipe() failed with %s\n", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start() failed with %s\n", err)
	}

	// Create a new scanner
	reader := bufio.NewReader(stdout)

	for {
		var dockerStats []DockerStat
		buffer := make([]byte, 4096)
		_, err := reader.Read(buffer)
		if err != nil {
			callback(nil, err)
		}
		lines := convertBufferToLines(buffer)
		for _, line := range lines {
			var dockerStat DockerStat
			err = json.Unmarshal([]byte(line), &dockerStat)
			if err != nil {
				callback(nil, err)
			}
			dockerStats = append(dockerStats, dockerStat)
		}
		callback(dockerStats, nil)
	}
}

func convertBufferToLines(buffer []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		start := strings.Index(line, "{")
		end := strings.LastIndex(line, "}")
		if start != -1 && end != -1 && start < end {
			line = line[start : end+1]
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Scanner error: %s", err)
	}
	return lines
}
