/*
Copyright 2018 BlackRock, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"github.com/argoproj/argo-events/common"
	"github.com/argoproj/argo-events/gateways"
	"github.com/argoproj/argo-events/gateways/core"
	"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/minio/minio-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"sync"
)

var (
	// gatewayConfig provides a generic configuration for a gateway
	gatewayConfig = gateways.NewGatewayConfiguration()
)

// S3Artifact contains information about an artifact in S3
type S3Artifact struct {
	// S3EventConfig contains configuration for bucket notification
	S3EventConfig S3EventConfig `json:"s3EventConfig" protobuf:"bytes,1,opt,name=s3EventConfig"`

	// Mode of operation for s3 client
	Insecure bool `json:"insecure,omitempty" protobuf:"bytes,2,opt,name=insecure"`

	// AccessKey
	AccessKey corev1.SecretKeySelector `json:"accessKey,omitempty" protobuf:"bytes,3,opt,name=accessKey"`

	// SecretKey
	SecretKey corev1.SecretKeySelector `json:"secretKey,omitempty" protobuf:"bytes,4,opt,name=secretKey"`
}

// S3EventConfig contains configuration for bucket notification
type S3EventConfig struct {
	Endpoint string                      `json:"endpoint,omitempty" protobuf:"bytes,1,opt,name=endpoint"`
	Bucket   string                      `json:"bucket,omitempty" protobuf:"bytes,2,opt,name=bucket"`
	Region   string                      `json:"region,omitempty" protobuf:"bytes,3,opt,name=region"`
	Event    minio.NotificationEventType `json:"event,omitempty" protobuf:"bytes,4,opt,name=event"`
	filter   S3Filter                    `json:"filter,omitempty" protobuf:"bytes,5,opt,name=filter"`
}

// S3Filter represents filters to apply to bucket nofifications for specifying constraints on objects
type S3Filter struct {
	Prefix string `json:"prefix" protobuf:"bytes,1,opt,name=prefix"`
	Suffix string `json:"suffix" protobuf:"bytes,2,opt,name=suffix"`
}

// getSecrets retrieves the secret value from the secret in namespace with name and key
func getSecrets(client kubernetes.Interface, namespace string, name, key string) (string, error) {
	secretsIf := client.CoreV1().Secrets(namespace)
	var secret *corev1.Secret
	var err error
	_ = wait.ExponentialBackoff(common.DefaultRetry, func() (bool, error) {
		secret, err = secretsIf.Get(name, metav1.GetOptions{})
		if err != nil {
			if !common.IsRetryableKubeAPIError(err) {
				return false, err
			}
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return "", err
	}
	val, ok := secret.Data[key]
	if !ok {
		return "", fmt.Errorf("secret '%s' does not have the key '%s'", name, key)
	}
	return string(val), nil
}

// Runs a configuration
func configRunner(config *gateways.ConfigContext) error {
	var err error
	var errMessage string

	// mark final gateway state
	defer gatewayConfig.GatewayCleanup(config, errMessage, err)

	gatewayConfig.Log.Info().Str("config-name", config.Data.Src).Msg("parsing configuration...")

	var artifact *S3Artifact
	err = yaml.Unmarshal([]byte(config.Data.Config), &artifact)
	if err != nil {
		errMessage = "failed to parse configuration"
		return err
	}

	gatewayConfig.Log.Debug().Str("config-key", config.Data.Config).Interface("artifact", *artifact).Msg("s3 artifact")

	// retrieve access key id and secret access key
	accessKey, err := getSecrets(gatewayConfig.Clientset, gatewayConfig.Namespace, artifact.AccessKey.Name, artifact.AccessKey.Key)
	if err != nil {
		errMessage = fmt.Sprintf("failed to retrieve access key id %s", artifact.AccessKey.Name)
		return err
	}
	secretKey, err := getSecrets(gatewayConfig.Clientset, gatewayConfig.Namespace, artifact.SecretKey.Name, artifact.SecretKey.Key)
	if err != nil {
		errMessage = fmt.Sprintf("failed to retrieve access key id %s", artifact.SecretKey.Name)
		return err
	}

	gatewayConfig.Log.Debug().Str("accesskey", accessKey).Str("secretaccess", secretKey).Msg("minio secrets")

	minioClient, err := minio.New(artifact.S3EventConfig.Endpoint, accessKey, secretKey, !artifact.Insecure)
	if err != nil {
		errMessage = "failed to get minio client"
		return err
	}

	// Create a done channel to control 'ListenBucketNotification' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	var wg sync.WaitGroup
	wg.Add(1)

	// waits till disconnection from client.
	go func() {
		<-config.StopCh
		config.Active = false
		gatewayConfig.Log.Info().Str("config", config.Data.Src).Msg("stopping the configuration...")
		wg.Done()
	}()

	gatewayConfig.Log.Info().Str("config-name", config.Data.Src).Msg("configuration is running...")
	config.Active = true

	event := gatewayConfig.GetK8Event("configuration running", v1alpha1.NodePhaseRunning, config.Data.Src)
	_, err = common.CreateK8Event(event, gatewayConfig.Clientset)
	if err != nil {
		gatewayConfig.Log.Error().Str("config-key", config.Data.Src).Err(err).Msg("failed to mark configuration as running")
		return err
	}

	// Listen for bucket notifications
	for notificationInfo := range minioClient.ListenBucketNotification(artifact.S3EventConfig.Bucket, artifact.S3EventConfig.filter.Prefix, artifact.S3EventConfig.filter.Suffix, []string{
		string(artifact.S3EventConfig.Event),
	}, doneCh) {
		if notificationInfo.Err != nil {
			err = notificationInfo.Err
			errMessage = "notification error"
			config.StopCh <- struct{}{}
			break
		} else {
			gatewayConfig.Log.Info().Str("config-key", config.Data.Src).Msg("dispatching event to gateway-processor")
			payload := []byte(fmt.Sprintf("%v", notificationInfo))
			gatewayConfig.DispatchEvent(&gateways.GatewayEvent{
				Src:     config.Data.Src,
				Payload: payload,
			})
		}
	}
	wg.Wait()
	return nil
}

func main() {
	_, err := gatewayConfig.WatchGatewayEvents(context.Background())
	if err != nil {
		gatewayConfig.Log.Panic().Err(err).Msg("failed to watch k8 events for gateway configuration state updates")
	}
	_, err = gatewayConfig.WatchGatewayConfigMap(context.Background(), configRunner, core.ConfigDeactivator)
	if err != nil {
		gatewayConfig.Log.Panic().Err(err).Msg("failed to watch gateway configuration updates")
	}
	select {}
}
