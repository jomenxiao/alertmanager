package plugin

func getValue(kv KV, key string) string {
	if val, ok := kv[key]; ok {
		return val
	}
	return ""
}
