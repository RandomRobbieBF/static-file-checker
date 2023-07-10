package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type StaticFiles struct {
	Paths   map[string]string `json:"paths"`
	Version string            `json:"version"`
}

func main() {
	// Define command-line flags
	urlFlag := flag.String("url", "", "URL of /static/staticfiles.json")
	flag.Parse()

	// Check if the URL flag is provided
	if *urlFlag == "" {
		fmt.Println("Please provide the URL flag.")
		flag.PrintDefaults()
		return
	}

	url := *urlFlag

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", ""+url+"/static/staticfiles.json", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Parse the JSON
	var staticFiles StaticFiles
	err = json.Unmarshal(bodyBytes, &staticFiles)
	if err != nil {
		fmt.Printf("Failed to parse the JSON: %v\n", err)
		return
	}

	// Create or truncate the file
	file, err := os.Create("tobechecked.txt")
	if err != nil {
		fmt.Printf("Failed to create the output file: %v\n", err)
		return
	}
	defer file.Close()

	// Write URLs to the file
	for _, filename := range staticFiles.Paths {
	if strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".gif") || strings.HasSuffix(filename, ".svg") || strings.HasSuffix(filename, ".woff2"){
		fmt.Println("Ignoring Image files.")
	} else {
		url := fmt.Sprintf("%s/static/%s", url, filename)
		file.WriteString(url + "\n")
	}
		
	}

	// Run the OS command
	cmd := exec.Command("nuclei", "-l", "tobechecked.txt", "-t", "http/exposures/tokens", "-o", "url.txt")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to execute the OS command: %v\n", err)
		return
	}

	fmt.Println("Script execution completed.")
}
