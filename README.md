# cloudtail
GNU tail-like Implementation for AWS CloudWatch Logs

[![Build Status](https://travis-ci.com/tinyzimmer/cloudtail.svg?branch=master)](https://travis-ci.com/tinyzimmer/cloudtail) [![Coverage Status](https://coveralls.io/repos/github/tinyzimmer/cloudtail/badge.svg?branch=master)](https://coveralls.io/github/tinyzimmer/cloudtail?branch=master)

Head to the [releases](https://github.com/tinyzimmer/cloudtail/releases) section to download pre-compiled binaries for **Linux** *(All Distributions)*, **macOS**, and **Windows**.

Only `amd64` binaries are provided, easy to add others if requested.

```bash
 OPTIONS

 -version
       display version and exit

  -f    follow the log group
  -n    number of lines to dump (default 10)
  -p    with -f, terminate after process ID, PID dies
  -s    interval (in seconds) to poll during a follow (default 3)
  -q    never output metadata for log events
  -v    always output metadata for log events (with -q, displays connect info)

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

 - [Environment Credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-environment.html)
 - EC2 IAM Instance Profile (not tested)
 - [Shared Credentials File](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/create-shared-credentials-file.html) (Linux/macOS: `$HOME/.aws/credentials`, Windows: `$env:HOME\.aws\credentials`)

## Build

```bash
$> go get -u github.com/tinyzimmer/cloudtail
```

#### TODO

- I know just from how it's written anything over n=50 will behave oddly, shouldn't affect -f
- pid poll can be threaded off probably to be more efficient. I only put it there anyway because original tail has it.
- bytes filters
- stream locking for follow
- date and search filters
- more inline comments
- custom output formats (json, yaml, etc.)
- multitail abilities - display multiple log groups side by side

```bash
-r                   keep trying to open a group even if it is
                     non-existant or permissions are denied at first
```
