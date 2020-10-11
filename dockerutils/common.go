package dockerutils

import (
	"encoding/base64"
	"fmt"
	"github.com/gookit/color"
	"github.com/cnlubo/go-docker-search/version"
)

type Environment struct {
	StorePath string
}

var logo = `%s

%s V%s
%s

`

func Displaylogo() {

	banner, _ := base64.StdEncoding.DecodeString(version.BannerBase64)
	fmt.Printf(color.FgLightGreen.Render(logo), banner, version.AppName, version.AppVersion, color.FgMagenta.Render(version.GitHub))
}
