//go:build !solution

package hogwarts

import "errors"

func TopologicalSort(adjList map[string][]string) []string {
	inDegrees := make(map[string]int)
	for node := range adjList {
		inDegrees[node] = 0
	}

	for _, neighbors := range adjList {
		for _, neighbor := range neighbors {
			inDegrees[neighbor]++
		}
	}

	var queue []string
	for node, degree := range inDegrees {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	var result []string
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)

		for _, v := range adjList[u] {
			inDegrees[v]--
			if inDegrees[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(result) < len(adjList) {
		return nil
	}

	return result
}

func GetCourseList(prereqs map[string][]string) []string {
	for _, value := range prereqs {
		for _, elem := range value {
			_, ok := prereqs[elem]
			if !ok {
				prereqs[elem] = []string{}
			}
		}
	}
	result := TopologicalSort(prereqs)
	if result == nil {
		panic(errors.New("there is cycle in the graph"))
	}

	reversed := make([]string, len(result))

	for i, s := range result {
		reversed[len(result)-1-i] = s
	}

	return reversed
}
