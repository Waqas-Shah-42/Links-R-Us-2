# Links-R-Us-2

# SLA
|Metric|Expectation|Measurement Period | Notes|
|------|-----------|-------------------|------|
|Links 'R' Us availability|99% uptime|Yearly|Tolerates up to 3d 15h 39m of downtime per year|
|Index service availability|99.9% uptime|Yearly|Tolerates up to 8h 45m of downtime per year|
PageRank calculator service availability|70% uptime|Yearly|Not a user-facing component of our system;the service can endure longer periods of downtime|
|Search response time|30% of requests answered in 0.5s| Monthly|-|
|Search response time|70% of requests answered in 1.2s|Monthly|-|
|Search response time|99% of requests answered in 2.0s|Monthly|-|
|CPU utilization for the PageRank calculator service|90%|Weekly|We shouldn't be paying for idle computing nodes|
|SRE team incident response time|90% of tickets resolved within 8h|Monthly|-


# UML
located page 170
```mermaid
flowchart LR
    A[hello] --> B;
    B(lols)
```
