set -e
# Install docker

for pkg in docker.io docker-doc docker-compose podman-docker containerd runc; do sudo apt-get remove $pkg; done
sudo apt-get update
sudo apt-get install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install docker-ce docker-ce-cli docker-buildx-plugin docker-compose-plugin

# Check if docker is installed
sudo docker run hello-world | grep "Hello from Docker!"

#Configure to run docker without root 
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker

# Install golang
sudo apt-get install git make wget

wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

tar xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:~/go/bin

#Install k8s
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl gpg
curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.29/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg

echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.29/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# Run docker registry
docker run -d -p 5000:5000 --name registry registry:2.8

#install crio
#https://github.com/cri-o/cri-o/blob/main/install.md


#Allow insecure registry docker
#Get your ip address with from ip addr show
#In file  /etc/docker/daemon.json
#{
#    "insecure-registries": ["<ip-address->:5000"]
#}
#Also add "<ip-address>:5000" in /etc/crio/crio.conf

kubeadm init --pod-network-cidr=192.168.0.0/16 --cri-socket=unix:///var/run/crio/crio.sock

bash build.sh
bash create_secrets.sh

kubectl apply -f Deployments/application.yaml

