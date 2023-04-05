package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var HREF []string = []string{
	"https://petropavlovsk.b2b.ivest.kz/zhivotnovodstvo-selskoe-hozyaystvo",
	//"https://petropavlovsk.b2b.ivest.kz/reklama-smi"
	//"https://petropavlovsk.b2b.ivest.kz/medicina"
	//"https://petropavlovsk.b2b.ivest.kz/gosudarstvo-i-pravo"
	//"https://petropavlovsk.b2b.ivest.kz/finansy"
	//"https://petropavlovsk.b2b.ivest.kz/telekommunikacii-i-svyaz"
	//"https://petropavlovsk.b2b.ivest.kz/dosug"
	//"https://petropavlovsk.b2b.ivest.kz/obrazovanie"
	//"https://petropavlovsk.b2b.ivest.kz/transport"
	//"https://petropavlovsk.b2b.ivest.kz/uslugi"
	//"https://petropavlovsk.b2b.ivest.kz/magaziny-i-torgovye-organizacii"
	//"https://petropavlovsk.b2b.ivest.kz/stroitelnye-materialy-i-instrumenty"
	//"https://petropavlovsk.b2b.ivest.kz/proizvodstvo/"
}

const (
	PARAM_VALUE_MAX int    = 3
	PARAM_KEY       string = "Company_page"
	OUT_FILE        string = "out11.json"
	SELECTOR        string = "div.order_com"
)

func GetUrlForScrape() []*HttpHelper {
	urls := make([]*HttpHelper, 0)
	for _, href := range HREF {
		for i := 1; i <= PARAM_VALUE_MAX; i++ {
			urls = append(urls,
				NewHttpHelper().URL(href).
					Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
					Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
					Header("Accept-Encoding", "gzip, deflate, br").
					Param(PARAM_KEY, strconv.Itoa(i)))
			//fmt.Println(i, href)
		}
	}
	return urls
}

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

func GetSpesialInfo(str string) string {
	spec := ""
	arr := strings.Split(str, "Специфика: ")
	if len(arr) > 1 && len(arr[1]) > 3 {
		spec = arr[1]
	}
	return spec
}

func (o *OrgJson) GetModifySend(in chan any, out chan []byte) {
	defer close(out)
	for b := range in {
		s := b.(*goquery.Selection)
		phones := GetMobileOnly(strings.TrimSpace(s.Find("div.phone_icon").Text()))
		if len(phones) > 0 {
			o.Name = s.Find("a.name_link").Text()
			o.Cat = s.Find("a.category_link").Text()
			o.Spec = GetSpesialInfo(s.Find("p").Text())
			o.Phones = phones
			jneworg, _ := json.Marshal(o)
			out <- jneworg
		}
	}
	//out <- []byte("out GetJsonOrg...")
	fmt.Println("<<<<< CLOSE", o.count)
}
