/*
Copyright 2020 VMware, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/suzerain-io/pinniped/generated/1.19/apis/idp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWebhookIdentityProviders implements WebhookIdentityProviderInterface
type FakeWebhookIdentityProviders struct {
	Fake *FakeIDPV1alpha1
	ns   string
}

var webhookidentityprovidersResource = schema.GroupVersionResource{Group: "idp.pinniped.dev", Version: "v1alpha1", Resource: "webhookidentityproviders"}

var webhookidentityprovidersKind = schema.GroupVersionKind{Group: "idp.pinniped.dev", Version: "v1alpha1", Kind: "WebhookIdentityProvider"}

// Get takes name of the webhookIdentityProvider, and returns the corresponding webhookIdentityProvider object, and an error if there is any.
func (c *FakeWebhookIdentityProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(webhookidentityprovidersResource, c.ns, name), &v1alpha1.WebhookIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookIdentityProvider), err
}

// List takes label and field selectors, and returns the list of WebhookIdentityProviders that match those selectors.
func (c *FakeWebhookIdentityProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WebhookIdentityProviderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(webhookidentityprovidersResource, webhookidentityprovidersKind, c.ns, opts), &v1alpha1.WebhookIdentityProviderList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.WebhookIdentityProviderList{ListMeta: obj.(*v1alpha1.WebhookIdentityProviderList).ListMeta}
	for _, item := range obj.(*v1alpha1.WebhookIdentityProviderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested webhookIdentityProviders.
func (c *FakeWebhookIdentityProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(webhookidentityprovidersResource, c.ns, opts))

}

// Create takes the representation of a webhookIdentityProvider and creates it.  Returns the server's representation of the webhookIdentityProvider, and an error, if there is any.
func (c *FakeWebhookIdentityProviders) Create(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.CreateOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(webhookidentityprovidersResource, c.ns, webhookIdentityProvider), &v1alpha1.WebhookIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookIdentityProvider), err
}

// Update takes the representation of a webhookIdentityProvider and updates it. Returns the server's representation of the webhookIdentityProvider, and an error, if there is any.
func (c *FakeWebhookIdentityProviders) Update(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.WebhookIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(webhookidentityprovidersResource, c.ns, webhookIdentityProvider), &v1alpha1.WebhookIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookIdentityProvider), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWebhookIdentityProviders) UpdateStatus(ctx context.Context, webhookIdentityProvider *v1alpha1.WebhookIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.WebhookIdentityProvider, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(webhookidentityprovidersResource, "status", c.ns, webhookIdentityProvider), &v1alpha1.WebhookIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookIdentityProvider), err
}

// Delete takes name of the webhookIdentityProvider and deletes it. Returns an error if one occurs.
func (c *FakeWebhookIdentityProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(webhookidentityprovidersResource, c.ns, name), &v1alpha1.WebhookIdentityProvider{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWebhookIdentityProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(webhookidentityprovidersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.WebhookIdentityProviderList{})
	return err
}

// Patch applies the patch and returns the patched webhookIdentityProvider.
func (c *FakeWebhookIdentityProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WebhookIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(webhookidentityprovidersResource, c.ns, name, pt, data, subresources...), &v1alpha1.WebhookIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WebhookIdentityProvider), err
}
