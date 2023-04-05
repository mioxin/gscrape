package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Converter interface {
	GetModifySend(in chan any, out chan []byte)
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

func Get_SaveData(arg ...any) error {
	fout := arg[0].(*bufio.Writer)
	request := arg[1].(*HttpHelper)
	mt := arg[2].(*sync.Mutex)
	count := arg[3].(int)

	var h_block Converter = NewOrgJson(count)

	r := request.Get()
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
	out_chan := make(chan []byte)

	go r.GetHtmlBySelector(SELECTOR, html_Block_chan)
	go h_block.GetModifySend(html_Block_chan, out_chan)

	for byt := range out_chan {
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
	f, err := os.OpenFile(OUT_FILE, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("Cant open ouput file", OUT_FILE, "\r\n", err)
	}
	defer f.Close()

	fout := bufio.NewWriter(f)
	defer fout.Flush()

	pool := Newpool(10)
	mutex := new(sync.Mutex)
	for ind, req := range GetUrlForScrape() {
		go pool.worker(Get_SaveData, fout, req, mutex, ind)
		//fmt.Printf("pool.Wait_g: %v\n", )
	}
	pool.Wait_g.Wait()
	log.Printf("END. OK... Time %v ms.\n WG %v", time.Since(time_start).Milliseconds(), pool.Wait_g)

}
