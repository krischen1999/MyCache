package GeeCache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*基于http提供被其他节点访问的能力*/
const defaltBasePath = "/_GeeCache/"

type HTTPPool struct {
	self     string
	basePath string
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaltBasePath,
	}
}

// log
func (p *HTTPPool) log(format string, v ...interface{}) {
	log.Printf("[Sever %s] %s", p.self, fmt.Sprintf(format, v...))
}

// handel all http request
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path :" + r.URL.Path)
	}
	p.log("%s %s", r.Method, r.URL.Path)

	//约定访问路径 /<basepath>/<groupname>/<key>
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group"+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content0Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
