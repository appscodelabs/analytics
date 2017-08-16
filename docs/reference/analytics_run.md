## analytics run

Run server

### Synopsis


Run server

```
analytics run [flags]
```

### Options

```
      --analytics                        Send analytical events to Google Analytics (default true)
      --cacert-file string               File containing CA certificate
      --cert-file string                 File container server TLS certificate
      --docker-hub-orgs stringToString   Map of Docker Hub organizations to Google spreadsheets (default [])
  -h, --help                             help for run
      --key-file string                  File containing server TLS private key
      --ops-addr string                  Address to listen on for web interface and telemetry. (default ":56790")
      --web-address string               Http server address (default ":9844")
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [analytics](analytics.md)	 - Analytics by AppsCode - Essential analytics for OSS

