# Chronos: The Keeper of Time Series Data

![license](https://img.shields.io/github/license/flohansen/chronos)
![server ci/cd](https://github.com/flohansen/chronos/actions/workflows/tests.yml/badge.svg)
![go report card](https://goreportcard.com/badge/github.com/flohansen/chronos)

## Quick Start

### Configuration

Chronos uses a YAML configuration to setup storage parameters and scrape
targets. Here is an example config containing all values.

```yaml
targets:
- url: http://localhost:3000/metrics
  interval: 1s
```

### Expose Metrics

Chronos scrapes metrics using HTTP protocol. To make Chronos being able to
store your specific metrics you have to expose them using a HTTP server and the
following format.

```
<metric_name> <value>
<metric_name> <value>
<metric_name> <value>
...
```

Now you can configure Chronos to scrape it. Take a look at the configuration
example to see how.
