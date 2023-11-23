package validator

func CheckValueIsAvailableInSlice(value interface{}, slice ...interface{}) bool {

	// iterate using the for loop
	for i := 0; i < len(slice); i++ {
		// check
		if slice[i] == value {
			// return true
			return true
		}
	}
	return false
}