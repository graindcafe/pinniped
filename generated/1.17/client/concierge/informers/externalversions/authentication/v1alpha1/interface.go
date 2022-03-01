// Copyright 2020-2022 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "go.pinniped.dev/generated/1.17/client/concierge/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// JWTAuthenticators returns a JWTAuthenticatorInformer.
	JWTAuthenticators() JWTAuthenticatorInformer
	// WebhookAuthenticators returns a WebhookAuthenticatorInformer.
	WebhookAuthenticators() WebhookAuthenticatorInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// JWTAuthenticators returns a JWTAuthenticatorInformer.
func (v *version) JWTAuthenticators() JWTAuthenticatorInformer {
	return &jWTAuthenticatorInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// WebhookAuthenticators returns a WebhookAuthenticatorInformer.
func (v *version) WebhookAuthenticators() WebhookAuthenticatorInformer {
	return &webhookAuthenticatorInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}
