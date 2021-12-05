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
