# k3s-experiments

Exploring [K3S](https://k3s.io/)

# Installation on Ubuntu 20.04 x86_64

Using the default options:

```bash
curl -sfL https://get.k3s.io | sh -
```

Create a configuration file`/etc/rancher/k3s/config.yaml` with the following values:

```yaml
write-kubeconfig-mode: "0644"
tls-san:
  - "foo.local"
node-label:
  - "foo=desktop"
```

Create the registries file at `/etc/rancher/k3s/registries.yaml`:

```
mirrors:
  "registry.local:5000":
    endpoint:
      - "http://registry.local:5000"
```

Restart k3s

```bash
sudo sytemctl restart k3s
```

Install [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/), if not present

```bash
snap install kubectl --classic
```

Install [Skaffold](https://skaffold.dev/docs/install/):

```bash
# For Linux on x86_64 (amd64)
curl -Lo skaffold https://storage.googleapis.com/skaffold/builds/latest/skaffold-linux-amd64 && \
sudo install skaffold /usr/local/bin/
```

# Local docker registry

Create a registry accessible in port 5000

```bash
# with a permanent volume
docker volume create "registry-volume"
docker container run -d --name "container-registry" -v "registry-volume":/var/lib/registry --restart unless-stopped -p 5000:5000 registry:2

# without a permanent volume
docker container run -d --name "container-registry" --restart unless-stopped -p 5000:5000 registry:2
```

Add `127.0.0.1 registry.local` to `/etc/hosts`


# Running Skaffold:

```bash
skaffold dev --kubeconfig /etc/rancher/k3s/k3s.yaml --default-repo registry.local:5000
```


# Useful commands

```bash
kubectl run temporary --image=radial/busyboxplus:curl -i --tty
```


# Cluster dashboard

Deploy the dashboard

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.4.0/aio/deploy/recommended.yaml
```

Create `admin-user` and role to access it:

```bash
kubectl create -f k8s/local/dashboard.admin-user.yaml -f k8s/local/dashboard.admin-user-role.yaml
```

Obtain bearer token:

```bash
kubectl -n kubernetes-dashboard describe secret admin-user-token | grep '^token'
```

Proxy the dashboard

```bash
kubectl proxy
```

It should be accessible at http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/


# Raspberry PI 4 Node configuration

Install the official Raspberry OS into a micro SD card. Enable SSH server following [these instructions](https://www.raspberrypi.com/documentation/computers/remote-access.html).

Access the RPi through SSH, here the IP address of the RPi is `192.168.1.63`

```bash
# default password: raspberry
ssh pi@192.168.1.63
```

Upgrade the system:

```bash
# upgrade the system
sudo apt update
sudo apt upgrade
```

Add `cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1` to the end of `/boot/cmdline.txt`:

```bash
sudo nano /boot/cmdline.txt
```

Reboot with `sudo reboot` and reconnect `ssh pi@192.168.1.63`.

Add DNS names to the k3s main server and the private registry running in the Ubuntu x86_64 machine in `/etc/hosts/`:

```bash
# K3S server node
192.168.1.65    k3s-server-1.local
192.168.1.65    registry.local
```

and verify the RPi can reach it:

```bash
ping -c 5 k3s-server-1.local

PING k3s-server-1.local (192.168.1.65) 56(84) bytes of data.
64 bytes from k3s-server-1.local (192.168.1.65): icmp_seq=1 ttl=64 time=2.76 ms
64 bytes from k3s-server-1.local (192.168.1.65): icmp_seq=2 ttl=64 time=2.36 ms
64 bytes from k3s-server-1.local (192.168.1.65): icmp_seq=3 ttl=64 time=2.32 ms
64 bytes from k3s-server-1.local (192.168.1.65): icmp_seq=4 ttl=64 time=3.41 ms
64 bytes from k3s-server-1.local (192.168.1.65): icmp_seq=5 ttl=64 time=7.53 ms

--- k3s-server-1.local ping statistics ---
5 packets transmitted, 5 received, 0% packet loss, time 10ms
rtt min/avg/max/mdev = 2.318/3.673/7.525/1.966 ms
```

On the Ubuntu x86_64 server machine, query the K3S node-token:

```bash
sudo more /var/lib/rancher/k3s/server/node-token

K10bab40d6f7d29afeb0eda2e53d46016cf28cf45a65fc9bdbb1aba2039b6e758d4::server:7e92645f5a0310e4bda7f2ceba57e1c9
```

On the RPi, install K3S in node mode:

```bash
curl -sfL https://get.k3s.io | K3S_URL=https://k3s-server-1.local:6443 K3S_TOKEN=K10bab40d6f7d29afeb0eda2e53d46016cf28cf45a65fc9bdbb1aba2039b6e758d4::server:7e92645f5a0310e4bda7f2ceba57e1c9 sh -
```

Add a `/etc/rancher/k3s/registries.yaml` file. This will enable the node to pull images over insecure HTTP:

```yaml
mirrors:
  "registry.local:5000":
    endpoint:
      - "http://registry.local:5000"
```

Reboot the RPi. After a while, the new node should become available. Run the following command on the Ubuntu machine:

```bash
kubectl get nodes

NAME          STATUS   ROLES                  AGE   VERSION
jadarve-gpu   Ready    control-plane,master   8d    v1.21.5+k3s2
raspberrypi   Ready    <none>                 30m   v1.21.7+k3s1
```

