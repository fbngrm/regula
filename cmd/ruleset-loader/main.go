// The ruleset-loader reads a json file containing a list of rulesets and sends them to etcd.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/heetch/rules-engine/rule"
	"github.com/heetch/rules-engine/store"
)

var (
	file = flag.String("f", "", "path to file containing rulesets informations")
	addr = flag.String("addr", "127.0.0.1:2379,etcd:2379", "coma separated list of the etcd endpoint addresses.")
)

func main() {
	flag.Parse()

	if *file == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	loadSnapshot(f)
}

func newClient() (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(*addr, ","),
		DialTimeout: 5 * time.Second,
	})
}

func loadSnapshot(r io.Reader) error {
	var snapshot map[string]*rule.Ruleset

	err := json.NewDecoder(r).Decode(&snapshot)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	for path, rs := range snapshot {
		path = strings.TrimSpace(path)
		if path == "" {
			return errors.New("empty path")
		}

		path = strings.Trim(path, "/")

		raw, err := json.Marshal(&store.RulesetEntry{
			Path:    path,
			Ruleset: rs,
		})
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = client.Put(ctx, path, string(raw))
		if err != nil {
			return err
		}

		fmt.Printf("Ruleset \"%s\" successfully saved.\n", path)
	}

	return nil
}
