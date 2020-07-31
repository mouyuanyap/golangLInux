package content

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dietsche/rfsnotify"
	"gopkg.in/fsnotify.v1"
)

type Page struct {
	Title     string
	Body      map[string][]string
	File      []byte
	Directory map[string]string
	FileType  map[string]string
}

func Read(filename string) *Page {
	if filename == "change" {
		file, err := os.Open(filename + ".txt")
		if err != nil {
			m := make(map[string][]string)
			return &Page{Title: "Change", Body: m}
		}
		defer file.Close()
		item := make(map[string][]string)
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
			if len(item[line[:b]]) > 0 {
				found := false
				for _, x := range item[line[:b]] {
					if x == line[b+1:a] {
						found = true
						break
					}
				}
				if found == false {
					item[line[:b]] = append(item[line[:b]], line[b+1:a])
				}
			} else {
				item[line[:b]] = append(item[line[:b]], line[b+1:a])
			}

		}
		for key, value := range item {
			fmt.Println("key", key, "value", value)
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		dir := make(map[string]string)
		for key := range item {
			var b int
			base := len("/home/mouyuan/go/src/github.com/mouyuan4598")
			for i := base + 1; i < len(key); i++ {
				//fmt.Print(i)
				//fmt.Println(string(key[i]))
				if key[i] == '/' {
					b = i
					break
				}
				if i == len(key)-1 {
					b = i + 1
				}
			}
			dir[key] = key[base+1 : b]
			//fmt.Println(base)
			//fmt.Println(b)
		}
		ty := make(map[string]string)
		for key, value := range item {
			var y, z int
			for i := 0; i < len(value); i++ {
				for j := 0; j < len(value[i]); j++ {

					if value[i][j] == '.' {
						y = i
						z = j
						break
					}
				}
			}
			ty[key] = value[y][z+1 : len(value[y])]
		}

		return &Page{Title: "Change", Body: item, Directory: dir, FileType: ty}

	} else if filename == "status" {
		body, err := ioutil.ReadFile(filename + ".txt")
		if err != nil {
			return &Page{Title: "Status", File: []byte(err.Error())}
		}
		return &Page{Title: "Status", File: body}
	} else {
		m := make(map[string][]string)
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
	watcher, err := rfsnotify.NewWatcher()
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

	err = watcher.AddRecursive("/home/mouyuan/go/src/github.com/mouyuan4598/gin")
	err = watcher.AddRecursive("/home/mouyuan/go/src/github.com/mouyuan4598/hello")
	err = watcher.AddRecursive("/home/mouyuan/go/src/github.com/mouyuan4598/py")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
