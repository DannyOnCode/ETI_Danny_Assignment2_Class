package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
	"encoding/json"
	"io/ioutil"
	"strings"
	"strconv"
	"bytes"

	"github.com/gorilla/mux"
)

// Initialising api routes
const classURL = "http://localhost:8041/api/v1/classes"

type User struct {
	UserID string `json:"UserID"`
}

var currentUserInfo User
var currentDate = time.Now()
var daysUntilMon = (1 - int(currentDate.Weekday()) + 7) % 7
var currentSemStartDate = currentDate.AddDate(0, 0, -(7 - daysUntilMon)).Format("02-01-2006")
var nextSemStartDate = currentDate.AddDate(0, 0, daysUntilMon).Format("02-01-2006")
var nextMon = currentDate.AddDate(0, 0, daysUntilMon).Format("02 Jan 2006")

type Class struct {
	ClassCode string `json:"ClassCode"`
	Schedule string `json:"Schedule"`
	Tutor    string `json:"Tutor"`
	Capacity int32 `json:"Capacity"`
	Students []string `json:"Students"`
}

type Module struct {
	ModuleCode string
	ModuleName string
	ModuleClasses []Class
}

type Semester struct {
	SemesterStartDate string
	SemesterModules []Module
}

func login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("tempLogin.html"))

	currentUserInfo = User{}

	details := User{
		UserID: r.FormValue("userid"),
	}

	if details.UserID != "" {
		currentUserInfo = details
		if currentUserInfo.UserID[0:1] == "T"{
			// Go to Tutor Page
			http.Redirect(w, r, "/TutorClassPage", http.StatusFound)
		}else{
			// Go to Student Page
			http.Redirect(w, r, "/StudentClassPage", http.StatusFound)
		}
	}

	if currentUserInfo.UserID != "" {
		if currentUserInfo.UserID[0:1] == "T"{
			// Go to Tutor Page
			http.Redirect(w, r, "/TutorClassPage", http.StatusFound)
		}else{
			// Go to Student Page
			http.Redirect(w, r, "/StudentClassPage", http.StatusFound)
		}
	}
	tmpl.Execute(w, nil)
}

func studentMain(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("studentClassMainPage.html"))
	if r.Method == http.MethodPost {
		if r.FormValue("semester") != ""{
			http.Redirect(w, r, "/viewClass/" + r.FormValue("classcode") + "?semester_start_date=" + r.FormValue("semester"), http.StatusFound)
		}else{
			http.Redirect(w, r, "/viewClass/" + r.FormValue("classcode"), http.StatusFound)
		}
		
		tmpl.Execute(w, nil)
		return
	}
	var url string
	url = classURL + "/" + currentSemStartDate

	var currentSemesterInfo Semester
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		semData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(semData))

		json.Unmarshal(semData, &currentSemesterInfo)
		fmt.Println(currentSemesterInfo)
		response.Body.Close()
	}

	data := map[string]interface{}{
		"UserID": currentUserInfo.UserID,
		"NextMon": nextMon,
		"SemInfo": currentSemesterInfo.SemesterModules,
		"CurrentSemesterStartDate" : currentSemStartDate,
	}

	tmpl.Execute(w, data)
}

func tutorMain(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("tutorClassMainPage.html"))
	var url string
	url = classURL + "/" + currentSemStartDate

	var currentSemesterInfo Semester
	response, err := http.Get(url)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		semData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(semData))

		json.Unmarshal(semData, &currentSemesterInfo)
		fmt.Println(currentSemesterInfo)
		response.Body.Close()
	}

	data := map[string]interface{}{
		"UserID": currentUserInfo.UserID,
		"NextMon": nextMon,
		"SemInfo": currentSemesterInfo.SemesterModules,
		"CurrentSemesterStartDate" : currentSemStartDate,
	}

	tmpl.Execute(w, data)
}

func viewClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	v := r.URL.Query()
	requestSemesterDate := v.Get("semester_start_date")
	requestedClassCode := params["classCode"]
	requestedModuleCode := strings.Split(requestedClassCode, "_")[0]
	var url string
	if requestSemesterDate != ""{
		// Send Url with semester date
		url = classURL + "/" + requestSemesterDate + "?moduleCode=" + requestedModuleCode + "&classCode=" + requestedClassCode
	}else{
		// Send Url with current sem date
		requestSemesterDate = currentSemStartDate
		url = classURL + "/" + requestSemesterDate + "?moduleCode=" + requestedModuleCode + "&classCode=" + requestedClassCode
	}
	var receivedClassDetails Class
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		classData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(classData))

		json.Unmarshal(classData, &receivedClassDetails)
		fmt.Println(receivedClassDetails)
		response.Body.Close()
	}

	tmpl := template.Must(template.ParseFiles("classDetails.html"))

	data := map[string]interface{}{
		"Class": receivedClassDetails,
		"CurrentSemesterStartDate" : requestSemesterDate,
	}

	tmpl.Execute(w, data)
}


func editClass(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("editClass.html"))
	if r.Method == http.MethodPost {
		convertCapacity, _ := strconv.ParseInt(r.FormValue("capacity"),10,32)
		newCapacity := int32(convertCapacity)
		details := Class{
			ClassCode: r.FormValue("classcode"),
			Schedule: r.FormValue("schedule"),
			Tutor: r.FormValue("tutor"),
			Capacity: newCapacity,
		}

		params := mux.Vars(r)
		v := r.URL.Query()
		requestedClassCode := params["classCode"]
		requestSemesterDate := v.Get("semester_start_date")

		if currentUserInfo.UserID[0:1] != "T"{
			fmt.Println("Youre not a tutor why are you here?")
		}else{
			url := classURL + "/" + requestSemesterDate + "?moduleCode=" + r.FormValue("modulecode") + "&classCode=" + requestedClassCode
			jsonValue, _ := json.Marshal(details)
			fmt.Println(url)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))

			request.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			response, err := client.Do(request)

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))
				response.Body.Close()
			}
		}

		tmpl.Execute(w, nil)
		return
	}

	params := mux.Vars(r)
	v := r.URL.Query()
	requestSemesterDate := v.Get("semester_start_date")
	requestedClassCode := params["classCode"]
	requestedModuleCode := strings.Split(requestedClassCode, "_")[0]
	var url string
	if requestSemesterDate != ""{
		// Send Url with semester date
		url = classURL + "/" + requestSemesterDate + "?moduleCode=" + requestedModuleCode + "&classCode=" + requestedClassCode
	}else{
		// Send Url with current sem date
		url = classURL + "/" + currentSemStartDate + "?moduleCode=" + requestedModuleCode + "&classCode=" + requestedClassCode
	}
	var receivedClassDetails Class
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		classData, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(classData))

		json.Unmarshal(classData, &receivedClassDetails)
		fmt.Println(receivedClassDetails)
		response.Body.Close()
	}

	tmpl.Execute(w, receivedClassDetails)
}

func createClass(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("createClass.html"))
	if r.Method == http.MethodPost {
		convertCapacity, _ := strconv.ParseInt(r.FormValue("capacity"),10,32)
		newCapacity := int32(convertCapacity)
		details := Class{
			ClassCode: r.FormValue("classcode"),
			Schedule: "",
			Tutor: currentUserInfo.UserID,
			Capacity: newCapacity,
		}
		if currentUserInfo.UserID[0:1] != "T"{
			fmt.Println("Youre not a tutor why are you here?")
		}else{
			url := classURL + "/" + currentSemStartDate + "?moduleCode=" + r.FormValue("modulecode") + "&classCode=" + details.ClassCode
			jsonValue, _ := json.Marshal(details)
			response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

			if err != nil {
				fmt.Printf("The HTTP request failed with error %s\n", err)
			} else {
				data, _ := ioutil.ReadAll(response.Body)
				fmt.Println(response.StatusCode)
				fmt.Println(string(data))
				response.Body.Close()
			}
		}

		tmpl.Execute(w, nil)
		return
	}

	tmpl.Execute(w, nil)
}

func deleteClass(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	v := r.URL.Query()
	requestSemesterDate := v.Get("semester_start_date")
	requestedClassCode := params["classCode"]
	requestedModuleCode := strings.Split(requestedClassCode, "_")[0]

	tmpl := template.Must(template.ParseFiles("deletedPage.html"))
	url := classURL + "/" + requestSemesterDate + "?moduleCode=" + requestedModuleCode + "&classCode=" + requestedClassCode

	request, err := http.NewRequest(http.MethodDelete, url, nil)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(response.StatusCode)
		fmt.Println(string(data))
		response.Body.Close()
	}
	tmpl.Execute(w, nil)
	return
}


func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", login)
	router.HandleFunc("/StudentClassPage", studentMain)
	router.HandleFunc("/TutorClassPage", tutorMain)
	router.HandleFunc("/viewClass/{classCode}", viewClass)
	router.HandleFunc("/editClass/{classCode}", editClass)
	router.HandleFunc("/createClass",createClass)
	router.HandleFunc("/deleteClass/{classCode}", deleteClass)
	fmt.Println("Listening on port 8040")
	http.ListenAndServe(":8040", router)
}
