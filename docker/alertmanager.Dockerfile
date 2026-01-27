FROM prom/alertmanager:latest

COPY observability/alertmanager/alertmanager.yaml /etc/alertmanager/alertmanager.yml