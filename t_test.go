package main

import (
	"fmt"
	"testing"
)

func Test_parsUrl(t *testing.T) {
	u := `"https://petropavlovsk.b2b.ivest.kz/zhivotnovodstvo-selskoe-hozyaystvo?Company_page=[1:4]"`
	arr, e := parseUrl(u)
	fmt.Println(arr, e)
	// ur, e := url.ParseRequestURI(u)
	// fmt.Println(ur, e)
}
