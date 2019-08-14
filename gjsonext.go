package main

import (
	"github.com/tidwall/gjson"
)

// GetStringIfExists returns json value or default
func GetStringIfExists(json gjson.Result, key string, def string) string {
	value := json.Get(key)
	if value.Exists() {
		return value.String()
	}
	return def
}

// GetFloatIfExists returns json value or default
func GetFloatIfExists(json gjson.Result, key string, def float64) float64 {
	value := json.Get(key)
	if value.Exists() {
		return value.Float()
	}
	return def
}
