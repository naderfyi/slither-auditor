# NB-Auditor

This is a Go program that uses the Gin framework to handle HTTP requests and file uploads. When a user makes a POST request to the "/uploadfile" endpoint, the program reads the uploaded file, saves it to the "files" directory, then calls the getCompilerVersion() and runSlither() functions with the file path. 

The `getCompilerVersion()` function reads the contents of the file and uses a regular expression to match and extract the compiler version. The `runSlither()` function runs the Slither tool on the file, which is a tool for analyzing smart contracts written in Solidity. 

The output of the `runSlither()` function is then returned to the user as a JSON response. The program also deletes the uploaded file after it is done with it.

## Prerequisites
* Go
* Gin
* Slither

## Running the program
1. Clone the repository
2. Make sure you have all the prerequisites installed
3. Run the command `make run`
4. Make a POST request to `http://localhost:8080/uploadfile` with a file attached
5. The output of the slither analysis will be returned to you as a JSON response

## Built With
* [Go](https://golang.org/)
* [Gin](https://github.com/gin-gonic/gin)
* [Slither](https://github.com/crytic/slither)
