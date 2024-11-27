//go:build !solution

package hotelbusiness

import (
	"sort"
)

type Guest struct {
	CheckInDate  int
	CheckOutDate int
}

type Load struct {
	StartDate  int
	GuestCount int
}

func ComputeLoad(guests []Guest) []Load {
	loadMap := make(map[int]int)
	for _, guest := range guests {
		for i := guest.CheckInDate; i <= guest.CheckOutDate; i++ {
			loadMap[i]++
		}
		loadMap[guest.CheckOutDate]--
	}

	if len(loadMap) == 0 {
		return []Load{}
	}

	var result []Load

	keys := make([]int, 0, len(loadMap))
	for key := range loadMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	var first int
	var flag bool
	for _, key := range keys {
		if !flag {
			flag = true
			first = loadMap[key]
			result = append(result, Load{key, loadMap[key]})
		} else {
			if first != loadMap[key] {
				result = append(result, Load{key, loadMap[key]})
			}
			first = loadMap[key]
		}
	}
	return result
}
