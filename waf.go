package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

type WAF struct {
	Public  *echo.Echo
	System  *echo.Echo
	Log     zerolog.Logger
	Handler chan bool
}

func NewWAF() (*WAF, error) {
	w := &WAF{
		Log:     G.Log.With().Str("module", "waf").Logger(),
		Public:  echo.New(),
		System:  echo.New(),
		Handler: make(chan bool, 1),
	}
	w.Public = w.configureWebserver(w.Public, "public")
	w.System = w.configureWebserver(w.System, "system")
	w.registerPublicRoutes()
	w.registerSystemRoutes()
	return w, nil
}

func (w *WAF) configureWebserver(e *echo.Echo, serverName string) *echo.Echo {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(newEchoContext(c, serverName))
		}
	})
	e.Use(ULIDMiddleware)
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:          true,
		LogStatus:       true,
		LogProtocol:     true,
		LogRemoteIP:     true,
		LogMethod:       true,
		LogResponseSize: true,
		LogValuesFunc:   w.logValuesFunc,
		Skipper:         healthcheckSkipFunc,
	}))
	e.HideBanner = true
	e.HidePort = true
	if G.Env != "prod" {
		e.Debug = true
	}
	return e
}

func (w *WAF) registerPublicRoutes() {
	w.Public.Use(CorazaMiddleware)
	upstreamHost, err := url.Parse(fmt.Sprintf("%s://%s:%d", G.Config.Upstream.Protocol, G.Config.Upstream.Host, G.Config.Upstream.Port))
	if err != nil {
		w.Log.Fatal().Err(err).Any("upstream", G.Config.Upstream).Msg("Error parsing upstream URL")
	}
	w.Public.Use(middleware.Proxy(middleware.NewRandomBalancer([]*middleware.ProxyTarget{{URL: upstreamHost}})))
}

func (w *WAF) registerSystemRoutes() {
	sys := w.System.Group("/system")
	sys.GET("/ping", RouteSysGetPing)
	sys.GET("/config", RouteSysGetConfig)
	sys.GET("/metrics", echoprometheus.NewHandlerWithConfig(echoprometheus.HandlerConfig{
		Gatherer: G.Metrics.Registry,
	}))
}

func (w *WAF) Run() {
	promHandler := echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: G.Metrics.Registry,
	})
	publicAddr := fmt.Sprintf("%s:%v",
		G.Config.WAF.Host,
		G.Config.WAF.Port,
	)
	w.Log.Debug().Str(this()).Str("server_name", "public").Str("host", G.Config.WAF.Host).Uint("port", G.Config.WAF.Port).Msg("Starting public WAF server")
	w.Public.Use(promHandler)
	go w.Public.Start(publicAddr) //nolint:all
	systemAddr := fmt.Sprintf("%s:%v",
		G.Config.System.Host,
		G.Config.System.Port,
	)
	w.Log.Debug().Str(this()).Str("server_name", "system").Str("host", G.Config.System.Host).Uint("port", G.Config.System.Port).Msg("Starting system server")
	w.System.Use(promHandler)
	go w.System.Start(systemAddr) //nolint:all
	w.HandleSignal()
}

func (w *WAF) HandleSignal() {
	for {
		select {
		case <-w.Handler:
			w.Log.Info().Str(this()).Msg("Recieved handler signal, closing connections")
			err := w.Public.Close()
			if err != nil {
				w.Log.Error().Err(err).Str(this()).Msg("Caught error closing down public server")
			}
			err = w.System.Close()
			if err != nil {
				w.Log.Error().Err(err).Str(this()).Msg("Caught error closing down system server")
			}
			return
		default:
			time.Sleep(time.Duration(G.Config.Core.SignalInterval) * time.Millisecond)
		}
	}
}

func (w *WAF) logValuesFunc(c echo.Context, v middleware.RequestLoggerValues) error {
	ctx := c.(*ServerContext)
	msg := w.Log.Info().
		Str("logtype", "access_log").
		Str("module", "proxy").
		Str("server_name", ctx.ServerName()).
		Str("URI", v.URI).
		Str("host", c.Request().Host).
		Str("protocol", v.Protocol).
		Str("ip", v.RemoteIP).
		Str("method", v.Method).
		Str("request_id", c.Response().Header().Get(echo.HeaderXRequestID)).
		Int64("size", v.ResponseSize).
		Int("status", v.Status)
	msg.Msg("WAF request")
	return nil
}

func healthcheckSkipFunc(c echo.Context) bool {
	skip := map[string]bool{
		"/system/ping?p=healthcheck": true,
		"/system/metrics":            true,
	}
	if val, ok := skip[c.Request().URL.String()]; ok {
		return val
	}
	return false
}
