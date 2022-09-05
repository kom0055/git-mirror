package config

import (
	"testing"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"gopkg.in/yaml.v3"
)

func TestParseConfig(t *testing.T) {
	cfg := Cfg{}
	b, _ := yaml.Marshal(cfg)
	_ = yaml.Unmarshal(b, &cfg)
	t.Log(cfg)
}

func TestParseUrl(t *testing.T) {
	e, err := transport.NewEndpoint("git@github.com:kom0055/go-flinx.git")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	e, err = transport.NewEndpoint("https://github.com/kom0055/go-flinx.git")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	e, err = transport.NewEndpoint("http://github.com/kom0055/go-flinx.git")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	e, err = transport.NewEndpoint("https://github.com/kom0055/go-flinx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	e, err = transport.NewEndpoint("http://github.com/kom0055/go-flinx")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)

	e, err = transport.NewEndpoint("http://github.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(e)
}
