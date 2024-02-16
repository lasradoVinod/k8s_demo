names=("lightfoot" "watcher")

for name in "${names[@]}"; do
    docker build -f Dockerfiles/$name/Dockerfile -t $name:v1 .
    docker tag $name:v1 localhost:5000/$name
done

