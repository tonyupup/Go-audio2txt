package lib

import (
	"time"
	"os"
	"io/ioutil"
	"net/http"
	"github.com/lestrrat-go/libxml2"
	"runtime"
	"path/filepath"
	"fmt"
	"sync"
)

type Signale interface {
	Getinstace()
}
type Manage struct {
}

var m *Manage
var once sync.Once

func GetInstance() *Manage {
	once.Do(func() {
		m = &Manage{}
	})
	return m
}
func (m *Manage) Do() {
	fmt.Println("ok")
}

const bingURL = "https://cn.bing.com"

func getPic() ([]byte, error) {
	body, err := http.Get(bingURL)
	if err != nil {
		return nil, err
	}
	defer body.Body.Close()
	html, err := libxml2.ParseHTMLReader(body.Body)
	if err != nil {
		return nil, err
	}
	defer html.Free()
	result, err := html.Find(`//div[@id="bgImgProgLoad"]/@data-ultra-definition-src`)
	if err != nil {
		return nil, err
	}
	pic, err := http.Get(bingURL + result.String())
	if err != nil {
		return nil, err
	}
	defer pic.Body.Close()

	return ioutil.ReadAll(pic.Body)
	// return nil,nil
}

func savePic(path string, pic []byte) (finepath string, err error) {
	var prepath string
	if runtime.GOOS == "linux" {
		prepath = os.Getenv("HOME")
	} else {
		prepath = os.Getenv("USERPROFILE")
	}
	finepath1 := filepath.Join(prepath, path, time.Now().Format("20060102")+".jpg")
	file, err := os.OpenFile(finepath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(pic)
	if err != nil {
		return "", err
	}
	return finepath1, nil
}