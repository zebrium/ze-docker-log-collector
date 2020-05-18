# DOCKER CONTAINER LOG COLLECTOR
Zebrium's docker container log collector collects container logs and and sends logs to Zebrium for automated Anomaly detection.
Our github repository is located [here](https://github.com/zebrium/ze-docker-log-collector).

# ze-docker-log-collector

## Getting Started
### Docker
Use the following command to create a docker log collector container:
```
sudo docker run -d --name="zdocker-log-collector" --restart=always \
                -v=/var/run/docker.sock:/var/run/docker.sock \
                -e ZE_LOG_COLLECTOR_URL="<ZE_LOG_COLLECTOR_URL>" \
                -e ZE_LOG_COLLECTOR_TOKEN="<ZE_LOG_COLLECTOR_TOKEN>" \
                -e ZE_HOSTNAME="<HOSTNAME>" \
                zebrium/docker-log-collector:latest
```

## Environment Variables
The following environment variables are supported by the collecotr:
<table>
  <tr>
    <th>Environment Variables</th>
    <th>Description</th>
    <th>Default value</th>
    <th>Note</th>
  </tr>
  <tr>
    <td>ZE_LOG_COLLECTOR_URL</td>
    <td>Zebrium log host server URL</td>
    <td>None. Must be set by user</td>
    <td>Provided by Zebrium once your account has been created.</td>
  </tr>
  <tr>
    <td>ZE_LOG_COLLECTOR_TOKEN</td>
    <td>Authentication token</td>
    <td>None. Must be set by user</td>
    <td>Provided by Zebrium once your account has been created.</td>
  </tr>
  <tr>
    <td>ZE_HOSTNAME</td>
    <td>Hostname of docker host</td>
    <td>Empty. Optional</td>
    <td>If ZE_HOSTNAME is not set, container hostname is used as source host for logs.</td>
  </tr>
  <tr>
    <td>ZE_MAX_INGEST_SIZE</td>
    <td>Maximum size of post request for Zebrium log server</td>
    <td>1048576 bytes. Optional</td>
    <td>Unit is in bytes</td>
  </tr>
  <tr>
    <td>ZE_FLUSH_TIMEOUT</td>
    <td>Interval between sending batches of log data to Zebrium log server.</td>
    <td>30 seconds. Optional</td>
    <td>Unit is in seconds. Please note Zebrium output plugin sends data immediately to log server when accumulated data reaches `ZE_MAX_INGEST_SIZE` bytes.</td>
  </tr>
</table>


## Testing your installation
Once the docker log collector software has been deployed in your environment, your container logs and anomaly detection will be available in the Zebrium UI.

## Contributors
* Brady Zuo (Zebrium)
