package transbytes

import (
	"fmt"
	"strconv"
)

func SizeToString(bytes float64) string {
	size, form := SizeToFloatAndString(bytes)
	return fmt.Sprintf("%+v %s", strconv.FormatFloat(size, 'f', 2, 64), form)
}

func SizeToFloatAndString(bytes float64) (float64, string) {
	size, sizeForm := transformBytes(bytes, 1)
	var form string
	switch sizeForm {
	case 1:
		form = "B"
	case 2:
		form = "KB"
	case 3:
		form = "MB"
	case 4:
		form = "GB"
	case 5:
		form = "TB"
	}
	return size, form
}

func transformBytes(size float64, sizeForm int) (float64, int) {
	if size/1024 >= 1 && sizeForm < 5 {
		return transformBytes(size/1024, sizeForm+1)
	} else {
		return size, sizeForm
	}
}
