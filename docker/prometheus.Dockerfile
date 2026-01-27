FROM prom/prometheus:latest

COPY observability/prometheus/prometheus.yaml /etc/prometheus/prometheus.yml
COPY observability/prometheus/rules /etc/prometheus/rules