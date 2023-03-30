package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func get_hresp(h, parm string) *HttpHelperResponse {
	resp := NewHttpHelper().URL(h).
		Header("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		Header("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		Header("accept-encoding", "gzip, deflate, br").
		Param(PARAM_KEY, parm).
		Get()
	return resp
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
	if clear_num[0] == rune('8') {
		clear_num = append([]rune("+7"), clear_num[1:]...)
	} else if clear_num[0] == rune('7') && len(clear_num) < 11 {
		clear_num = append([]rune("+7"), clear_num...)
	} else {
		clear_num = append([]rune("+"), clear_num...)
	}
	return string(clear_num)
}

func GetMobileOnly(phones string) []string {
	ph := make([]string, 0)
	arr_ph := strings.Split(phones, ";")
	arr_ph2 := strings.Split(phones, ",")
	if len(arr_ph) < len(arr_ph2) {
		arr_ph = arr_ph2
	}
	for _, p := range arr_ph {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "8") || strings.HasPrefix(p, "+7") {
			if !(strings.HasPrefix(p, "8715") || strings.HasPrefix(p, "8 715") || strings.HasPrefix(p, "8(715") ||
				strings.HasPrefix(p, "8 (715") || strings.HasPrefix(p, "8-(715") || strings.HasPrefix(p, "8-715") ||
				strings.HasPrefix(p, "+7715") || strings.HasPrefix(p, "+7 715") || strings.HasPrefix(p, "+7(715") ||
				strings.HasPrefix(p, "+7 (715") || strings.HasPrefix(p, "+7-(715") || strings.HasPrefix(p, "+7-715")) {
				ph = append(ph, CleanPhonNumber(p))
			}
		}
	}
	return ph
}

func GetSpesialInfo(str string) string {
	spec := ""
	arr := strings.Split(str, "Специфика: ")
	if len(arr) > 1 && len(arr[1]) > 3 {
		spec = arr[1]
	}
	return spec
}

func GetJsonOrg(in chan any, out chan []byte, num int) {
	defer close(out)
	for b := range in {
		s := b.(*goquery.Selection)
		phones := GetMobileOnly(strings.TrimSpace(s.Find("div.phone_icon").Text()))
		spec := GetSpesialInfo(s.Find("p").Text())
		if len(phones) > 0 {
			neworg := Org{s.Find("a.name_link").Text(), s.Find("a.category_link").Text(), spec, phones}
			jneworg, _ := json.Marshal(neworg)
			out <- jneworg
		}
	}
	//out <- []byte("out GetJsonOrg...")
	fmt.Println("<<<<< CLOSE", num)
}

func Get_SaveData(arg ...any) error {
	fout := arg[0].(*bufio.Writer)
	i := arg[1].(int)
	mt := arg[2].(*sync.Mutex)
	wg := arg[3].(*sync.WaitGroup)
	defer wg.Done()

	r := get_hresp(HREF, strconv.Itoa(i))
	if r.err != nil {
		log.Println("Get_SaveData: error END. ", r.err)
		return r.err
	}
	if !r.OK() {
		log.Println("Get_SaveData: error END. ", r.Status)
		return r.err
	}
	for _, c := range r.response.Cookies() {
		fmt.Println(c.Name, c.Value)
	}

	html_Block_chan := make(chan any)
	json_org := make(chan []byte)

	go GetJsonOrg(html_Block_chan, json_org, i)
	go r.GetHtmlBySelector(SELECTOR, html_Block_chan)

	for byt := range json_org {
		mt.Lock()
		_, err := fout.Write(byt)
		if err != nil {
			log.Println(err, string(byt))
		}
		fout.Write([]byte("\r\n"))
		mt.Unlock()
	}

	log.Println(r.response.Request.URL)
	return nil
}

func main() {
	time_start := time.Now()
	//http://petropavlovsk.b2b.ivest.kz/proizvodstvo/?Company_page=2
	f, err := os.OpenFile(OUT_FILE, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Cant open ouput file", OUT_FILE, "\r\n", err)
	}
	defer f.Close()

	fout := bufio.NewWriter(f)
	defer fout.Flush()

	pool := Newpool(3)
	mutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
LOOP:
	for i := 1; i <= PARAM_VALUE_MAX; i++ {
		wg.Add(1)
		go pool.worker(Get_SaveData, fout, i, mutex, wg)
		select {
		case e := <-pool.err:
			if e != nil {
				fmt.Println("error from pool", e)
				break LOOP
			}
		default:
		}
	}
	// go func() {
	for {
		if e, ok := <-pool.err; ok {
			fmt.Println("error from pool IN END", e)

		} else {
			break
		}

	}
	//}()
	wg.Wait()
	log.Printf("END. OK... Time %v ms.\n", time.Since(time_start).Milliseconds())

}
