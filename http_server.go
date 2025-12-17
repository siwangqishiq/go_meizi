package main

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const URL_INDEX = "/index"
const URL_ASSETS = "/assets/"
const URL_ABLUMS = "/ablums"
const URL_WEB = "/web/"

const HTTP_CODE_SUCCESS = 200
const HTTP_CODE_ERROR = 404

const DEFAULT_PAGESIZE = 20

type HttpResp struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data any `json:"data"`
}

type HttpRespAblums struct {
	Title  string   `json:"title"`
	Href   string   `json:"href"`
	Id     int      `json:"id"`
	Cover  string   `json:"cover"`
	Images []string `json:"images"`
}

func NewSuccessHttpResp(data any) HttpResp{
	return HttpResp{
		Code: HTTP_CODE_SUCCESS,
		Msg:"",
		Data: data,
	}
}

func NewErrorHttpResp(errCode int , msg string) HttpResp{
	return HttpResp{
		Code: errCode,
		Msg:msg,
		Data: nil,
	}
}

type HttpServer struct {
	crawler *Crawler
	port    int
}

func NewHttpServer(craw *Crawler) *HttpServer {
	return &HttpServer{
		crawler: craw,
		port:    8910,
	}
}

func (h *HttpServer) StartServer() {
	mux := http.NewServeMux()

	mux.HandleFunc(URL_INDEX, func(w http.ResponseWriter,req *http.Request){
		HandlerHome(w, req, h)
	})

	mux.HandleFunc(URL_ASSETS,func(w http.ResponseWriter,req *http.Request){
		HandleAsset(w, req)
	})

	mux.HandleFunc(URL_WEB,func(w http.ResponseWriter,req *http.Request){
		PutDirToCloud(w, req, URL_WEB, true)
	})

	mux.HandleFunc(URL_ABLUMS, func(w http.ResponseWriter,req *http.Request){
		HandleAblums(w, req, h)
	})

	var portStr string = strconv.Itoa(h.port)
	fmt.Println("Server started at http://localhost:" + portStr)
    http.ListenAndServe(":" + portStr, mux)
}

func HandlerHome(w http.ResponseWriter, r *http.Request, httpServer *HttpServer) {
	fmt.Fprintln(w, "Home Page ablumsszie", len(httpServer.crawler.Ablums))
}

func HandleAblums(w http.ResponseWriter, req *http.Request, httpServer *HttpServer){
	queryValues := req.URL.Query()
	// fmt.Println(queryValues)
	pageSize := DEFAULT_PAGESIZE
	pageSizeValue,ok := queryValues["pagesize"]
	if(ok){
		pageSize , _= strconv.Atoi(pageSizeValue[0])
	}
	fmt.Println("pagesize",pageSize)
	var startIndex int = 0
	idStrs,ok := queryValues["id"]
	if(ok){
		id , err := strconv.Atoi(idStrs[0])
		if(err == nil){
			fmt.Println("id",id)
			for i := 0; i < len(httpServer.crawler.Ablums); i++ {
				// fmt.Println("query id",httpServer.crawler.Ablums[i].Id)
				if(httpServer.crawler.Ablums[i].Id == id){
					startIndex = i + 1
					break
				}
			}//end for i
		}
	}

	startIndex = Min(startIndex, len(httpServer.crawler.Ablums) - 1)
	fmt.Println("find start index",startIndex)
	queryList := httpServer.crawler.Ablums[startIndex:startIndex + pageSize]
	list := ReplaceImageResource(queryList,httpServer)
	data := NewSuccessHttpResp(list)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}

func ReplaceImageResource(list []Ablum, h *HttpServer)[]HttpRespAblums{
	ret := make([]HttpRespAblums, len(list))
	
	for i, abl := range list {
		ret[i].Id = abl.Id
		ret[i].Title = abl.Title
		ret[i].Href = abl.Href

		ret[i].Cover = GenAssetFullPath(abl.Cover, h.port)
		ret[i].Images = make([]string, len(list[i].Images))
		for j, img := range list[i].Images {
			ret[i].Images[j] = GenAssetFullPath(img, h.port)
		} //end for j
	}//end for i
	return ret
}

func GenAssetFullPath(input string, port int) string{
	return fmt.Sprintf("%s:%d/assets/%s",FindPathUrl(),port,input)
}

func PutDirToCloud(resp http.ResponseWriter, req *http.Request, dir string,isLocalPath bool){
	fmt.Println("handle PutDirToCloud")
	fmt.Println("req path->", req.URL.Path)
	// fmt.Println("req query->", r.URL.Query())

	path := req.URL.Path
	if len(path) <= 0 {
		http.Error(resp, "file not found", http.StatusNotFound)
		return
	}

	subpath, success := strings.CutPrefix(path, dir)
	if !success {
		http.Error(resp, "file not found", http.StatusNotFound)
		return
	}
	fmt.Println("req subpath->", subpath)
	
	var filepath string = ""
	if(isLocalPath){
		filepath = dir[1:] + subpath
	}else{
		filepath = FindPathByOs("") + subpath
	}
	fmt.Println("file path =", filepath)
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Println("file not found->", filepath)
		http.Error(resp, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		http.Error(resp, "cannot stat file", http.StatusInternalServerError)
		return
	}
	fmt.Println("file IsDir", fi.IsDir(),"size",fi.Size())

	// 获取 MIME 类型
	ext := StringExt(filepath)
	fmt.Println("file ext " ,ext)
	mimeType := mime.TypeByExtension(ext)
	fmt.Println("file mimeType " ,mimeType)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	// 设置响应头
	resp.Header().Set("Content-Type", mimeType)
	resp.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	// resp.Header().Set("Content-Disposition", "attachment; filename=\""+fi.Name()+"\"")

	http.ServeContent(resp, req, fi.Name(), fi.ModTime(), f)
}

func HandleAsset(resp http.ResponseWriter, req *http.Request){
	PutDirToCloud(resp, req, URL_ASSETS, false)
}


