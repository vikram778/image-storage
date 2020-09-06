# image-storage
To run in Local
1. Would need to update the ENV values in env.local file accordinlgy for DB and kafka brokers.
2. Application dependencies are included via go mod can execute go mod tidy before stating off the application.
3. make sure the db is created before starting the application,on start of the application the tables would be migrated via goose.
4. Once updated the env.local file execute sh run.sh script.
5. Application by default would start at localhost:8080 can choose to update the same in env.local file.
6. hit the swaggger endpoint to check the API contract.
7. Check the database if the tables were migrated/created successfully after the application start.
8. On start of application can see the logs on console and also logs would be written to a file for every request in /log/image-storage directory can update the path in env.local if required.



Application Emdpoints : 
basepath : http://localhost:8080/
SwaggerDocs Path : /api/documentation

Add Image :/add/image -> To add new image to album
Add Album :/add/album -> To create new Album
Get Image by ID : /get/image/{id} -> Displays image from album
Delete Image : /del/image/{id}  -> Deletes Image from Album
Delete album : /del/album/{tittle} -> Deletes entire album

This Application uses filesytem approcah to store and serve images.
While the corresponding Path of the image and its association with album is also Stored in Db.
can checkout the Db schema in migration/migration.go

This application uses different kafka topic to produce notification for CREATE/DELETE image and album.
All the topics would be created when the message is being produced in those topics no need to explicitly create topics.
This is for connecting to local kafka cluster. connecting to AWS kafka cluster would require different config.

#DOCKER
Have Include the Docker file which be used to create docker image and run the application in docker container.
can execute the below cmd to run application in docker container in Local :

1. navigate to application directory
2. run - docker build -t image-storage .
3. image created with name image-storage
4. run - docker run image-storage.