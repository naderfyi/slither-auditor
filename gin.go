package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type Detector struct {
	Check       string `json:"Check"`
	Description string `json:"Description"`
	Impact      string `json:"Impact"`
	Confidence  string `json:"Confidence"`
}

func main() {
	// Create Gin engine
	r := gin.Default()

	// Register handlers
	r.POST("/uploadfile", uploadFile)

	r.Static("/", "template")
	// Start server
	r.Run(":8080")
}

func uploadFile(c *gin.Context) {
	// Maximum upload of 10 MB
	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["selectedFile"]

	for _, file := range files {
		log.Printf("Uploaded File: %+v\n", file.Filename)
		log.Printf("File Size: %+v\n", file.Size)
		log.Printf("MIME Header: %+v\n", file.Header)

		src, err := file.Open()
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("open file err: %s", err.Error()))
			return
		}
		defer src.Close()

		dstFilePath := "./files/" + file.Filename
		dst, err := os.Create(dstFilePath)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("create file err: %s", err.Error()))
			return
		}
		if _, err := io.Copy(dst, src); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("copy file err: %s", err.Error()))
			return
		}

		if err := dst.Close(); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("close file err: %s", err.Error()))
			return
		}

		// Call the function with the file path
		compilerVersion, err := getCompilerVersion(dstFilePath)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(">>>>>>>> COMPILER VERSION: ", compilerVersion)
		}

		// Call the function with the file path
		output, err := runSlither(dstFilePath, compilerVersion)
		if err != nil {
			fmt.Println(err)
		}

		// Parse the output as a string
		results := string(output)

		// Delete the uploaded file
		if err := os.Remove(dstFilePath); err != nil {
			log.Printf("Error deleting file: %v", err)
		}

		//c.JSON(http.StatusOK, results)

		jsonIssues, err := getIssues(results)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error getting issues: %s", err.Error()))
			return
		}
		//println(string(jsonIssues))
		//c.JSON(http.StatusOK, jsonIssues)
		var result map[string][]Detector
		json.Unmarshal(jsonIssues, &result)

		response := map[string]interface{}{
			"compilerVersion": compilerVersion,
			"issues":          result,
		}
		c.JSON(http.StatusOK, response)
	}
}

func getCompilerVersion(filePath string) (string, error) {
	// Read the entire file into a single string
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("Error reading file: %s", err)
	}
	fileString := string(fileBytes)

	// Define the regular expression to match the compiler version
	versionRegex := regexp.MustCompile(`pragma solidity\s*\^?([0-9]+\.[0-9]+\.[0-9]+)`)

	// Check if the file matches the regular expression
	match := versionRegex.FindStringSubmatch(fileString)
	if match != nil {
		// Extract the compiler version from the match
		compilerVersion := match[1]
		return compilerVersion, nil
	}

	return "", fmt.Errorf("Error: Compiler version not found in file")
}

func runSlither(filePath string, compilerVersion string) (string, error) {
	scriptName := "./run_slither.sh"
	// Define the command to run the shell script
	cmd := exec.Command(scriptName, filePath, compilerVersion)

	// Execute the command and capture the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		if !strings.Contains(string(output), "Switched global version to") {
			return "", fmt.Errorf("Error running slither:\n %s", output)
		}
	}

	// Parse the output as a string
	outputStr := string(output)

	return outputStr, nil
}

func getIssues(output string) ([]byte, error) {
	detectors := extractSlitherDetectors()

	issues := make(map[string][]Detector)
	for _, detector := range detectors {
		check := detector.Check
		if _, ok := issues[check]; !ok {
			issues[check] = []Detector{}
		}
		if strings.Contains(output, check) {
			issues[check] = append(issues[check], detector)
		}
	}

	jsonIssues, err := json.Marshal(issues)
	if err != nil {
		return nil, err
	}

	return jsonIssues, nil
}

func extractSlitherDetectors() []Detector {
	detectors := []Detector{}

	output, err := exec.Command("slither", "--list-detectors").Output()
	if err != nil {
		fmt.Println("Error running slither --list-detectors")
	}
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	for i, line := range lines {
		if i < 3 || i > len(lines)-3 {
			continue
		}

		parts := strings.Fields(line)
		oneline := strings.Join(parts[2:len(parts)-1], " ")
		// Split the line by "|"
		words := strings.Split(oneline, "|")

		// Trim whitespace from the parts
		for i := range words {
			words[i] = strings.TrimSpace(words[i])
		}

		check := words[1]
		description := words[2]
		impact := words[3]
		confidence := words[4]

		detectors = append(detectors, Detector{
			Check:       check,
			Description: description,
			Impact:      impact,
			Confidence:  confidence,
		})
	}
	return detectors
}
