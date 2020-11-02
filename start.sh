# Possibly export process id to environ
# Possibly write pids to file so log.go can change a pid if it restarts a program
  # If done, make the kill command read from a file of pids using a while loop
build() {
  echo "Started building"
  cd server
  go build .
  gfortran -c indicators.f90 indicators.so
  cd ../logs
  go build .
  cd ..
  echo "Done building"
}
run() {
  echo "Running"
  logs/logs &
  logs_pid=$!
  server/server &
  server_pid=$!
  server/chat.py &
  chat_pid=$!
}
build
run
endprogram=""
while [ "$endprogram" != "kill" ]
do
  read endprogram
  if [ "$endprogram" == "kill" ]
  then
    echo "Killing"
    kill "$logs_pid"
    kill "$server_pid"
    kill "$chat_pid"
    echo "Killed"
  fi
done