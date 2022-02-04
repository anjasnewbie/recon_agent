package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type CommandPayload struct {
	Command string `json:"command"`
	Timeout int    `json:"timeout"`
	Param   string `json:"param"`
	Authkey string `json:"authkey"`
}

func main() {
	// Hello world, the web server

	executeTask := func(w http.ResponseWriter, req *http.Request) {
		var commandPayload CommandPayload

		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.

		err := json.NewDecoder(req.Body).Decode(&commandPayload)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		fmt.Println(commandPayload.Command)

		if commandPayload.Authkey != "mekans" {
			io.WriteString(w, "invalid key "+commandPayload.Authkey)
			return
		}
		param := commandPayload.Param
		command := commandPayload.Command
		fmt.Println("Command " + command)
		fmt.Println("Param " + param)
		parameters := strings.Split(param, " ")
		fmt.Println(parameters[0])
		fmt.Println(parameters)
		output := executeCommand(command, 2, parameters...)
		io.WriteString(w, output)
	}

	http.HandleFunc("/execute", executeTask)
	log.Println("Listing for requests at http://localhost:8000/hello")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func executeCommand(commandLine string, timeout int, param ...string) string {

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel() // The cancel should be deferred so resources are cleaned up

	// Create the command with our context
	cmd := exec.CommandContext(ctx, commandLine, param...)

	// This time we can simply use Output() to get the result.
	out, err := cmd.Output()

	// We want to check the context error to see if the timeout was executed.
	// The error returned by cmd.Output() will be OS specific based on what
	// happens when a process is killed.
	if ctx.Err() == context.DeadlineExceeded {
		return "Command timed out"

	}

	// If there's no context error, we know the command completed (or errored).
	if err != nil {
		fmt.Println(err)
		return "Non-zero exit code:"
	}
	return string(out)

}
