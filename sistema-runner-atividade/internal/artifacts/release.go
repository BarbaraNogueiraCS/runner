package artifacts

import (
	"encoding/json"
	"os"
)

type ReleaseMetadata struct {
	Artifact       string `json:"artifact"`
	Version        string `json:"version"`
	URL            string `json:"url"`
	ChecksumSHA256 string `json:"checksumSha256"`
}

type ReleaseJSON struct {
	Jar struct {
		URL     string `json:"url"`
		Version string `json:"version"`
		SHA256  string `json:"sha256"`
	} `json:"jar"`
}

func ReadReleaseJSON(path string) (ReleaseMetadata, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return ReleaseMetadata{}, err
	}
	var data ReleaseJSON
	if err := json.Unmarshal(content, &data); err != nil {
		return ReleaseMetadata{}, err
	}
	return ReleaseMetadata{Artifact: "simulador.jar", Version: data.Jar.Version, URL: data.Jar.URL, ChecksumSHA256: data.Jar.SHA256}, nil
}
