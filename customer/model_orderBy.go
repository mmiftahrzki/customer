package customer

import "strings"

type OrderBy struct {
	Column    orderColumn
	Direction orderDirection
}

func (ob OrderBy) String() string {
	if ob.Column == Name {
		nameStr := string(Name)
		nameStrs := strings.Split(nameStr, ", ")

		return nameStrs[0] + " " + string(ob.Direction) + ", " + nameStrs[1] + " " + string(ob.Direction)
	}

	return string(ob.Column) + " " + string(ob.Direction)
}
