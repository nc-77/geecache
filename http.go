package geecache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/nc-77/geecache/geecachepb"
	"github.com/nc-77/geecache/hash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/geecache"
	defaultReplicas = 50
)

type HTTPPool struct {
	self        string
	basePath    string // eg /geecache
	mu          sync.Mutex
	peers       *hash.Map
	httpGetters map[string]*httpGetter
}

type httpGetter struct {
	baseUrl string // eg http://example.com/geecache/
}

var (
	_ PeerGetter = (*httpGetter)(nil)
	_ PeerPicker = (*HTTPPool)(nil)
)

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:        self,
		basePath:    defaultBasePath,
		httpGetters: make(map[string]*httpGetter),
	}
}

func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, path)
	// /<basepath>/<groupname>/<key> required
	parts := strings.Split(path[len(p.basePath):], "/")

	if len(parts) != 3 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName, key := parts[1], parts[2]
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

	body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = hash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseUrl: peer + p.basePath}
	}
}

func (p *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

func (h *httpGetter) Get(req *pb.Request, resp *pb.Response) error {
	u := fmt.Sprintf("%v/%v/%v", h.baseUrl, url.QueryEscape(req.Group), url.QueryEscape(req.Key))

	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %v", res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body : %v", err)
	}

	if err = proto.Unmarshal(data, resp); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}
