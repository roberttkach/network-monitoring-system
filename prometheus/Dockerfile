FROM prom/prometheus

COPY prometheus.yaml /prometheus/prometheus.yaml
EXPOSE 9090

CMD ["--config.file=/prometheus/prometheus.yaml"]