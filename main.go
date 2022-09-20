package main

import (
	"fmt"
	"strings"
)

func main() {
	var userList []string = []string{"wonk", "bob", "martin"}
	var reportList []string = []string{"wonk bob", "bob martin", "shi wonk", "bob martin", "wonk bob", "martin bob"}
	k := 2
	fmt.Println(findBadUser(userList, reportList, k))
}

func isUser(userList []string, userName string) bool {
	for _, v := range userList {
		if v == userName {
			return true
		}
	}
	return false
}

func isDuplicatedReport(history map[string]struct{}, report string) bool {
	if _, ok := history[report]; ok {
		return true
	}
	return false
}

func findBadUser(userList []string, reportList []string, k int) []string {
	history := make(map[string]struct{})
	agg := make(map[string]int)

	validReportCount := 0
	for _, v := range reportList {
		a := strings.Split(v, " ")
		if len(a) != 2 {
			continue
		}

		reporter, reported := a[0], a[1]
		if !isUser(userList, reporter) || !isUser(userList, reported) {
			continue
		}

		if isDuplicatedReport(history, v) {
			continue
		}
		history[v] = struct{}{}

		if count, ok := agg[reported]; ok {
			agg[reported] = count + 1
		} else {
			agg[reported] = 1
		}
		validReportCount++
	}

	result := make([]string, 0, validReportCount)
	for userName, reportedCount := range agg {
		if reportedCount < k {
			continue
		}
		result = append(result, userName)
	}

	return result
}
