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
	t, err := template.ParseFiles("/home/mouyuan/go/src/github.com/mouyuan4598/microservice/" + tmpl + ".html")
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

	switch r.Method {
	case "GET":
		p := content.Read("change")
		fmt.Println(p.FileType)
		if len(p.Body) == 0 {
			renderTemplate(w, "view1", p)
		} else {
			renderTemplate(w, "view2", p)
		}
	default:
		fmt.Fprint(w, "Only GET method")
	}

}
func buildHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.Redirect(w, r, "/view/", http.StatusFound)
	case "POST":
		var err = os.Remove("status.txt")
		if err != nil {
			fmt.Println(err.Error())
		}
		p := content.Read("change")

		for key, value := range p.Body {
			// var b int
			// base := len("/home/mouyuan/go/src/github.com/mouyuan4598")
			// for i := base + 1; i < len(key); i++ {
			// 	//fmt.Print(i)
			// 	//fmt.Println(string(key[i]))
			// 	if key[i] == '/' {
			// 		b = i
			// 		break
			// 	}
			// 	if i == len(key)-1 {
			// 		b = i + 1
			// 	}
			// }
			//fmt.Println(base)
			//fmt.Println(b)

			folder := "/home/mouyuan/go/src/github.com/mouyuan4598/" + p.Directory[key]
			fmt.Println(folder)

			var temp string
			err = os.Remove(folder + "/Dockerfile")
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println(r.FormValue(p.Directory[key] + "buildName"))
			command.DockerCreate(folder, p.Directory[key], r.FormValue(p.Directory[key]+"buildName"), p.FileType[key], value)
			newDir := r.FormValue(p.Directory[key] + "imageName")
			temp = command.DockerBuild(folder, newDir)
			save("status.txt", temp)
			// temp = command.DockerSave(folder, newDir)
			// save("status.txt", temp)
			// // temp = command.DockerTransfer(folder, newDir, "149.28.11.6:~/dockerimage/")
			// // save("status.txt", temp)
			// save("status.txt", "\n\nFinished successfully")
		}
		err = os.Remove("change.txt")
		if err != nil {
			fmt.Println(err.Error())
		}

		/////////////////////////////////////////////////////////////////////
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

	default:
		fmt.Fprint(w, "Only GET or POST method.")
	}

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
