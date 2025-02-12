# DataHow

DataHow is a simple microservice that processes structured JSON log messages to count unique IP addresses in memory and exposes the count as a Prometheus metric.

## Challenge Description

For the full challenge details, see [CHALLENGE](CHALLENGE.md).

## Requirements

- Go (version 1.23.x)
- Docker
- Docker Compose
- Graphana k6 (for load testing)

## Getting Started

To run the `ipcounter` microservice, build and start the Docker environment using the following step:

```sh
docker compose up -d --build
```

## Prometheus

Connect to Prometheus at [http://localhost:9102/metrics](http://localhost:9102/metrics) and search for the `unique_ip_addresses` at the bottom of the page.

```sh
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 40
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP unique_ip_addresses No. of unique IP addresses
# TYPE unique_ip_addresses counter
unique_ip_addresses 119
```

## Grafana

Login to Grafana at [http://localhost:3000](http://localhost:3000) using the following dummy credentials: `admin` and `secret`. Then configure the Prometheus data source and explore the `unique_ip_addresses` metric.

## Cleanup

To stop and remove the Docker containers, use the following command:

```sh
docker compose down
```

## Load Testing

In the CHALLENGE document, it was suggested as a bonus point to run benchmark tests using [siege](https://github.com/JoeDog/siege) or [gobench](https://github.com/gobench-io/gobench). I chose to use [Graphana k6](https://k6.io) instead, as I had experience with it in a previous Go project, [Alessandrina](https://github.com/rotiroti/alessandrina).
After installing k6, you can run the load test using the following command:

```sh
k6 run ./tests/k6/script.js
...
...
running (15.1s), 00/10 VUs, 150 complete and 0 interrupted iterations
```

## TODO

- Replace the current map implementation with [sync.Map](https://pkg.go.dev/sync#Map) to improve performance in concurrent environments.
- Explore alternative data structures, such as probabilistic counters (e.g., HyperLogLog), to reduce memory usage while maintaining accuracy.
- Add more unit tests, particularly for the `run()` function of the ipcounter service, to ensure robustness.
- Develop additional benchmark scenarios to evaluate performance under different workloads.

## License

[MIT](LICENSE)
