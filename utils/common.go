package utils

func StrArrayToMap(input []string) map[string]bool {
	output := make(map[string]bool, len(input))
	for idx := range input {
		output[input[idx]] = true
	}

	return output
}
