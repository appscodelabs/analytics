[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/client-ip)](https://goreportcard.com/report/github.com/appscode/client-ip)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Forum](https://discuss.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# client-ip
Display Client IP

## Motivation
Cloud providers like Google Cloud Platform does not assign public ip to an interface on the VM. This tool & hosted service can be used to detect public ip address from inside a VM on those cloud providers. 

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

# update image
kubectl set image deployment/client-ip client-ip=appscode/client-ip:<updated-tag>
```
### Hosted Service
To use a hosted version of this tool, check below:
* Get Client IP: https://my-ip.space/index.json
* Get IP & request headers: https://my-ip.space/index.json?include_headers=true
