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
api_key: my_api_key
api_secret: my_api_secret
labels: 
  label1: value1
  label2: value2
rate: 100
timeout: 30s
format: logfmt
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
  config      Configuration commands
  help        Help about any command
  run         Run the generator in cli-mode
  server      Run the generator in server-mode

Flags:
  -h, --help             help for loggen

Use "loggen [command] --help" for more information about a command.
```

#### Configure
```shell
Configuration commands

Usage:
  loggen config [command]

Available Commands:
  get         Show current default configuration settings
  reset       Reset configuration to initial defaults
  set         Set current default configuration setting

Flags:
  -h, --help   help for config

Use "loggen config [command] --help" for more information about a command.
```

##### Get current configuration
```shell
Show current default configuration settings

Usage:
  loggen config get [setting name] [flags]

Flags:
  -h, --help                   help for get
  -o, --output-format string   output format of config; yaml, or json (default "yaml")
```

##### Reset configuration
```shell
Reset configuration to initial defaults

Usage:
  loggen config reset [flags]

Flags:
  -h, --help   help for reset
```

##### Set new configuration setting
```shell
Set current default configuration setting

Usage:
  loggen config set [setting name] [setting value] [flags]

Flags:
  -h, --help   help for set
```

#### CLI mode
```shell
Run the generator in cli-mode

Usage:
  loggen run [flags]

Flags:
  -k, --api-key string          API key to use for authenticating with qryn Cloud (default "")
  -s, --api-secret string       API key to use for authenticating with qryn Cloud (default "")
  -f, --format string           format to use when sending logs (default "logfmt")
  -h, --help                    help for run
  -l, --labels stringToString   labels for each log (default [])
  -r, --rate int                number of logs to generate per second (default 100)
  -d, --timeout duration        length of time to run the generator before exiting (default 30s)
  -m, --enable-metrics          enable collection of Prometheus metrics (default true)
  -t, --enable-traces           enable collection of OpenTelemetry traces (default true)
```

#### Server mode
```shell
Run the generator in server-mode

Usage:
  loggen server [flags]

Flags:
  -h, --help   help for server
```
