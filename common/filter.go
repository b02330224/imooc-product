package common

import (
	"net/http"
	"strings"
)

type FilterHandle func(rw http.ResponseWriter, req *http.Request) error

type Filter struct {
	//用来存储需要拦截的URI
	filterMap map[string]FilterHandle
}

func NewFilter() *Filter {
	return &Filter{filterMap: map[string]FilterHandle{}}
	//return &Filter{filterMap:make(map[string]FilterHandle)} //都可以
}


func (f *Filter) RegisterFilterUri(uri string, handler FilterHandle) {
	f.filterMap[uri] = handler
}

func (f *Filter) GetFilterHandle(uri string) FilterHandle {
	return f.filterMap[uri]
}

type WebHandle func (rw http.ResponseWriter, req *http.Request)

func (f *Filter) Handle(webHandle WebHandle) func(rw http.ResponseWriter, r *http.Request) {

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		for path, handle := range f.filterMap {
			if strings.Contains(r.RequestURI, path) {
				err := handle(rw, r)
				if err != nil {
					rw.Write([]byte(err.Error()))
					return
				}

				break
			}
		}
		webHandle(rw, r)
	}

}