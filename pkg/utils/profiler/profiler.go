package profiler

import (
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/pkg/logger"
	"log"
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	googleProfiler "cloud.google.com/go/profiler"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func StartProfiler(conf *config.Config) {
	log.Println("Starting profiler ddog...")
	if err := profiler.Start(
		profiler.WithService(conf.Core.Name),
		profiler.WithEnv(conf.Core.Environment),
		profiler.WithVersion(conf.Core.Version),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
		),
	); err != nil {
		log.Fatal(err)
	}
}

func StopProfiler() {
	log.Println("Stopping profiler ddog...")
	profiler.Stop()
}

func StartGoProfiler(conf *config.Config) {
	if !conf.GoProfiler.Enable {
		return
	}

	runtime.SetBlockProfileRate(1)

	log.Println("Starting go profiler, running on port:", conf.GoProfiler.Port)

	// to prevent other http port to have pprof endpoints.
	http.DefaultServeMux = http.NewServeMux()

	go func() {
		serverMuxA := http.NewServeMux()
		serverMuxA.HandleFunc("/debug/pprof/", pprof.Index)
		serverMuxA.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		serverMuxA.HandleFunc("/debug/pprof/profile", pprof.Profile)
		serverMuxA.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		serverMuxA.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// to overcome issue with socketmaster respawning 2nd app.
		for {
			_ = http.ListenAndServe(":"+conf.GoProfiler.Port, serverMuxA) // if open socket successful, it does not go to next line.

			// otherwise sleep 10s and try again until the port is available (i.e. app1 is down).
			time.Sleep(time.Second * 10)
		}
	}()
}

func StartGoCloudProfiler(conf *config.Config) {
	if !conf.GoCloudProfiler.Enable {
		return
	}

	log.Println("Starting go cloud profiler")

	profilerCfg := googleProfiler.Config{
		Service:        conf.Core.Name,
		ServiceVersion: conf.Core.Version,
		// ProjectID must be set if not running on GCP.
		ProjectID: conf.GCP.ProjectId,

		// For OpenCensus users:
		// To see Profiler agent spans in APM backend,
		// set EnableOCTelemetry to true
		EnableOCTelemetry: conf.GoCloudProfiler.OpenTelemetry,
		DebugLogging:      conf.GoCloudProfiler.DebugLogging,
	}

	// Profiler initialization, best done as early as possible.
	if err := googleProfiler.Start(profilerCfg); err != nil {
		logger.Log.Fatal(err)
	}
}
