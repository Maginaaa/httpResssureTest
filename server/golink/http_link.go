package golink

import (
	"bytes"
	"httpPressureTest/model"
	"httpPressureTest/server/client"
	"io/ioutil"
	"net/http"
	"sync"
)

type ReqListMany struct {
	list []*model.Request
}

func (r *ReqListMany) getCount() int {
	return len(r.list)
}

var (
	clientList *ReqListMany
)

func init() {
	clientList = &ReqListMany{}
}

func Http(chanId uint64, ch chan<- *model.RequestResults, totalNumber uint64, wg *sync.WaitGroup, request *model.Request) {

	defer func() {
		wg.Done()
	}()

	// fmt.Printf("启动协程 编号:%05d \n", chanId)
	for i := uint64(0); i < totalNumber; i++ {

		list := getRequestList(request)

		isSucceed, errCode, requestTime, contentLength := sendList(list)

		requestResults := &model.RequestResults{
			Time:          requestTime,
			IsSucceed:     isSucceed,
			ErrCode:       errCode,
			ReceivedBytes: contentLength,
		}

		requestResults.SetId(chanId, i)

		ch <- requestResults
	}

	return
}

func getRequestList(request *model.Request) []*model.Request {

	if len(clientList.list) <= 0 {

		return []*model.Request{request}
	}

	return clientList.list
}

// 多个接口分步压测
func sendList(requestList []*model.Request) (isSucceed bool, errCode int, requestTime uint64, contentLength int64) {

	errCode = model.HttpOk
	for _, request := range requestList {
		succeed, code, u, length := send(request)
		isSucceed = succeed
		errCode = code
		requestTime = requestTime + u
		contentLength = contentLength + length
		if succeed == false {

			break
		}
	}

	return
}

var buf = make([]byte, 1024*1024)

// send 发送一次请求
func send(request *model.Request) (bool, int, uint64, int64) {
	var (
		isSucceed     = false
		errCode       = model.HttpOk
		contentLength = int64(0)
	)

	newRequest := request

	resp, requestTime, err := client.HttpRequest(newRequest.Method, newRequest.Url, newRequest.GetBody(), newRequest.Headers, newRequest.Timeout)
	if err != nil {
		errCode = model.RequestErr // 请求错误
	} else {

		contentLength = 0
		for {
			n, err := resp.Body.Read(buf)
			resp.Body = ioutil.NopCloser(bytes.NewReader(buf))
			contentLength += int64(n)
			if err != nil {
				break
			}
		}


		// 验证请求是否成功
		if resp.StatusCode == http.StatusOK {
			isSucceed = true
		} else {
			errCode = resp.StatusCode
		}

	}
	return isSucceed, errCode, requestTime, contentLength
}


