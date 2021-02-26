package adapter

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gliderlabs/logspout/router"
)

const (
	collectorVers = "1.47.0-docker"
)

type Adapter struct {
	MaxIngestSize  int
	ZapiUrl        string
	ZapiToken      string
	VerifySsl      bool
	DeploymentName string
	Hostname       string
	FlushTimeout   time.Duration
	Queue        chan ContainerLogMessage
}

// Message structure:
type ContainerLogMessage struct {
	Message          string        `json:"message"`
	Source           string        `json:"source"`
	// Timestamp on log message from container log
	EpochNanos       int64         `json:"epoch_nanos"`
	Collector        string        `json:"collector"`
	ZeDeploymentName string        `json:"ze_deployment_name"`
	Container        ContainerMeta `json:"container"`
}

type ContainerMeta struct {
	Name      string             `json:"name"`
	Id        string             `json:"id"`
	Image     string             `json:"image"`
	Hostname  string             `json:"hostname"`
	Labels    map[string]string  `json:"labels"`
}

type IngestRequest struct {
	LogType   string             `json:"log_type"`
	Messages  []string           `json:"messages"`
}

func New(zapiUrl string, zapiToken string, verifySsl bool,
	 deploymentName string, hostname string,
	 maxIngestSize int, flushTimeout int) *Adapter {
	log.Println("hostname=" + hostname)
	adapter := &Adapter{
		MaxIngestSize:  maxIngestSize,
		ZapiUrl:        zapiUrl,
		ZapiToken:      zapiToken,
		VerifySsl:      verifySsl,
		DeploymentName: deploymentName,
		Hostname:       hostname,
		FlushTimeout:   time.Duration(flushTimeout) * time.Second,
		Queue:      make(chan ContainerLogMessage),
	}

	log.Printf("zapiUrl=%s maxIngestSize=%d flushTimeout=%d\n", zapiUrl, maxIngestSize, flushTimeout)

	go adapter.readQueue()
	return adapter
}

func (a *Adapter) Stream(logstream chan *router.Message) {
	for m := range logstream {
		hostname := a.Hostname
		if hostname == "" {
			hostname = m.Container.Config.Hostname
		}
		clm := ContainerLogMessage{
			Message:           m.Data,
			Source:            m.Source,
			EpochNanos:        m.Time.UnixNano(),
			Collector:         collectorVers,
			ZeDeploymentName:  a.DeploymentName,
                        Container:  ContainerMeta {
				Name:     strings.Trim(m.Container.Name, "/"),
				Id:       m.Container.ID,
				Image:    m.Container.Config.Image,
				Hostname: hostname,
				Labels:   m.Container.Config.Labels,
			},
		}
		a.Queue <- clm
	}
}

func (a *Adapter) readQueue() {
	msgData := make([]string, 0)
	dataLen := 0
	timeout := time.NewTimer(a.FlushTimeout)

	for {
		select {
		case m := <-a.Queue:
			if dataLen >= a.MaxIngestSize {
				timeout.Stop()
				a.send(msgData)
				dataLen = 0
				msgData = make([]string, 0)
				timeout.Reset(a.FlushTimeout)
			}
			mstr, err := json.Marshal(m)
			if err == nil {
				msgData = append(msgData, string(mstr))
				dataLen += len(mstr)
			} else {
				log.Println("Failed to marsal container log message: " + err.Error())
			}

		case <-timeout.C:
			if len(msgData) > 0 {
				a.send(msgData)
				msgData = make([]string, 0)
			}
			timeout.Reset(a.FlushTimeout)
		}
	}
}

func (a *Adapter) send(msgData []string) {
	var data bytes.Buffer

	ir := IngestRequest{
			LogType: "docker_container",
			Messages: msgData,
	}

	err := json.NewEncoder(&data).Encode(ir)
	if err != nil {
		log.Println("Failed to marshal ingest request: " + err.Error())
		return
	}

	req, _ := http.NewRequest(http.MethodPost, a.ZapiUrl, &data)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", "Token " + a.ZapiToken)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : !a.VerifySsl},
	}
	httpClient := &http.Client{Transport: tr}
	resp, err := httpClient.Do(req)

	if err != nil {
		log.Println("Post request failed: " + err.Error())
		return
	}

	if resp != nil {
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)
			log.Printf("Received Status Code: %d\n", resp.StatusCode)
			log.Printf("             Message: %s\n", string(body))
		}
		defer resp.Body.Close()
        }
}
