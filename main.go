package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/kkyr/fig"
	"github.com/prometheus/client_golang/prometheus"
)

var G *Global
var VERSION = "0.0.0-local" // defined in ldflags at build time

func main() {
	env, ok := os.LookupEnv("APP_PROFILE")
	if !ok {
		env = "local"
	}
	G = &Global{
		Env:      env,
		Config:   &Config{},
		Handlers: []chan bool{},
	}
	err := fig.Load(
		G.Config,
		fig.File(fmt.Sprintf("profiles/%s.yaml", env)),
		fig.UseEnv("APP"),
	)
	if err != nil {
		panic(err)
	}
	G.Metrics = &Metrics{
		Registry: prometheus.NewRegistry(),
	}
	RegisterMetrics()
	G.Log, err = configureLogging(env, G.Config.Log.Level)
	if err != nil {
		G.Log.Error().Err(err).Str(this()).Msg("Error configuring logging facility; defaulting to INFO")
	}
	logmain := G.Log.With().Str("module", "main").Logger()
	if G.Config.Testing.ProfileCPU {
		logmain.Info().Str(this()).Msg("Starting CPU profiling...")
		f, err := os.Create("cpu.pprof")
		if err != nil {
			logmain.Fatal().Err(err).Msg("Could not create CPU profile file.")
		}
		defer f.Close() //nolint:all
		if err := pprof.StartCPUProfile(f); err != nil {
			logmain.Fatal().Err(err).Msg("Could not take CPU profile")
		}
	}
	if G.Config.Testing.ProfileRAM {
		logmain.Info().Str(this()).Msg("Starting RAM profiling...")
		f, err := os.Create("memory.pprof")
		if err != nil {
			logmain.Fatal().Err(err).Msg("Could not create memory profile file.")
		}
		defer f.Close() //nolint:all
		runtime.GC()
		if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
			logmain.Fatal().Err(err).Msg("Could not take memory profile")
		}
	}
	logmain.Info().Str(this()).Msgf("WAF v%s", VERSION)
	logmain.Debug().Str(this()).Msg("Configured the logger; starting the WAF!")
	logmain.Info().Str(this()).Msgf("PROFILE: %s", env)
	logmain.Trace().Any("config", G.Config).Str(this()).Msg("CONFIG")

	/// DO STUFF

	logmain.Info().Str(this()).Msg("Startup complete!")
	os.Exit(catchSignals(G.Handlers))
}
