FROM node:21-alpine3.18

WORKDIR /app

COPY helloworld/node_client/* ./
COPY helloworld/helloworld/helloworld.proto ./

RUN npm install
RUN npm install @grpc/grpc-js

CMD ["node", " "--target", "server-service.default.svc.cluster.local:50051"]
