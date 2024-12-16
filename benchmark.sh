#!/bin/bash

cleanup() {
    echo "Cleanup...."
    [ -e "$INSERT_FILE" ] && unlink "$INSERT_FILE"
    [ -e "$FIND_FILE" ] && unlink "$FIND_FILE"
    [ -e "post.lua" ] && unlink "post.lua"
    [ -e "get.lua" ] && unlink "get.lua"
    [ -e "rwget.lua" ] && unlink "rwget.lua"
    echo "Cleanup completed"
}

# trap for cleanup on signals
trap cleanup EXIT SIGINT SIGTERM

# Generate random numbers for insert and find operations
INSERT_FILE="inserts.json"
FIND_FILE="finds.txt"
NUM_ENTRIES=100000

# Create random insert JSON
echo "[" > $INSERT_FILE
for ((i=0; i<$NUM_ENTRIES; i++)); do
  VALUE=$((RANDOM))
  if [ $i -eq $((NUM_ENTRIES-1)) ]; then
    echo "{\"index\": $i, \"value\": $VALUE}" >> $INSERT_FILE
    echo "$VALUE" >> $FIND_FILE
  else
    echo "{\"index\": $i, \"value\": $VALUE}," >> $INSERT_FILE
    echo "$VALUE" >> $FIND_FILE
  fi
done
echo "]" >> $INSERT_FILE

# Insert the numbers into the server
echo "Inserting numbers into the server..."
while IFS= read -r line; do
  index=$(echo $line | sed 's/.*"index": \([0-9]*\).*/\1/')
  value=$(echo $line | sed 's/.*"value": \([0-9]*\).*/\1/')
  if ! curl -s -X POST -H "Content-Type: application/json" -d "$line" http://localhost:8080/v2/numbers/$index/$value > /dev/null; then
    echo "Error inserting $line"
    exit 1
  fi
done < <(tail -n +2 $INSERT_FILE | head -n -1 | sed 's/,$//')

# Create a script for wrk to use for POST requests
cat <<EOF > post.lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
request = function()
    local index = math.random(0, $NUM_ENTRIES-1)
    local value = math.random(0, 32767)
    local body = string.format('{"index": %d, "value": %d}', index, value)
    return wrk.format(nil, "/v2/numbers/" .. index .. "/" .. value, nil, body)
end
EOF

# Create a script for wrk to use for GET requests
cat <<EOF > get.lua
wrk.method = "GET"
values = {}
index = 1
$(awk '{print "table.insert(values, \"" $0 "\")"}' $FIND_FILE)
request = function()
    local path = "/v2/numbers/value/" .. values[index]
    index = index + 1
    if index > #values then
        index = 1
    end
    return wrk.format(nil, path)
end
EOF

# Create a script for wrk to use for GET requests for RWMutex
cat <<EOF > rwget.lua
wrk.method = "GET"
values = {}
index = 1
$(awk '{print "table.insert(values, \"" $0 "\")"}' $FIND_FILE)
request = function()
    local path = "/v2/numbers/rwmutex/value/" .. values[index]
    index = index + 1
    if index > #values then
        index = 1
    end
    return wrk.format(nil, path)
end
EOF

# Benchmark the insert operation
echo "Benchmarking insert operation..."
wrk -t12 -c100 -d30s -s post.lua http://localhost:8080 &
pid1=$!

# Benchmark the find operation
echo "Benchmarking find operation..."
wrk -t12 -c100 -d30s -s get.lua http://localhost:8080 &
pid2=$!

# Benchmark the RWMutex find operation
echo "Benchmarking RWMutex find operation..."
wrk -t12 -c100 -d30s -s rwget.lua http://localhost:8080 &
pid3=$!

# Wait for all background processes to finish and capture their exit statuses
wait $pid1
status1=$?
wait $pid2
status2=$?
wait $pid3
status3=$?

# Check the exit statuses and exit with an error if any command failed
if [ $status1 -ne 0 ]; then
    echo "Error: Benchmarking insert operation failed"
    exit 1
fi

if [ $status2 -ne 0 ]; then
    echo "Error: Benchmarking find operation failed"
    exit 1
fi

if [ $status3 -ne 0 ]; then
    echo "Error: Benchmarking RWMutex find operation failed"
    exit 1
fi

echo "Completed successfully"
