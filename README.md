Network Monitoring System
This is a network monitoring system built using Go, Prometheus, and Grafana. The system consists of three main components: a Go application that collects network metrics, Prometheus for storing and querying the metrics, and Grafana for visualizing the metrics.

Components
Go Application
The Go application is responsible for collecting network metrics such as open ports, DNS resolution times, and TCP connection times. It exposes these metrics via an HTTP endpoint at /metrics for Prometheus to scrape.

Prometheus
Prometheus is a open-source monitoring and alerting system. It scrapes metrics from the Go application and stores them in a time-series database. Prometheus is configured to scrape the Go application's /prometheus endpoint at regular intervals.

Grafana
Grafana is a open-source visualization and analytics platform. It connects to Prometheus and allows you to create dashboards and visualizations for the collected metrics.

Deployment
The system is deployed using Kubernetes. The provided YAML files define three Deployments:

go-app-deployment: Deploys the Go application with three replicas.
grafana-deployment: Deploys Grafana with a ConfigMap for dashboards.
prometheus-deployment: Deploys Prometheus with a ConfigMap for configuration.
Additionally, there are Dockerfiles for building the Go application, Prometheus, and Grafana images.

Usage
Build the Docker images:

Go Application: docker build -t hallexz/src .
Grafana: docker build -t hallexz/grafana -f Dockerfile.grafana .
Prometheus: docker build -t hallexz/prometheus -f Dockerfile.prometheus .
Push the Docker images to a registry (e.g., Docker Hub).

Deploy the system to a Kubernetes cluster:

kubectl apply -f go-app-deployment.yaml
kubectl apply -f grafana-deployment.yaml
kubectl apply -f prometheus-deployment.yaml
Access Grafana by forwarding the port:

kubectl port-forward <grafana-pod-name> 3030:3030
Open http://localhost:3030 in your browser.
Configure Grafana to connect to Prometheus and create dashboards for visualizing the collected metrics.
