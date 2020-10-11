package registry

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

func SearchContainer(containerName string, limit uint8) (error, SearchResults) {
	// GET /v1/search?q=search_term HTTP/1.1
	// Host: example.com
	// Accept: application/json
	var result SearchResults
	utils.PrintN(utils.Info, "Query Docker registry for: "+containerName)
	fmt.Println()
	setOfDigits := spinner.GenerateNumberSequence(20)
	s := spinner.New(setOfDigits, 200*time.Millisecond)
	s.Prefix = "running ..... "
	s.HideCursor = true
	s.Start()
	request := gorequest.New()
	url := HostURL + "/v1" + "/search" + "?q=" + containerName + "&n=" + strconv.Itoa(int(limit))
	fmt.Println(url)
	resp, body, err := request.Get(url).End()
	if err != nil {
		return errors.WithMessage(err[0], "search container error"), result
	}
	defer resp.Body.Close()

	e := json.Unmarshal([]byte(body), &result)
	s.Stop()
	if e != nil {
		return e, result
	}
	return nil, result

}
