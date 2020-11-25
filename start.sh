# Possibly export process id to environ
# Possibly write pids to file so log.go can change a pid if it restarts a program
  # If done, make the kill command read from a file of pids using a while loop

# Used to build/compile all the necessary files
build() {
  echo "Building..."
  cd server
  go build .
  chmod a+x chat.py
  cd others
  # f95 -c utils.f90
  # f95 -c indicators.f90 -o indicators.so
  cd ../../logs
  go build .
  cd ..
  echo "Done building"
}

# Used to start the servers
# Possibly move code outside the function
run() {
  echo "Running..."
  logs/logs & logs_pid=$!
  server/server & server_pid=$!
  server/chat.py & chat_pid=$!
}

# Used to restart the Go server
# Possibly won't be used
restartserver() {
  echo "Restarting server..."
  kill "$server_pid"
  cd server
  go build .
  cd ..
  server/server & server_pid="$!"
  echo "Server restarted"
}

export IP="192.168.1.125"
export PORT="8000"
export STOCKS_PORT="8080"
export BOT_PORT="8008"
export LOG_IP="localhost"
export LOG_PORT="7000"

build
run
action=""

while [ "$action" != "kill" ]
do
  read -p 'Action: ' action
  if [ "$action" == "kill" ]; then
    echo "Killing..."
    kill "$logs_pid"
    kill "$server_pid"
    kill "$chat_pid"
    echo "Killed"
  elif [ "$action" == "restart server" ]; then
    restartserver
  else
    echo "Invalid action: $action"
  fi
done