package main

import (
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"
)

func catchSignals(Handlers []chan bool) int {
	log := G.Log.With().Str("service", "signal").Logger()
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)
	for {
		//METRICserviceUp.WithLabelValues("ticker", "signal").Inc()
		s := <-sigchnl
		if s == syscall.SIGTERM || s == syscall.SIGINT {
			//METRICsignalCatch.WithLabelValues("token", resolveSignal(s)).Inc()
			log.Info().Str("signal", resolveSignal(s)).Msg("Got signal. Draining log buffer and exiting.")
			for _, h := range Handlers {
				h <- true
			}
			time.Sleep(100 * time.Millisecond) // drain the log buffer
			if G.Config.Testing.ProfileCPU {
				pprof.StopCPUProfile()
			}
			return 0
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func resolveSignal(s os.Signal) string {
	switch s.String() {
	case "no signal":
		return "SIGNONE"
	case "hangup":
		return "SIGHUP"
	case "interrupt":
		return "SIGINT"
	case "quit":
		return "SIGQUIT"
	case "illegal instruction":
		return "SIGILL"
	case "trace/breakpoint trap":
		return "SIGTRAP"
	case "abort":
		return "SIGABRT"
	case "bus error":
		return "SIGBUS"
	case "floating point exception":
		return "SIGFPE"
	case "killed":
		return "SIGKILL"
	case "user defined signal 1":
		return "SIGUSR1"
	case "segmentation fault":
		return "SIGSEGV"
	case "user defined signal 2":
		return "SIGUSR2"
	case "broken pipe":
		return "SIGPIPE"
	case "alarm clock":
		return "SIGALRM"
	case "terminated":
		return "SIGTERM"
	case "child exited":
		return "SIGCHLD"
	case "continued":
		return "SIGCONT"
	case "stopped (signal)":
		return "SIGSTOP"
	case "stopped":
		return "SIGTSTP"
	case "stopped (tty input)":
		return "SIGTTIN"
	case "stopped (tty output)":
		return "SIGTTOU"
	case "urgent I/O condition":
		return "SIGURG"
	case "CPU time limit exceeded":
		return "SIGXCPU"
	case "file size limit exceeded":
		return "SIGXFSZ"
	case "virtual timer expired":
		return "SIGVTARLM"
	case "profiling timer expired":
		return "SIGPROF"
	case "window changed":
		return "SIGWINCH"
	case "I/O possible":
		return "SIGPOLL"
	case "power failure":
		return "SIGPWR"
	case "bad system call":
		return "SIGSYS"
	default:
		return "UNKNOWN"
	}
}
