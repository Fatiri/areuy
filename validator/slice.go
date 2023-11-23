package validator

func CheckValueIsAvailableInSlice(slice []interface{}, value interface{}) bool {

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
