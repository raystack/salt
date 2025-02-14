package utils

func StringFoundInArray(item string, arr []string) bool {
	for _, curr := range arr {
		if curr == item {
			return true
		}
	}
	return false
}
