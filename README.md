## Table of Contents
1. [Introduction](#introduction)
2. [Design Consideration of Microservice](#consideration)
3. [Microservice APIs of Application : Class](#microservices_explain)
4. [Instructions to setting up and running of microservices](#instructions)


## Introduction  <a name="introduction"></a>
Good Day, My name is Danny Chan currently a year 3 student studying in Information Technology from Ngee Ann Polytechnic

This assignment is about an education and financial application which has been split into several microservices. The package that I have been assigned with is the **Management of Classes**
This readme first covers the design considerations that has been made, which will go through the architecture diagram that are connected to the package that I am working on. Following that, the microservice, class, will be explained along with all the API calls. To end off, instructions on how to set up and deploy the database, frontend, and backend api will be mentioned.

This repository contains the source code to the web frontend to call the API as well as the backend API code.


The task breakdown for the package that I have done is as follows:

- 3.8.	Management of classes
    - 3.8.1.	Create, update, delete classes. Info includes
        - 3.8.1.1.	Class code
        - 3.8.1.2.	Schedule of the class
        - 3.8.1.3.	Capacity of class
    - 3.8.2.	View class info and ratings
    - 3.8.3.	List all classes
    - 3.8.4.	Search for classes
    - 3.8.5.	List all students in a class

## Design Consideration of Microservice <a name="consideration"></a>

### Architecture diagram
Management of class will consist of a web service (a frontend), 1 microservice, 
as well as one database that will be used by 3 connecting microservices, from the package 3.14 Bidding Dashboard and 3.15 Timetable. 
However in the diagram below, only the main microservice connection will be shown.

![Architecture of Application](/images/Architecture.png)<br>

To quickly go through the process flow. 
As seen from the diagram, the user will interact with the frontend and the frontend will send the request to the backend to get the data.
However, as other packages may require the use of data from the package that I develop, frontend is not strictly required and the API can be called directly through the curl command.

### Scalability
As class is the only microservice and there is only one database connected to the microservice, there will be no issue if the microservices needs to be scaled.

#### Database
Mongodb has been considered as my database as it allows for scalability and flexibility.

### Security
As for security, since the package that I have worked on is Class, not much security is necessary as the information is open to students as well as tutors. 
However, for creation, deletion and editing of the classes. ID has been checked to ensure the user editing the class is a tutor.

## Microservice APIs of Application : Class <a name="microservices_explain"></a>
Under this section, I will be discussing the class microservice and the resources that it provides along with the routes to access the resources.

**Base URL :**
```
http://10.31.11.11:8041/api/v1/classes
```
This is the base URL to the API that will needed before any specifications.
 

**Get all classes of certain semester:**
```
GET http://10.31.11.11:8041/api/v1/classes/{semester_start_date}
```  
This is a GET request route that will look into the database and fetch all the classes under the specified semester start date<br>
Format of semester_start_date : 24-01-2022<br>
Returns: json array of all modules and classes under the module<br>

**Get specific class of certain semester:**
```
GET http://10.31.11.11:8041/api/v1/classes/{semester_start_date}?classCode=...
```
This is a GET request route that will look into the database and fetch the class specified under the specified semester start date<br>
Format of semester_start_date : 24-01-2022<br>
Format of classCode : IS_01 (ModuleCode_ClassNumber)<br>
Returns: json array of all modules and classes under the module<br>

**Add Semester, Modules and Empty Classes**
```
POST http://10.31.11.11:8041/api/v1/classes/{semester_start_date}
```
This is a POST request route that will be called automatically where it creates the shell with all the modules.<br>
This call will also call from *Myron API call* to retreive the available modules for the semester (not yet implemented)<br>

**Add/Update/Delete of certain Class**
```
POST http://10.31.11.11:8041/api/v1/classes/{semester_start_date}?moduleCode=...&classCode=... \ 
--header 'Content-Type: application/json' \
--data '{
    "ClassCode":"...",
    "Schedule": "...", 
    "Tutor": "...", 
    "Capacity": "...", 
    "Students":["...","...","..."]
    }'
```
This is a POST request route that will add a new class under the module and semester specified<br>
Format of semester_start_date : 24-01-2022<br>
Format of classCode : IS_01 (ModuleCode_ClassNumber)<br>
Format of moduleCode : IS (ModuleCode)<br>

```
PUT http://10.31.11.11:8041/api/v1/classes/{semester_start_date}?moduleCode=...&classCode=... \ 
--header 'Content-Type: application/json' \
--data '{
    "ClassCode":"...",
    "Schedule": "...", 
    "Tutor": "...", 
    "Capacity": "...", 
    "Students":["...","...","..."]
    }'
```
This is a PUT request route that will edit the specified class under the module and semester specified<br>
Format of semester_start_date : 24-01-2022<br>
Format of classCode : IS_01 (ModuleCode_ClassNumber)<br>
Format of moduleCode : IS (ModuleCode)<br>

```
DELETE http://10.31.11.11:8041/api/v1/classes/{semester_start_date}?moduleCode=...&classCode=...
```
This is a DELETE request route that will delete the specified class under the module and semester specified<br>
Format of semester_start_date : 24-01-2022 (dd-mm-yyyy)<br>
Format of classCode : IS_01 (ModuleCode_ClassNumber)<br>
Format of moduleCode : IS (ModuleCode)<br>

**Base Front-End URL**
```
http://10.31.11.11:8040
```

**Access the temp login page**
```
http://10.31.11.11:8040
```

**Student Main page after login**
```
http://10.31.11.11:8040/StudentClassPage
```

**Tutor Main Page after login**
```
http://10.31.11.11:8040/TutorClassPage
```

**Create Class Page**
```
http://10.31.11.11:8040/createClass
```

**View Specific Class Page**
```
http://10.31.11.11:8040/viewClass/{classCode}
```

**Edit Class Page**
```
http://10.31.11.11:8040/editClass/{classCode}
```

**Delete Class Page**
```
http://10.31.11.11:8040/deleteClass/{classCode}
```


## Database Structure

semester_start
- module_code (string)
- moduleName (string)
- classes (array)
  - class_code (string)  
  - schedule (datetime)  
  - tutor (string)  
  - students (array)    
    - student_id (string)  
  - class_capacity (int)

The below image shows an example of the data in the database<br>

![Database Structure](/images/database_structure.png)<br>

## Link to your container image

Each service has been publicly published onto Docker Hub.

Front-end web view: https://hub.docker.com/repository/docker/nihilitydas/frontend_class

Class microservice: https://hub.docker.com/repository/docker/nihilitydas/microservice_class

MongoDB Service : https://hub.docker.com/repository/docker/nihilitydas/db_mongo_class

## Instructions for setting up and running your microservices <a name="instructions"></a>

After setting up the services, the applications would be hosted from http://10.31.11.11:8040

The Class Director will first do a curl command to the API to create the shell for all the classes
```
curl --request POST "http://10.31.11.11:8041/api/v1/classes/:semester_start_date" -H "Content-type:application/json"
```
This curl command will create the shell for all classes after retrieving the available classes from Myron API (Management of Modules)

**Automatic deployment**

The deployment is done through the `docker-compose.yml` file where it will automatically build/update and run the application containers.

Prerequisite:

- Downloaded git repository to local storage
- Docker Destop installed

Steps:

1. Open a command terminal and navigate to project ROOT directory under Danny/04
2. Run command `docker-compose up --build`

**Single Docker image run**<br>
Running class microservice
```
docker run --name class_microservice -dp 8041:8041 nihilitydas/microservice_class
```

Running class frontend
```
docker run --name class_frontend -dp 8040:8040 nihilitydas/frontend_class
```

Running class mongo_db
```
docker run --name mongo_db -dp 8049:27017 nihilitydas/db_mongo_class
```

**Manual deployment**

In order to run the project without docker, pull the project from the local branch under github and
run each of the following commands in seperate command prompts.

```
# run front-end application
go run web/main.go

# run add credits microservice
go run Class/main.go
```
