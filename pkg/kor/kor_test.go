package kor

import (
	"os"
	"sort"
	"testing"
)

func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Sort the slices before comparing
	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestRemoveDuplicatesAndSort(t *testing.T) {
	// Test case 1: Test removing duplicates and sorting the slice
	slice := []string{"b", "a", "c", "b", "a"}
	expected := []string{"a", "b", "c"}
	result := RemoveDuplicatesAndSort(slice)

	if !stringSlicesEqual(result, expected) {
		t.Errorf("RemoveDuplicatesAndSort failed, expected: %v, got: %v", expected, result)
	}

	// Test case 2: Test removing duplicates and sorting an empty slice
	emptySlice := []string{}
	emptyExpected := []string{}
	emptyResult := RemoveDuplicatesAndSort(emptySlice)

	if !stringSlicesEqual(emptyResult, emptyExpected) {
		t.Errorf("RemoveDuplicatesAndSort failed for empty slice, expected: %v, got: %v", emptyExpected, emptyResult)
	}
}

func TestCalculateResourceDifference(t *testing.T) {
	usedResourceNames := []string{"resource1", "resource2", "resource3"}
	allResourceNames := []string{"resource1", "resource2", "resource3", "resource4", "resource5"}

	expectedDifference := []string{"resource4", "resource5"}
	difference := CalculateResourceDifference(usedResourceNames, allResourceNames)

	if len(difference) != len(expectedDifference) {
		t.Errorf("Expected %d difference items, but got %d", len(expectedDifference), len(difference))
	}

	for i, item := range difference {
		if item != expectedDifference[i] {
			t.Errorf("Difference item at index %d should be %s, but got %s", i, expectedDifference[i], item)
		}
	}
}

func getFakeConfigContent() string {
	fakeContent := `
apiVersion: v1
clusters:
- cluster:
    server: https://localhost:8080
    extensions:
    - name: client.authentication.k8s.io/exec
      extension:
        audience: foo
        other: bar
  name: foo-cluster
contexts:
- context:
    cluster: foo-cluster
    namespace: bar
  name: foo-context
current-context: foo-context
kind: Config
`
	return fakeContent
}

func TestGetKubeClientFromEnvVar(t *testing.T) {
	configFile, err := os.CreateTemp("", "kubeconfig-")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(configFile.Name())
	if err := os.WriteFile(configFile.Name(), []byte(getFakeConfigContent()), 0666); err != nil {
		t.Error(err)
	}

	originalKCEnv := os.Getenv("KUBECONFIG")
	defer os.Setenv("KUBECONFIG", originalKCEnv)
	os.Setenv("KUBECONFIG", configFile.Name())

	kcs := GetKubeClient("")
	if kcs == nil {
		t.Errorf("Expected valid clientSet")
	}
}

func TestGetKubeClientFromInput(t *testing.T) {
	configFile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(configFile.Name())
	if err := os.WriteFile(configFile.Name(), []byte(getFakeConfigContent()), 0666); err != nil {
		t.Error(err)
	}

	kcs := GetKubeClient(configFile.Name())
	if kcs == nil {
		t.Errorf("Expected valid clientSet")
	}
}
