// Code generated by "stringer -type=race"; DO NOT EDIT.

package internal

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[orcs-0]
	_ = x[humans-1]
}

const _race_name = "orcshumans"

var _race_index = [...]uint8{0, 4, 10}

func (i race) String() string {
	if i < 0 || i >= race(len(_race_index)-1) {
		return "race(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _race_name[_race_index[i]:_race_index[i+1]]
}
