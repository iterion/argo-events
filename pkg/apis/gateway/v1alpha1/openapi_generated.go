// +build !ignore_autogenerated

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
// Code generated by main. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.Gateway":       schema_pkg_apis_gateway_v1alpha1_Gateway(ref),
		"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewayList":   schema_pkg_apis_gateway_v1alpha1_GatewayList(ref),
		"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewaySpec":   schema_pkg_apis_gateway_v1alpha1_GatewaySpec(ref),
		"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewayStatus": schema_pkg_apis_gateway_v1alpha1_GatewayStatus(ref),
	}
}

func schema_pkg_apis_gateway_v1alpha1_Gateway(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Gateway is the definition of a gateway resource",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewayStatus"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewaySpec"),
						},
					},
				},
				Required: []string{"metadata", "status", "spec"},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewaySpec", "github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.GatewayStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_gateway_v1alpha1_GatewayList(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GatewayList is the list of Gateway resources",
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"),
						},
					},
					"items": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.Gateway"),
									},
								},
							},
						},
					},
				},
				Required: []string{"metadata", "items"},
			},
		},
		Dependencies: []string{
			"github.com/argoproj/argo-events/pkg/apis/gateway/v1alpha1.Gateway", "k8s.io/apimachinery/pkg/apis/meta/v1.ListMeta"},
	}
}

func schema_pkg_apis_gateway_v1alpha1_GatewaySpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GatewaySpec represents gateway specifications",
				Properties: map[string]spec.Schema{
					"deploySpec": {
						SchemaProps: spec.SchemaProps{
							Description: "DeploySpec is description of gateway",
							Ref:         ref("k8s.io/api/core/v1.PodSpec"),
						},
					},
					"configMap": {
						SchemaProps: spec.SchemaProps{
							Description: "ConfigMap is name of the configmap for gateway-processor",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"type": {
						SchemaProps: spec.SchemaProps{
							Description: "Type is type of the gateway",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"version": {
						SchemaProps: spec.SchemaProps{
							Description: "Version is used for marking event version",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"serviceSpec": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceSpec is the specifications of the service to expose the gateway",
							Ref:         ref("k8s.io/api/core/v1.ServiceSpec"),
						},
					},
					"sensors": {
						SchemaProps: spec.SchemaProps{
							Description: "Sensors are list of sensors to dispatch events to",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"rpcPort": {
						SchemaProps: spec.SchemaProps{
							Description: "RPCPort if provided is used to communicate between gRPC gateway client and gRPC gateway server",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"httpServerPort": {
						SchemaProps: spec.SchemaProps{
							Description: "HTTPServerPort if provided is used to communicate between gateway client and server over http",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"deploySpec", "type", "version", "sensors"},
			},
		},
		Dependencies: []string{
			"k8s.io/api/core/v1.PodSpec", "k8s.io/api/core/v1.ServiceSpec"},
	}
}

func schema_pkg_apis_gateway_v1alpha1_GatewayStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "GatewayStatus contains information about the status of a gateway.",
				Properties: map[string]spec.Schema{
					"phase": {
						SchemaProps: spec.SchemaProps{
							Description: "Phase is the high-level summary of the gateway",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"startedAt": {
						SchemaProps: spec.SchemaProps{
							Description: "StartedAt is the time at which this gateway was initiated",
							Ref:         ref("k8s.io/apimachinery/pkg/apis/meta/v1.Time"),
						},
					},
					"message": {
						SchemaProps: spec.SchemaProps{
							Description: "Message is a human readable string indicating details about a gateway in its phase",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"phase"},
			},
		},
		Dependencies: []string{
			"k8s.io/apimachinery/pkg/apis/meta/v1.Time"},
	}
}
