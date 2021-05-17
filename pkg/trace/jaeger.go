package trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Config struct {
	Agent       string  `yaml:"agent"`
	Sampling    string  `yaml:"sampling"`
	ServiceName string  `yaml:"serviceName"`
	LogSpan     bool    `yaml:"logSpan"`
	Type        string  `yaml:"type"`
	Param       float64 `yaml:"param"`
}

func New(config *Config) (opentracing.Tracer, func(), error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			SamplingServerURL: config.Sampling,
			Type:              jaeger.SamplerTypeConst,
			Param:             1,
		},

		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: config.Agent,
			LogSpans:           config.LogSpan,
		},
	}

	if config.Type != "" {
		cfg.Sampler.Type = config.Type
		cfg.Sampler.Param = config.Param
	}

	cfg.ServiceName = config.ServiceName
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, func() {
		closer.Close()
	}, nil
}
