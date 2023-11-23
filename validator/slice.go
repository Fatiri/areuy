package validator

func CheckValueIsAvailableInSliceInt(value interface{}, slice []int) bool {
	// iterate using the for loop
	for _, sliceValue := range slice {
		// check
		if sliceValue == value {
			// return true
			return true
		}
	}
	return false
}

func CheckValueIsAvailableInSliceStr(value interface{}, slice []string) bool {
	// iterate using the for loop
	for _, sliceValue := range slice {
		// check
		if sliceValue == value {
			// return true
			return true
		}
	}
	return false
}

func CheckValueIsAvailableInSliceInterface(value interface{}, slice []interface{}) bool {
	// iterate using the for loop
	for _, sliceValue := range slice {
		// check
		if sliceValue == value {
			// return true
			return true
		}
	}
	return false
}
