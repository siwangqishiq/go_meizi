package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func FindPathByOs(pathName string) string{
	if runtime.GOOS == "linux" {
		return "/root/assets/" + pathName
	}
	return "assets/"+pathName
}

func FindPathUrl() string{
	if runtime.GOOS == "linux" {
		return HOST_REL
	}
	return HOST_TEST
}


func PrepareDirs() {
	fmt.Println("Prepare dirs.")
	os.Mkdir(FindPathByOs("data"), 0777)
	os.Mkdir(FindPathByOs("imgs"), 0777)
}

func DownloadFile(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// 发起 HTTP 请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 判断服务器返回状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// 将响应内容写入文件
	_, err = io.Copy(out, resp.Body)
	return err
}

func StringExt(str string) string {
	if len(str) <= 0 {
		return ""
	}
	
	lstIndex := strings.LastIndex(str,".")
	if(lstIndex < 0){
		return ""
	}
	return str[lstIndex: len(str)]
}

func Min(a,b int) int{
	if a < b {
		return a
	} else{
		return b
	}
}

func Max(a,b int) int{
	if a >= b {
		return a
	} else{
		return b
	}
}

