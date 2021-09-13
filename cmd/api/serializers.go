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
	data []*data.Currency
	base string
}

func (rs *RatesSerializer) dump() interface{} {
	date := ""
	if len(rs.data) > 0 {
		date = rs.data[0].CreatedAt.Format("2006-01-02")
	}
	rates := make(map[string]float64)
	for _, c := range rs.data {
		if c.Code != rs.base {
			rates[c.Code], _ = c.Rate.Round(5).Float64()
		}
	}
	// Drop the base rate (it will be always equal to 1)
	delete(rates, rs.base)

	return struct {
		Base  string             `json:"base"`
		Date  string             `json:"date"`
		Rates map[string]float64 `json:"rates"`
	}{
		Base:  rs.base,
		Date:  date,
		Rates: rates,
	}
}
