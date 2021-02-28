package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v33/github"
)

func untarHelmBin(body []byte) ([]byte, error) {
	byteReader := bytes.NewReader(body)
	gzr, err := gzip.NewReader(byteReader)
	if err != nil {
		return nil, fmt.Errorf(`creating gzip reader: %w`, err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil, fmt.Errorf(`no file with name "helm" found in tar.gz file`)
		case err != nil:
			return nil, fmt.Errorf(`parsing file in tar.gz file: %w`, err)
		case header == nil:
			continue
		default:
			filename := filepath.Base(header.Name)
			if filename == "helm" && header.Typeflag == tar.TypeReg {
				helmBytes, err := io.ReadAll(tr)
				if err != nil {
					return nil, fmt.Errorf(`reading bytes from %s in tar.gz file`, header.Name)
				}
				return helmBytes, nil
			}
		}
	}
}

func unzipHelmBin(body []byte) ([]byte, error) {
	r := bytes.NewReader(body)
	rdr, err := zip.NewReader(r, int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf(`creating zip reader: %s`, err)
	}

	for _, zipFile := range rdr.File {
		filename := filepath.Base(zipFile.Name)
		if filename == "helm.exe" {
			f, err := zipFile.Open()
			if err != nil {
				return nil, err
			}
			helmBytes, err := io.ReadAll(f)
			if err != nil {
				return nil, err
			}
			return helmBytes, err
		}
	}

	return nil, fmt.Errorf(`no file named "helm.exe" found in zip file`)
}

func downloadLatestHelm() ([]byte, error) {
	// get the latest github release
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "helm", "Helm")
	if err != nil {
		return nil, fmt.Errorf(`getting latest release from github for helm/helm: %w`, err)
	}
	if len(strings.TrimSpace(*release.Body)) == 0 {
		return nil, fmt.Errorf(`getting latest release from github for helm/helm: empty release body was found for release %s`, *release.Name)
	}

	// get the correct compressed extension
	var compressExt string
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "linux":
		compressExt = "tar.gz"
	case "windows":
		compressExt = "zip"
	default:
		return nil, fmt.Errorf(`downloading helm binary: unsupported host %s`, runtime.GOOS)
	}

	// download the os specific release
	downloadURL := fmt.Sprintf(`https://get.helm.sh/helm-%s-%s-amd64.%s`, *release.TagName, runtime.GOOS, compressExt)
	resp, err := http.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf(`downloading helm from %s: %w`, downloadURL, err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(`reading body of helm download response from %s: %w`, downloadURL, err)
	}

	// decompress the file and get the helm bin bytes
	var helmBinBytes []byte
	switch runtime.GOOS {
	case "darwin":
		fallthrough
	case "linux":
		helmBinBytes, err = untarHelmBin(bodyBytes)
	case "windows":
		helmBinBytes, err = unzipHelmBin(bodyBytes)
	default:
		return nil, fmt.Errorf(`unsupported os for decompressing downloaded helm binary`)
	}
	if err != nil {
		return nil, fmt.Errorf(`decompressing downloaded helm binary: %w`, err)
	}
	if helmBinBytes == nil {
		return nil, fmt.Errorf(`empty byte slice found for "helm" file in compressed download file`)
	}

	return helmBinBytes, nil
}

func main() {
	b, err := downloadLatestHelm()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d\n", len(b))
}
