// Copyright 2020-2022 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "go.pinniped.dev/generated/1.17/apis/concierge/identity/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	testing "k8s.io/client-go/testing"
)

// FakeWhoAmIRequests implements WhoAmIRequestInterface
type FakeWhoAmIRequests struct {
	Fake *FakeIdentityV1alpha1
}

var whoamirequestsResource = schema.GroupVersionResource{Group: "identity.concierge.pinniped.dev", Version: "v1alpha1", Resource: "whoamirequests"}

var whoamirequestsKind = schema.GroupVersionKind{Group: "identity.concierge.pinniped.dev", Version: "v1alpha1", Kind: "WhoAmIRequest"}

// Create takes the representation of a whoAmIRequest and creates it.  Returns the server's representation of the whoAmIRequest, and an error, if there is any.
func (c *FakeWhoAmIRequests) Create(whoAmIRequest *v1alpha1.WhoAmIRequest) (result *v1alpha1.WhoAmIRequest, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(whoamirequestsResource, whoAmIRequest), &v1alpha1.WhoAmIRequest{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WhoAmIRequest), err
}
