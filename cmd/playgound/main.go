package main

import (
	"fmt"
	"log"

	"github.com/microsoft/fabrikate/internal/helm"
)

type BuildInfo struct {
	Version      string
	GitCommit    string
	GitTreeState string
	GoVersion    string
}

func main() {
	// out, err := helm.TemplateWithCRDs(helm.TemplateOptions{
	// 	Chart:   "traefik",
	// 	Repo:    "https://helm.traefik.io/traefik",
	// 	Version: "9.14.3",
	// 	Release: "traefik",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(out)

	// if err := yaml.NewDecoder(bytes.NewReader([]byte(out))).Decode(&maps); err != nil {
	// 	log.Fatal(err)
	// }

	// maps, err := yaml.DecodeMaps([]byte(out))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// b, err := json.MarshalIndent(maps, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))

	buildInfo, err := helm.Version()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", buildInfo)
	fmt.Println(buildInfo.IsHelm2())
	fmt.Println(buildInfo.IsHelm3())
	parsed, err := buildInfo.Parse()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", parsed)

	var foo []string
	foo = append(foo, "foo", "bar", "baz")
	fmt.Println(foo[10])
}
