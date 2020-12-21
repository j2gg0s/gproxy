package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/j2gg0s/gproxy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{
		Use: "gproxy",
	}
	root.PersistentFlags().StringVar(&source, "source", "", "source image url")
	root.PersistentFlags().StringVar(&dest, "dest", defaultDest, "dest acr's addr, {domain}/{namespace}/{name}:{tag}, name and tag is optional")
	root.PersistentFlags().StringVar(&username, "username", defaultUsername, "username for auth")
	root.PersistentFlags().StringVar(&password, "password", defaultPassword, "passowrd for auth")
	root.PersistentFlags().BoolVar(&debug, "debug", false, "log level")

	root.RunE = func(cmd *cobra.Command, args []string) error {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}

		token := ""
		if dest == defaultDest || username != defaultUsername || password != defaultPassword {
			authConfig, err := json.Marshal(types.AuthConfig{Username: username, Password: password})
			if err != nil {
				return fmt.Errorf("invalid username/password: %v", err)
			}
			token = base64.URLEncoding.EncodeToString(authConfig)
		}

		return gproxy.ImageCopy(context.Background(), source, dest, token)
	}

	if err := root.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

var (
	source   string
	dest     string
	username string
	password string
	debug    bool

	defaultDest     string = "registry.cn-huhehaote.aliyuncs.com/gproxy"
	defaultUsername string = "gproxy@j2gg0s"
	defaultPassword string = "Aliyun123456"
)
