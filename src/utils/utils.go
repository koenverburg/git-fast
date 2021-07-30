package utils

import (
	"log"

	"github.com/koenverburg/git-fast/types"
)

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

func IsEmpty(str string) bool {
  if str == "" && str != " " {
    return true
  }
  return false
}

func CreateSegment(value string, part string) types.Segment {
  var s types.Segment

  s.Value = value
  s.Part = part

  return s
}
