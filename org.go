package main

var HREF string = "https://petropavlovsk.b2b.ivest.kz/zhivotnovodstvo-selskoe-hozyaystvo"

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

const (
	PARAM_VALUE_MAX int    = 6
	PARAM_KEY       string = "Company_page"
	OUT_FILE        string = "out11.json"
	SELECTOR        string = "div.order_com"
)

// func GetUrlForScrape() []*HttpHelper {
// 	urls := make([]*HttpHelper, 0)
// 	for i := 0; i <= PARAM_VALUE_MAX; i++ {
// 		urls = append(urls,
// 			NewHttpHelper().URL(href).
// 				Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
// 				Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
// 				Header("accept-encoding", "gzip, deflate, br").
// 				Param(PARAM_KEY, strconv.Itoa(i)))
// 	}
// 	return urls
// }

type Org struct {
	Name, Cat, Spec string
	Phones          []string
}
