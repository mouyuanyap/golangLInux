package content

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	
	"github.com/mouyuan4598/server/command"

	"github.com/fsnotify/fsnotify"
)

func Read(filename string) {

	file, err := os.Open(filename + ".txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()
	item := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		length := len(scanner.Text())
		var a, b int
		for i := length - 1; i > -1; i-- {
			if line[i] == ',' {
				a = i
			}
			if line[i] == '/' {
				b = i
				break
			}
		}
		if item[line[:b]] == "" {
			item[line[b+1:a-3]] = line[a+1:]
		}

	}
	for key, value := range item {
		fmt.Println("key", key, "value", value)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func Save(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		err := ioutil.WriteFile(filename, []byte(content+"\n"), 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		defer file.Close()
		if _, err := file.WriteString(content + "\n"); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return err

}

func action(filename string) {
	command.DockerStop("/root/dockerimage", "/root/dockerimage/server/stop.sh")
//	command.DockerUnzip("/root/dockerimage", filename)
//	command.DockerRun("/root/dockerimage", "8080", filename)
}

func Change() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				//event.Name event.Op.String()
				if event.Op&fsnotify.Write == fsnotify.Write {
					var eventDetail [2]string
					eventDetail[0] = event.Name
					eventDetail[1] = event.Op.String()
					log.Println("modified file:", event.Name)
					length := len(eventDetail[0])
					var b int
					for i := length - 1; i > -1; i-- {
//						fmt.Println(i)
//						fmt.Println(eventDetail[0][i])
						if eventDetail[0][i] == '/' {
							b = i
							break
						}

					}
					//fmt.Println(eventDetail[0][b+1 : length-4])
					Save("change.txt", eventDetail[0][b+1:length-4])
					action(eventDetail[0][b+1 : length-4])
				
				}// else {
				//	var eventDetail [2]string
				//	eventDetail[0] = event.Name
				//	eventDetail[1] = event.Op.String()
					//log.Println("modified file:", event.Name)
				//	length := len(eventDetail[0])
				//	var b int
				//	for i := length - 1; i > -1; i-- {

				//		if eventDetail[0][i] == '/' {
				//			b = i
				//			break
				//		}
						//fmt.Println(eventDetail[0][b+1 : length-4])
				//		Save("change.txt", eventDetail[0][b+1:length-4])
				//		action(eventDetail[0][b+1 : length-4])
				//	}
//				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/root/dockerimage")
	if err != nil {
		fmt.Println("wtf")
		log.Fatal(err)
	}
	<-done
}
