# cloudtail
Tail Implementation for CloudWatch Logs

[![Build Status](https://travis-ci.com/tinyzimmer/cloudtail.svg?branch=master)](https://travis-ci.com/tinyzimmer/cloudtail)

Head to the [releases](https://github.com/tinyzimmer/cloudtail/releases) section for pre-compiled binaries

```bash
 OPTIONS
  -f    Follow the log group
  -n int
        Number of lines to dump (default 10)
  -s int

        Interval (in seconds) to poll during a follow (default 3)
$> cloudtail [OPTIONS] logGroup # accepts substring
```

## TODO (stolen from real tail)
```bash
-r                   keep trying to open a group even if it is
                     non-existant or permissions denied at first
-p                   with -f, terminate after process ID, PID dies
-q                   never output metadata for log events
-v                   always output metadata for log events
--version            output version information and exit
```
