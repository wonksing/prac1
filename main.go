package main

import (
	"fmt"
	"strings"
)

func main() {
	var userList []string = []string{"wonk", "bob", "martin"}
	var reportList []string = []string{
		"wonk bob", "bob martin", "shi wonk", "bob martin", "wonk bob", "martin bob"}

	k := 2

	fmt.Println(findBadUserRefactored(userList, reportList, k))
}

func findBadUserRefactored(userList []string, reportList []string, k int) []string {
	// 사용자 확인용 저장소
	var userRepo UserRepo = NewUserRepoMap(userList)
	// 이력 확인용 저장소
	var history ReportHistoryRepo = NewReportHistoryRepoMap()
	// 신고당한 사용자 저장소
	var badUserReport BadUserRepo = NewBadUserRepoMap()

	// 신고 목록 순회
	for _, report := range reportList {
		a := strings.Split(report, " ")
		if len(a) != 2 {
			continue
		}

		// 사용자 유효성 체크
		reporter, reported := a[0], a[1]
		if !userRepo.IsValidUser(reporter) || !userRepo.IsValidUser(reported) {
			continue
		}

		if history.IsDuplicatedReport(report) {
			continue
		}
		history.AddReport(report)

		badUserReport.AccumulateReport(reporter, reported)
	}

	return badUserReport.ReportBadUserList(k)
}

type UserRepo interface {
	// 유효한 사용자인지 확인
	IsValidUser(userName string) bool
}
type ReportHistoryRepo interface {
	// 중복 신고인지 확인
	IsDuplicatedReport(report string) bool
	// 신고 추가
	AddReport(report string)
}
type BadUserRepo interface {
	// 신고당한 사용자 건수 누적
	AccumulateReport(reporter, reported string)
	// 슬라이스 타입의 보고서를 반환(신고당한 건수가 minReportCount 이상인 사용자)
	ReportBadUserList(minReportCount int) []string
}

type UserRepoMap struct {
	userMap map[string]struct{}
}

func NewUserRepoMap(userList []string) *UserRepoMap {
	userMap := make(map[string]struct{})
	for _, v := range userList {
		userMap[v] = struct{}{}
	}

	return &UserRepoMap{
		userMap: userMap,
	}
}

func (r *UserRepoMap) IsValidUser(userName string) bool {
	if _, ok := r.userMap[userName]; ok {
		return true
	}
	return false
}

type ReportHistoryRepoMap struct {
	history map[string]struct{}
}

func NewReportHistoryRepoMap() *ReportHistoryRepoMap {
	return &ReportHistoryRepoMap{
		history: make(map[string]struct{}),
	}
}

func (r *ReportHistoryRepoMap) IsDuplicatedReport(report string) bool {
	if _, ok := r.history[report]; ok {
		return true
	}
	return false
}

func (r *ReportHistoryRepoMap) AddReport(report string) {
	r.history[report] = struct{}{}
}

type BadUserRepoMap struct {
	countAggregation map[string]int
}

func NewBadUserRepoMap() *BadUserRepoMap {
	return &BadUserRepoMap{
		countAggregation: make(map[string]int),
	}
}

func (r *BadUserRepoMap) AccumulateReport(reporter, reported string) {
	if count, ok := r.countAggregation[reported]; ok {
		r.countAggregation[reported] = count + 1
	} else {
		r.countAggregation[reported] = 1
	}
}

func (r *BadUserRepoMap) ReportBadUserList(minReportCount int) []string {
	result := make([]string, 0, len(r.countAggregation))
	for userName, reportedCount := range r.countAggregation {
		if reportedCount < minReportCount {
			continue
		}
		result = append(result, userName)
	}

	return result
}
