FROM prom/prometheus:latest
COPY observability/prometheus/prometheus.yaml /etc/prometheus/prometheus.yml