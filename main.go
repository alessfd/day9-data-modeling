package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"personal-web/connection"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addproject", addProject).Methods("GET")
	route.HandleFunc("/submitproject", submitProject).Methods("POST")
	route.HandleFunc("/editproject/{id}", editProject).Methods("GET")
	route.HandleFunc("/submitedit", submitEdit).Methods("POST")
	route.HandleFunc("/deleteproject/{id}", deleteProject).Methods("GET")

	port := "5000"

	fmt.Print("Server sedang berjalan di port " + port + "\n")
	http.ListenAndServe("localhost:"+port, route)
}

// Home
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf=8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title, content, start_date, end_date FROM tb_project")

	if errQuery != nil {
		w.Write([]byte("Message: " + errQuery.Error()))
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.Id, &each.Title, &each.Content, &each.Start_date, &each.End_date)

		if err != nil {
			w.Write([]byte("Message: " + err.Error()))
			return
		}

		each.Duration = ProjectDuration(each.Start_date, each.End_date)

		result = append(result, each)
	}

	// fmt.Println(result)

	data := map[string]interface{}{
		"Projects": result,
	}

	tmpt.Execute(w, data)
}

// Contact
func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf=8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

// Add Project
func addProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf=8")
	tmpt, err := template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

// Edit Project
func editProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf=8")
	tmpt, err := template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte("Message: " + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

// Project struct
type Project struct {
	Id         int
	Title      string
	Content    string
	Image      string
	Start_date time.Time
	End_date   time.Time
	Tech       string
	Duration   string
}

var s1, _ = time.Parse(timeLayout, "2022-11-25")
var s2, _ = time.Parse(timeLayout, "2022-12-02")

// var projects = []
var projects = []Project{
	{
		Title:      "Judul",
		Content:    "Halo Dumbways",
		Start_date: s1,
		End_date:   s2,
	},
}

// Project Form Submit
func submitProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	var title = r.PostForm.Get("title")
	var content = r.PostForm.Get("content")

	var newProject = Project{
		Title:   title,
		Content: content,
	}

	// fmt.Println(
	// 	"Title: "+r.PostForm.Get("title"),
	// 	"\nContent: "+r.PostForm.Get("content"),
	// 	"\nDate Start: "+r.PostForm.Get("datestart"),
	// 	"\nDate Start: "+r.PostForm.Get("datestart"),
	// 	"\nTechnologies: ",
	// 	"\n Node Js: "+r.PostForm.Get("nodejs"),
	// 	"\n React Js: "+r.PostForm.Get("reactjs"),
	// 	"\n Next Js: "+r.PostForm.Get("nextjs"),
	// 	"\n TypeScript: "+r.PostForm.Get("typescript"),
	// )

	// projects.push(newProject)
	projects = append(projects, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Submit Edit
func submitEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var title = r.PostForm.Get("title")
	var content = r.PostForm.Get("content")

	projects[id].Title = title
	projects[id].Content = content

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

// Project Detail
func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf=8")
	tmpt, err := template.ParseFiles("views/project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	// Id
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var projectInfo = Project{}

	for index, data := range projects {
		if index == id {
			projectInfo = Project{
				Title:    data.Title,
				Content:  data.Content,
				Duration: data.Duration,
			}
		}
	}

	dataDetail := map[string]interface{}{
		"Project": projectInfo,
	}

	tmpt.Execute(w, dataDetail)
}

// Delete Project
func deleteProject(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// fmt.Println(id)

	projects = append(projects[:id], projects[id+1:]...)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

var timeLayout = "2006-01-02"

// Duration
func ProjectDuration(datestart time.Time, dateend time.Time) string {
	var duration string
	// var t1, _ = time.Parse(timeLayout, datestart)
	// var t2, _ = time.Parse(timeLayout, dateend)

	days := dateend.Sub(datestart).Hours() / 24
	weeks := days / 7
	months := days / 30
	years := months / 12

	if int(years) > 0 {
		duration = strconv.Itoa(int(years)) + " tahun"
	} else if int(months) > 0 {
		duration = strconv.Itoa(int(months)) + " bulan"
	} else if int(weeks) > 0 {
		duration = strconv.Itoa(int(weeks)) + " minggu"
	} else {
		duration = strconv.Itoa(int(years)) + " hari"
	}
	return duration
}
