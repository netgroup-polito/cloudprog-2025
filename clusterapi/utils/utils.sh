#!/usr/bin/env bash

# Define the retry function
function waitandretry() {
  local waittime="$1"
  local retries="$2"
  local command="$3"
  local options="$-" # Get the current "set" options

  sleep "${waittime}"

  echo "Running command: ${command} (retries left: ${retries})"

  # Disable set -e
  if [[ $options == *e* ]]; then
    set +e
  fi

  # Run the command, and save the exit code
  $command
  local exit_code=$?

  # restore initial options
  if [[ $options == *e* ]]; then
    set -e
  fi

  # If the exit code is non-zero (i.e. command failed), and we have not
  # reached the maximum number of retries, run the command again
  if [[ $exit_code -ne 0 && $retries -gt 0 ]]; then
    waitandretry "$waittime" $((retries - 1)) "$command"
  else
    # Return the exit code from the command
    return $exit_code
  fi
}

function forgename() {
  local index=$1
  local cni=$2
  local podcidrtype=$3
  local template=$4
  local image=$5

  filename=$(basename -- "$template")
  template="${filename%.*}"

  if [[ $image == *ubuntu* ]]; then
    image="ubuntu"
  elif [[ $image == *rocky* ]]; then
    image="rocky"
  else
    image="unknown"
  fi

  echo "${template}-${image}-${cni}-${podcidrtype}-${index}"
}
