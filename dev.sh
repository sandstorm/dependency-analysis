#!/bin/bash
############################## DEV_SCRIPT_MARKER ##############################
# This script is used to document and run recurring tasks in development.     #
#                                                                             #
# You can run your tasks using the script `./dev some-task`.                  #
# You can install the Sandstorm Dev Script Runner and run your tasks from any #
# nested folder using `dev some-task`.                                        #
# https://github.com/sandstorm/Sandstorm.DevScriptRunner                      #
###############################################################################

source ./dev_utilities.sh

set -e

######### TASKS #########

# Starts the project from source
function run() {
  go run . $@
}

# Compiles the project into a local binary
function build() {
  go build .
  _log_green "done"
}

# Auto-formats all source files
function format() {
  go fmt main.go
  go fmt ./analysis
  go fmt ./cmd
  go fmt ./dataStructures
  go fmt ./parsing
  go fmt ./rendering
}

# runs all tests
function test() {
  go test ./analysis
  go test ./cmd
  go test ./dataStructures
  go test ./parsing
  go test ./rendering
  self_check
}

# runs the dependency analysis against this project
function self_check() {
  go run . validate .
  go run . visualize .
}

# Builds the project locall and creates a release
function release() {
  test
  docker build -t sandstormmedia/dependency-analysis:latest .
  goreleaser release --clean
  _log_green "done"
  _log_yellow "Manual steps:"
  _log_yellow " - docker tag sandstormmedia/dependency-analysis:latest sandstormmedia/dependency-analysis:version-goes-here"
  _log_yellow " - docker push sandstormmedia/dependency-analysis:latest"
  _log_yellow " - docker push sandstormmedia/dependency-analysis:the-version"
}

# Docs in the browser at localhost:8080
function docs() {
  which godoc || go install golang.org/x/tools/cmd/godoc@latest
  _log_green "Please open http://localhost:8080/pkg/github.com/sandstorm/dependency-analysis/"
  _log_green "Terminate with Ctrl + C"
  godoc -http=:8080
}

_log_green "---------------------------- RUNNING TASK: $1 ----------------------------"

# THIS NEEDS TO BE LAST!!!
# this will run your tasks
"$@"
