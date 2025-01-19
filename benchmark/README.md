# Benchmark Results

This document details the performance characteristics of the tree parsing implementations, comparing ASCII tree and JSON parsing methods across various node counts and input methods.

The benchmark results are in `./benchmark_results.txt`

## Environment

- **OS**: Darwin (macOS)
- **Architecture**: ARM64
- **CPU**: Apple M3 Pro

## Methodology

Benchmarks were conducted using Go's built-in testing framework with the following parameters:

- 100, 500, 1000, and 5000 nodes (files and dirs)
- 2 second runs
- 3 runs per test
- Single core and quad core
- File and String input
- Metrics measured: 
  - Time (ns/op)
  - Memory allocation (B/op)
  - Allocation count (allocs/op)

## Key Findings

### Performance Comparison

#### Time performance

| Nodes | ASCII (ms) | JSON (ms) | Difference |
|-------|------------|-----------|------------|
| 100   | 8.62      | 9.00      | +4.4%      |
| 500   | 35.55     | 36.76     | +3.4%      |
| 1000  | 64.48     | 66.23     | +2.7%      |
| 5000  | 428.16    | 438.29    | +2.4%      |

#### Memory Usage

| Nodes | ASCII (KB) | JSON (KB) | Difference |
|-------|------------|-----------|------------|
| 100   | 13.89     | 13.95     | +0.4%      |
| 500   | 17.09     | 17.14     | +0.3%      |
| 1000  | 24.82     | 25.11     | +1.2%      |
| 5000  | 235.83    | 235.89    | +0.03%     |

### Input Method Comparison (500 nodes)

- **String Input**
  - Time: 35.82ms
  - Memory: 40.62KB
  - Allocations: 77/op

- **File Input**
  - Time: 35.57ms
  - Memory: 16.17KB
  - Allocations: 77/op

## Analysis

1. **Parsing Performance**
   - ASCII tree parsing consistently outperforms JSON parsing
   - Performance gap decreases with larger node counts
   - Both methods show linear scaling with node count

2. **Memory Efficiency**
   - Both parsers show similar memory patterns
   - JSON consistently uses marginally more memory
   - Differences in memory usage become negligible at larger node counts

3. **Input Methods**
   - File input shows significant memory advantages (~60% reduction)
   - No performance penalty for file-based input
   - Consistent allocation patterns across both input methods

## Practical Implications

1. **For Performance-Critical Applications**
   - ASCII tree parsing offers a slight but consistent performance advantage
   - Benefits are most noticeable with smaller node counts
   - Consider using file-based input for better memory efficiency

2. **For Memory-Constrained Systems**
   - File input method is strongly recommended
   - Both parsing methods show similar memory characteristics
   - Memory usage scales linearly with node count

3. **For General Usage**
   - Both methods are viable for typical use cases
   - Performance differences are minimal for most applications
   - Choose based on your specific needs for data format and interoperability

## Running the Benchmarks

To run the benchmarks yourself:

```bash
# Standard
make benchmark
# More intensive (time consuming)
make benchmark:full
```
## Notes

- All benchmarks were run with Go's default settings
- Results may vary based on hardware and system load
- Memory statistics include both heap and stack allocations
- Each benchmark includes multiple runs to account for variance
