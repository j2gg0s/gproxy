package pkg

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ImageCopy(ctx context.Context, r *Request) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("new docker client with error: %v", err)
	}

	log.Info().Msgf("docker pull %s", r.Source)
	err = readMessage(cli.ImagePull(ctx, r.Source, types.ImagePullOptions{}))
	if err != nil {
		return err
	}

	log.Info().Msgf("docker image tag %s %s", r.Source, r.Dest)
	if err := cli.ImageTag(ctx, r.Source, r.Dest); err != nil {
		return fmt.Errorf("tag image with error: %v", err)
	}

	log.Info().Msgf("docker push %s", r.Dest)
	return readMessage(cli.ImagePush(ctx, r.Dest, types.ImagePushOptions{RegistryAuth: r.Token}))
}

func readMessage(reader io.ReadCloser, err error) error {
	if err != nil {
		return err
	}
	defer reader.Close()

	buf := bytes.Buffer{}
	if zerolog.GlobalLevel() == zerolog.DebugLevel {
		_, err = io.Copy(os.Stderr, io.TeeReader(reader, &buf))
		if err != nil {
			return err
		}
	} else {
		b, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		buf = *bytes.NewBuffer(b)
	}

	scanner := bufio.NewScanner(&buf)
	for scanner.Scan() {
		var msg interface{}
		err = json.Unmarshal(scanner.Bytes(), &msg)
		if err != nil {
			log.Warn().Err(err).Msgf("unmarshal docker response with error: %s", string(scanner.Bytes()))
			continue
		}
		if v, ok := msg.(map[string]interface{}); ok {
			if message, ok := v["message"]; ok {
				return fmt.Errorf("request docker api with error message: %s", message)
			}
		}
	}

	return nil
}
