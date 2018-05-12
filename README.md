# cloudtail
Tail Implementation for AWS CloudWatch Logs

[![Build Status](https://travis-ci.com/tinyzimmer/cloudtail.svg?branch=master)](https://travis-ci.com/tinyzimmer/cloudtail)

Head to the [releases](https://github.com/tinyzimmer/cloudtail/releases) section to download pre-compiled binaries for **Linux** *(All Distributions)*, **macOS**, and **Windows**.

Only `amd64` binaries are provided.

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
 - EC2 IAM Instance Profile (not tests)
 - [Shared Credentials File](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/create-shared-credentials-file.html) (Linux/macOS: `$HOME/.aws/credentials`, Windows: `$env:HOME\.aws\credentials`)

## Build

```bash
$> go get -u github.com/tinyzimmer/cloudtail
```

#### TODO

- I need to write tests...
- I know just from how it's written anything over n=50 will behave oddly, shouldn't affect -f
- more inline comments

```bash
-r                   keep trying to open a group even if it is
                     non-existant or permissions are denied at first
```
