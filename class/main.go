package main

import (
	"context"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"time"
	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// Get all the classes in a semester and putting it into a struct
func GetSemesterClasses(collection *mongo.Collection, context context.Context, input_semester string) Semester {
	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var modules []bson.M
	if err = cursor.All(context, &modules); err != nil {
		log.Fatal(err)
	}
	var array_of_modules []Module
	for _, i_module := range modules {
		totalClasses := i_module["moduleClasses"].(primitive.A)
		var array_of_classes []Class
		for _, k_class := range totalClasses {
			var single_class Class
			bsonBytes, _ := bson.Marshal(k_class)
			bson.Unmarshal(bsonBytes,&single_class)
			array_of_classes = append(array_of_classes, single_class)
		}
		current_module := Module{
			ModuleCode: i_module["moduleCode"].(string),
			ModuleName: i_module["moduleName"].(string),
			ModuleClasses: array_of_classes,
		}
		array_of_modules = append(array_of_modules,current_module)
	}
	semester := Semester{
		SemesterStartDate : input_semester,
		SemesterModules : array_of_modules,
	}
	return semester
}


func GetSingleClass(collection *mongo.Collection, context context.Context, specifiedClassCode string) Class {
	cursor, err := collection.Find(context, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var modules []bson.M
	if err = cursor.All(context, &modules); err != nil {
		log.Fatal(err)
	}
	for _, i_module := range modules {
		totalClasses := i_module["moduleClasses"].(primitive.A)
		var array_of_classes []Class
		for _, k_class := range totalClasses {
			var single_class Class
			bsonBytes, _ := bson.Marshal(k_class)
			bson.Unmarshal(bsonBytes,&single_class)
			array_of_classes = append(array_of_classes, single_class)
			if single_class.ClassCode == specifiedClassCode {
				return single_class
			}
		}
	}
	return Class{}
}

// Adding a new collection to the database
func AddNewSemester(input_semester string, arr_of_modules map[string]string) {
	// TO CHANGE : when myron code is done and the receive format is confirm
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:rootpassword@mongo_db:27017"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    classDatabase := client.Database("classes")
    newSemester_startCollection := classDatabase.Collection(input_semester)

	for i, j := range arr_of_modules{
		newDocument, _ := newSemester_startCollection.InsertOne(ctx, bson.D{
			{Key: "moduleCode", Value: i},
			{Key: "moduleName", Value: j},
			{Key: "moduleClasses", Value: bson.A{}},
		})
		fmt.Println(newDocument.InsertedID)
	}
	return
}

// Used by POST and PUT
func UpdateOrInsertClassInSemester(collection *mongo.Collection, context context.Context, classToUpdate Class, requestModuleCode string, requestClassCode string) { 
	cursor, err := collection.Find(context, bson.M{"moduleClasses.classCode":requestClassCode})
	if err != nil {
		log.Fatal(err)
	}

	var foundSpecifiedClass []bson.M
	if err = cursor.All(context, &foundSpecifiedClass); err != nil {
		log.Fatal(err)
	}
	fmt.Println(foundSpecifiedClass)
	if foundSpecifiedClass != nil{
		retrieveClassInfo := GetSingleClass(collection,context,requestClassCode)
		if classToUpdate.Tutor == "" && retrieveClassInfo.Tutor != ""{
			classToUpdate.Tutor = retrieveClassInfo.Tutor
		}
		if classToUpdate.Schedule == "" && retrieveClassInfo.Schedule != ""{
			classToUpdate.Schedule = retrieveClassInfo.Schedule
		}
		if classToUpdate.Capacity == 0 && retrieveClassInfo.Capacity != 0{
			classToUpdate.Capacity = retrieveClassInfo.Capacity
		}
		if classToUpdate.Students == nil && retrieveClassInfo.Students != nil{
			classToUpdate.Students = retrieveClassInfo.Students
		}
		result, err := collection.UpdateOne(
			context,
			bson.M{"moduleClasses.classCode":requestClassCode},
			bson.D{
				{"$set", bson.M{"moduleClasses.$.schedule": classToUpdate.Schedule}}, 
				{"$set", bson.M{"moduleClasses.$.tutor": classToUpdate.Tutor}}, 
				{"$set", bson.M{"moduleClasses.$.capacity": classToUpdate.Capacity}},  
				{"$set", bson.M{"moduleClasses.$.students": classToUpdate.Students}},  
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	}else{ // Inserting new class into existing document
		result, err := collection.UpdateOne(
			context,
			bson.M{"moduleCode":requestModuleCode},
			bson.D{
				{"$push", bson.M{"moduleClasses": bson.D{
					{Key :"classCode", Value: classToUpdate.ClassCode},
					{Key :"schedule", Value: classToUpdate.Schedule},
					{Key :"tutor", Value: classToUpdate.Tutor},
					{Key :"capacity", Value: classToUpdate.Capacity},
					{Key :"students", Value: classToUpdate.Students},
				},	
				}},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
	}
	return
}

func RemoveClassFromSemester(collection *mongo.Collection, context context.Context,requestClassCode string){
	result, err := collection.UpdateOne(
		context,
		bson.M{"moduleClasses.classCode":requestClassCode},
		bson.D{
			{"$pull", bson.M{"moduleClasses": bson.M{"classCode" : requestClassCode}}}, 
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", result.ModifiedCount)
}

func classes(w http.ResponseWriter, r *http.Request) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:rootpassword@mongo_db:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	classDatabase := client.Database("classes")
	params := mux.Vars(r)

	if r.Method == "GET" {
		v := r.URL.Query()
		requestClassCode := v.Get("classCode")
		if requestClassCode != ""{
			// Get Specific Class
			semester_startCollection := classDatabase.Collection(params["semester_start_date"])
			retrievedClass := GetSingleClass(semester_startCollection, ctx, requestClassCode)
			json.NewEncoder(w).Encode(retrievedClass)
		}else{
			semester_startCollection := classDatabase.Collection(params["semester_start_date"])
			semesterClasses := GetSemesterClasses(semester_startCollection, ctx,params["semester_start_date"])
			json.NewEncoder(w).Encode(semesterClasses)
		}
		return
	}

	if r.Method == "DELETE"{
		v := r.URL.Query()
		requestClassCode := v.Get("classCode")
		if requestClassCode != ""{
			// Get Specific Class
			semester_startCollection := classDatabase.Collection(params["semester_start_date"])
			RemoveClassFromSemester(semester_startCollection, ctx, requestClassCode)
		}
		return
	}

	if r.Header.Get("Content-type") == "application/json" {
		var newClass Class
		v := r.URL.Query()
		if r.Method == "POST" {
			reqBody, _ := ioutil.ReadAll(r.Body)
			fmt.Println(string(reqBody))
			json.Unmarshal(reqBody, &newClass)
			fmt.Println(newClass)
			requestModuleCode := v.Get("moduleCode")
			requestClassCode := v.Get("classCode")
			fmt.Println(requestClassCode)
			if requestClassCode != "" && requestModuleCode != ""{
				if newClass.Tutor != "" && newClass.Capacity != 0 {
					semester_startCollection := classDatabase.Collection(params["semester_start_date"])
					UpdateOrInsertClassInSemester(semester_startCollection, ctx, newClass, requestModuleCode, requestClassCode)
				}
			}else{
				//arr_of_modules = get modules from myron api
				arr_of_modules := map[string]string{
					"CM" : "Computing Mathematics",
					"CSF" : "Cybersecurity Fundamentals",
					"DP" : "Design Principles",
					"PRG1" : "Programming 1",
					"DB" : "Databases",
					"ID" : "Interactive Development",
					"OSNF" : "Operating Systems and Networking Fundamentals",
					"PRG2" : "Programming 2",
					"OOAD" : "Object-Oriented Analysis and Design",
					"WEB" : "Web Application Development",
					"PFD" : "Portfolio Development",
					"SDD" : "Solution Design & Development",
				}
				AddNewSemester(params["semester_start_date"],arr_of_modules)
			}
		}

		if r.Method == "PUT" {
			reqBody, _ := ioutil.ReadAll(r.Body)
			fmt.Println(string(reqBody))
			json.Unmarshal(reqBody, &newClass)
			fmt.Println(newClass)
			requestModuleCode := v.Get("moduleCode")
			requestClassCode := v.Get("classCode")
			fmt.Println(requestClassCode)
			if requestClassCode != "" && requestModuleCode != ""{
				semester_startCollection := classDatabase.Collection(params["semester_start_date"])
				UpdateOrInsertClassInSemester(semester_startCollection, ctx, newClass, requestModuleCode, requestClassCode)
			}
		}

		return
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/classes/{semester_start_date}", classes).Methods(
		"GET", "PUT", "POST", "DELETE")

	fmt.Println("Listening at port 8041")
	log.Fatal(http.ListenAndServe(":8041", router))

	fmt.Println("Database opened")

}
