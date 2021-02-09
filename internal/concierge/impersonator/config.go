// Copyright 2021 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package impersonator

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type Mode string

const (
	// Explicitly enable the impersonation proxy.
	ModeEnabled Mode = "enabled"

	// Explicitly disable the impersonation proxy.
	ModeDisabled Mode = "disabled"

	// Allow the proxy to decide if it should be enabled or disabled based upon the cluster in which it is running.
	ModeAuto Mode = "auto"
)

const (
	ConfigMapDataKey = "config.yaml"
)

// When specified, both CertificateAuthoritySecretName and TLSSecretName are required. They may be specified to
// both point at the same Secret or to point at different Secrets.
type TLSConfig struct {
	// CertificateAuthoritySecretName contains the name of a namespace-local Secret resource. The corresponding Secret
	// must contain a key called "ca.crt" whose value is the CA certificate which clients should trust when connecting
	// to the impersonation proxy.
	CertificateAuthoritySecretName string `json:"certificateAuthoritySecretName"`

	// TLSSecretName contains the name of a namespace-local Secret resource. The corresponding Secret must be of type
	// "kubernetes.io/tls" and contain keys called "tls.crt" and "tls.key" whose values are the TLS certificate and
	// private key that will be used by the impersonation proxy to serve its endpoints.
	TLSSecretName string `json:"tlsSecretName"`
}

type Config struct {
	// Enable or disable the impersonation proxy. Optional. Defaults to ModeAuto.
	Mode Mode `json:"mode,omitempty"`

	// The HTTPS URL of the impersonation proxy for clients to use from outside the cluster. Used when creating TLS
	// certificates and for clients to discover the endpoint. Optional. When not specified, if the impersonation proxy
	// is started, then it will automatically create a LoadBalancer Service and use its ingress as the endpoint.
	Endpoint string `json:"endpoint,omitempty"`

	// The TLS configuration of the impersonation proxy's endpoints. Optional. When not specified, a CA and TLS
	// certificate will be automatically created based on the Endpoint setting.
	TLS *TLSConfig `json:"tls,omitempty"`
}

func FromConfigMap(configMap *v1.ConfigMap) (*Config, error) {
	stringConfig, ok := configMap.Data[ConfigMapDataKey]
	if !ok {
		return nil, fmt.Errorf(`ConfigMap is missing expected key "%s"`, ConfigMapDataKey)
	}
	var config Config
	if err := yaml.Unmarshal([]byte(stringConfig), &config); err != nil {
		return nil, fmt.Errorf("decode yaml: %w", err)
	}
	if config.Mode == "" {
		config.Mode = ModeAuto // set the default value
	}
	return &config, nil
}
