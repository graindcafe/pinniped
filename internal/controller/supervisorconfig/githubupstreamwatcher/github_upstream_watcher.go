// Copyright 2020-2024 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package githubupstreamwatcher implements a controller which watches GitHubIdentityProviders.
package githubupstreamwatcher

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1informers "k8s.io/client-go/informers/core/v1"

	"go.pinniped.dev/generated/latest/apis/supervisor/idp/v1alpha1"
	supervisorclientset "go.pinniped.dev/generated/latest/client/supervisor/clientset/versioned"
	idpinformers "go.pinniped.dev/generated/latest/client/supervisor/informers/externalversions/idp/v1alpha1"
	pinnipedcontroller "go.pinniped.dev/internal/controller"
	"go.pinniped.dev/internal/controller/conditionsutil"
	"go.pinniped.dev/internal/controllerlib"
	"go.pinniped.dev/internal/federationdomain/upstreamprovider"
	"go.pinniped.dev/internal/upstreamgithub"
)

const (
	// Setup for the name of our controller in logs.
	controllerName = "github-upstream-observer"

	// Constants related to the client credentials Secret.
	gitHubClientSecretType corev1.SecretType = "secrets.pinniped.dev/github-client"

	// // fixes lint to split from above group where const has explicit type
	// clientIDDataKey     = "clientID"
	// clientSecretDataKey = "clientSecret"
	//
	// // Constants related to conditions.
	// typeClientCredentialsValid = "ClientCredentialsValid" //nolint:gosec // this is not a credential.
)

// UpstreamGitHubIdentityProviderICache is a thread safe cache that holds a list of validated upstream GitHub IDP configurations.
type UpstreamGitHubIdentityProviderICache interface {
	SetGitHubIdentityProviders([]upstreamprovider.UpstreamGithubIdentityProviderI)
}

type gitHubWatcherController struct {
	cache                          UpstreamGitHubIdentityProviderICache
	log                            logr.Logger
	client                         supervisorclientset.Interface
	gitHubIdentityProviderInformer idpinformers.GitHubIdentityProviderInformer
	secretInformer                 corev1informers.SecretInformer
}

// New instantiates a new controllerlib.Controller which will populate the provided UpstreamGitHubIdentityProviderICache.
func New(
	idpCache UpstreamGitHubIdentityProviderICache,
	client supervisorclientset.Interface,
	gitHubIdentityProviderInformer idpinformers.GitHubIdentityProviderInformer,
	secretInformer corev1informers.SecretInformer,
	log logr.Logger,
	withInformer pinnipedcontroller.WithInformerOptionFunc,
) controllerlib.Controller {
	c := gitHubWatcherController{
		cache:                          idpCache,
		client:                         client,
		log:                            log.WithName(controllerName),
		gitHubIdentityProviderInformer: gitHubIdentityProviderInformer,
		secretInformer:                 secretInformer,
	}

	return controllerlib.New(
		controllerlib.Config{Name: controllerName, Syncer: &c},
		withInformer(
			gitHubIdentityProviderInformer,
			pinnipedcontroller.MatchAnythingFilter(pinnipedcontroller.SingletonQueue()),
			controllerlib.InformerOption{},
		),
		withInformer(
			secretInformer,
			pinnipedcontroller.MatchAnySecretOfTypeFilter(gitHubClientSecretType, pinnipedcontroller.SingletonQueue()),
			controllerlib.InformerOption{},
		),
	)
}

// Sync implements controllerlib.Syncer.
func (c *gitHubWatcherController) Sync(ctx controllerlib.Context) error {
	actualUpstreams, err := c.gitHubIdentityProviderInformer.Lister().List(labels.Everything())
	if err != nil {
		return fmt.Errorf("failed to list GitHubIdentityProviders: %w", err)
	}

	requeue := false
	validatedUpstreams := make([]upstreamprovider.UpstreamGithubIdentityProviderI, 0, len(actualUpstreams))
	for _, upstream := range actualUpstreams {
		valid := c.validateUpstream(ctx, upstream)
		if valid == nil {
			requeue = true
		} else {
			validatedUpstreams = append(validatedUpstreams, upstreamprovider.UpstreamGithubIdentityProviderI(valid))
		}
	}
	c.cache.SetGitHubIdentityProviders(validatedUpstreams)
	if requeue {
		return controllerlib.ErrSyntheticRequeue
	}
	return nil
}

func (c *gitHubWatcherController) validateUpstream(ctx controllerlib.Context, upstream *v1alpha1.GitHubIdentityProvider) *upstreamgithub.ProviderConfig {
	result := upstreamgithub.ProviderConfig{
		Name: upstream.Name,
	}

	conditions := []*metav1.Condition{
		// TODO: once we firm up the proposal doc & merge, then firm up the CRD & merge, we can
		// fill out these validations.
		// c.validateHost(),
		// c.validateTLS(),
		// c.validateAllowedOrganizations(),
		// c.validateOrganizationLoginPolicy(),
		// c.validateClient(),
	}

	c.updateStatus(ctx.Context, upstream, conditions)

	return &result
}

func (c *gitHubWatcherController) updateStatus(ctx context.Context, upstream *v1alpha1.GitHubIdentityProvider, conditions []*metav1.Condition) {
	log := c.log.WithValues("namespace", upstream.Namespace, "name", upstream.Name)
	updated := upstream.DeepCopy()

	hadErrorCondition := conditionsutil.MergeIDPConditions(conditions, upstream.Generation, &updated.Status.Conditions, log)

	updated.Status.Phase = v1alpha1.GitHubPhaseReady
	if hadErrorCondition {
		updated.Status.Phase = v1alpha1.GitHubPhaseError
	}

	if equality.Semantic.DeepEqual(upstream, updated) {
		return
	}

	_, err := c.client.
		IDPV1alpha1().
		GitHubIdentityProviders(upstream.Namespace).
		UpdateStatus(ctx, updated, metav1.UpdateOptions{})
	if err != nil {
		log.Error(err, "failed to update status")
	}
}
