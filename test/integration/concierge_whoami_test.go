// Copyright 2020-2024 the Pinniped contributors. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0
package integration

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	authenticationv1 "k8s.io/api/authentication/v1"
	certificatesv1 "k8s.io/api/certificates/v1"
	certificatesv1beta1 "k8s.io/api/certificates/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/certificate/csr"
	"k8s.io/client-go/util/keyutil"

	identityv1alpha1 "go.pinniped.dev/generated/latest/apis/concierge/identity/v1alpha1"
	"go.pinniped.dev/internal/testutil"
	"go.pinniped.dev/test/testlib"
)

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_Kubeadm_Parallel(t *testing.T) {
	// use the cluster signing key being available as a proxy for this being a kubeadm cluster
	// we should add more robust logic around skipping clusters based on vendor
	_ = testlib.IntegrationEnv(t).WithCapability(testlib.ClusterSigningKeyIsAvailable)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	adminClient := testlib.NewKubernetesClientset(t)

	whoAmI, err := testlib.NewConciergeClientset(t).IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	var wantGroups []string
	if testutil.KubeServerMinorVersionInBetweenInclusive(t, adminClient.Discovery(), 0, 28) {
		wantGroups = []string{"system:masters", "system:authenticated"}
	} else {
		// See https://github.com/kubernetes/enhancements/issues/4214. Admin kubeconfigs from kubeadm
		// which previously had system:masters now have kubeadm:cluster-admins instead.
		wantGroups = []string{"kubeadm:cluster-admins", "system:authenticated"}
	}

	// this user info is based off of the bootstrap cert user created by kubeadm
	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "kubernetes-admin",
						Groups:   wantGroups,
					},
				},
			},
		},
		whoAmI,
	)
}

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_ServiceAccount_Legacy_Parallel(t *testing.T) {
	_ = testlib.IntegrationEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeClient := testlib.NewKubernetesClientset(t).CoreV1()

	ns, err := kubeClient.Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-whoami-",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, kubeClient.Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{}))
	})

	sa, err := kubeClient.ServiceAccounts(ns.Name).Create(ctx, &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-whoami-",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	secret, err := kubeClient.Secrets(ns.Name).Create(ctx, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-whoami-",
			Annotations: map[string]string{
				corev1.ServiceAccountNameKey: sa.Name,
			},
		},
		Type: corev1.SecretTypeServiceAccountToken,
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	testlib.RequireEventuallyWithoutError(t, func() (bool, error) {
		secret, err = kubeClient.Secrets(ns.Name).Get(ctx, secret.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return len(secret.Data[corev1.ServiceAccountTokenKey]) > 0, nil
	}, time.Minute, time.Second)

	saConfig := testlib.NewAnonymousClientRestConfig(t)
	saConfig.BearerToken = string(secret.Data[corev1.ServiceAccountTokenKey])

	whoAmI, err := testlib.NewKubeclient(t, saConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err)

	// legacy service account tokens do not have any extra info
	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "system:serviceaccount:" + ns.Name + ":" + sa.Name,
						UID:      "", // aggregation drops UID: https://github.com/kubernetes/kubernetes/issues/93699
						Groups: []string{
							"system:serviceaccounts",
							"system:serviceaccounts:" + ns.Name,
							"system:authenticated",
						},
					},
				},
			},
		},
		whoAmI,
	)
}

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_ServiceAccount_TokenRequest_Parallel(t *testing.T) {
	env := testlib.IntegrationEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeClient := testlib.NewKubernetesClientset(t)
	coreV1client := kubeClient.CoreV1()

	ns, err := coreV1client.Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-whoami-",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, coreV1client.Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{}))
	})

	sa, err := coreV1client.ServiceAccounts(ns.Name).Create(ctx, &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "test-whoami-",
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	_, tokenRequestProbeErr := coreV1client.ServiceAccounts(ns.Name).CreateToken(ctx, sa.Name, &authenticationv1.TokenRequest{}, metav1.CreateOptions{})
	if errors.IsNotFound(tokenRequestProbeErr) && tokenRequestProbeErr.Error() == "the server could not find the requested resource" {
		return // stop test early since the token request API is not enabled on this cluster - other errors are caught below
	}

	pod := testlib.CreatePod(ctx, t, "whoami", ns.Name,
		corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "sleeper",
					Image:           env.ShellContainerImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command:         []string{"sh", "-c", "sleep 3600"},
					// Use a restrictive security context just in case the test cluster has PSAs enabled.
					SecurityContext: testlib.RestrictiveSecurityContext(),
				},
			},
			ServiceAccountName: sa.Name,
		})

	tokenRequestBadAudience, err := coreV1client.ServiceAccounts(ns.Name).CreateToken(ctx, sa.Name, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences: []string{"should-fail-because-wrong-audience"}, // anything that is not an API server audience
			BoundObjectRef: &authenticationv1.BoundObjectReference{
				Kind:       "Pod",
				APIVersion: "",
				Name:       pod.Name,
				UID:        pod.UID,
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	saBadAudConfig := testlib.NewAnonymousClientRestConfig(t)
	saBadAudConfig.BearerToken = tokenRequestBadAudience.Status.Token

	_, badAudErr := testlib.NewKubeclient(t, saBadAudConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.True(t, errors.IsUnauthorized(badAudErr), testlib.Sdump(badAudErr))

	tokenRequest, err := coreV1client.ServiceAccounts(ns.Name).CreateToken(ctx, sa.Name, &authenticationv1.TokenRequest{
		Spec: authenticationv1.TokenRequestSpec{
			Audiences: []string{},
			BoundObjectRef: &authenticationv1.BoundObjectReference{
				Kind:       "Pod",
				APIVersion: "",
				Name:       pod.Name,
				UID:        pod.UID,
			},
		},
	}, metav1.CreateOptions{})
	require.NoError(t, err)

	saTokenReqConfig := testlib.NewAnonymousClientRestConfig(t)
	saTokenReqConfig.BearerToken = tokenRequest.Status.Token

	whoAmITokenReq, err := testlib.NewKubeclient(t, saTokenReqConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	whoAmIUser := whoAmITokenReq.Status.KubernetesUserInfo.User

	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "system:serviceaccount:" + ns.Name + ":" + sa.Name,
						UID:      "", // aggregation drops UID: https://github.com/kubernetes/kubernetes/issues/93699
						Groups: []string{
							"system:serviceaccounts",
							"system:serviceaccounts:" + ns.Name,
							"system:authenticated",
						},
						Extra: whoAmIUser.Extra, // This will be a dynamic assertion below based on the version of K8s
					},
				},
			},
		},
		whoAmITokenReq,
	)

	require.Equal(t, whoAmIUser.Extra["authentication.kubernetes.io/pod-name"], identityv1alpha1.ExtraValue{pod.Name})
	require.Equal(t, whoAmIUser.Extra["authentication.kubernetes.io/pod-uid"], identityv1alpha1.ExtraValue{string(pod.UID)})

	if testutil.KubeServerMinorVersionAtLeastInclusive(t, kubeClient.Discovery(), 30) {
		// Starting in K8s 1.30, three additional `Extra` fields were added with unpredictable values.
		// This is because the following three feature gates were enabled by default in 1.30.
		// https://kubernetes.io/docs/reference/command-line-tools-reference/feature-gates/
		// - ServiceAccountTokenJTI
		// - ServiceAccountTokenNodeBindingValidation
		// - ServiceAccountTokenPodNodeInfo
		// These were added in source code in 1.29 but not enabled by default until 1.30.
		// <1.29: https://pkg.go.dev/k8s.io/apiserver@v0.28.7/pkg/authentication/serviceaccount
		// 1.29+: https://pkg.go.dev/k8s.io/apiserver@v0.29.0/pkg/authentication/serviceaccount

		require.Equal(t, 5, len(whoAmIUser.Extra))
		require.NotEmpty(t, whoAmIUser.Extra["authentication.kubernetes.io/credential-id"])
		require.NotEmpty(t, whoAmIUser.Extra["authentication.kubernetes.io/node-name"])
		require.NotEmpty(t, whoAmIUser.Extra["authentication.kubernetes.io/node-uid"])
	} else {
		require.Equal(t, 2, len(whoAmIUser.Extra))
	}
}

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_CSR_Parallel(t *testing.T) {
	// use the cluster signing key being available as a proxy for this not being an EKS cluster
	// we should add more robust logic around skipping clusters based on vendor
	_ = testlib.IntegrationEnv(t).WithCapability(testlib.ClusterSigningKeyIsAvailable)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	kubeClient := testlib.NewKubernetesClientset(t)

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	der, err := x509.MarshalECPrivateKey(privateKey)
	require.NoError(t, err)

	keyPEM := pem.EncodeToMemory(&pem.Block{Type: keyutil.ECPrivateKeyBlockType, Bytes: der})

	csrPEM, err := cert.MakeCSR(privateKey, &pkix.Name{
		CommonName:   "panda-man",
		Organization: []string{"living-the-dream", "need-more-sleep"},
	}, nil, nil)
	require.NoError(t, err)

	csrName, csrUID, err := csr.RequestCertificate(
		kubeClient,
		csrPEM,
		"",
		certificatesv1.KubeAPIServerClientSignerName,
		nil,
		[]certificatesv1.KeyUsage{certificatesv1.UsageClientAuth},
		privateKey,
	)
	require.NoError(t, err)

	useCertificatesV1API := testutil.KubeServerSupportsCertificatesV1API(t, kubeClient.Discovery())

	t.Cleanup(func() {
		if useCertificatesV1API {
			require.NoError(t, kubeClient.CertificatesV1().CertificateSigningRequests().
				Delete(context.Background(), csrName, metav1.DeleteOptions{}))
		} else {
			// On old clusters use v1beta1
			require.NoError(t, kubeClient.CertificatesV1beta1().CertificateSigningRequests().
				Delete(context.Background(), csrName, metav1.DeleteOptions{}))
		}
	})

	if useCertificatesV1API {
		_, err = kubeClient.CertificatesV1().CertificateSigningRequests().UpdateApproval(ctx, csrName, &certificatesv1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: csrName,
			},
			Status: certificatesv1.CertificateSigningRequestStatus{
				Conditions: []certificatesv1.CertificateSigningRequestCondition{
					{
						Type:   certificatesv1.CertificateApproved,
						Status: corev1.ConditionTrue,
						Reason: "WhoAmICSRTest",
					},
				},
			},
		}, metav1.UpdateOptions{})
		require.NoError(t, err)
	} else {
		// On old Kubernetes clusters use CertificatesV1beta1
		_, err = kubeClient.CertificatesV1beta1().CertificateSigningRequests().UpdateApproval(ctx, &certificatesv1beta1.CertificateSigningRequest{
			ObjectMeta: metav1.ObjectMeta{
				Name: csrName,
			},
			Status: certificatesv1beta1.CertificateSigningRequestStatus{
				Conditions: []certificatesv1beta1.CertificateSigningRequestCondition{
					{
						Type:   certificatesv1beta1.CertificateApproved,
						Status: corev1.ConditionTrue,
						Reason: "WhoAmICSRTest",
					},
				},
			},
		}, metav1.UpdateOptions{})
		require.NoError(t, err)
	}

	crtPEM, err := csr.WaitForCertificate(ctx, kubeClient, csrName, csrUID)
	require.NoError(t, err)

	csrConfig := testlib.NewAnonymousClientRestConfig(t)
	csrConfig.CertData = crtPEM
	csrConfig.KeyData = keyPEM

	whoAmI, err := testlib.NewKubeclient(t, csrConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "panda-man",
						Groups: []string{
							"need-more-sleep",
							"living-the-dream",
							"system:authenticated",
						},
					},
				},
			},
		},
		whoAmI,
	)
}

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_Anonymous_Parallel(t *testing.T) {
	_ = testlib.IntegrationEnv(t).WithCapability(testlib.AnonymousAuthenticationSupported)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	anonymousConfig := testlib.NewAnonymousClientRestConfig(t)

	whoAmI, err := testlib.NewKubeclient(t, anonymousConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	// this also asserts that all users, even unauthenticated ones, can call this API when anonymous is enabled
	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "system:anonymous",
						Groups: []string{
							"system:unauthenticated",
						},
					},
				},
			},
		},
		whoAmI,
	)
}

// whoami requests are non-mutating and safe to run in parallel with serial tests, see main_test.go.
func TestWhoAmI_ImpersonateDirectly_Parallel(t *testing.T) {
	_ = testlib.IntegrationEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	impersonationConfig := testlib.NewClientConfig(t)
	impersonationConfig.Impersonate = rest.ImpersonationConfig{
		UserName: "solaire",
		// need to impersonate system:authenticated directly to support older clusters otherwise we will get RBAC errors below
		Groups: []string{"astora", "lordran", "system:authenticated"},
		Extra: map[string][]string{
			"covenant": {"warrior-of-sunlight"},
			"loves":    {"sun", "co-op"},
		},
	}

	whoAmI, err := testlib.NewKubeclient(t, impersonationConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "solaire",
						UID:      "", // no way to impersonate UID: https://github.com/kubernetes/kubernetes/issues/93699
						Groups: []string{
							"astora",
							"lordran",
							"system:authenticated", // impersonation will add this implicitly but only in newer clusters
						},
						Extra: map[string]identityv1alpha1.ExtraValue{
							"covenant": {"warrior-of-sunlight"},
							"loves":    {"sun", "co-op"},
						},
					},
				},
			},
		},
		whoAmI,
	)

	impersonationAnonymousConfig := testlib.NewClientConfig(t)
	impersonationAnonymousConfig.Impersonate.UserName = "system:anonymous"
	// need to impersonate system:unauthenticated directly to support older clusters otherwise we will get RBAC errors below
	impersonationAnonymousConfig.Impersonate.Groups = []string{"system:unauthenticated"}

	whoAmIAnonymous, err := testlib.NewKubeclient(t, impersonationAnonymousConfig).PinnipedConcierge.IdentityV1alpha1().WhoAmIRequests().
		Create(ctx, &identityv1alpha1.WhoAmIRequest{}, metav1.CreateOptions{})
	require.NoError(t, err, testlib.Sdump(err))

	require.Equal(t,
		&identityv1alpha1.WhoAmIRequest{
			Status: identityv1alpha1.WhoAmIRequestStatus{
				KubernetesUserInfo: identityv1alpha1.KubernetesUserInfo{
					User: identityv1alpha1.UserInfo{
						Username: "system:anonymous",
						Groups: []string{
							"system:unauthenticated", // impersonation will add this implicitly but only in newer clusters
						},
					},
				},
			},
		},
		whoAmIAnonymous,
	)
}
