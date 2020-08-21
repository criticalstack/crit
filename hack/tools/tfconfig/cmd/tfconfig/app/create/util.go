package create

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/blang/semver"
)

var (
	latestKubernetesURL = "https://storage.googleapis.com/kubernetes-release/release/stable.txt"
	releaseURL          = "https://storage.googleapis.com/kubernetes-release/release/stable-%s.txt"
)

// getKubernetesVersions returns a list of the last n minor versions of
// Kubernetes.
func getKubernetesVersions(n int) ([]string, error) {
	client := http.Client{Timeout: 1 * time.Second}
	versions := make([]string, 0)
	resp, err := client.Get(latestKubernetesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	v, err := semver.Parse(strings.TrimPrefix(strings.TrimSpace(string(data)), "v"))
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		resp, err := client.Get(fmt.Sprintf(releaseURL, fmt.Sprintf("1.%d", v.Minor-uint64(i))))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			break
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		v := strings.TrimPrefix(strings.TrimSpace(string(data)), "v")
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	return versions, nil
}
