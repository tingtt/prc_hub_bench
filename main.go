package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tingtt/prc_hub_bench/application/benchmark"
	"github.com/tingtt/prc_hub_bench/infrastructure/externalapi/backend"
	"github.com/tingtt/prc_hub_bench/presentation/marshaler_json"
	"github.com/tingtt/prc_hub_bench/presentation/marshaler_yaml"

	flags "github.com/jessevdk/go-flags"
)

type options struct {
	Target       string `short:"t" long:"target" description:"Benchmark target" default:"http://localhost:1323"`
	OutputFormat string `short:"o" long:"output" description:"Output format" default:"json"`
}

func main() {
	// Options
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// Marshaler
	marshaler := marshaler_json.New("	")
	switch opts.OutputFormat {
	case "json", "JSON":
		break
	case "yaml", "YAML", "yml":
		marshaler = marshaler_yaml.New()
	default:
		b, _ := marshaler.Marshal(benchmark.Result{
			Score: 0,
			Error: fmt.Sprintf("invalid format option '%s'", opts.OutputFormat)},
		)
		fmt.Println(string(b))
		os.Exit(1)
	}

	// Run benchmark
	c, err := backend.NewClient(opts.Target)
	if err != nil {
		r := benchmark.Result{
			Score: 0,
			Error: err.Error(),
		}
		b, _ := marshaler.Marshal(r)
		fmt.Println(string(b))
		os.Exit(1)
	}
	r := benchmark.Run(c, time.Minute)
	b, _ := marshaler.Marshal(r)
	fmt.Println(string(b))
	if r.Error != "" {
		os.Exit(1)
	}
}
