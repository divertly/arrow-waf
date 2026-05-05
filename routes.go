package main

import (
	"bytes"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

func CorazaMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(e echo.Context) error {
		timer := prometheus.NewTimer(wafProcessLatency)
		defer timer.ObserveDuration()
		tx := G.CZA.NewTransactionWithID(e.Get("RequestID").(string))
		defer func() {
			tx.ProcessLogging()
			_ = tx.Close()
		}()
		G.Log.Trace().Str("id", tx.ID()).Str("ip", e.RealIP()).Msg("Processing WAF request")
		tx.ProcessConnection(e.RealIP(), int(G.Config.WAF.Port), G.Config.Upstream.Host, int(G.Config.Upstream.Port))
		if tx.IsInterrupted() {
			i := tx.Interruption()
			rule := (*G.Rules)[i.RuleID]
			G.Log.Warn().Int("rule_id", i.RuleID).Str("msg", rule.Msg).Str("action", rule.Action).Str("severity", rule.Severity).Str("version", rule.Version).Uint("phase", rule.Phase).Str("tx_id", tx.ID()).Str("ip", e.RealIP()).Strs("tags", rule.Tag).Msg("Request matched WAF rule")
			G.Log.Trace().Str("phase", "process_connection").Str("action", i.Action).Str("data", i.Data).Int("rule_id", i.RuleID).Int("status", i.Status).Str("ip", e.RealIP()).Msg("Connection interrupted")
			wafOperationsCount.WithLabelValues("process_connection", i.Action).Inc()
			wafRuleCount.WithLabelValues(strconv.Itoa(i.RuleID)).Inc()
			return e.NoContent(http.StatusForbidden)
		}
		tx.ProcessURI(e.Request().RequestURI, e.Request().Method, e.Request().Proto)
		if tx.IsInterrupted() {
			i := tx.Interruption()
			rule := (*G.Rules)[i.RuleID]
			G.Log.Warn().Int("rule_id", i.RuleID).Str("msg", rule.Msg).Str("action", rule.Action).Str("severity", rule.Severity).Str("version", rule.Version).Uint("phase", rule.Phase).Str("tx_id", tx.ID()).Str("ip", e.RealIP()).Strs("tags", rule.Tag).Msg("Request matched WAF rule")
			G.Log.Trace().Str("phase", "process_uri").Str("action", i.Action).Str("data", i.Data).Int("rule_id", i.RuleID).Int("status", i.Status).Str("uri", e.Request().RequestURI).Msg("Connection interrupted")
			wafOperationsCount.WithLabelValues("process_uri", i.Action).Inc()
			wafRuleCount.WithLabelValues(strconv.Itoa(i.RuleID)).Inc()
			return e.NoContent(http.StatusForbidden)
		}
		for k, v := range e.Request().Header {
			for _, h := range v {
				tx.AddRequestHeader(k, h)
				G.Log.Trace().Str(k, h).Msg("Request header")
			}
		}
		inter := tx.ProcessRequestHeaders()
		if inter != nil {
			i := tx.Interruption()
			rule := (*G.Rules)[i.RuleID]
			G.Log.Warn().Int("rule_id", i.RuleID).Str("msg", rule.Msg).Str("action", rule.Action).Str("severity", rule.Severity).Str("version", rule.Version).Uint("phase", rule.Phase).Str("tx_id", tx.ID()).Str("ip", e.RealIP()).Strs("tags", rule.Tag).Msg("Request matched WAF rule")
			G.Log.Trace().Str("phase", "process_headers").Str("action", i.Action).Str("data", i.Data).Int("rule_id", i.RuleID).Int("status", i.Status).Any("headers", e.Request().Header).Msg("Connection interrupted")
			wafOperationsCount.WithLabelValues("process_headers", i.Action).Inc()
			wafRuleCount.WithLabelValues(strconv.Itoa(i.RuleID)).Inc()
			return e.NoContent(http.StatusForbidden)
		}
		b := bytes.NewBuffer(nil)
		_, err := b.ReadFrom(e.Request().Body)
		if err != nil {
			e.Logger().Error(err)
		}
		inter, _, err = tx.WriteRequestBody(b.Bytes())
		if err != nil {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, "There was an error processing the request; try again later")
		}
		if inter != nil {
			i := tx.Interruption()
			rule := (*G.Rules)[i.RuleID]
			G.Log.Warn().Int("rule_id", i.RuleID).Str("msg", rule.Msg).Str("action", rule.Action).Str("severity", rule.Severity).Str("version", rule.Version).Uint("phase", rule.Phase).Str("tx_id", tx.ID()).Str("ip", e.RealIP()).Strs("tags", rule.Tag).Msg("Request matched WAF rule")
			G.Log.Trace().Str("phase", "write_body").Str("action", i.Action).Str("data", i.Data).Int("rule_id", i.RuleID).Int("status", i.Status).Bytes("body", b.Bytes()).Msg("Connection interrupted")
			wafOperationsCount.WithLabelValues("write_body", i.Action).Inc()
			wafRuleCount.WithLabelValues(strconv.Itoa(i.RuleID)).Inc()
			return e.NoContent(http.StatusForbidden)
		}
		inter, err = tx.ProcessRequestBody()
		if err != nil {
			e.Logger().Error(err)
			return e.String(http.StatusInternalServerError, "There was an error processing the request; try again later")
		}
		if inter != nil {
			i := tx.Interruption()
			rule := (*G.Rules)[i.RuleID]
			G.Log.Warn().Int("rule_id", i.RuleID).Str("msg", rule.Msg).Str("action", rule.Action).Str("severity", rule.Severity).Str("version", rule.Version).Uint("phase", rule.Phase).Str("tx_id", tx.ID()).Str("ip", e.RealIP()).Strs("tags", rule.Tag).Msg("Request matched WAF rule")
			G.Log.Trace().Str("phase", "process_body").Str("action", i.Action).Str("data", i.Data).Int("rule_id", i.RuleID).Int("status", i.Status).Bytes("body", b.Bytes()).Msg("Connection interrupted")
			wafOperationsCount.WithLabelValues("process_body", i.Action).Inc()
			wafRuleCount.WithLabelValues(strconv.Itoa(i.RuleID)).Inc()
			return e.NoContent(http.StatusForbidden)
		}
		G.Log.Trace().Str("phase", "process_success").Str("action", "pass").Msg("Request passed inspection")
		wafOperationsCount.WithLabelValues("process_success", "pass").Inc()
		return next(e)
	}
}
