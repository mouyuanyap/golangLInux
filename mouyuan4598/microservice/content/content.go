package content

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type Page struct {
	Title string
	Body  map[string]string
	File  []byte
}

func Read(filename string) *Page {
	if filename == "change" {
		file, err := os.Open(filename + ".txt")
		if err != nil {
			m := make(map[string]string)
			return &Page{Title: "Change", Body: m}
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
				item[line[:b]] = line[a+1:]
			}

		}
		for key, value := range item {
			fmt.Println("key", key, "value", value)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return &Page{Title: "Change", Body: item}

	} else if filename == "status" {
		body, err := ioutil.ReadFile(filename + ".txt")
		if err != nil {
			return &Page{Title: "Status", File: []byte(err.Error())}
		}
		return &Page{Title: "Status", File: body}
	} else {
		m := make(map[string]string)
		return &Page{Title: "Error", Body: m}
	}

}

func save(filename string, content [2]string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		err := ioutil.WriteFile(filename, []byte(content[0]+","+content[1]+"\n"), 0644)
		if err != nil {
			log.Fatal(err)
			return err
		}
	} else {
		defer file.Close()
		if _, err := file.WriteString(content[0] + "," + content[1] + "\n"); err != nil {
			log.Fatal(err)
			return err
		}
	}
	return err

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
					save("change.txt", eventDetail)
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("created file:", event.Name)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/home/mouyuan/go/src/github.com/mouyuan4598/gin")
	err = watcher.Add("/home/mouyuan/go/src/github.com/mouyuan4598/hello")
	err = watcher.Add("/home/mouyuan/go/src/github.com/mouyuan4598/server")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
