// Copyright 2020-2022 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "go.pinniped.dev/generated/1.22/apis/supervisor/idp/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeActiveDirectoryIdentityProviders implements ActiveDirectoryIdentityProviderInterface
type FakeActiveDirectoryIdentityProviders struct {
	Fake *FakeIDPV1alpha1
	ns   string
}

var activedirectoryidentityprovidersResource = schema.GroupVersionResource{Group: "idp.supervisor.pinniped.dev", Version: "v1alpha1", Resource: "activedirectoryidentityproviders"}

var activedirectoryidentityprovidersKind = schema.GroupVersionKind{Group: "idp.supervisor.pinniped.dev", Version: "v1alpha1", Kind: "ActiveDirectoryIdentityProvider"}

// Get takes name of the activeDirectoryIdentityProvider, and returns the corresponding activeDirectoryIdentityProvider object, and an error if there is any.
func (c *FakeActiveDirectoryIdentityProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(activedirectoryidentityprovidersResource, c.ns, name), &v1alpha1.ActiveDirectoryIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActiveDirectoryIdentityProvider), err
}

// List takes label and field selectors, and returns the list of ActiveDirectoryIdentityProviders that match those selectors.
func (c *FakeActiveDirectoryIdentityProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ActiveDirectoryIdentityProviderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(activedirectoryidentityprovidersResource, activedirectoryidentityprovidersKind, c.ns, opts), &v1alpha1.ActiveDirectoryIdentityProviderList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ActiveDirectoryIdentityProviderList{ListMeta: obj.(*v1alpha1.ActiveDirectoryIdentityProviderList).ListMeta}
	for _, item := range obj.(*v1alpha1.ActiveDirectoryIdentityProviderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested activeDirectoryIdentityProviders.
func (c *FakeActiveDirectoryIdentityProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(activedirectoryidentityprovidersResource, c.ns, opts))

}

// Create takes the representation of a activeDirectoryIdentityProvider and creates it.  Returns the server's representation of the activeDirectoryIdentityProvider, and an error, if there is any.
func (c *FakeActiveDirectoryIdentityProviders) Create(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.CreateOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(activedirectoryidentityprovidersResource, c.ns, activeDirectoryIdentityProvider), &v1alpha1.ActiveDirectoryIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActiveDirectoryIdentityProvider), err
}

// Update takes the representation of a activeDirectoryIdentityProvider and updates it. Returns the server's representation of the activeDirectoryIdentityProvider, and an error, if there is any.
func (c *FakeActiveDirectoryIdentityProviders) Update(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(activedirectoryidentityprovidersResource, c.ns, activeDirectoryIdentityProvider), &v1alpha1.ActiveDirectoryIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActiveDirectoryIdentityProvider), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeActiveDirectoryIdentityProviders) UpdateStatus(ctx context.Context, activeDirectoryIdentityProvider *v1alpha1.ActiveDirectoryIdentityProvider, opts v1.UpdateOptions) (*v1alpha1.ActiveDirectoryIdentityProvider, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(activedirectoryidentityprovidersResource, "status", c.ns, activeDirectoryIdentityProvider), &v1alpha1.ActiveDirectoryIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActiveDirectoryIdentityProvider), err
}

// Delete takes name of the activeDirectoryIdentityProvider and deletes it. Returns an error if one occurs.
func (c *FakeActiveDirectoryIdentityProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(activedirectoryidentityprovidersResource, c.ns, name), &v1alpha1.ActiveDirectoryIdentityProvider{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeActiveDirectoryIdentityProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(activedirectoryidentityprovidersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ActiveDirectoryIdentityProviderList{})
	return err
}

// Patch applies the patch and returns the patched activeDirectoryIdentityProvider.
func (c *FakeActiveDirectoryIdentityProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ActiveDirectoryIdentityProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(activedirectoryidentityprovidersResource, c.ns, name, pt, data, subresources...), &v1alpha1.ActiveDirectoryIdentityProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ActiveDirectoryIdentityProvider), err
}
