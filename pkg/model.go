package pkg

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
)

type Request struct {
	Source string
	Dest   string
	Token  string
}

var (
	DefaultDest     string = "registry.cn-huhehaote.aliyuncs.com/gproxy"
	DefaultUsername string = "gproxy@j2gg0s"
	DefaultPassword string = "Aliyun123456"

	DefaultPort int = 80
)

func NewRequest(source, dest, username, password string) (*Request, error) {
	token := ""
	if dest == DefaultDest || username != DefaultUsername || password != DefaultPassword {
		authConfig, err := json.Marshal(types.AuthConfig{Username: username, Password: password})
		if err != nil {
			return nil, fmt.Errorf("invalid username/password: %v", err)
		}
		token = base64.URLEncoding.EncodeToString(authConfig)
	}

	namedSource, err := reference.ParseNormalizedNamed(source)
	if err != nil {
		return nil, fmt.Errorf("invalid source image: %s, %v", source, err)
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
		return nil, fmt.Errorf("invalid dest image: %s, %v", dest, err)
	}

	return &Request{
		Source: source,
		Dest:   dest,
		Token:  token,
	}, nil
}
