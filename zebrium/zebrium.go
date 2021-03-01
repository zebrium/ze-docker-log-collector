package zebrium

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gliderlabs/logspout/router"
	"github.com/zebrium/ze-docker-log-collector/zebrium/adapter"
)

const (
	adapterName          = "zebrium"
	filterNameEnvVar     = "ZE_FILTER_NAME"
	filterLabelsEnvVar   = "ZE_FILTER_LABELS"
	ZapiUrlEnvVar        = "ZE_LOG_COLLECTOR_URL"
	ZapiTokenEnvVar      = "ZE_LOG_COLLECTOR_TOKEN"
	VerifySslEnvVar      = "ZE_VERIFY_SSL"
	DeploymentNameEnvVar = "ZE_DEPLOYMENT_NAME"
	HostnameEnvvar       = "ZE_HOSTNAME"
	MaxIngestSizeEnvVar  = "ZE_MAX_INGEST_SIZE"
	FlushTimeoutEnvVar   = "ZE_FLUSH_TIMEOUT"
)

func init() {
	log.Println("zebrium: init() called")
	router.AdapterFactories.Register(NewZebriumAdapter, adapterName)

	filterName := os.Getenv(filterNameEnvVar)
	filterLabels := make([]string, 0)
	filterLabelsVal := os.Getenv(filterLabelsEnvVar)
	if filterLabelsVal != "" {
		filterLabels = strings.Split(filterLabelsVal, ",")
	}

	r := &router.Route{
		Adapter:      adapterName,
		FilterName:   filterName,
		FilterLabels: filterLabels,
	}
	if err := router.Routes.Add(r); err != nil {
		log.Fatal("Failed to add New Route: ", err.Error())
	}
	log.Println("zebrium: init() done")
}

func NewZebriumAdapter(route *router.Route) (router.LogAdapter, error) {
	url := os.Getenv(ZapiUrlEnvVar)
	url = strings.Trim(url, " \t\"'")
	if url == "" {
		log.Fatal("Environment variable ", ZapiUrlEnvVar, " is not set")
	}
	if strings.HasSuffix(url, "zebrium.com") {
		url = url + "/log/api/v2/ingest"
	}

	token := os.Getenv(ZapiTokenEnvVar)
	token = strings.Trim(token, " \t\"'")
	if token == "" {
		log.Fatal("Environment variable ", ZapiTokenEnvVar, " is not set")
	}

	verifySsl := true
	verifySslStr := os.Getenv(VerifySslEnvVar)
	if strings.EqualFold(verifySslStr, "false") {
		verifySsl = false
	}

	deploymentName := os.Getenv(DeploymentNameEnvVar)
	deploymentName = strings.Trim(deploymentName, " \t\"'")
	if deploymentName == "" {
		deploymentName = "default"
		log.Println("Use default deployment name ", deploymentName)
	}

	ingestSizeStr := os.Getenv(MaxIngestSizeEnvVar)
	if ingestSizeStr == "" {
		ingestSizeStr = "1048576"
	}
	ingestSize, err := strconv.Atoi(ingestSizeStr)
	if err != nil {
		log.Fatal("Max ingest size setting is invalid: ", err.Error())
	}

	flushTimeoutStr := os.Getenv(FlushTimeoutEnvVar)
	if flushTimeoutStr == "" {
		flushTimeoutStr = "30"
	}
	flushTimeout, err := strconv.Atoi(flushTimeoutStr)
	if err != nil {
		log.Fatal("Flush timeout setting is invalid: ", err.Error())
	}

	log.Println("zebrium: create new adapter")
	return adapter.New(
		url,
		token,
		verifySsl,
		deploymentName,
		os.Getenv(HostnameEnvvar),
		ingestSize,
		flushTimeout,
	), nil
}
