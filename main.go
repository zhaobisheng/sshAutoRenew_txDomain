// txDomain project main.go
package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"

	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type DomainRecord struct {
	Id    int64  `json:"id"`
	Value string `json:"value"`
	Name  string `json:"name"`
	Type  string `json:"type"`
}

type Records struct {
	Records []*DomainRecord `json:"records"`
}

type DomainList struct {
	Data *Records `json:"data"`
}

const apiUrl = "cns.api.qcloud.com/v2/index.php?"

func main() {
	var action, subDomain, recordType string
	flag.StringVar(&action, "a", "add", "-a=add/clean 添加记录或者删除记录,默认是add")
	flag.Parse()
	if (len(os.Args) < 8 && strings.EqualFold(action, "add")) || (len(os.Args) < 5 && strings.EqualFold(action, "clean")) {
		//fmt.Println("Usage:\r\n./txDomain AccessKeyID AccessKeySecret funwan.cn blog A 127.0.0.1")
		flag.Usage()
		fmt.Println("eg.\r\nAddRecord:./txDomain (-a=add)添加时可以省略-a AccessKeyID AccessKeySecret funwan.cn blog A 127.0.0.1")
		fmt.Println("\r\nCleanRecord:./txDomain -a=Clean AccessKeyID AccessKeySecret funwan.cn www A")
		os.Exit(0)
	}
	rand.Seed(time.Now().UnixNano())
	AccessKeyID := os.Args[2]
	AccessKeySecret := os.Args[3]
	domain := os.Args[4]
	if strings.EqualFold(action, "add") {
		subDomain = os.Args[5]
		recordType = strings.ToUpper(os.Args[6])
		value := os.Args[7]
		rs := TXRecordCreate(AccessKeyID, AccessKeySecret, domain, subDomain, recordType, value)
		fmt.Println("\r\n" + rs)
	} else if strings.EqualFold(action, "clean") {
		if len(os.Args) == 7 {
			subDomain = os.Args[5]
			recordType = strings.ToUpper(os.Args[6])
		}
		if subDomain == "" {
			subDomain = "_acme-challenge"
		}
		if recordType == "" {
			recordType = "TXT"
		}
		rs := TXGetRecordList(AccessKeyID, AccessKeySecret, domain, subDomain, recordType)
		fmt.Println("\r\n" + rs)
		domainList := &DomainList{}
		err := json.Unmarshal([]byte(rs), domainList)
		if err != nil {
			fmt.Println("json-error:", err)
			return
		}
		for _, val := range domainList.Data.Records {
			fmt.Println("\r\n", val.Id, val.Name, val.Type, val.Value)
			result := TXRecordDel(AccessKeyID, AccessKeySecret, domain, val.Id)
			fmt.Println("\r\n" + result)
		}
	}
}

func TXGetRecordList(AccessKeyID, AccessKeySecret, domain, subDomain, recordType string) string {
	data := map[string]string{
		"SecretId":   AccessKeyID,
		"Action":     "RecordList",
		"domain":     domain,
		"recordType": recordType,
		"subDomain":  subDomain,
	}
	data = MakePublicParam(data)
	rs := SendRequest(data, AccessKeySecret)
	return rs
}

func TXRecordDel(AccessKeyID, AccessKeySecret, domain string, recordId int64) string {
	data := map[string]string{
		"SecretId": AccessKeyID,
		"Action":   "RecordDelete",
		"domain":   domain,
		"recordId": fmt.Sprintf("%d", recordId),
	}
	data = MakePublicParam(data)
	rs := SendRequest(data, AccessKeySecret)
	return rs
}

func TXRecordCreate(AccessKeyID, AccessKeySecret, domain, subDomain, recordType, value string) string {
	data := map[string]string{
		"SecretId":   AccessKeyID,
		"Action":     "RecordCreate",
		"domain":     domain,
		"subDomain":  subDomain,
		"recordType": recordType,
		"recordLine": "默认",
		"value":      value,
	}
	if recordType == "MX" {
		data["mx"] = "10"
	}
	data = MakePublicParam(data)
	rs := SendRequest(data, AccessKeySecret)
	return rs
}

func MakePublicParam(params map[string]string) map[string]string {
	if params["SignatureMethod"] == "" {
		params["SignatureMethod"] = "HmacSHA1"
	}
	if params["Timestamp"] == "" {
		params["Timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	}
	if params["Nonce"] == "" {
		rand.Seed(time.Now().UnixNano())
		params["Nonce"] = fmt.Sprintf("%d", rand.Int31())
	}
	if params["Region"] == "" {
		params["Region"] = "ap-guangzhou"
	}
	return params
}

func SendRequest(data map[string]string, AccessKeySecret string) string {
	sortQueryString := SortedString(data)
	stringToSign := "GET" + apiUrl + sortQueryString[1:] //UrlEncode(sortQueryString[1:])
	data["Signature"] = Sign(stringToSign, AccessKeySecret)
	sortQueryString = SortedString(data)
	sortQueryString = sortQueryString[1:]
	url := fmt.Sprintf("https://%s%s", apiUrl, sortQueryString)
	r, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(b)
}

func SortedString(data map[string]string) string {
	var sortQueryString string
	for _, v := range Keys(data) {
		sortQueryString = fmt.Sprintf("%s&%s=%s", sortQueryString, v, data[v])
		//fmt.Sprintf("%s&%s=%s", sortQueryString, v, UrlEncode(data[v]))
	}
	return sortQueryString
}

func UrlEncode(in string) string {
	r := strings.NewReplacer("+", "%20", "*", "%2A", "%7E", "~")
	return r.Replace(url.QueryEscape(in))
}

func Sign(stringToSign, AccessKeySecret string) string {
	h := hmac.New(sha1.New, []byte(AccessKeySecret))
	h.Write([]byte(stringToSign))
	return UrlEncode(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func Keys(data map[string]string) []string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
