package registry

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/cnlubo/go-docker-search/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"strings"
)

// var (
// 	ErrNoMorePages = errors.New("no more pages")
// )

// This function parses the Www-Authenticate header provided in the challenge
// It has the following format
// Bearer realm="https://gitlab.com/jwt/auth",service="container_registry",scope="repository:andrew18/container-test:pull"
func parseBearer(bearer []string) map[string]string {
	out := make(map[string]string)
	for _, b := range bearer {
		for _, s := range strings.Split(b, " ") {
			if s == "Bearer" {
				continue
			}
			for _, params := range strings.Split(s, ",") {
				fields := strings.Split(params, "=")
				key := fields[0]
				val := strings.Replace(fields[1], "\"", "", -1)
				out[key] = val
			}
		}
	}
	return out
}

func getLastN(s []*semver.Version, number uint8) (semVerTags []*semver.Version) {

	var semVersions []*semver.Version
	if len(s) <= int(number) {
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] != nil && s[i].Original() != "" {
				semVersions = append(semVersions, s[i])
			}
		}
	} else {
		for i := len(s) - 1; i >= len(s)-int(number); i-- {
			if s[i] != nil && s[i].Original() != "" {
				semVersions = append(semVersions, s[i])
			}
		}
	}
	if len(semVersions) > 0 {
		return semVersions
	} else {
		return nil
	}
}

func getMajorVersions(s []*semver.Version, major int64, minor int64) (semVerTags []*semver.Version) {

	var semVersions []*semver.Version
	for i := len(s) - 1; i >= 0; i-- {
		fmt.Println(s[i].Patch())
		if s[i] != nil && s[i].Original() != "" {
			if (s[i].Major() == major && minor == 0) || (s[i].Major() == major && s[i].Minor() == minor) {
				semVersions = append(semVersions, s[i])
			}
		}
	}
	if len(semVersions) > 0 {
		return semVersions
	} else {
		return nil
	}
}

// Parse a list of tags and return a slice of SemVer versions.
func tags2SemVer(tags []string) (semVerTags []*semver.Version, errTags []string) {

	if len(tags) == 0 {
		return nil, nil
	}
	for _, t := range tags {
		v, err := semver.NewVersion(t)
		if err == nil {
			semVerTags = append(semVerTags, v)
		} else {
			errTags = append(errTags, t)
		}
	}
	return semVerTags, errTags
}

func listAllTags(repository *RepositoryRequest) (tags []string, err error) {
	var url string
	// var buf bytes.Buffer
	// logger := log.New(&buf, "[gorequest]", log.LstdFlags)
	// logger.SetOutput(os.Stdout)
	// request := gorequest.New().SetDebug(true).SetLogger(logger)

	request := gorequest.New()
	url += fixSuffixPrefix(repository.RepositoryUrl)          // https://index.docker.io/v2/
	url += fixSuffixPrefix(fixOfficialRepos(repository.Repo)) // "_/mongo/" or "foo/bar/"
	url += fixSuffixPrefix(repository.Endpoint)               // "tags/list/
	utils.PrintN(utils.Info, fmt.Sprintf("registry.tags url=%s repository=%s", url, repository.Repo))
	// fmt.Println()
	// First step is to tagFetch the endpoint where we'll be authenticating
	resp, _, _ := request.Get(url).End()
	// fmt.Println(resp.Status)
	// fmt.Println(resp.StatusCode)
	// fmt.Println("ddd")
	// Get the token
	// This has the various things we'll need to parse and use in the request
	// fmt.Println(resp.Header["Www-Authenticate"])
	utils.PrintN(utils.Info,"ok2")
	params := parseBearer(resp.Header["Www-Authenticate"])
	// fmt.Println(params["scope"])
	paramsJSON, _ := json.Marshal(&params)
	// fmt.Println(string(paramsJSON))

	challenge := gorequest.New()
	resp, body, _ := challenge.Get(params["realm"]).
		SetBasicAuth(repository.User, repository.Password).
		Query(string(paramsJSON)).
		End()
	token := make(map[string]string)
	_ = json.Unmarshal([]byte(body), &token)
	utils.PrintN(utils.Info,"ok1")
	// Now reissue the challenge with the token in the Header
	// curl -IL https://index.docker.io/v2/odewahn/image/tags/list
	authenticatedRequest := gorequest.New()

	resp, body, _ = authenticatedRequest.Get(url).
		Set("Authorization", "Bearer "+token["token"]).
		End()
	var tagResponse TagResponse
	err = json.Unmarshal([]byte(body), &tagResponse)
	utils.PrintN(utils.Info,"ok")
	if err != nil {
		return nil, err
	}
	return tagResponse.Tags, err

	// var response TagResponse
	// for {
	// 	url, err = getPaginatedJSON(url, token, &response)
	// 	switch err {
	// 	case ErrNoMorePages:
	// 		tags = append(tags, response.Tags...)
	// 		return tags, nil
	// 	case nil:
	// 		tags = append(tags, response.Tags...)
	// 		continue
	// 	default:
	// 		return nil, err
	// 	}
	// }
	return nil, err
}

// getPaginatedJSON accepts a string and a pointer, and returns the
// next page URL while updating pointed-to variable with a parsed JSON
// value. When there are no more pages it returns `ErrNoMorePages`.
// func getPaginatedJSON(url string, token map[string]string, response interface{}) (string, error) {
//
// 	// Now reissue the challenge with the token in the Header
// 	// curl -IL https://index.docker.io/v2/odewahn/image/tags/list
// 	authenticatedRequest := gorequest.New()
// 	resp, body, err := authenticatedRequest.Get(url).
// 		Set("Authorization", "Bearer "+token["token"]).
// 		End()
// 	if err != nil {
// 		return "", err[0]
// 	}
// 	defer resp.Body.Close()
// 	// var tagResponse TagResponse
// 	_ = json.Unmarshal([]byte(body), &response)
// 	// fmt.Println(resp)
// 	return getNextLink(resp)
// }

// Matches an RFC 5988 (https://tools.ietf.org/html/rfc5988#section-5)
// Link header. For example,
//
//    <http://registry.example.com/v2/_catalog?n=5&last=tag5>; type="application/json"; rel="next"
//
// The URL is _supposed_ to be wrapped by angle brackets `< ... >`,
// but e.g., quay.io does not include them. Similarly, params like
// `rel="next"` may not have quoted values in the wild.
// var nextLinkRE = regexp.MustCompile(`^ *<?([^;>]+)>? *(?:;[^;]*)*; *rel="?next"?(?:;.*)?`)
//
// func getNextLink(resp *http.Response) (string, error) {
// 	for _, link := range resp.Header[http.CanonicalHeaderKey("Link")] {
// 		parts := nextLinkRE.FindStringSubmatch(link)
// 		if parts != nil {
// 			fmt.Println("ok")
// 			return parts[1], nil
// 		}
// 	}
// 	return "", ErrNoMorePages
// }

func ListRepoTags(repository *RepositoryRequest) (semVerTags []*semver.Version, errTags []string, err error) {

	// setOfDigits := spinner.GenerateNumberSequence(25)
	// s := spinner.New(setOfDigits, 1000*time.Millisecond)
	// s.Prefix = "running ..... "
	// s.HideCursor = true
	// utils.PrintN(utils.Info, fmt.Sprintf("registry.tags url=%s repository=%s", repository.RepositoryUrl, repository.Repo))
	// fmt.Println()
	// s.Start()
	var semTags []*semver.Version
	allTags, err := listAllTags(repository)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "get repo all tags failure")
	}
	allSemVer, errT := tags2SemVer(allTags)
	if allSemVer != nil {
		if repository.Major == 0 {
			semTags = getLastN(allSemVer, repository.PageSize)
		} else {
			semTags = getMajorVersions(allSemVer, repository.Major, repository.Minor)
		}
	}
	// s.Stop()
	return semTags, errT, nil
}
