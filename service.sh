#!/bin/sh
#
# service script

# Check the application status
#
# This function checks if the application is running
check_status() {

  # Running ps with some arguments to check if the PID exists
  # -C : specifies the command name
  # -o : determines how columns must be displayed
  # h : hides the data header
  s=`ps -C 'java -jar /path/to/your.jar' -o pid h`

  echo $s

  # If somethig was returned by the ps command, this function returns the PID
  if [ $s ] ; then
    return $s
  fi

  # In any another case, return 0
  return 0

}

# Starts the application
start() {

  # At first checks if the application is already started calling the check_status
  # function
  pid=$(check_status)
  
  if [[ "$pid" -eq "" ]]
  then
    pid=0
  fi

  if [ $pid -ne 0 ] ; then
    echo "The application is already started"
    exit 1
  fi

  # If the application isn't running, starts it
  echo -n "Starting application: "

  # Redirects default and error output to a log file
  #java -jar /path/to/application.jar >> /path/to/logfile 2>&1 &
  java -jar /path/to/your.jar >> /path/to/your.log 2>&1 &
  echo "OK"

}

# Stops the application
stop() {

  # Like as the start function, checks the application status
  pid=$(check_status)

  if [[ "$pid" -eq "" ]]
  then
    pid=0
  fi  

  if [ $pid -eq 0 ] ; then
    echo "Application is already stopped"
    exit 1
  fi

  # Kills the application process
  echo -n "Stopping application: "
  kill -9 $pid &
  echo "OK"

}

# Redeploys the application
redeploy() {
  
  stop
  currentDir=`pwd`
  #echo $currentDir
  cd /path/to/your/repo
  git pull
  mvn clean install
  cd $currentDir
  #echo `pwd`
  start

}

# Show the application status
status() {

  # The check_status function, again...
  pid=$(check_status)

  # If the PID was returned means the application is running
  if [ $pid -ne 0 ] ; then
    echo "Application is started: $pid"
  else
    echo "Application is stopped"
  fi

}

# Main logic, a simple case to call functions
case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  status)
    status
    ;;
  redeploy)
    redeploy
    ;;
  restart)
    stop
    start
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|reload|status}"
    exit 1
esac

exit 0