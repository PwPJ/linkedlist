#!/bin/bash

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
      rm -f $INSERT_FILE $FIND_FILE post.lua get.lua
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

# Benchmark the insert operation
echo "Benchmarking insert operation..."
if ! wrk -t12 -c100 -d30s -s post.lua http://localhost:8080; then
      rm -f $INSERT_FILE $FIND_FILE post.lua get.lua
  exit 1
fi

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

# Benchmark the find operation
echo "Benchmarking find operation..."
if ! wrk -t12 -c100 -d30s -s get.lua http://localhost:8080; then
  rm -f $INSERT_FILE $FIND_FILE post.lua get.lua
  exit 1
fi

# Cleanup
rm -f $INSERT_FILE $FIND_FILE post.lua get.lua
