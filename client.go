package gproxy

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ImageCopy(ctx context.Context, source, dest string, token string) error {
	namedSource, err := reference.ParseNormalizedNamed(source)
	if err != nil {
		return fmt.Errorf("invalid source image: %s, %v", source, err)
	}
	if strings.Count(dest, "/") < 2 {
		// default acr repo under registry.cn-huhehaote.aliyuncs.com/gproxy
		dest += "/default"
	}
	if !strings.Contains(dest, ":") {
		path := reference.Path(namedSource)
		tag := source[strings.LastIndex(source, ":")+1:]
		// translate path to tag
		dest += ":" + strings.ReplaceAll(path, "/", ".") + "-" + tag
	}
	_, err = reference.ParseNormalizedNamed(dest)
	if err != nil {
		return fmt.Errorf("invalid dest image: %s, %v", dest, err)
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("new docker client with error: %v", err)
	}

	log.Info().Msgf("docker pull %s", source)
	sourceImage, err := cli.ImagePull(ctx, source, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pull image with error: %v", err)
	}
	defer sourceImage.Close()
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		_, _ = io.Copy(os.Stdout, sourceImage)
	} else {
		_, _ = io.Copy(ioutil.Discard, sourceImage)
	}

	if err := cli.ImageTag(ctx, source, dest); err != nil {
		return fmt.Errorf("tag image with error: %v", err)
	}
	log.Info().Msgf("docker image tag %s %s", source, dest)

	log.Info().Msgf("docker push %s", dest)
	destImage, err := cli.ImagePush(ctx, dest, types.ImagePushOptions{RegistryAuth: token})
	if err != nil {
		return fmt.Errorf("push image with error: %s, %v", dest, err)
	}
	defer destImage.Close()
	if zerolog.GlobalLevel() <= zerolog.DebugLevel {
		_, _ = io.Copy(os.Stdout, destImage)
	} else {
		_, _ = io.Copy(ioutil.Discard, destImage)
	}

	return nil
}
