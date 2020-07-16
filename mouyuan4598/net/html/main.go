package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/mouyuan4598/net/html/command"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func loadFile() (*Page, error) {
	var files []string
	root := "/home/mouyuan/go/src/github.com/mouyuan4598/net/html"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".txt" {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		panic(err)
	}
	var filename string
	for _, file := range files {
		filename = filename + file + "\n"
	}
	body := []byte(filename)
	return &Page{Title: "file", Body: body}, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, err := template.ParseFiles(tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/", http.StatusFound)
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	p, err := loadFile()
	if err != nil {

		return
	}
	renderTemplate(w, "read", p)
}

func buildHandler(w http.ResponseWriter, r *http.Request) {
	command.DockerStop("../../gin", "test")
	message := ""
	message = message + command.DockerBuild("../../gin", "project")
	m := &Page{Title: "Build", Body: []byte(message)}
	err := m.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+m.Title, http.StatusFound)
	_ = command.DockerRun("../../gin", "6060", "test", "project")

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
					log.Println("modified file:", event.Name)
					file, err := os.OpenFile("change.txt", os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						log.Println(err)
						err := ioutil.WriteFile("change.txt", []byte(event.Name+"\n"), 0644)
						if err != nil {
							log.Fatal(err)
						}
					} else {
						defer file.Close()
						if _, err := file.WriteString(event.Name); err != nil {
							log.Fatal(err)
						}

					}

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

	err = watcher.Add("../../gin")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func main() {
	go Change()
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/read/", readHandler)
	http.HandleFunc("/build/", buildHandler)
	log.Fatal(http.ListenAndServe(":7070", nil))

}
