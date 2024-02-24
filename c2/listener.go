package c2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type HTTPListener struct {
	IP         string
	Port       string
	errorLog   *log.Logger
	requestLog *log.Logger
}

func (h *HTTPListener) taskHandler(w http.ResponseWriter, r *http.Request) {
	// when someone makes a request this function will be called
	h.requestLog.Println("Request Recv:", r.Method, r.URL.Path, r.RemoteAddr)
	task := map[string]any{
		"cmd":  "echo",
		"args": []string{"hello"},
	}

	// convert our task to json
	jsonData, err := json.Marshal(task)
	if err != nil {
		h.errorLog.Println("Error Marshalling data: ", err)
	}

	// send response back to user
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(jsonData)
}

// Listen functions our struct can do
func (h *HTTPListener) Listen() {
	// create logging dir
	err := os.Mkdir("c2/logs", 0755)
	if err != nil {
		fmt.Println("Error creating log directory: ", err)
	}

	logFile, err := os.OpenFile("c2/logs/listener.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	h.errorLog = log.New(logFile, "ERROR", log.LstdFlags)
	h.requestLog = log.New(logFile, "REQUEST", log.LstdFlags)

	fmt.Printf("Listening on: %s:%s\n", h.IP, h.Port)

	// create new server
	mux := http.NewServeMux()

	// register routes and handlers
	mux.HandleFunc("/task", h.taskHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(h.IP + ":" + h.Port),
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		h.errorLog.Println("Error starting server:", err)
	}
}
