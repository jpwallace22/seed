#!/usr/bin/env bash
set -euo pipefail

# Add debug mode
if [[ "${TRACE-0}" == "1" ]]; then
  set -x
fi

# Archives benchmark results with an incremented number
archive_benchmark() {
  local current_file="./benchmark/benchmark_results.txt"
  local historical_dir="./benchmark/historical"

  echo "Checking for current file: $current_file"
  if [ ! -f "$current_file" ]; then
    echo "No benchmark results file found at $current_file"
    return 0 # This might be the issue - returning 0 when file doesn't exist
  fi

  # Create historical directory if it doesn't exist
  echo "Checking historical directory: $historical_dir"
  if [ ! -d "$historical_dir" ]; then
    echo "Creating historical directory at $historical_dir"
    mkdir -p "$historical_dir"
  fi

  # Find the next number - add error checking
  echo "Finding next number..."
  next_num=$(ls -1 "$historical_dir" 2>/dev/null | grep -E '^[0-9]+_benchmark_results\.txt$' | sed 's/_.*$//' | sort -n | tail -1 || echo "0")

  if [ -z "$next_num" ]; then
    next_num=0
    echo "No existing benchmark archives found, starting at 1"
  else
    echo "Last benchmark archive number: $next_num"
  fi
  next_num=$((next_num + 1))

  # Archive the file
  local archive_file="$historical_dir/${next_num}_benchmark_results.txt"
  mv "$current_file" "$archive_file" || {
    echo "Move failed with error code $?"
    return 1
  }

  if [ -f "$archive_file" ]; then
    echo "Successfully archived benchmark results to $archive_file"
    return 0
  else
    echo "ERROR: Failed to archive benchmark results"
    return 1
  fi
}

# Run the archiving function
archive_benchmark || {
  echo "Archive function failed with error code $?"
  exit 1
}
