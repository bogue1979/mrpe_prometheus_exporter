# mrpe exporter

An local monitoring agent exporting output of nagios plugins.
The configuration is parsed from mrpe config directory.


## Status

===   WORK IN PROGRESS   ===



## How it works

The daemon will read the config directory of mrpe at startup and periodically run the defined checks. The output of these checks are collected and exposed as metrics for prometheus.
When nagios performance data is found in the output of the executed plugin, these metrics will be added to the defaults.
These defaults are:

exit ( exit code of the plugin )
duration ( how long the plugin took to execute in ns )

MRPE configuration /etc/mrpe/conf.d/foo.cfg
```
# Interval 60
fooplugin echo "Test | baz=1;;; foo=2;;; bar=0.4;;;"
```

Will run by default every minute the fooplugin. You can change the check interval in comment above the check definition.

Start daemon:
```
Usage of ./mrpe_prometheus_exporter:
  -conf.dir string
        directory with mrpe config files (default "./conf.d")
  -env.key string
        environment differentiator (default "stage")
  -env.val string
        environment name (default "dev")
  -log.json
        log as json
  -log.lvl string
        loglevel from [debug,info,warn,error,fatal,panic] (default "info")
  -version
        show version and exit

```

Metrics:
```
# HELP cmk_fooplugin Check_MK metrics for fooplugin
# TYPE cmk_fooplugin gauge
cmk_fooplugin{metric="baz",stage="dev"} 1
cmk_fooplugin{metric="duration",stage="dev"} 0.007444337
cmk_fooplugin{metric="exit",stage="dev"} 0
cmk_fooplugin{metric="foo",stage="dev"} 2
```

## Todo

- [x] configuration via flag
- [x] logging
- [ ] initscript for OS
- [ ] ...
