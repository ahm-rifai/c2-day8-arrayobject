package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	route := mux.NewRouter()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/project-detail/{index}", projectDetail).Methods("GET")
	
	route.HandleFunc("/add-project", addProject).Methods("GET")	
	route.HandleFunc("/create-project", createProject).Methods("POST") // CREATE PROJECT
	
	route.HandleFunc("/edit-project/{index}", editProject).Methods("GET") // EDIT PROJECT

	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET") // DELETE PROJECT

	fmt.Println("Server running on port 8000")
	http.ListenAndServe("localhost:8000", route)

}


type Project struct {
	ProjectName string
	StartDate string
	EndDate string
	Description string
	Duration string
	Technologies []string
}

var dataProject = []Project {
	{
		ProjectName:  "Grow Your Business With Mobile",
		StartDate:    "20 October 2022",
		EndDate:      "01 November 2022",
		Duration:     "3 Months",
		Description:  "mobile app proses pengembangan aplikasi yang dibuat untuk perangkat genggam, atau yang lebih dikenal dengan smartphone.",
		Technologies: []string{"node", "react", "next", "type"},
	},
	{
		ProjectName:  "Portofolio",
		StartDate:    "31 October 2022",
		EndDate:      "29 November 2022",
		Duration:     "2 Months",
		Description:  "Hi, I'm Rifa'i. A student of S1 Informatics Engineering University of Raharja, who focuses on Frontend Development. Often use HTML5, CSS3, JavaScript(ES6), Bootstrap, Node JS, React Native, and Golang in some projects. I also have communication skills, design thinking, and problem solving. I am currently attending a coding bootcamp at DumbWays Indonesia.",
		Technologies: []string{"node", "react", "next", "type"},
	},
}



func home(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	response := map[string]interface{}{
		"DataProject": dataProject,
	}

	tmpl.Execute(w, response)
}

func contact(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func projectDetail(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/project-detail.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var RenderProjectDetail = Project{}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	
	for i, data := range dataProject {
		if i == index {
			RenderProjectDetail = Project{
				ProjectName: data.ProjectName,
				StartDate: data.StartDate,
				EndDate: data.EndDate,
				Duration: data.Duration,
				Description: data.Description,
				Technologies: data.Technologies,
			}
		}
	}

	data := map[string]interface{}{
		"RenderProjectDetail": RenderProjectDetail,
	}

	// fmt.Println(data)
	tmpl.Execute(w, data)
}

func addProject(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/add-project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	tmpl.Execute(w, nil)
}

func createProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	projectName := r.PostForm.Get("input-project")
	startDate  := r.PostForm.Get("input-start")
	endDate  := r.PostForm.Get("input-end")
	description := r.PostForm.Get("input-desc")
	iconNodeJS := r.PostForm.Get("node")
	iconReactJS := r.PostForm.Get("react")
	iconNextJS := r.PostForm.Get("next")
	iconTypescript := r.PostForm.Get("type")
	

	
	var newProject = Project {
		ProjectName: projectName,
		StartDate: formatDate(startDate),
		EndDate: formatDate(endDate),
		Duration: getDuration(startDate, endDate),
		Description: description,
		Technologies: []string{iconNodeJS, iconReactJS, iconNextJS, iconTypescript},
	}

	dataProject = append(dataProject, newProject)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func editProject(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "Text/html; charset=utp-8")
	var tmpl, err = template.ParseFiles("views/edit-project.html")

	if err != nil {
		w.Write([]byte("Message : " + err.Error()))
		return
	}

	var updateProject = Project{}
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	for i, data := range dataProject {
		if index == i {
			updateProject = Project {
				ProjectName: data.ProjectName,
				StartDate: returnDate(data.StartDate),
				EndDate: returnDate(data.EndDate),
				Description: data.Description,
				Technologies: data.Technologies,
			}
			dataProject = append(dataProject[:index], dataProject[index+1:]...)
		}
	}

	data := map[string]interface{} {
		"updateProject" : updateProject,
	}

	tmpl.Execute(w, data)

}

func deleteProject(w http.ResponseWriter, r *http.Request)  {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	dataProject = append(dataProject[:index], dataProject[index+1:]...)

	http.Redirect(w, r, "/", http.StatusFound)
}











func getDuration(startDate string, endDate string) string {

	layout := "2006-01-02"

	projectPost, _ := time.Parse(layout, startDate)
	currenTime, _ := time.Parse(layout, endDate)

	distance := currenTime.Sub(projectPost).Hours() / 24
	var duration string

	if distance > 30 {
		if (distance / 30) <= 1 {
			duration = "1 Month"
		} else {
			duration = strconv.Itoa(int(distance)/30) + " Months"
		}
	} else {
		if distance <= 1 {
			duration = "1 Day"
		} else {
			duration = strconv.Itoa(int(distance)) + " Days"
		}
	}

	return duration
}

func formatDate(InputDate string) string {

	layout := "2006-01-02"
	t, _ := time.Parse(layout, InputDate)

	formated := t.Format("02 January 2006")

	return formated
}

func returnDate(InputDate string) string {

	layout := "02 January 2006"
	t, _ := time.Parse(layout, InputDate)

	formated := t.Format("2006-01-02")

	return formated
}
