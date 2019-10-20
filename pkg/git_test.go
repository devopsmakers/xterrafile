package xterrafile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsGitSourceAddr(t *testing.T) {
	assert.True(t, IsGitSourceAddr("git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"))
	assert.True(t, IsGitSourceAddr("git::https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"))
	assert.True(t, IsGitSourceAddr("https://github.com/terraform-digitalocean-modules/terraform-digitalocean-droplet.git"))
	assert.False(t, IsGitSourceAddr("/some/absolute/path"), "absolute path should be false")
	assert.False(t, IsGitSourceAddr("https://something"), "http source should be false")
}

func TestGetGitTags(t *testing.T) {
	assert.Contains(t, getGitTags("droplet", "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"), "v0.1.7")
	assert.Contains(t, getGitTags("droplet", "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git"), "v0.0.2")
}

func TestGetGitSource(t *testing.T) {
	module1Src, module1Version := GetGitSource("droplet", "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git", "> 0.1.2 <= 0.1.7")
	assert.Equal(t, "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git", module1Src)
	assert.Equal(t, "v0.1.7", module1Version)

	module2Src, module2Version := GetGitSource("droplet", "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git", "39bda6c7aabac9226ec6628339463aa1708bef85")
	assert.Equal(t, "git@github.com:terraform-digitalocean-modules/terraform-digitalocean-droplet.git", module2Src)
	assert.Equal(t, "39bda6c7aabac9226ec6628339463aa1708bef85", module2Version)
}
