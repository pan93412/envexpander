package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pan93412/envexpander/v2"
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	env := make(map[string]string)

	for r.Scan() {
		s := r.Text()

		// none
		if s == "" {
			continue
		}

		// comment
		if s[0] == '#' {
			continue
		}

		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			log.Fatalf("not a valid line: %s", kv)
		}

		k := kv[0]
		v := kv[1]

		// add to env list
		env[k] = v
	}

	// resolve env variables
	resolvedEnv := envexpander.ResolveEnvVariable(env)

	// print
	for k, v := range resolvedEnv {
		fmt.Printf("%s=%s\n", k, v)
	}
}
