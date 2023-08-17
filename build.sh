eval $(minikube docker-env)
docker build -f Dockerfiles/client/Dockerfile -t client:v1 .
docker build -f Dockerfiles/server/Dockerfile -t server:v1 .
