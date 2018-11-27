package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func provisionDemos(dir string, demos map[string]DemoConfig) (map[string][]string, error) {
	m := map[string][]string{}

	for name, dc := range demos {
		if dc.Build != "" {
			cmd := exec.Command("docker", []string{"build", "-t", "demo-" + name, dc.Build}...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = append(os.Environ(), "DOCKER_BUILDKIT=1")
			cmd.Dir = dir
			if err := cmd.Run(); err != nil {
				return nil, err
			}
		}
		containerName := name + fmt.Sprintf(".%d", rand.Intn(1e5))
		cmd := exec.Command("docker", append(append([]string{"run", "-d", "--name", containerName}, dc.Flags...), "demo-"+name)...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			return nil, err
		}

		m[name] = []string{containerName, dc.Command}
	}

	return m, nil
}
