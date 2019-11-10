package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cnlubo/go-docker-search/version"
	"net"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// file to check to determine Operating System
const etcOsRelease = "/etc/os-release"

// GetUsername return current username
func GetUsername() string {
	username := ""
	u, err := user.Current()
	if err == nil {
		username = u.Username
	}
	return username
}

// GetNodeIP fetches node ip via command hostname.
// If it fails to get this, return empty string directly.
func GetNodeIP() string {
	output, err := exec.Command("hostname", "-i").CombinedOutput()
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		ip := scanner.Text()
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	return ""
}

// GetOSName gets data in /etc/os-release and gets OS name.
// For example, in a Ubuntu host, fetched data are like:
// root@i-8brpbc9t:~# cat /etc/os-release
// NAME="Ubuntu"
// VERSION="16.04.2 LTS (Xenial Xerus)"
// ID=ubuntu
// ID_LIKE=debian
// PRETTY_NAME="Ubuntu 16.04.2 LTS"
// VERSION_ID="16.04"
// HOME_URL="http://www.ubuntu.com/"
// SUPPORT_URL="http://help.ubuntu.com/"
// BUG_REPORT_URL="http://bugs.launchpad.net/ubuntu/"
// VERSION_CODENAME=xenial
// UBUNTU_CODENAME=xenial
func GetOSName() (string, error) {
	etcOsReleaseFile, err := os.Open(etcOsRelease)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("failed to open %s: %v", etcOsRelease, err)
		}
	}
	defer etcOsReleaseFile.Close()

	var prettyName string

	scanner := bufio.NewScanner(etcOsReleaseFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "PRETTY_NAME=") {
			continue
		}

		data := strings.SplitN(line, "=", 2)
		prettyName = data[1]
		return prettyName, nil
	}

	return "Linux", nil

}

type SystemVersion struct {
	Name string `json:"Name,omitempty"`

	// Arch type of underlying hardware
	Arch string `json:"Arch,omitempty"`

	// The time when this binary of daemon is built
	BuildTime string `json:"BuildTime,omitempty"`

	// Commit ID held by the latest commit operation
	GitCommit string `json:"GitCommit,omitempty"`

	// version of Go runtime
	GoVersion string `json:"GoVersion,omitempty"`

	// Operating system kernel version
	// KernelVersion string `json:"KernelVersion,omitempty"`

	// Operating system type of underlying system
	Os string `json:"Os,omitempty"`

	Version string `json:"Version,omitempty"`
}

func Version() (SystemVersion, error) {

	return SystemVersion{
		Name:      version.AppName,
		Version:   version.AppVersion,
		GoVersion: runtime.Version(),
		Arch:      runtime.GOARCH,
		Os:        runtime.GOOS,
		BuildTime: version.BuildTime,
		GitCommit: version.GitCommit,
	}, nil
}

func Root() bool {
	u, err := user.Current()
	CheckAndExit(err)
	return u.Uid == "0" || u.Gid == "0"
}
