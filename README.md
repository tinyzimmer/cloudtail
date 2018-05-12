# cloudtail
Tail Implementation for AWS CloudWatch Logs

[![Build Status](https://travis-ci.com/tinyzimmer/cloudtail.svg?branch=master)](https://travis-ci.com/tinyzimmer/cloudtail)

Head to the [releases](https://github.com/tinyzimmer/cloudtail/releases) section to download pre-compiled binaries for **Linux** *(All Distributions)*, **macOS**, and **Windows**.

Only `amd64` binaries are provided.

```bash
 OPTIONS

 -version
       Display version and exit

  -f    Follow the log group
  -n    Number of lines to dump (default 10)
  -p    with -f, terminate after process ID, PID dies
  -s    Interval (in seconds) to poll during a follow (default 3)

$> cloudtail [OPTIONS] logGroup # accepts substring
```

## Docker

For whatever reason, there is a docker image you can use also.

```bash
$> alias ctail='docker run --rm -it tinyzimmer/cloudtail:latest  /cloudtail'
$> ctail --help
```

## AWS Credentials

See the AWS documentation for configuring an SDK client. The order in which `cloudtail` checks credentials is:

 - Environment Credentials
 - EC2 IAM Instance Profile
 - Shared Credentials File (~/.aws/credentials)

## Build

```bash
$> go get -u github.com/tinyzimmer/cloudtail
```

#### TODO (stolen from real GNU tail)
```bash
-r                   keep trying to open a group even if it is
                     non-existant or permissions denied at first
-q                   never output metadata for log events
-v                   always output metadata for log events
```
