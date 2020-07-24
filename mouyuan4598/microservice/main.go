package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/mouyuan4598/microservice/command"
	"github.com/mouyuan4598/microservice/content"
)

func save(filename string, content string) error {
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

func renderTemplate(w http.ResponseWriter, tmpl string, p *content.Page) {
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

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
	p := content.Read("change")
	if len(p.Body) == 0 {
		renderTemplate(w, "view1", p)
	} else {
		renderTemplate(w, "view2", p)
	}

}
func buildHandler(w http.ResponseWriter, r *http.Request) {
	var err = os.Remove("status.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	p := content.Read("change")
	var b int
	j := 1

	for key := range p.Body {
		fmt.Println(j)
		for i := len(key) - 1; i > -1; i-- {
			if key[i] == '/' {
				b = i
				break
			}
		}
		var temp string
		err = os.Remove(key + "/Dockerfile")
		if err != nil {
			fmt.Println(err.Error())
		}
		command.DockerCreate(key, key[b+1:])

		temp = command.DockerBuild(key, key[b+1:])
		save("status.txt", temp)
		temp = command.DockerSave(key, key[b+1:])
		save("status.txt", temp)
		temp = command.DockerTransfer(key, key[b+1:], "149.28.11.6:~/dockerimage/")
		save("status.txt", temp)
		save("status.txt", "\n\nFinished successfully")
		j++
	}
	err = os.Remove("change.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	//command.DockerStop("../../gin", "test")
	//command.DockerBuild("../../gin", "project")
	//message := ""
	//message = message + command.DockerBuild("../../gin", "project")
	//m := &Page{Title: "Build", Body: []byte(message)}
	//err := m.save()
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	http.Redirect(w, r, "/view/", http.StatusFound)
	//_ = command.DockerRun("../../gin", "6060", "test", "project")

}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
	p := content.Read("status")
	renderTemplate(w, "status", p)

}

func main() {
	go content.Change()
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/build/", buildHandler)
	http.HandleFunc("/status/", statusHandler)
	log.Fatal(http.ListenAndServe(":7070", nil))
}
