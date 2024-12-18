package internal

import "strconv"

// Helper to parse integer with a default fallback
func defaultInt(valStr string, defaultValue int) int {
	val, err := strconv.Atoi(valStr)
	if err != nil || val < 1 {
		return defaultValue
	}
	return val
}
