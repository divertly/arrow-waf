package main

import (
	"github.com/corazawaf/coraza/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type Global struct {
	Env      string
	Handlers []chan bool
	Log      zerolog.Logger
	Config   *Config
	Metrics  *Metrics
	WAF      *WAF
	CZA      coraza.WAF
	Rules    *CRSRuleList
}

type Config struct {
	Core struct {
		TestMode         bool   `json:"testMode" fig:"test_mode"`
		SignalInterval   uint64 `json:"signalInterval" fig:"signal_interval_ms" default:"100"`
		RuleListLocation string `json:"rule_list_location" fig:"rule_list_location" default:"conf/idmap.yaml"`
	} `json:"core"`
	WAF struct {
		Host string `json:"host" fig:"host" default:"0.0.0.0"`
		Port uint   `json:"port" fig:"port" default:"6080"`
	} `json:"waf"`
	Log struct {
		Format string `json:"log" fig:"format" default:"text"`
		Level  string `json:"level" fig:"level" default:"info"`
		Color  bool   `json:"color" fig:"color"`
	} `json:"log"`
	System struct {
		Host string `json:"host" fig:"host" default:"0.0.0.0"`
		Port uint   `json:"port" fig:"port" default:"6081"`
	} `json:"system"`
	Testing struct {
		ProfileCPU bool `json:"profile_cpu" fig:"profile_cpu"`
		ProfileRAM bool `json:"profile_ram" fig:"profile_ram"`
	} `json:"testing" fig:"testing"`
	Upstream struct {
		Host     string `json:"host" fig:"host"`
		Port     uint   `json:"port" fig:"port" default:"80"`
		Protocol string `json:"protocol" fig:"protocol" default:"http"`
	} `json:"upstream"`
}

type Metrics struct {
	Registry *prometheus.Registry
}

type GenericAPIResponse struct {
	Status int
	Body   any `json:"body,omitempty"`
	Meta   any `json:"meta,omitempty"`
	Error  any `json:"error,omitempty"`
}

type RedirectResponse struct {
	Found bool   `json:"found"`
	Type  string `json:"type"`
	Body  string `json:"body"`
}

type CRSRule struct {
	Msg      string   `yaml:"msg"`
	Severity string   `yaml:"severity"`
	Phase    uint     `yaml:"phase"`
	Action   string   `yaml:"action"`
	Version  string   `yaml:"ver"`
	Tag      []string `yaml:"tag"`
}

type CRSRuleList map[int]CRSRule
