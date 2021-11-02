package models

import "strconv"

type Money float64

func (m Money) MarshalJSON() ([]byte, error) {
	if float64(m) == float64(int(m)) {
		return []byte(strconv.FormatFloat(float64(m), 'f', 1, 32)), nil
	}
	return []byte(strconv.FormatFloat(float64(m), 'f', -1, 32)), nil
}
