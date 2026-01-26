.PHONY: prometheus-build prometheus-run prometheus-up prometheus-down

PROM_IMAGE=sentinel-prometheus

prometheus-build:
	docker build -t $(PROM_IMAGE) -f docker/prometheus.Dockerfile .

prometheus-run:
	docker run -p 9090:9090 $(PROM_IMAGE)

prometheus-up: prometheus-build prometheus-run

prometheus-down:
	- docker stop $$(docker ps -q --filter ancestor=$(PROM_IMAGE))