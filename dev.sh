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
}

# Builds the project locall and creates a release
function release() {
  test
  goreleaser release --clean
  _log_green "done"
  _log_yellow "TODO: implement docker image releases"
}

_log_green "---------------------------- RUNNING TASK: $1 ----------------------------"

# THIS NEEDS TO BE LAST!!!
# this will run your tasks
"$@"
