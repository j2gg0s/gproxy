package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/j2gg0s/gproxy/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func main() {
	root := cobra.Command{
		Use: "gproxy",
	}
	addCommonFlags(root.PersistentFlags())
	root.PersistentPreRunE = func(*cobra.Command, []string) error {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		return nil
	}
	root.RunE = func(*cobra.Command, []string) error {
		r, err := pkg.NewRequest(source, dest, username, password)
		if err != nil {
			log.Warn().Err(err).Send()
			return err
		}

		return pkg.ImageCopy(context.Background(), r)
	}

	httpServer := &cobra.Command{
		Use: "http",
	}
	httpServer.PersistentFlags().IntVar(&port, "port", 8080, "http port")
	httpServer.RunE = func(cmd *cobra.Command, args []string) error {
		return pkg.RunGin(port)
	}

	remote := &cobra.Command{
		Use: "remote",
	}
	addCommonFlags(remote.PersistentFlags())
	remote.PersistentFlags().StringVar(&addr, "addr", "https://gproxy.j2gg0s.com", "server url")
	remote.RunE = func(*cobra.Command, []string) error {
		values := url.Values{}
		values.Add("source", source)
		values.Add("dest", dest)
		values.Add("username", username)
		values.Add("password", password)

		resp, err := http.Get(fmt.Sprintf("%s?%s", addr, values.Encode()))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Info().Int("status", resp.StatusCode).Str("body", string(body)).Msg("response")

		return nil
	}

	root.AddCommand(httpServer)
	root.AddCommand(remote)
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

func addCommonFlags(set *pflag.FlagSet) {
	set.StringVar(&source, "source", "", "source image url")
	set.StringVar(&dest, "dest", pkg.DefaultDest, "dest acr's addr, {domain}/{namespace}/{name}:{tag}, name and tag is optional")
	set.StringVar(&username, "username", pkg.DefaultUsername, "username for auth")
	set.StringVar(&password, "password", pkg.DefaultPassword, "passowrd for auth")
	set.BoolVar(&debug, "debug", false, "log level")
}

var (
	source   string
	dest     string
	username string
	password string
	debug    bool

	addr string

	port int
)
