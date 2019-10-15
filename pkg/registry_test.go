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

	module1Src, module1Version := GetRegistrySource("droplet", "example.com/test-versions/name/provider", "1.2.x", test.Disco(server))
	assert.IsType(t, "string", module1Src, "download URL should be a string")
	assert.Equal(t, "1.2.2", module1Version)

	module2Src, module2Version := GetRegistrySource("droplet", "example.com/test-versions/name/provider", "2.1.0", test.Disco(server))
	assert.IsType(t, "string", module2Src, "download URL should be a string")
	assert.Equal(t, "2.1.0", module2Version)

	module3Src, module3Version := GetRegistrySource("droplet", "example.com/test-versions/name/provider", "", test.Disco(server))
	assert.IsType(t, "string", module3Src, "download URL should be a string")
	assert.Equal(t, "2.2.0", module3Version)

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

func TestGetRegistryModuleVersions(t *testing.T) {
	server := test.Registry()
	defer server.Close()

	testClient := registry.NewClient(test.Disco(server), nil)

	modSrc, _ := getModSrc("example.com/test-versions/name/provider")
	versions := getRegistryModuleVersions(testClient, modSrc)
	assert.Equal(t, []string{"2.2.0", "2.1.1", "1.2.2", "1.2.1"}, versions)
}
