# DOCKER CONTAINER LOG COLLECTOR

Zebrium's docker container log collector collects container logs and sends logs to Zebrium for automated Incident detection.  This is achieved by using the [Fluentd logging driver for Docker](https://docs.docker.com/config/containers/logging/fluentd/) and the Zebrium Fluentd [output plugin.](https://github.com/zebrium/fluentd-output-zebrium)

## Getting Started

When sending your logs from your docker daemon to Zebrium, there are two configuration options for where your log collector can be installed in configured.  The collector can be installed within the docker daemon context that you are sending all the logs from, or it could be installed on an external host, and have the logs routed to it by each docker daemon.

### Deploying the Collector

Regardless on the installation method, you will start the collector using the following command, substituting the token and URL in for the values found in your Zebrium Integration and Collectors page.  Additional ENVS listed [below](#environment-variables) can be specified to the collector to further extend the functionality.

```bash
docker run -p 24224:24224 -e ZE_LOG_COLLECTOR_URL=<URL> -e ZE_LOG_COLLECTOR_TOKEN=<TOKEN> --restart always zebrium/docker-log-collector:latest
```

### Configuring the Docker daemon

Once our collector has been deployed and configured, we need to modify the docker daemon configuration to start sending logs to the collector.  For a complete list of configuration options, please see the [official docker documentation](https://docs.docker.com/config/containers/logging/fluentd/).  The docker daemon is located in `/etc/docker/daemon.json` on Linux host and `C:\ProgramData\docker\config\daemon.json` on windows host.  For more about the docker daemon.json, see the [official documentation](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file)

Add the following configuration to your daemon.json file, substituting `<fluentd-address>` for the address of your log collector.  If your log collector is deployed in the same docker daemon, then use `127.0.0.1:24224` as your address.  

```bash
{
"log-driver": "fluentd",
  "log-opts": {
    "fluentd-address": "<fluentd-address>",
    "fluentd-async": "true"
  }
}
```

Once the daemon file is updated, restart the docker daemon for the new changes to take effect.  After this, your should be able to view the logs of the log collector and verify that it is receiving and forwarding logs to Zebrium.

### Environment Variables

Below is a list of environment variables that are available for configuration of the Fluentd container.

| Environment Variables | Default | Description | Required |
|-------------------|-------------------|-------------------| ---|
| ZE_LOG_COLLECTOR_URL | "" | Zebrium URL Endpoint for log ingestion| yes|
| ZE_LOG_COLLECTOR_TOKEN | "" | Zebrium ZAPI token for log ingestion| yes|
| ZE_DEPLOYMENT_NAME | "default" | Zebrium Service Group Name.  Read more [here](https://docs.sciencelogic.com/zebrium/latest/Content/Web_Zebrium/Key_Concepts.html#service-groups)| no|
| FLUSH_INTERVAL | "60s" | Buffer Flush Interval| no|
| ZE_LOG_LEVEL | "info" | Sets the log level for the output plugin | no |
| VERIFY_SSL | "true" | Enables or disables SSL verification on endpoint| no|

## Additional Resources

* [Github repository](https://github.com/zebrium/ze-docker-log-collector)
* [Fluentd Documentation](https://www.fluentd.org/guides/recipes/docker-logging)
* [Docker Plugin Documentation](https://docs.docker.com/config/containers/logging/fluentd/)
