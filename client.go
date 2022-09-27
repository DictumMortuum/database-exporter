package main

import (
	"github.com/DictumMortuum/servus/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Client struct {
	httpClient http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}
}

func (c *Client) Metrics() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := c.getStatistics()
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		promhttp.Handler().ServeHTTP(writer, request)
	}
}

func (c *Client) getStatistics() error {
	db, err := sqlx.Connect("mysql", Cfg.Databases["mariadb"])
	if err != nil {
		return err
	}
	defer db.Close()

	var tmp []models.KeyVal
	err = db.Select(&tmp, `select * from tkeyval`)
	if err != nil {
		return err
	}

	for _, modem := range tmp {
		var stats models.Modem
		err = modem.Unmarshal(&stats)
		if err != nil {
			return err
		}

		Uptime.WithLabelValues(stats.Host).Set(float64(stats.Uptime))
		CurrentUp.WithLabelValues(stats.Host).Set(float64(stats.CurrentUp))
		CurrentDown.WithLabelValues(stats.Host).Set(float64(stats.CurrentDown))
		CRCUp.WithLabelValues(stats.Host).Set(float64(stats.CRCUp))
		CRCDown.WithLabelValues(stats.Host).Set(float64(stats.CRCDown))
		MaxUp.WithLabelValues(stats.Host).Set(float64(stats.MaxUp))
		MaxDown.WithLabelValues(stats.Host).Set(float64(stats.MaxDown))
		DataUp.WithLabelValues(stats.Host).Set(float64(stats.DataUp))
		DataDown.WithLabelValues(stats.Host).Set(float64(stats.DataDown))
		FECUp.WithLabelValues(stats.Host).Set(float64(stats.FECUp))
		FECDown.WithLabelValues(stats.Host).Set(float64(stats.FECDown))
		SNRUp.WithLabelValues(stats.Host).Set(float64(stats.SNRUp))
		SNRDown.WithLabelValues(stats.Host).Set(float64(stats.SNRDown))

		var isEnabled int = 0
		if stats.Status == true {
			isEnabled = 1
		}

		Status.WithLabelValues(stats.Host).Set(float64(isEnabled))

		var isVoipEnabled int = 0
		if stats.VoipStatus == true {
			isVoipEnabled = 1
		}

		VoipStatus.WithLabelValues(stats.Host).Set(float64(isVoipEnabled))
	}

	return nil
}
