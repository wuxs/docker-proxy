package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var Mirrors = map[string]string{
	"docker.io":   "dockerproxy.com",
	"gcr.io":      "gcr.dockerproxy.com",
	"quay.io":     "quay.dockerproxy.com",
	"k8s.gcr.io":  "k8s.dockerproxy.com",
	"ghcr.gcr.io": "ghcr.dockerproxy.com",
}

func init() {
	os.Getenv("DOCKER_MIRROR")
}

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	action := os.Args[1]
	if err := CheckAction(action); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	originImage := os.Args[2]
	mirrorImage := FullImageName(originImage)
	Run(originImage, mirrorImage)
}

func Run(originImage, mirrorImage string) {
	fmt.Printf("\033[1;32mPull image <%s> with mirror image <%s>\033[0m\n", originImage, mirrorImage)
	var err error
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "pull", mirrorImage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		panic(err.Error())
	}
	err = exec.Command("docker", "tag", mirrorImage, originImage).Run()
	if err != nil {
		panic(err.Error())
	}
	err = exec.Command("docker", "rmi", mirrorImage).Run()
	if err != nil {
		panic(err.Error())
	}
}

func CheckAction(action string) error {
	if action == "pull" {
		return nil
	}
	return errors.New("invalid action")
}

func FullImageName(image string) string {
	var domain = "docker.io"
	var project = "library"
	var name string
	items := strings.Split(image, "/")
	if len(items) == 3 {
		domain = items[0]
		project = items[1]
		name = items[2]
	} else if len(items) == 2 {
		project = items[0]
		name = items[1]
	} else if len(items) == 1 {
		name = items[0]
	}
	if val, ok := Mirrors[domain]; ok {
		domain = val
	}

	return fmt.Sprintf("%s/%s/%s", domain, project, name)
}
