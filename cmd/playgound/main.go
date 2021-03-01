package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/microsoft/fabrikate/internal/helm"
)

func main() {
	out, err := helm.TemplateWithCRDs(helm.TemplateOptions{
		Chart:   "traefik",
		Repo:    "https://helm.traefik.io/traefik",
		Version: "9.14.3",
		Release: "traefik",
	})
	if err != nil {
		log.Fatal(err)
	}

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

	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	v, err := helm.Version()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", v)

	fmt.Println(v.IsHelm3())

	foo := map[string]interface{}{
		"foo": 123,
		// "bar": map[string]interface{}{},
	}

	bar, ok := foo["bar"].(map[string]interface{})
	fmt.Println(ok)
	fmt.Println(bar)
	bar = map[string]interface{}{}
	foo["bar"] = bar
	// bar["baz"] = 123

	fmt.Printf("%+v\n", foo)
	fmt.Printf("%+v\n", foo["baz"])
}
