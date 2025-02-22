/*
Copyright 2017 The Kubernetes Authors.

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

package constants

import (
	"path/filepath"
	"testing"

	"k8s.io/apimachinery/pkg/util/version"
	apimachineryversion "k8s.io/apimachinery/pkg/version"
)

func TestGetStaticPodDirectory(t *testing.T) {
	expected := "/etc/kubernetes/manifests"
	actual := GetStaticPodDirectory()

	if actual != expected {
		t.Errorf(
			"failed GetStaticPodDirectory:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetAdminKubeConfigPath(t *testing.T) {
	expected := filepath.Join(KubernetesDir, AdminKubeConfigFileName)
	actual := GetAdminKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetAdminKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetBootstrapKubeletKubeConfigPath(t *testing.T) {
	expected := "/etc/kubernetes/bootstrap-kubelet.conf"
	actual := GetBootstrapKubeletKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetBootstrapKubeletKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetKubeletKubeConfigPath(t *testing.T) {
	expected := "/etc/kubernetes/kubelet.conf"
	actual := GetKubeletKubeConfigPath()

	if actual != expected {
		t.Errorf(
			"failed GetKubeletKubeConfigPath:\n\texpected: %s\n\t  actual: %s",
			expected,
			actual,
		)
	}
}

func TestGetStaticPodFilepath(t *testing.T) {
	var tests = []struct {
		componentName, manifestsDir, expected string
	}{
		{
			componentName: "kube-apiserver",
			manifestsDir:  "/etc/kubernetes/manifests",
			expected:      "/etc/kubernetes/manifests/kube-apiserver.yaml",
		},
		{
			componentName: "kube-controller-manager",
			manifestsDir:  "/etc/kubernetes/manifests/",
			expected:      "/etc/kubernetes/manifests/kube-controller-manager.yaml",
		},
		{
			componentName: "foo",
			manifestsDir:  "/etc/bar/",
			expected:      "/etc/bar/foo.yaml",
		},
	}
	for _, rt := range tests {
		t.Run(rt.componentName, func(t *testing.T) {
			actual := GetStaticPodFilepath(rt.componentName, rt.manifestsDir)
			if actual != rt.expected {
				t.Errorf(
					"failed GetStaticPodFilepath:\n\texpected: %s\n\t  actual: %s",
					rt.expected,
					actual,
				)
			}
		})
	}
}

func TestEtcdSupportedVersion(t *testing.T) {
	var supportedEtcdVersion = map[uint8]string{
		13: "3.2.24",
		14: "3.3.10",
		15: "3.3.10",
		16: "3.3.17-0",
		17: "3.4.3-0",
		18: "3.4.3-0",
	}
	var tests = []struct {
		kubernetesVersion string
		expectedVersion   *version.Version
		expectedWarning   bool
		expectedError     bool
	}{
		{
			kubernetesVersion: "1.x.1",
			expectedVersion:   nil,
			expectedWarning:   false,
			expectedError:     true,
		},
		{
			kubernetesVersion: "1.10.1",
			expectedVersion:   version.MustParseSemantic("3.2.24"),
			expectedWarning:   true,
			expectedError:     false,
		},
		{
			kubernetesVersion: "1.99.0",
			expectedVersion:   version.MustParseSemantic("3.4.3-0"),
			expectedWarning:   true,
			expectedError:     false,
		},
		{
			kubernetesVersion: "v1.16.0",
			expectedVersion:   version.MustParseSemantic("3.3.17-0"),
			expectedWarning:   false,
			expectedError:     false,
		},
		{
			kubernetesVersion: "1.17.2",
			expectedVersion:   version.MustParseSemantic("3.4.3-0"),
			expectedWarning:   false,
			expectedError:     false,
		},
	}
	for _, rt := range tests {
		t.Run(rt.kubernetesVersion, func(t *testing.T) {
			actualVersion, actualWarning, actualError := EtcdSupportedVersion(supportedEtcdVersion, rt.kubernetesVersion)
			if (actualError != nil) != rt.expectedError {
				t.Fatalf("expected error %v, got %v", rt.expectedError, actualError != nil)
			}
			if (actualWarning != nil) != rt.expectedWarning {
				t.Fatalf("expected warning %v, got %v", rt.expectedWarning, actualWarning != nil)
			}
			if actualError == nil && actualVersion.String() != rt.expectedVersion.String() {
				t.Errorf("expected version %s, got %s", rt.expectedVersion.String(), actualVersion.String())
			}
		})
	}
}

func TestGetKubernetesServiceCIDR(t *testing.T) {
	var tests = []struct {
		svcSubnetList string
		isDualStack   bool
		expected      string
		expectedError bool
		name          string
	}{
		{
			svcSubnetList: "192.168.10.0/24",
			isDualStack:   false,
			expected:      "192.168.10.0/24",
			expectedError: false,
			name:          "valid: valid IPv4 range from single-stack",
		},
		{
			svcSubnetList: "fd03::/112",
			isDualStack:   false,
			expected:      "fd03::/112",
			expectedError: false,
			name:          "valid: valid IPv6 range from single-stack",
		},
		{
			svcSubnetList: "192.168.10.0/24,fd03::/112",
			isDualStack:   true,
			expected:      "192.168.10.0/24",
			expectedError: false,
			name:          "valid: valid <IPv4,IPv6> ranges from dual-stack",
		},
		{
			svcSubnetList: "fd03::/112,192.168.10.0/24",
			isDualStack:   true,
			expected:      "fd03::/112",
			expectedError: false,
			name:          "valid: valid <IPv6,IPv4> ranges from dual-stack",
		},
		{
			svcSubnetList: "192.168.10.0/24,fd03:x::/112",
			isDualStack:   true,
			expected:      "",
			expectedError: true,
			name:          "invalid: failed to parse subnet range for dual-stack",
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			actual, actualError := GetKubernetesServiceCIDR(rt.svcSubnetList, rt.isDualStack)
			if rt.expectedError {
				if actualError == nil {
					t.Errorf("failed GetKubernetesServiceCIDR:\n\texpected error, but got no error")
				}
			} else if !rt.expectedError && actualError != nil {
				t.Errorf("failed GetKubernetesServiceCIDR:\n\texpected no error, but got: %v", actualError)
			} else {
				if actual.String() != rt.expected {
					t.Errorf(
						"failed GetKubernetesServiceCIDR:\n\texpected: %s\n\t  actual: %s",
						rt.expected,
						actual.String(),
					)
				}
			}
		})
	}
}

func TestGetSkewedKubernetesVersionImpl(t *testing.T) {
	tests := []struct {
		name                string
		versionInfo         *apimachineryversion.Info
		n                   int
		isRunningInTestFunc func() bool
		expectedResult      *version.Version
	}{
		{
			name:           "invalid versionInfo; running in test",
			versionInfo:    &apimachineryversion.Info{},
			expectedResult: defaultKubernetesVersionForTests,
		},
		{
			name:                "invalid versionInfo; not running in test",
			versionInfo:         &apimachineryversion.Info{},
			isRunningInTestFunc: func() bool { return false },
			expectedResult:      nil,
		},
		{
			name:           "valid skew of -1",
			versionInfo:    &apimachineryversion.Info{Major: "1", GitVersion: "v1.23.0"},
			n:              -1,
			expectedResult: version.MustParseSemantic("v1.22.0"),
		},
		{
			name:           "valid skew of 0",
			versionInfo:    &apimachineryversion.Info{Major: "1", GitVersion: "v1.23.0"},
			n:              0,
			expectedResult: version.MustParseSemantic("v1.23.0"),
		},
		{
			name:           "valid skew of +1",
			versionInfo:    &apimachineryversion.Info{Major: "1", GitVersion: "v1.23.0"},
			n:              1,
			expectedResult: version.MustParseSemantic("v1.24.0"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isRunningInTestFunc == nil {
				tc.isRunningInTestFunc = func() bool { return true }
			}
			result := getSkewedKubernetesVersionImpl(tc.versionInfo, tc.n, tc.isRunningInTestFunc)
			if (tc.expectedResult == nil) != (result == nil) {
				t.Errorf("expected result: %v, got: %v", tc.expectedResult, result)
			}
			if result == nil {
				return
			}
			if cmp, _ := result.Compare(tc.expectedResult.String()); cmp != 0 {
				t.Errorf("expected result: %v, got %v", tc.expectedResult, result)
			}
		})
	}
}
