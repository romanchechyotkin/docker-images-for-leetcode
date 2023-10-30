package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"docker-images-for-leetcode/docker"
)

func main() {
	http.HandleFunc("/submit", submit)
	http.HandleFunc("/task", task)
	log.Println("server running on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func submit(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	log.Printf(task)

	if task == "" {
		w.Write([]byte("what tasK?"))
		return
	}

	fs := os.DirFS(task)
	if fs == nil {
		fmt.Println("Failed")
		return
	}

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, r.Body)
	if err != nil {
		fmt.Println("failed to read build response:", err)
		os.Exit(1)
	}
	r.Body.Close()

	file, err := os.Create(task + "/math_operations.py")
	if err != nil {
		log.Println("failed to create file:", err)
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		log.Println("failed to write to file", err)
	}

	logsBuffer := docker.RunSubmit()

	w.Header().Set("Content-Type", "application/text")
	w.Write(logsBuffer.Bytes())
}

func task(w http.ResponseWriter, r *http.Request) {
	task := r.URL.Query().Get("task")
	log.Printf(task)

	file, err := os.Open(task + "/base.py")
	if err != nil {
		log.Println("failed to open file:", err)
		w.Write([]byte("no such task"))
		return
	}
	defer file.Close()

	bufferCode := make([]byte, 1024)
	n, err := file.Read(bufferCode)
	if err != nil {
		log.Println("failed to read file:", err)
	}

	w.Write(bufferCode[:n])
}
