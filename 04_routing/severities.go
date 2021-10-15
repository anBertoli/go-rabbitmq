package main

import "strings"

// The list of possible severity values.
var SEVERITIES = [3]string{"info", "warn", "error"}

// Helper functions to validate a list of '-' delimited severities.
func validateSeverities(sevsInput string) ([]string, bool) {
	sevs := strings.Split(sevsInput, "-")
	if len(sevs) == 0 || len(sevs) > 3 {
		return nil, false
	}

	for _, sev := range sevs {
		valid := false
		for _, s := range SEVERITIES {
			if sev == s {
				valid = true
				break
			}
		}
		if !valid {
			return nil, false
		}
	}

	return sevs, true
}
