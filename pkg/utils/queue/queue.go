package queue

import (
	"ar_exhibition/pkg/utils"
	"log"
	"net/http"
	"sync"
	"time"
)

type MyRequest struct {
	req     *http.Request
	timeout time.Time
	//out     chan *http.Response
}

type Queue struct {
	sync.Mutex
	requests []*MyRequest
}

func NewQueue() *Queue {
	return &Queue{
		requests: make([]*MyRequest, 0),
	}
}

func (q *Queue) Push(r *http.Request, t time.Time) {
	q.Lock()
	q.requests = append(q.requests, &MyRequest{r, t})
	q.Unlock()
}

func (q *Queue) Pop() (*http.Request, time.Time) {
	if !q.Empty() {
		q.Lock()
		defer q.Unlock()
		imin := 0
		for i := range q.requests {
			if q.requests[i].timeout.Before(q.requests[imin].timeout) {
				imin = i
			}
		}
		result := q.requests[imin]
		q.requests = append(q.requests[:imin], q.requests[imin+1:]...)
		return result.req, result.timeout
	}
	return nil, time.Time{}
}

func (q *Queue) Empty() bool {
	return len(q.requests) == 0
}

func (q *Queue) CheckStatus(status int, service string, repeat ...*MyRequest) bool {
	if status >= http.StatusInternalServerError {
		var req *http.Request
		if len(repeat) > 0 {
			req = repeat[0].req
		} else {
			req, _ = http.NewRequest(http.MethodGet, service+utils.ApiPing, nil)
		}
		q.Push(req, time.Now().Add(utils.RequestLimit*time.Second))
		return false
	}
	return true
}

func Execute(q *Queue) {
	for {
		if q.Empty() {
			time.Sleep(time.Second)
		} else {
			req, timeout := q.Pop()
			service := req.URL.Scheme + "://" + req.URL.Host
			time.Sleep(time.Until(timeout))
			if resp, err := http.DefaultClient.Do(req); err == nil {
				q.CheckStatus(resp.StatusCode, service, &MyRequest{req, timeout})
			} else {
				log.Println(err)
				q.Push(req, time.Now().Add(utils.RequestLimit*time.Second))
			}
		}
	}
}
