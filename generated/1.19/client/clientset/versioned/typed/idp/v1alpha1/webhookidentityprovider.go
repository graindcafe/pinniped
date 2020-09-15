/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/suzerain-io/pinniped/generated/1.19/apis/idp/v1alpha1"
	scheme "github.com/suzerain-io/pinniped/generated/1.19/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WebhookIdentityProvidersGetter has a method to return a WebhookIdentityProviderInterface.
// A group's client should implement this interface.
type WebhookIdentityProvidersGetter interface {
	WebhookIdentityProviders(namespace string) WebhookIdentityProviderInterface
}

// WebhookIdentityProviderInterface has methods to work with WebhookIdentityProvider resources.
type WebhookIdentityProviderInterface interface {
	Create(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.CreateOptions) (*v1alpha1.WebhookIdentityProvider, error)
	Update(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.WebhookIdentityProvider, error)
	UpdateStatus(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.WebhookIdentityProvider, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.WebhookIdentityProvider, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.WebhookIdentityProviderList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookIdentityProvider, err error)
	WebhookIdentityProviderExpansion
}

// webhookIdentityProviders implements WebhookIdentityProviderInterface
type webhookIdentityProviders struct {
	client rest.Interface
	ns     string
}

// newWebhookIdentityProviders returns a WebhookIdentityProviders
func newWebhookIdentityProviders(c *IDPV1alpha1Client, namespace string) *webhookIdentityProviders {
	return &webhookIdentityProviders{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the webhookIdentityProvider, and returns the corresponding webhookIdentityProvider object, and an error if there is any.
func (c *webhookIdentityProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	result = &v1alpha1.WebhookIdentityProvider{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of WebhookIdentityProviders that match those selectors.
func (c *webhookIdentityProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WebhookIdentityProviderList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.WebhookIdentityProviderList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested webhookIdentityProviders.
func (c *webhookIdentityProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a webhookIdentityProvider and creates it.  Returns the server's representation of the webhookIdentityProvider, and an error, if there is any.
func (c *webhookIdentityProviders) Create(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.CreateOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	result = &v1alpha1.WebhookIdentityProvider{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a webhookIdentityProvider and updates it. Returns the server's representation of the webhookIdentityProvider, and an error, if there is any.
func (c *webhookIdentityProviders) Update(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	result = &v1alpha1.WebhookIdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		Name(webhookIdentityProvider.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *webhookIdentityProviders) UpdateStatus(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	result = &v1alpha1.WebhookIdentityProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		Name(webhookIdentityProvider.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(webhookIdentityProvider).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the webhookIdentityProvider and deletes it. Returns an error if one occurs.
func (c *webhookIdentityProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *webhookIdentityProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched webhookIdentityProvider.
func (c *webhookIdentityProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookIdentityProvider, err error) {
	result = &v1alpha1.WebhookIdentityProvider{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("webhookidentityproviders").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
