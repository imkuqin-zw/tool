package throttle

import (
	"testing"
	"net/http"
	"fmt"
	"log"
	"github.com/imkuqin-zw/tool/cache"
)

func TestThrottle_Handle(t *testing.T) {
	cache.Register("127.0.0.1:6379", nil)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		filter := NewThrottle(request, writer, NewRedisCatch("127.0.0.1:6379"), 20, 1)
		pass, err := filter.Handle()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
			return
		}
		if !pass {
			return
		}
		fmt.Fprintf(writer, "%v", "hello world")
	})
	log.Fatal(http.ListenAndServe(":8092", nil))
}
