# Backend challenge (advanced)

## DataHow backend coding challenge

In this coding challenge, for which you should **not invest more than 3 hours**, you will be developing a micro-service.

Your team has been asked by business to provide them with a dashboard that shows how many unique IP addresses have been visiting the company's website. The DevOps engineer in your team says, he can easily use Prometheus and Grafana to provide the dashboard, but you need to provide a micro-service to count the unique IP addresses. The DevOps engineer can configure the ingress controller to send you a structured log message in JSON format for each page request in the following format:

```log
{ "timestamp": "2020-06-24T15:27:00.123456Z", "ip": "83.150.59.250", "url": ... }
```

Each log entry is sent to your web-service via HTTP/1.1 POST to  
[http://your-service:5000/logs](http://your-service:5000/logs) . For simplicity, you are only supposed to count how many unique IP addresses have been logged since the start of the service. Keep in mind that there might be potentially thousands of visitors per second. Provide the cumulative count as custom metric to a Prometheus server. Prometheus will periodically scrape metrics from your service at [http://your-service:9102/metrics](http://your-service:9102/metrics) .

You have 3 hours to implement, benchmark, and document your solution. Decide on the order in which you implement the requirements. We value clean and tested code over fully implemented requirements. Also, we appreciate the use of smart algorithms and 3rd party libraries to achieve your goal. Finally, upload your solution to a public Git repository like GitHub.

### Requirements

- **Time limit** for assignment: 3 hours
- Listen on ports :5000 and :9102
- Receive JSON logs on :5000/logs
- Serve Prometheus metrics on :9102/metrics
- Compute number of unique IP addresses in logs since service start
- Create custom Prometheus metric "unique_ip_addresses"
- Publish your result in public Git repository

### NOT required

- Persistence
- Validate inputs (assume the logs are well formatted)

### Bonus

- Benchmark your API with a tool like siege or gobench
- Preferred languages Go, Rust, TypeScript
- Develop your code in logical increments and document those via git commit messages
- Test driven development

### Evaluation criteria

- Clean code
- Memory usage
- Performance
- Unit tests
- Git history and commit messages
