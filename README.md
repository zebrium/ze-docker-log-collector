# DOCKER CONTAINER LOG COLLECTOR
Zebrium's docker container log collector collects container logs and sends logs to Zebrium for automated Incident detection.
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
                -e ZE_DEPLOYMENT_NAME="YOUR_DEPLOYMENT_NAME_HERE" \
                zebrium/docker-log-collector:latest
```

### Docker Compose
Use the following configuration file to deploy via docker-compose command:
```
version: '3.5'

services:
  zdocker-log-collector:
    image: zebrium/docker-log-collector:latest
    restart: always
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      ZE_LOG_COLLECTOR_URL: "<ZE_LOG_COLLECTOR_URL>"
      ZE_LOG_COLLECTOR_TOKEN: "<ZE_LOG_COLLECTOR_TOKEN>"
      ZE_DEPLOYMENT_NAME: "<YOUR_DEPLOYMENT_NAME_HERE>"
      ZE_HOSTNAME: "<HOSTNAME>"
```
### AWS Elastic Container Service (ECS)

Add the following service to ECS on EC2 cluster configuration.
```
services:
  zdocker-log-collector:
    image: zebrium/docker-log-collector:latest
    restart: always
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      ZE_LOG_COLLECTOR_URL: "<ZE_LOG_COLLECTOR_URL>"
      ZE_LOG_COLLECTOR_TOKEN: "<ZE_LOG_COLLECTOR_TOKEN>"
      ZE_DEPLOYMENT_NAME: "<YOUR_DEPLOYMENT_NAME_HERE>"
```
To collect container logs from all nodes in an ECS cluster, zdocker-log-collector service must be configured to run as an ECS daemon task. Please follow the steps below to configure the daemon task:
1. Log in to the AWS console and navigate to the ECS Clusters section. Click into your cluster you run the Agent on.
2. Choose Service tab, click on the Create button.
3. For launch type, select EC2, for service type, select DAEMON, type a service name, and click on Next step.
4. For Load balance type option, select None, and click on Next step. On next page, select Next step without configuring Auto Scaling.
5. Review and click on Create Service.

## Environment Variables
The following environment variables are supported by the collector:
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
    <td>Unit is in seconds. Please note Zebrium output plugin sends data immediately to log server when accumulated data reaches ZE_MAX_INGEST_SIZE bytes.</td>
  </tr>
  <tr>
    <td>ZE_FILTER_NAME</td>
    <td>Collect logs for containers whose names match filter name pattern. These can include wildcards, for example, <i>my_container1*</i></td>
    <td>Empty. Optional</td>
    <td></td>
  </tr>
  <tr>
    <td>ZE_FILTER_LABELS</td>
    <td>Collect logs for containers whose labels match the labels as defined in ZE_FILTER_LABELS. The format is: <i>label1:label1_value,label2:label2_value</i> These can include wildcards, for example, <i>my_label:xyz*</i></td>
    <td>Empty. Optional</td>
    <td></td>
  </tr>

</table>


## Testing your installation
Once the docker log collector software has been deployed in your environment, your container logs and incident detection will be available in the Zebrium UI.

## Contributors
* Brady Zuo (Zebrium)
