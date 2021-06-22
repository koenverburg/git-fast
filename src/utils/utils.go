package utils

import "log"

// @Summary CheckIfError
// @Description
// @Param  eError
func CheckIfError(e error) {
  if e != nil {
    log.Fatal(e)
    panic(e)
  }
}


// @Summary FilterEmptyString
// @Description 
// @Param  s[]string 
func FilterEmptyString(s []string ) []string {
  var cleaned []string
	for _, str := range s {
    if str != "" {
      cleaned = append(cleaned, str)
    }
	}
	return cleaned
}
