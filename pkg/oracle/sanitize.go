package oracle

// SanitizePrices drops non-positive readings so invalid oracle data never affects storage.
func SanitizePrices(in map[string]int64) map[string]int64 {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]int64, len(in))
	for name, price := range in {
		if price > 0 {
			out[name] = price
		}
	}
	return out
}
