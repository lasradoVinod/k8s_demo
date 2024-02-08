FROM gradle:7.6 AS build

# Set the working directory to /app
WORKDIR /app


COPY helloworld/java_client/* /app/
COPY helloworld/helloworld/helloworld.proto /app/src/main/proto/

RUN gradle shadowJar --no-daemon 

FROM openjdk:11-jre-slim
# Run the JAR file when the container launches
CMD ["java", "-jar", "build/libs/helloworld-all.jar","--target", "server-service.default.svc.cluster.local:50051", 
    "--Cert" , "/etc/secret/java-client-cert.pem" , "--Key",
    "/etc/secret/java-client-key.pem", "--caKey", "/etc/secret/ca.pem"]
