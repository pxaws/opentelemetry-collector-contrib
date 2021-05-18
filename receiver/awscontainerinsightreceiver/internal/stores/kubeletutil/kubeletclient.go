// Copyright  OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeletutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awscontainerinsightreceiver/internal/stores/tls"
)

type KubeClient struct {
	Port            string
	BearerToken     string
	KubeIP          string
	responseTimeout time.Duration
	roundTripper    http.RoundTripper
	tls.ClientConfig
}

var ErrKubeClientAccessFailure = errors.New("KubeClinet Access Failure")

func (k *KubeClient) ListPods() ([]corev1.Pod, error) {
	var result []corev1.Pod
	var req *http.Request
	url := fmt.Sprintf("https://%s:%s/pods", k.KubeIP, k.Port)

	req, _ = http.NewRequest("GET", url, nil)
	var resp *http.Response

	k.InsecureSkipVerify = true
	tlsCfg, err := k.ClientConfig.TLSConfig()
	if err != nil {
		return result, err
	}

	if k.roundTripper == nil {
		// Set default values
		if k.responseTimeout < time.Second {
			k.responseTimeout = time.Second * 5
		}
		k.roundTripper = &http.Transport{
			TLSHandshakeTimeout:   5 * time.Second,
			TLSClientConfig:       tlsCfg,
			ResponseHeaderTimeout: k.responseTimeout,
		}
	}

	if k.BearerToken != "" {
		var token []byte
		if token, err = ioutil.ReadFile(k.BearerToken); err != nil {
			return result, err
		}
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(string(token)))
	}
	req.Header.Add("Accept", "application/json")

	resp, err = k.roundTripper.RoundTrip(req)
	if err != nil {
		log.Printf("E! error making HTTP request to %s: %s", url, err)
		return result, ErrKubeClientAccessFailure
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("E! %s returned HTTP status %s", url, resp.Status)
		return result, ErrKubeClientAccessFailure
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("E! Fail to read request %s body: %s", req.URL.String(), err)
		return result, err
	}

	pods := corev1.PodList{}
	err = json.Unmarshal(b, &pods)
	if err != nil {
		log.Printf("E! parsing response: %s", err)
		return result, err
	}

	return pods.Items, nil
}
