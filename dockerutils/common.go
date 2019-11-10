package dockerutils

import (
	"encoding/base64"
	"fmt"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/gookit/color"
)

type Environment struct {
	StorePath string
}

var logo = `%s

%s V%s
%s

`

func Displaylogo() {

	banner, _ := base64.StdEncoding.DecodeString(utils.BannerBase64)
	fmt.Printf(color.FgLightGreen.Render(logo), banner, utils.Appname, utils.Version, color.FgMagenta.Render(utils.GitHub))
}
