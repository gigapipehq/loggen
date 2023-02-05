# Loggen
## Fake log, metric and trace generator for qryn Cloud

### Install

```shell
git clone git@github.com:gigapipehq/loggen.git
cd loggen

go install
```

### Configuration

A YAML configuration file will be generated automatically
on first run, that will be stored in `~/.loggen/config.yaml`. 

The config file defines the default parameters, but can be 
overridden using any of the flags defined below. 

### Example config file

```yaml
url: https://qryn.gigapipe.com
api_key: abc
api_secret: cba
labels: 
  label1: value1
  label2: value2
rate: 100
timeout: 30s
enable_metric: true
enable_traces: true
```

### Usage
```shell
A fake log, metric and trace generator for qryn Cloud

Usage:
  loggen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run the generator in cli-mode
  server      Run the generator in server-mode

Flags:
  -m, --enable-metrics   Enable collection of Prometheus metrics (default true)
  -o, --enable-traces    Enable collection of OpenTelemetry traces (default true)
  -h, --help             help for loggen

Use "loggen [command] --help" for more information about a command.
```

#### CLI mode

```shell
Run the generator in cli-mode

Usage:
  loggen run [flags]

Flags:
  -k, --api-key string          API key to use for authenticating with qryn Cloud (default "abc")
  -s, --api-secret string       API key to use for authenticating with qryn Cloud (default "cba")
  -h, --help                    help for run
  -l, --labels stringToString   labels for each log (default [])
  -r, --rate int                number of logs to generate per second (default 100)
  -t, --timeout duration        length of time to run the generator before exiting (default 30s)

Global Flags:
  -m, --enable-metrics   Enable collection of Prometheus metrics (default true)
  -o, --enable-traces    Enable collection of OpenTelemetry traces (default true)
```

#### Server mode

```shell
Run the generator in server-mode

Usage:
  loggen server [flags]

Flags:
  -h, --help   help for server

Global Flags:
  -m, --enable-metrics   Enable collection of Prometheus metrics (default true)
  -o, --enable-traces    Enable collection of OpenTelemetry traces (default true)
```
