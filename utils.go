package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/globalsign/mgo/bson"
	"github.com/revel/revel"
)

//AsOptions convert [a] to [a,a]
func AsOptions(list []string) [][]interface{} {
	output := [][]interface{}{}
	for _, item := range list {
		output = append(output, []interface{}{item, item})
	}
	return output
}

// Split by ",", ";"
func Split(s string) []string {
	out := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case ',', ';':
			return true
		}
		return false
	})
	return out
}

func SliceToMap(colNames []string, rows [][]string) ([]map[string]string, error) {
	if len(rows) == 0 {
		return nil, nil
	}
	colNum := len(colNames)
	rowNum := len(rows)

	// create maps
	results := []map[string]string{}
	for i := 0; i < rowNum; i++ {
		if len(rows[i]) != colNum {
			return nil, errors.New("Error: ColNames length and record length not consistent.")
		}
		t_map := map[string]string{}
		for j := 0; j < colNum; j++ {
			t_map[colNames[j]] = rows[i][j]
		}
		results = append(results, t_map)
	}

	return results, nil
}

func ReadUrl(apiUrl, user, pwd string, proxyUrl string) ([]byte, error) {
	var client = &http.Client{}
	if proxyUrl != "" {
		proxy, _ := url.Parse(proxyUrl)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
		}
	}

	//read body string
	req, err := http.NewRequest("GET", apiUrl, nil)
	req.SetBasicAuth(user, pwd)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	// Retrieve the body of the response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, err
	}
	return []byte(body), nil
}

//StrInSlice check if string is in string list
func StrInSlice(str string, list []string) bool {
	for _, item := range list {
		if str == item {
			return true
		}
	}
	return false
}

// Remove string from a string list
func RemoveStrFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func SplitByPipe(s string) []string {
	out := strings.FieldsFunc(s, func(r rune) bool {
		switch r {
		case '|':
			return true
		}
		return false
	})
	return out
}

func RenderTpl(t_path string, viewArgs map[string]interface{}) string {
	// Get the Template.
	t, err := revel.MainTemplateLoader.Template(t_path)
	if err != nil {
		revel.AppLog.Errorf("Load template err %s", err)
	}

	var outbuf bytes.Buffer
	out := bufio.NewWriter(&outbuf)
	t.Render(out, viewArgs)
	out.Flush()

	return outbuf.String()
}

//UniqueInts returns a unique subset of the int slice provided.
func UniqueInts(input []int) []int {
	u := make([]int, 0, len(input))
	m := make(map[int]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

//UniqueStrings returns a unique subset of the int slice provided.
func UniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

//RemoveBlankStrings
func RemoveBlankStrings(input []string) []string {
	u := make([]string, 0, len(input))
	for _, val := range input {
		if strings.TrimSpace(val) != "" {
			u = append(u, val)
		}
	}

	return u
}

func FilterUnicodeSymbol(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case 0x000A, 0x000B, 0x000C, 0x000D, 0x0085, 0x2028, 0x2029, 0xFFFD:
			return -1
		default:
			return r
		}
	}, s)
}

// parse url query string return to bson.M
func ParseQueryToBson(reqStr string) bson.M {
	args, _ := url.ParseQuery(reqStr)
	out := bson.M{}
	for k, item := range args {
		out[k] = fmt.Sprintf("%v", item[0])
	}
	return out
}
