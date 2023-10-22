package utils

func In[T comparable](element T, array []T) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}
	return false
}
