# cron coordinator

cron coordinator will run kube-jarvis as a cron job
this coordinator will make kube-jarvis as a http server with following handlers
* "/coordinator/cron/run" : run a diagnostic immediately

  curl http://127.0.0.1:9005/coordinator/cron/run

* "/coordinator/cron/update":  update cron config 

  curl http://127.0.0.1:9005/coordinator/cron/update -d '1 * * * *'

* "/coordinator/cron/state" : return current coordinator state

  running : diagnostic job is running now

  pendding:  diagnostic job is not running now 

# config
```yaml
coordinate:
  type: "cron"
  config:
    cron: "1 * * * *"
```