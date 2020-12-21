package gproxy

import (
	"strconv"
	"testing"

	"github.com/docker/distribution/reference"
	"github.com/stretchr/testify/require"
)

func TestImageCopy(t *testing.T) {
	fixtures := []struct {
		source string
		domain string
		path   string
	}{
		{
			"ubuntu",
			"docker.io",
			"library/ubuntu",
		},
		{
			"gcr.io/google-containers/busybox:1.27",
			"gcr.io",
			"google-containers/busybox",
		},
		{
			"kubernetesui/dashboard:latest",
			"docker.io",
			"kubernetesui/dashboard",
		},
	}

	for i, f := range fixtures {
		fixture := f
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			named, err := reference.ParseNormalizedNamed(fixture.source)
			require.NoError(t, err)
			require.Equal(t, fixture.domain, reference.Domain(named))
			require.Equal(t, fixture.path, reference.Path(named))
		})
	}
}
