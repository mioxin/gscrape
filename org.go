package main

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	SELECTOR string = "div.order_com"
)

type Org struct {
	Name, Cat, Spec string
	Phones          []string
}

type OrgJson struct {
	count int
	Org
}

func NewOrgJson(c int) *OrgJson {
	return &OrgJson{c, Org{}}
}

func (o *OrgJson) Scrape(in io.Reader, out chan []byte) {
	defer close(out)
	doc, err := goquery.NewDocumentFromReader(in)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(SELECTOR).Each(func(i int, s *goquery.Selection) {
		phones := GetMobileOnly(splitStr2Arr(strings.TrimSpace(s.Find("div.phone_icon").Text())))
		if len(phones) > 0 {
			o.Name = s.Find("a.name_link").Text()
			o.Cat = s.Find("a.category_link").Text()
			o.Spec = GetSpesialInfo(s.Find("p").Text())
			o.Phones = phones
			jneworg, _ := json.Marshal(o)
			out <- jneworg
		}
	})
}

func splitStr2Arr(str string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return []string{}
	}
	arr_ph := strings.Split(str, ";")
	arr_ph2 := strings.Split(str, ",")
	if len(arr_ph) < len(arr_ph2) {
		arr_ph = arr_ph2
	}
	return arr_ph
}

func GetSpesialInfo(str string) string {
	spec := ""
	arr := strings.Split(str, "Специфика: ")
	if len(arr) > 1 && len(arr[1]) > 3 {
		spec = arr[1]
	}
	return spec
}

func CleanPhonNumber(phone_num string) string {
	clear_num := []rune("")
	cache := make(map[rune]struct{})
	for _, ch := range phone_num {
		_, inMap := cache[ch]
		if inMap || strings.ContainsRune("0123456789", ch) {
			clear_num = append(clear_num, ch)
			cache[ch] = struct{}{}
		}
	}
	if len(clear_num) < 2 {
		return ""
	}
	if clear_num[0] == rune('8') {
		clear_num = append([]rune("+7"), clear_num[1:]...)
	} else if clear_num[0] == rune('7') && len(clear_num) < 11 {
		clear_num = append([]rune("+7"), clear_num...)
	} else {
		clear_num = append([]rune("+"), clear_num...)
	}
	return string(clear_num)
}

func GetMobileOnly(phones []string) []string {
	ph := make([]string, 0)

	for _, p := range phones {
		if len(strings.TrimSpace(p)) < 6 {
			continue
		}
		p = CleanPhonNumber(strings.TrimSpace(p))
		if strings.HasPrefix(p, "+7") && !strings.HasPrefix(p, "+7715") {
			ph = append(ph, p)
		}
	}
	return ph
}
