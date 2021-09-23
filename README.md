# Documentation
This is the repository of fikriprayoga1/go-mongo in docker hub. Function of this repository is guide to create image of fikriprayoga1/go-mongo in docker hub, so you can learn how to create image and run the docker system image easily. This repository is part of https://hub.docker.com/repository/docker/fikriprayoga1/go-mongo project

## Best Practice
Use this command in your CMD(Windows) or Terminal(Mac or Linux). Before start the command, you must be in the directory of the folder you pulled.

### Step 1. Create fikriprayoga1/go-mongo image
```
docker build -t fikriprayoga1/go-mongo:1.0 .
```

### Step 2. Pull mongodb image
```
docker pull mongo
```

### Step 3. Run docker compose
```
docker-compose up
```

### Step 4. Test in Postman app
You can use Postman app to test server with this configuration. Please import postman collection file in this folder to your postman app collection, and then try the API

## Warning
Ensure the system running before you try API. You can run 'docker logs go-mongo' command to see the log, if you have seen 'Server listener started.' log, you can try the API

## Utility Command
This part is command line to help you something

### Show running container list
```
docker container ls
```

### Show image list
```
docker images
```

### Show  container list
```
docker container ls -a
```

### Show  container log for debugging
```
docker logs go-mongo
```

### Show network list
```
docker network ls
```
