package main

// Para usar el modulo compile daemon usar en el terminal $ CompileDaemon -command=".exe"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type alTask []task

var tasks = alTask{
	{
		ID:      1,
		Name:    "Task one",
		Content: "some content",
	},
}

// enviar tareas
func getTasks(w http.ResponseWriter, r *http.Request) {
	// Codificar el archivo a formato json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	// Obtener parametro de la ruta
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "ID invalid")
		return
	}

	// Buscar en la lista
	for _, task := range tasks {
		if task.ID == taskID {
			// Codificar el archivo a formato json
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	// Obtener parametro de la ruta
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])

	if err != nil {
		fmt.Fprintf(w, "ID invalid")
		return
	}

	// Buscar en la lista
	for index, task := range tasks {
		if task.ID == taskID {
			// eliminar un elemento del slice
			tasks = append(tasks[:index], tasks[index+1:]...)
			fmt.Fprintf(w, "The task with Id %v has been remove successfully", taskID)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	// Obtener parametro de la ruta
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updateTask task

	if err != nil {
		fmt.Fprintf(w, "ID invalid")
		return
	}

	// leer datos del Body
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Enter valid data")
		return
	}

	// Guardar los nuevos datos
	json.Unmarshal(reqBody, &updateTask)

	for index, task := range tasks {
		if task.ID == taskID {
			// elimino la tarea
			tasks = append(tasks[:index], tasks[index+1:]...)
			// guardar id en la nueva tarea
			updateTask.ID = taskID
			tasks = append(tasks, updateTask)
		}
	}
	fmt.Fprintf(w, "The task with Id %v has been updated successfully", taskID)
}

// crear una tarea
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Insert a valid task")
	}

	// Decodificar json y asignarlo a una nueva variable
	json.Unmarshal(reqBody, &newTask)

	// almaceno en el slices tasks
	newTask.ID = len(tasks) + 1
	tasks = append(tasks, newTask)

	// Respondo al cliente la tarea que se creo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the jungle!!!")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// rutas api
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks", createTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	// Start server
	log.Fatal(http.ListenAndServe(":3000", router))

}
