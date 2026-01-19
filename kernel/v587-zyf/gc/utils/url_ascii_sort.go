package utils

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

type (
	UrlParam struct {
		Key   string
		Value string
	}
	UrlAsciiByKey []UrlParam
)

func (a UrlAsciiByKey) Len() int           { return len(a) }
func (a UrlAsciiByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a UrlAsciiByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

func UrlParamParse(data string) (url.Values, error) {
	return url.ParseQuery(data)
}
func UrlParamCheckExist(urlQuery url.Values, check []string) (bool, string) {
	if len(check) > 0 {
		for _, c := range check {
			if urlQuery.Get(c) == "" {
				return false, c
			}
		}
	}

	return true, ""
}

func UrlParamSort(urlQuery url.Values, check []string, noEmpty bool, filter ...string) (string, error) {
	if len(check) > 0 {
		if res, c := UrlParamCheckExist(urlQuery, check); !res {
			return "", errors.New(fmt.Sprintf("%s is empty", c))
		}
	}

	var params []UrlParam
	for k, vs := range urlQuery {
		if len(filter) > 0 && InCollection(k, filter) {
			continue
		}
		//fmt.Println(k, vs[0])
		if noEmpty && (vs[0] == "" || vs[0] == "0" || vs[0] == "omitempty") {
			continue
		}
		params = append(params, UrlParam{Key: k, Value: vs[0]})
	}
	sort.Sort(UrlAsciiByKey(params))

	var sortedQuery strings.Builder
	for _, p := range params {
		if sortedQuery.Len() > 0 {
			sortedQuery.WriteRune('\n')
		}
		sortedQuery.WriteString(p.Key)
		sortedQuery.WriteRune('=')
		sortedQuery.WriteString(p.Value)
	}
	return sortedQuery.String(), nil
}
