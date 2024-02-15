FROM gradle:7.6 AS build

# Set the working directory to /app
WORKDIR /app


COPY helloworld/java_client/* /app/
COPY helloworld/helloworld/helloworld.proto /app/src/main/proto/

RUN gradle shadowJar --no-daemon 

FROM openjdk:16-slim
WORKDIR /app
COPY --from=0 /app/build/libs/helloworld-all.jar ./
# Run the JAR file when the container launches
#CMD ["java", "-jar", "./helloworld-all.jar","--target", "server:50051", \
#    "--Cert" , "/etc/secret/java-client-cert.pem" , "--Key", \
#    "/etc/secret/java-client-key.pem", "--caKey", "/etc/secret/ca.pem"]

CMD ["tail", "-f", "/dev/null"]
