// Copyright 2020 Red Hat, Inc. and/or its affiliates
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

package kubernetes

import (
	"io/ioutil"

	"github.com/kiegroup/kogito-cloud-operator/pkg/client"
	corev1 "k8s.io/api/core/v1"
)

// PodInterface has functions that interacts with pod object in the Kubernetes cluster
type PodInterface interface {
	// Immediately return pod log
	GetLogs(namespace, podName, containerName string) (string, error)
	// Wait until pod is terminated and then return pod log
	GetLogsWithFollow(namespace, podName, containerName string) (string, error)
}

type pod struct {
	client *client.Client
}

func newPod(c *client.Client) PodInterface {
	client.MustEnsureClient(c)
	return &pod{
		client: c,
	}
}

func (pod *pod) GetLogs(namespace, podName, containerName string) (string, error) {
	return pod.getLogs(namespace, podName, containerName, false)
}

func (pod *pod) GetLogsWithFollow(namespace, podName, containerName string) (string, error) {
	return pod.getLogs(namespace, podName, containerName, true)
}

func (pod *pod) getLogs(namespace, podName, containerName string, follow bool) (string, error) {
	log.Debugf("About to fetch log of pod with name '%s' in namespace %s from cluster with follow '%t'", podName, namespace, follow)
	podLogOpts := corev1.PodLogOptions{
		Follow:    follow,
		Container: containerName,
	}
	req := pod.client.KubernetesExtensionCli.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)
	readCloser, err := req.Stream()
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(readCloser)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
