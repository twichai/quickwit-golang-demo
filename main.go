package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	jsonData := []byte(`
	[
		{"message":"hello world","level":"info","timestamp":"2026-02-25T10:00:00Z"},
		{"message":"error happened","level":"error","timestamp":"2026-02-25T10:01:00Z"}
	]`)

	url := "http://localhost:7280/api/v1/logs/ingest"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status:", resp.Status)
}
