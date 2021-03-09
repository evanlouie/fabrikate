package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/microsoft/fabrikate/internal/cmd"
)

type build struct {
	GOOS   string
	GOARCH string
}

func (b build) binName() string {
	outfile := "fab"
	if b.GOOS == "windows" {
		outfile = outfile + ".exe"
	}

	return outfile
}

// build fabrikate for the target build and return the binary as []byte.
func (b build) build() ([]byte, error) {
	// ensure that the path to main.go is correct
	fabCore := filepath.Join("cmd", "fab", "main.go")
	if _, err := os.Stat(fabCore); err != nil {
		return nil, fmt.Errorf(`unable to find Fabrikate core main package %s: %w`, fabCore, err)
	}

	// create a temp file to build to.
	tempFile, err := os.CreateTemp("", "fabrikate")
	if err != nil {
		return nil, fmt.Errorf(`creating temporary build file: %w`, err)
	}
	defer tempFile.Close()

	// Write to the temp file
	// NOTE we only use the tempFile to generate a random safe name. After this
	// command runs, the tempFile we created no longer points to the correct
	// file header as the `go build` command created a new file in its place.
	cmd := exec.Command("go", "build", "-o", tempFile.Name(), fabCore)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		fmt.Sprintf(`GOOS=%s`, b.GOOS),
		fmt.Sprintf(`GOARCH=%s`, b.GOARCH),
	)
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(`running build commands "%s": %w`, cmd, err)
	}

	// Have to open the newly created build file as the tempfile was overwritten
	// by the build command.
	fab, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf(`opening fabrikate build %s: %w`, tempFile.Name(), err)
	}
	if err := fab.Chmod(os.ModePerm); err != nil {
		return nil, fmt.Errorf(`making fabrikate binary executable: %w`, err)
	}
	defer fab.Close()
	defer os.Remove(fab.Name())

	fabBytes, err := io.ReadAll(fab)
	if err != nil {
		return nil, fmt.Errorf(`reading build bytes from temp file %s: %w`, tempFile.Name(), err)
	}

	return fabBytes, nil
}

// buildandZip fabrikate to the target zip file.
func (b build) buildAndZip(zipName string) error {
	outZip, err := os.Create(zipName)
	if err != nil {
		return fmt.Errorf(`creating zip file %s: %w`, outZip.Name(), err)
	}
	defer outZip.Close()

	zWriter := zip.NewWriter(outZip)
	defer zWriter.Close()
	_, filename := filepath.Split(b.binName())
	// do not use writer.Create to create the file -- manually create a FileHeader
	// so we can SetMode to keep the file executable.
	fabHeader := &zip.FileHeader{
		Name:   filename,
		Method: zip.Deflate,
	}
	fabHeader.SetMode(os.ModePerm)
	fabFile, err := zWriter.CreateHeader(fabHeader)
	if err != nil {
		return fmt.Errorf(`creating fab binary zip header: %w`, err)
	}
	fabBytes, err := b.build()
	if err != nil {
		return fmt.Errorf(`building fabrikate: %w`, err)
	}
	if _, err := fabFile.Write(fabBytes); err != nil {
		return fmt.Errorf(`writing fabrikate binary to zip file: %w`, err)
	}

	return nil
}

func main() {
	version := flag.String("version", cmd.Version, "specify the version tag for this release")
	flag.Parse()

	var builds = []build{
		{"darwin", "amd64"},
		{"linux", "amd64"},
		{"windows", "amd64"},
	}

	// ensure that the path to main.go is correct
	fabCore := filepath.Join("cmd", "fab", "main.go")
	if _, err := os.Stat(fabCore); err != nil {
		log.Fatal(fmt.Errorf(`unable to find fabrikate core main package %s: %w`, fabCore, err))
	}

	const releasesDir = "_releases"
	for _, build := range builds {
		log.Printf("Building %s-%s...\n", build.GOOS, build.GOARCH)
		// buildDir, _ := filepath.Split(build.outPath())
		if err := os.MkdirAll(releasesDir, os.ModePerm); err != nil {
			log.Fatal(fmt.Errorf(`creating release directory: %w`, err))
		}
		zipName := fmt.Sprintf(`fab-%s-%s-%s.zip`, *version, build.GOOS, build.GOARCH)
		zipPath := filepath.Join(releasesDir, zipName)
		if err := build.buildAndZip(zipPath); err != nil {
			log.Fatal(fmt.Errorf(`building and zipping fabrikate: %w`, err))
		}

		log.Printf("%s complete\n", zipPath)
	}

	log.Println("Done!")
}
