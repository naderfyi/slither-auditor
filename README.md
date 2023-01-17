# NB-Auditor

This is a Go program that uses the Gin framework to handle HTTP requests and file uploads. When a user makes a POST request to the "/uploadfile" endpoint, the program reads the uploaded file, saves it to the "files" directory, then calls the getCompilerVersion() and runSlither() functions with the file path. 

The `getCompilerVersion()` function reads the contents of the file and uses a regular expression to match and extract the compiler version. The `runSlither()` function runs the Slither tool on the file, which is a tool for analyzing smart contracts written in Solidity. The output of the `runSlither()` function is then parsed and returned to the user as a JSON response containing the compiler version used and the issues found by the tool.

## Prerequisites
* Go
* Gin
* Slither
* Solc-select

## Running the program
1. Clone the repository
2. Make sure you have all the prerequisites installed
3. Run the command `make run` to start the server
4. Make a POST request to `http://localhost:8080/uploadfile` with a file attached
5. The output of the slither analysis, including the compiler version used and the issues found, will be returned to you as a JSON response

## Built With
* [Go](https://golang.org/)
* [Gin](https://github.com/gin-gonic/gin)
* [Slither](https://github.com/crytic/slither)
