package main

import (
	"github.com/Luqqk/goexrates/internal/data"
)

type ErrorSerializer struct {
	data interface{}
}

func (es *ErrorSerializer) dump() interface{} {
	return struct {
		Error interface{} `json:"error"`
	}{
		Error: es.data,
	}
}

type HealthCheckSerializer struct {
	data map[string]string
}

func (hs *HealthCheckSerializer) dump() interface{} {
	type SystemInfo struct {
		Environment string `json:"environment"`
		Version     string `json:"version"`
	}

	return struct {
		Status     string     `json:"status"`
		SystemInfo SystemInfo `json:"system_info"`
	}{
		Status: hs.data["status"],
		SystemInfo: SystemInfo{
			Environment: hs.data["environment"],
			Version:     hs.data["version"],
		},
	}
}

type RatesSerializer struct {
	data   []*data.Currency
	source string
}

func (rs *RatesSerializer) dump() interface{} {
	date := ""
	if len(rs.data) > 0 {
		date = rs.data[0].CreatedAt.Format("2006-01-02")
	}
	rates := make(map[string]interface{})
	for _, c := range rs.data {
		rates[c.Code] = c.Rate
	}

	return struct {
		Date   string                 `json:"date"`
		Rates  map[string]interface{} `json:"rates"`
		Source string                 `json:"source"`
	}{
		Date:   date,
		Rates:  rates,
		Source: rs.source,
	}
}
