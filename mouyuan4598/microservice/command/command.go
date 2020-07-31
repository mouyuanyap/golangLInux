package command

import (
	"fmt"
	"io/ioutil"
	"log"
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

func DockerCreate(dir string, imageName string, buildName string, fileType string, fileChange []string) error {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	file, err := os.OpenFile(dir+"/Dockerfile", os.O_APPEND|os.O_WRONLY, 0644)
	if fileType == "go" {

		if err != nil {
			log.Println(err)
			err := ioutil.WriteFile(dir+"/Dockerfile", []byte(`
FROM golang AS builder

WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/`+imageName+`

COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o `+buildName+`

FROM scratch
COPY --from=builder `+buildName+` `+buildName+` 
ENTRYPOINT `+`["`+buildName+`"]`), 0644)
			if err != nil {
				log.Fatal(err)
				return err
			}
		} else {
			defer file.Close()
			if _, err := file.WriteString(`
FROM golang AS builder
		
WORKDIR /home/mouyuan/go/src/github.com/mouyuan4598/` + imageName + `
		
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 go build -o ` + buildName + `
		
FROM scratch
COPY --from=builder ` + buildName + ` ` + buildName + ` 
ENTRYPOINT ` + `["` + buildName + `"]`); err != nil {
				log.Fatal(err)
				return err
			}
		}

	} else if fileType == "py" {
		if err != nil {
			log.Println(err)
			err := ioutil.WriteFile(dir+"/Dockerfile", []byte(`
FROM python:3.8-slim-buster

WORKDIR `+dir+`

COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
 
CMD `+`["python", "./`+fileChange[0][:len(fileChange[0])]+`"]`), 0644)
			if err != nil {
				log.Fatal(err)
				return err
			}
		} else {
			defer file.Close()
			if _, err := file.WriteString(`
FROM python:3.8-slim-buster
			
WORKDIR ` + dir + `
			
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt
			
COPY . .
			 
CMD ` + `["python", "./` + fileChange[0][:len(fileChange[0])] + `"]`); err != nil {
				log.Fatal(err)
				return err
			}
		}

	}

	return err
}

func DockerBuild(dir string, imageName string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	var output string
	output = "Docker Build: \n"
	output = output + command(dir, "docker", "build", "-t", imageName, ".")
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

func DockerSave(dir string, name string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	command(dir, "docker", "save", "-o", dir+"/"+name+".tar", name)
	output := "\nSaved to .tar file\n"
	return output
}

func DockerTransfer(dir string, name string, ip string) string {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	command(dir, "scp", name+".tar", "root@"+ip)
	output := "\nSent to Destination Machine\n"
	return output
}
