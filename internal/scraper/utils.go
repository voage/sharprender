package scraper

func SafeString(val interface{}) string {
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func SafeInt(val interface{}) int {
	if num, ok := val.(float64); ok {
		return int(num)
	}
	return 0
}
