# client-ip
Display Client IP

## Usage
### CLI
```bash
go run main.go
```
### Docker
Link: https://hub.docker.com/r/appscode/client-ip/
```bash
docker pull appscode/client-ip
```
### Kubernetes
```bash
kubectl run client-ip --image=appscode/client-ip:1.1.0 --port=60010
kubectl expose deployment client-ip --port=60010
```
### Hosted Service
To use a hosted version of this tool, check here: https://my-ip.space/index.json
