# Running Benchmark Guide

This guide provides instructions on how to use the provided benchmark bash script with `wrk`.
The `benchmark.sh` script generates random numbers, inserts them into server, and benchmarks the `insert` and `find` operations.

### Prerequisites

Ensure you have the following installed on your system:

- `wrk`: HTTP benchmarking tool
- `curl`: Command-line tool for transferring data with URLs

### Installing `wrk`

To install wrk on a Linux system, use the following command:

```bash
sudo apt-get install wrk
```

### Making the Script Executable

To make the `benchmark.sh` script executable, use the following command:

```bash
chmod +x benchmark.sh
```

### Running the Script

To run the benchmark script, execute:

```bash
./benchmark.sh
```

### Run Server

Run the server by the following command:

```bash
go run main.go
```

### Script Overview

The `benchmark.sh` script performs the following actions:

1. **Generate Random Numbers**:
   - Creates `100,000` random numbers and saves them into `inserts.json` for insert operations and `finds.txt` for find operations.

2. **Insert Numbers into Server**:
   - Reads the `inserts.json` file and uses `curl` to send POST requests to insert each number into the server at `http://localhost:8080/v2/numbers/{index}/{value}`.

3. **Create Lua Script for `wrk` POST Requests**:
   - Generates a `post.lua` script to simulate random insert operations using `wrk`.

4. **Benchmark Insert Operation**:
   - Uses `wrk` to benchmark the insert operation with 12 threads, 100 connections, for 30 seconds.

5. **Create Lua Script for `wrk` GET Requests**:
   - Generates a `get.lua` script to simulate find operations using values from `finds.txt`.

6. **Benchmark Find Operation**:
   - Uses `wrk` to benchmark the find operation with 12 threads, 100 connections, for 30 seconds.

7. **Cleanup**:
   - Removes the generated files `inserts.json`, `finds.txt`, `post.lua`, and `get.lua`.