package ingester

func mp(arr []string) map[string]struct{} {
	res := make(map[string]struct{}, len(arr))

	for _, v := range arr {
		res[v] = struct{}{}
	}
	return res
}
