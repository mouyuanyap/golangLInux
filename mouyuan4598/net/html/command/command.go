package command

import (
	"fmt"
	"os"
	"os/exec"
)

func command(dir string, exe ...string) string {
	var cmd *exec.Cmd
	if len(exe) == 1 {
		cmd = exec.Command(exe[0])
	} else if len(exe) == 2 {
		cmd = exec.Command(exe[0], exe[1])
	} else if len(exe) == 3 {
		cmd = exec.Command(exe[0], exe[1], exe[2])
	} else if len(exe) == 4 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3])
	} else if len(exe) == 5 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4])
	} else if len(exe) == 6 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4], exe[5])
	} else if len(exe) == 7 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4], exe[5], exe[6])
	} else if len(exe) == 8 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4], exe[5], exe[6], exe[7])
	} else if len(exe) == 9 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4], exe[5], exe[6], exe[7], exe[8])
	} else if len(exe) == 10 {
		cmd = exec.Command(exe[0], exe[1], exe[2], exe[3], exe[4], exe[5], exe[6], exe[7], exe[8], exe[9])
	}
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return ""
	}
	return string(output)

}

func Ll(dir string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	pwd := command(dir, "pwd")
	ll := command(dir, "ls", "-l")
	output := pwd + ll
	return output
}
func DockerBuild(dir string, imageName string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	output := command(dir, "docker", "build", "-t", imageName, ".")
	fmt.Println(output)
	return output
}

func DockerRun(dir string, port string, name string, image string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	output := command(dir, "docker", "run", "--publish", port+":8080", "--name", name, "--rm", image)
	fmt.Println(output)
	return output
}

func DockerStop(dir string, name string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	output := command(dir, "docker", "stop", name)
	fmt.Println(output)
	return output
}
