[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/analytics)](https://goreportcard.com/report/github.com/appscode/analytics)

[Website](https://appscode.com) • [Slack](https://slack.appscode.com) • [Twitter](https://twitter.com/AppsCodeHQ)

# analytics
Essential analytics for OSS

## Motivation
Cloud providers like Google Cloud Platform does not assign public ip to an interface on the VM. This tool & hosted service can be used to detect public ip address from inside a VM on those cloud providers.

## Usage
### CLI
```bash
go run main.go
```
### Docker
Link: https://hub.docker.com/r/appscode/analytics/
```bash
docker pull appscode/analytics
```
### Kubernetes
```bash
kubectl run analytics --image=appscode/analytics:1.1.0 --port=60010
kubectl expose deployment analytics --port=60010

# update image
kubectl set image deployment/analytics analytics=appscode/analytics:<updated-tag>
```
### Hosted Service
To use a hosted version of this tool, check below:
* Get Client IP: https://my-ip.space/index.json
* Get IP & request headers: https://my-ip.space/index.json?include_headers=true
