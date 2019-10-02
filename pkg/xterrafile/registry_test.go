package xterrafile

import (
	"testing"

	"github.com/hashicorp/terraform/registry"
	"github.com/hashicorp/terraform/registry/test"
	"github.com/stretchr/testify/assert"
)

func TestIsRegistrySourceAddr(t *testing.T) {
	assert.False(t, IsRegistrySourceAddr("./some/path"), "non-registry path should be false")
	assert.True(t, IsRegistrySourceAddr("terraform-digitalocean-modules/droplet/digitalocean"),
		"registry path should be true")
	assert.True(t, IsRegistrySourceAddr(
		"terraform-digitalocean-modules/droplet/digitalocean//examples/test"),
		"registry path with sub-path should be true")
	assert.True(t, IsRegistrySourceAddr("app.terraform.io/terraform-digitalocean-modules/droplet/digitalocean"),
		"private registry path should be true")
	assert.True(t, IsRegistrySourceAddr(
		"app.terraform.io/terraform-digitalocean-modules/droplet/digitalocean//examples/test"),
		"private registry path with sub-path should be true")
}

func TestGetRegistrySource(t *testing.T) {
	server := test.Registry()
	defer server.Close()

	module1Src := GetRegistrySource("droplet", "example.com/test-versions/name/provider", "2.1.x", test.Disco(server))
	assert.IsType(t, "string", module1Src, "download URL should be a string")
}

func TestGetModSrc(t *testing.T) {
	module1String := "terraform-digitalocean-modules/droplet/digitalocean"
	module1Src, err := getModSrc(module1String)
	assert.Equal(t, nil, err, "error should be nil")
	assert.Equal(t, "registry.terraform.io", module1Src.Host().Normalized(), "host should be public repo")

	module2String := "app.terraform.io/terraform-digitalocean-modules/droplet/digitalocean"
	module2Src, err := getModSrc(module2String)
	assert.Equal(t, nil, err, "error should be nil")
	assert.Equal(t, "app.terraform.io", module2Src.Host().Normalized(), "host should be private repo")

	module3String := "---.io/terraform-digitalocean-modules/droplet/digitalocean"
	module3Src, err := getModSrc(module3String)
	assert.NotEqual(t, nil, err, "error should be present")
	assert.Panics(t, func() { module3Src.Host().Normalized() }, "accessing host should panic")
}

func TestGetRegistryVersion(t *testing.T) {
	server := test.Registry()
	defer server.Close()

	testClient := registry.NewClient(test.Disco(server), nil)

	modSrc, _ := getModSrc("example.com/test-versions/name/provider")
	version, _ := getRegistryVersion(testClient, modSrc, ">= 2.0.0 < 2.2.0")
	assert.Equal(t, "2.1.1", version, "version should be >= 2.0.0 < 2.2.0")

	_, err := getRegistryVersion(testClient, modSrc, ">= 3.0.0")
	assert.Error(t, err, "should have returned an error")

	_, err = getRegistryVersion(testClient, modSrc, "not.a.version")
	assert.Error(t, err, "should return an error")

	modSrc, _ = getModSrc("invalid.com/test-versions/name/provider")
	_, err = getRegistryVersion(testClient, modSrc, ">= 3.0.0")
	assert.Error(t, err, "should return an error")
}
