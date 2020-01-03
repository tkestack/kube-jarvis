# cron coordinator

cron coordinator will run kube-jarvis as a cron job
this coordinator will make kube-jarvis as a http server with following handlers
* POST "/coordinator/cron/run" : run a diagnostic immediately

  curl http://127.0.0.1:9005/coordinator/cron/run -X POST

* POST "/coordinator/cron/period" : set cron period

  curl http://127.0.0.1:9005/coordinator/cron/period -d '1 * * * * *'

* GET "/coordinator/cron/period":  get cron period 

  curl http://127.0.0.1:9005/coordinator/cron/period
  "1 * * * * *"  

* GET "/coordinator/cron/state" : return current coordinator state

  curl http://127.0.0.1:9005/coordinator/cron/state 

  {"State":"running","Progress":{"IsDone":false,"Steps":{"diagnostic":{"Title":"Diagnosing...","Total":4,"Percent":0,"Current":0},"init_components":{"Title":"Fetching all components..","Total":10,"Percent":10,"Current":1},"init_env":{"Title":"Preparing environment","Total":2,"Percent":100,"Current":2},"init_k8s_resources":{"Title":"Fetching k8s resources..","Total":20,"Percent":100,"Current":20},"init_machines":{"Title":"Fetching all machines..","Total":4,"Percent":0,"Current":0}},"CurStep":"init_components","Total":40,"Current":23,"Percent":57.49999999999999}}  

# config
```yaml
coordinate:
  type: "cron"
  config:
    cron: "1 * * * * *"
    # this is the path that will used to save wal file,
    # the wal file is  used to auto retry if If the process is restarted at diagnostic time
    walpath: "/tmp/" 
```