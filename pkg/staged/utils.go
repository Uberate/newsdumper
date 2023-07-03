package staged

func TryToMap(m interface{}) (map[string]any, bool) {
	res, ok := m.(map[string]any)
	return res, ok
}

func GetFromMap[T any](m map[string]any, key string, defaultValue T) (T, bool) {
	if len(m) == 0 {
		return defaultValue, false
	}

	if res, ok := m[key]; ok {
		if result, canConvert := res.(T); canConvert {
			return result, true
		}
	}

	return defaultValue, false
}
