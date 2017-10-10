# mrpe exporter

An local monitoring agent exporting output of nagios plugins.
The configuration is parsed from mrpe config directory.


## Status

===   WORK IN PROGRESS   ===



## How it works

The daemon will read the config directory of mrpe at startup and periodically run the defined checks. The output of these checks are collected and exposed as metrics for prometheus.

MRPE configuration /etc/mrpe/conf.d/foo.cfg
```
# Interval 60
fooplugin /usr/lib/nagios/plugins/fooplugin -w 2 -c 3
```

Will run by default every minute the fooplugin. You can change the check interval in comment above the check definition.

Metrics:
```
# HELP cmk_fooplugin_duration_ns runtime in ns for fooplugin
# TYPE cmk_fooplugin_duration_ns gauge
cmk_fooplugin_duration_ns{stage="dev"} 89

# HELP cmk_fooplugin_exit_code check exitcode for fooplugin
# TYPE cmk_fooplugin_exit_code gauge
cmk_fooplugin_exit_code{stage="dev"} 0
```

## Todo

- [x] configuration via flag
- [ ] initscript for OS
- [ ] ...
