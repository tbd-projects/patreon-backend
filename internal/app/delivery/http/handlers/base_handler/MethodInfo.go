package base_handler

type MethodInfo struct {
	queries []string
}

func (mi *MethodInfo) AddQuery(key string, value string) *MethodInfo {
	mi.queries = append(mi.queries, key, value)
	return mi
}