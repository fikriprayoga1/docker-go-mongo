# Documentation
This is repository of fikriprayoga1/go-mongo in docker hub. The function of this repository is create image of fikriprayoga1/go-mongo in docker hub, so you can learn how to create image and run the docker system image easily. This repository is part of https://hub.docker.com/repository/docker/fikriprayoga1/go-mongo project

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
docker-compose up -d
```

### Step 4. Test in Postman app
You can use Postman app to test server with this configuration
```
Request Method = POST
Address = localhost:8080/profile/create
Body Type = form-data
Key = ProfileImage & Value = https://www.google.com/image-profile.png
Key = Name & Value = Budi
Key = Email & Value = budi@example.com
Key = Password & Value = Test123
```

```
Request Method = POST
Address = localhost:8080/profile/read
Body Type = form-data
Key = Id & Value = 614775ff8d156f54ebb02be0
```

```
Request Method = POST
Address = localhost:8080/profile/update
Body Type = form-data
Key = Id & Value = 614775ff8d156f54ebb02be0
Key = ProfileImage & Value = https://www.google.com/image-profile.png
Key = Name & Value = Budi
Key = Email & Value = budi@example.com
Key = Password & Value = Test123
```

```
Request Method = POST
Address = localhost:8080/profile/delete
Body Type = form-data
Key = Id & Value = 614775ff8d156f54ebb02be0
```

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
