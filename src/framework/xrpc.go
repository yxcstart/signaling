package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"signaling/src/framework/xrpc"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var xrpcClients map[string]*xrpc.Client = make(map[string]*xrpc.Client)

func SetupXprc() error {
	sections := configFile.GetSectionList()
	for _, section := range sections {
		if !strings.HasPrefix(section, "xrpc.") {
			continue
		}

		mSection, err := configFile.GetSection(section)
		if err != nil {
			return err
		}
		values, ok := mSection["server"]
		if !ok {
			return fmt.Errorf("no server field in config file")
		}
		arrServer := strings.Split(values, ",")
		for i, server := range arrServer {
			arrServer[i] = strings.TrimSpace(server)
		}

		client := xrpc.NewClient(arrServer)

		if values, ok := mSection["connectTimeout"]; ok {
			if connectTimeout, err := strconv.Atoi(values); err == nil {
				client.ConnectTimeout = time.Duration(connectTimeout) * time.Millisecond
			}
		}

		if values, ok := mSection["readTimeout"]; ok {
			if readTimeout, err := strconv.Atoi(values); err == nil {
				client.ReadTimeout = time.Duration(readTimeout) * time.Millisecond
			}
		}

		if values, ok := mSection["writeTimeout"]; ok {
			if writeTimeout, err := strconv.Atoi(values); err == nil {
				client.WriteTimeout = time.Duration(writeTimeout) * time.Millisecond
			}
		}

		xrpcClients[section] = client
	}

	return nil
}

func Call(serviceName string, req interface{}, rsp interface{}, logId uint32) error {
	log.Infof("xrpc call %s req %+v", serviceName, req)

	client, ok := xrpcClients["xrpc."+serviceName]
	if !ok {
		return fmt.Errorf("[%s] service not found", serviceName)
	}
	content, err := json.Marshal(req)
	if err != nil {
		return err
	}

	request := xrpc.NewRequest(bytes.NewReader(content), logId)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	log.Infof("xrpc call %s rsp %+v", serviceName, resp)

	return nil
}
