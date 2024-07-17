# Commands:
docker build -t bookserver . \n
docker run -p 1000:8080 bookserver start -u username -s password  //1000 is local host port and 8080 is container port 
