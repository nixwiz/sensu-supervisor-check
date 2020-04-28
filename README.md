[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-supervisor-check)
![Go Test](https://github.com/nixwiz/sensu-supervisor-check/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-supervisor-check/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/nixwiz/sensu-supervisor-check)](https://goreportcard.com/report/github.com/nixwiz/sensu-supervisor-check)

# Supervisor Check for Sensu Go

## Table of Contents
- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The Supervisor Check for Sensu Go is a [Sensu Check][2] that reports on the status of processes
managed by [Supervisor][5].

## Files

N/A

## Usage examples

```
Supervisor Check for Sensu Go

Usage:
  sensu-supervisor-check [flags]
  sensu-supervisor-check [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical string   Supervisor states to consider critical (default "FATAL")
  -h, --help              help for sensu-supervisor-check
  -H, --host string       Host running Supervisor (default "localhost")
  -P, --port int          Supervisor listening port (default 9001)
  -s, --socket string     Supervisor listening UNIX domain socket

Use "sensu-supervisor-check [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][4] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add nixwiz/sensu-supervisor-check
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][3]

### Check definition

**NOTE:**  This plugin does not currently support the use of username/password authentication to the
Supervisor URL or UNIX domain socket.

#### Using HTTP
```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-supervisor-check
  namespace: default
spec:
  command: sensu-supervisor-check --host super.example.com --port 9001 -c STARTING,BACKOFF,STOPPING,EXITED,FATAL,UNKNOWN,STOPPED
  subscriptions:
  - system
  runtime_assets:
  - nixwiz/sensu-supervisor-check
```

#### Using a UNIX domain socket
```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: sensu-supervisor-check
  namespace: default
spec:
  command: sensu-supervisor-check -s /var/run/supervisor.sock -c STARTING,BACKOFF,STOPPING,EXITED,FATAL,UNKNOWN,STOPPED
  subscriptions:
  - system
  runtime_assets:
  - nixwiz/sensu-supervisor-check
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-supervisor-check repository:

```
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[3]: https://bonsai.sensu.io/assets/nixwiz/sensu-supervisor-check
[4]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[5]: http://supervisord.org/
