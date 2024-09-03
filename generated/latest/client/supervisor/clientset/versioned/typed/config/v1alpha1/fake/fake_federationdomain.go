// Copyright 2020-2024 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/latest/apis/supervisor/config/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFederationDomains implements FederationDomainInterface
type FakeFederationDomains struct {
	Fake *FakeConfigV1alpha1
	ns   string
}

var federationdomainsResource = v1alpha1.SchemeGroupVersion.WithResource("federationdomains")

var federationdomainsKind = v1alpha1.SchemeGroupVersion.WithKind("FederationDomain")

// Get takes name of the federationDomain, and returns the corresponding federationDomain object, and an error if there is any.
func (c *FakeFederationDomains) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FederationDomain, err error) {
	emptyResult := &v1alpha1.FederationDomain{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(federationdomainsResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// List takes label and field selectors, and returns the list of FederationDomains that match those selectors.
func (c *FakeFederationDomains) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FederationDomainList, err error) {
	emptyResult := &v1alpha1.FederationDomainList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(federationdomainsResource, federationdomainsKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FederationDomainList{ListMeta: obj.(*v1alpha1.FederationDomainList).ListMeta}
	for _, item := range obj.(*v1alpha1.FederationDomainList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested federationDomains.
func (c *FakeFederationDomains) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(federationdomainsResource, c.ns, opts))

}

// Create takes the representation of a federationDomain and creates it.  Returns the server's representation of the federationDomain, and an error, if there is any.
func (c *FakeFederationDomains) Create(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.CreateOptions) (result *v1alpha1.FederationDomain, err error) {
	emptyResult := &v1alpha1.FederationDomain{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(federationdomainsResource, c.ns, federationDomain, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// Update takes the representation of a federationDomain and updates it. Returns the server's representation of the federationDomain, and an error, if there is any.
func (c *FakeFederationDomains) Update(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.UpdateOptions) (result *v1alpha1.FederationDomain, err error) {
	emptyResult := &v1alpha1.FederationDomain{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(federationdomainsResource, c.ns, federationDomain, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFederationDomains) UpdateStatus(ctx context.Context, federationDomain *v1alpha1.FederationDomain, opts v1.UpdateOptions) (result *v1alpha1.FederationDomain, err error) {
	emptyResult := &v1alpha1.FederationDomain{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(federationdomainsResource, "status", c.ns, federationDomain, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}

// Delete takes name of the federationDomain and deletes it. Returns an error if one occurs.
func (c *FakeFederationDomains) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(federationdomainsResource, c.ns, name, opts), &v1alpha1.FederationDomain{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFederationDomains) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(federationdomainsResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FederationDomainList{})
	return err
}

// Patch applies the patch and returns the patched federationDomain.
func (c *FakeFederationDomains) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FederationDomain, err error) {
	emptyResult := &v1alpha1.FederationDomain{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(federationdomainsResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.FederationDomain), err
}
