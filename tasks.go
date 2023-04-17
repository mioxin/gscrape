package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

type uri struct {
	href   string
	params []pair
}

type Scraper interface {
	Scrape(in io.Reader, out chan []byte)
}

type Tasks chan *Task

type Task struct {
	reqst *HttpHelper
	num   int
	Scraper
}

func NewTasks(href *bufio.Reader) Tasks {
	urls_chan := make(chan *Task)
	go func(out chan *Task) {
		j := 0
		for {
			href, err := href.ReadString('\n')
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(err)
			}
			j++
			u, err := parseUrl(href)
			if err != nil {
				log.Printf("parseUrl: SKIP. %v", err)
				continue
			}
			for i, uri_src := range u {
				url := NewHttpHelper().URL(uri_src.href)
				for _, param := range uri_src.params {
					url = url.Param(param.k, param.v)
				}
				//fmt.Println(i, href)
				out <- &Task{reqst: url.
					Header("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
					Header("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
					Header("Accept-Encoding", "*"),
					num: i + j}
			}
		}
		close(out)
	}(urls_chan)
	return Tasks(urls_chan)
}

// func (t *Task) parse(html_block Scraper, fout *bufio.Writer, mt *sync.Mutex) error {
func (t *Task) parse(arg ...any) error {
	html_block := arg[0].(Scraper)
	fout := arg[1].(*bufio.Writer)
	mt := arg[2].(*sync.Mutex)

	r := t.reqst.Get()
	if r.err != nil {
		log.Println("parse: error GET. ", r.err)
		return r.err
	}
	if !r.OK() {
		log.Println("parse: error GET. ", r.response.Request.URL, r.Status)
		return r.err
	}

	out_chan := make(chan []byte)
	rd := r.response.Body
	go (html_block).Scrape(rd, out_chan)
	defer r.response.Body.Close()

	for byt := range out_chan {
		mt.Lock()
		_, err := fout.Write(byt)
		if err != nil {
			log.Println("parse: ", err, string(byt))
		}
		fout.Write([]byte("\r\n"))
		mt.Unlock()
	}

	log.Println("parse: ", r.response.Request.URL)
	return nil
}

func parseUrl(s_url string) ([]uri, error) {
	s_url = strings.Trim(strings.TrimSpace(s_url), "\"")
	if len(s_url) < 6 {
		return nil, fmt.Errorf("%s The line is too short", s_url)
	}

	if s_url[:2] == "//" {
		return nil, fmt.Errorf("%s The line is commented", s_url)
	}
	_, err := url.Parse(s_url)
	if err != nil {
		return nil, fmt.Errorf("%s %s", s_url, err)
	}
	arr_uri := make([]uri, 0)
	arr_url := strings.Split(s_url, "?")
	href_path := arr_url[0]
	if len(arr_url) < 2 {
		return []uri{{href: href_path}}, nil
	}

	params := strings.Split(arr_url[1], "&")
	list_prms := make([][]pair, 0)
	for _, pr := range params {
		arr_pr := strings.Split(pr, "=")
		if len(arr_pr) == 2 {
			if prms, err := parseMaskParam(arr_pr); err != nil {
				return nil, fmt.Errorf("%s ", err)
			} else {
				list_prms = append(list_prms, prms)
			}
		} else {
			log.Printf("parseUrl: SKIP PARAM %v", arr_pr)
		}
	}
	prms := make([][]pair, 0)
	getAllParam([]pair{}, &prms, &list_prms, 0)

	for _, p := range prms {
		arr_uri = append(arr_uri, uri{href_path, p})
	}

	return arr_uri, nil
}

func getAllParam(pre_param []pair, res_list, list_par *[][]pair, ind int) {
	if ind >= len(*list_par) {
		*res_list = append(*res_list, pre_param)
		return
	}
	for _, p := range (*list_par)[ind] {
		new_param := append(pre_param, p)
		new_ind := ind + 1
		getAllParam(new_param, res_list, list_par, new_ind)
	}
}

func parseMaskParam(params []string) ([]pair, error) {
	if params[1][0] == '[' && params[1][len(params[1])-1] == ']' {
		s := strings.Split(strings.Trim(params[1], "[]"), ":")
		prms := make([]pair, 0)
		if len(s) == 2 {
			start, err := strconv.Atoi(s[0])
			if err != nil {
				return nil, fmt.Errorf("parsMaskParam: error AtoI mask %v", params)
			}
			max, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, fmt.Errorf("parsMaskParam: error AtoI mask %v", params)
			}
			for i := start; i < max; i++ {
				prms = append(prms, pair{params[0], strconv.Itoa(i)})
			}
			return prms, nil
		} else if len(s) > 2 {
			return nil, fmt.Errorf("parsMaskParam: error mask %v", params)
		}
		s = strings.Split(strings.Trim(params[1], "[]"), ";")
		if len(s) > 0 {
			for _, p_val := range s {
				prms = append(prms, pair{params[0], p_val})
			}
			return prms, nil
		} else {
			return nil, fmt.Errorf("parsMaskParam: error mask %v", params)
		}
	}
	return []pair{{params[0], params[1]}}, nil
}
