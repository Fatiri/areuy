package manipulator

//FindStringOnStringSlice ...
func FindStringOnStringSlice(slice []string, val string) (index int, found bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
