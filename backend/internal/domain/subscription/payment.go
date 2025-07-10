package subscription

import (
	"sort"
)

func sortPayments(payments []Payment) {
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].startDate.After(payments[j].startDate)
	})

	var previous *Payment
	//var next *Payment
	for idx := 0; idx < len(payments); idx++ {
		current := &payments[idx]
		if idx == 0 {
			previous = nil
		} else {
			previous = &payments[idx-1]
		}
		//if idx+1 < len(payments) {
		//	next = &payments[idx+1]
		//}

		if previous != nil {
			if current.endDate == nil {
				current.endDate = &previous.startDate
			}
			if current.endDate.After(previous.startDate) {
				current.endDate = &previous.startDate
			}
		}

		payments[idx] = *current
	}

}
