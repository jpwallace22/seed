#!/usr/bin/env bash
set -euo pipefail

# Add debug mode
if [[ "${TRACE-0}" == "1" ]]; then
  set -x
fi

historical_dir="./benchmark/historical"
current_file="./benchmark/benchmark_results.txt"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Benchmark Comparison Tool${NC}"
echo "--------------------------------"

if [ ! -f "$current_file" ]; then
  echo -e "${RED}Error: No current benchmark file found at $current_file${NC}"
  exit 1
fi

# Get the latest historical file
latest_historical=$(ls -1 "$historical_dir" | grep -E '^[0-9]+_benchmark_results\.txt$' | sort -n | tail -1)
if [ -z "$latest_historical" ]; then
  echo -e "${RED}Error: No historical benchmark file found to compare against${NC}"
  exit 1
fi

echo -e "${YELLOW}Comparing:${NC}"
echo "Old: $historical_dir/$latest_historical"
echo "New: $current_file"
echo "--------------------------------"

# Capture benchstat output for analysis
bench_output=$(benchstat \
  -alpha=0.05 \
  -confidence=0.95 \
  -format=text \
  "$historical_dir/$latest_historical" \
  "$current_file")

# Display the benchstat output
echo "$bench_output"

# Extract geomean lines
geomean_lines=$(echo "$bench_output" | grep "geomean" || true)

# Count improvements and regressions
improvements=$(echo "$geomean_lines" | grep -c '\-[0-9.]*%' || true)
regressions=$(echo "$geomean_lines" | grep -c '\+[0-9.]*%' || true)

echo -e "\n${YELLOW}Changes:${NC}"
if [ "$improvements" -gt 0 ] || [ "$regressions" -gt 0 ]; then
  # Extract and display geomean changes
  echo -e "${BLUE}Overall changes (geomean):${NC}"
  while IFS= read -r line; do
    if [[ $line =~ .*-[0-9.].*% ]]; then
      echo -e "${GREEN}→ $(echo $line | awk '{print $NF}') improvement${NC}"
    elif [[ $line =~ .*\+[0-9.].*% ]]; then
      echo -e "${RED}→ $(echo $line | awk '{print $NF}') regression${NC}"
    fi
  done <<<"$geomean_lines"
else
  echo "→ No significant changes detected"
fi

# Show timing details
echo -e "\n${YELLOW}Date Information:${NC}"
echo "Old benchmark: $(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "$historical_dir/$latest_historical")"
echo "New benchmark: $(stat -f "%Sm" -t "%Y-%m-%d %H:%M:%S" "$current_file")"

# Add notes about interpretation
echo -e "\n${BLUE}Notes:${NC}"
echo "* Delta shows percentage change from old to new"
echo "* p-value < 0.05 indicates statistically significant change"
echo "* Negative delta means improvement (faster/less memory)"
echo "* '~' indicates the difference is likely due to noise"
echo "* Use TRACE=1 for debug output"
