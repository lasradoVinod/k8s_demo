FROM gradle:latest AS build

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY helloworld/java_client/* /app/
COPY helloworld/helloworld/helloworld.proto /app/src/main/proto/
# Build the Gradle project
RUN gradle build --no-daemon 

FROM openjdk:11-jre-slim
# Run the JAR file when the container launches
CMD ["java", "-jar", "build/libs/helloworld.jar","--target","server-service.default.svc.cluster.local:50051"]
