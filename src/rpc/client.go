package rpc

import (
	"bytes"
	"encoding/json"
	"github.com/dubuqingfeng/explorer-parser/producer/config"
	"github.com/gabstv/httpdigest"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type RpcClient struct {
	address    string
	user       string
	password   string
	ssl        bool
	authType   string
	httpClient *http.Client
}

type rpcResponse struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *error          `json:"error"`
}

type rpcRequest struct {
	// ID: time.Now().UnixNano()
	ID      int64       `json:"id"`
	JSONRPC string      `json:"jsonrpc"` // value: "2.0"
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type rpcResult struct {
	response *http.Response
	err      error
}

func newHttpClient(nodeConfig config.NodeConfig) *http.Client {
	if nodeConfig.AuthType == "digest" {
		// if auth type is 'digest-2617' initial the client
		transport := httpdigest.New(nodeConfig.User, nodeConfig.Password)
		httpClient := &http.Client{Transport: transport}
		return httpClient
	} else {
		httpClient := &http.Client{}
		return httpClient
	}
}

// Need to force open authentication
func newRPCClient(nodeConfig config.NodeConfig) *RpcClient {
	httpClient := newHttpClient(nodeConfig)
	client := &RpcClient{
		address:    nodeConfig.Address,
		user:       nodeConfig.User,
		password:   nodeConfig.Password,
		authType:   nodeConfig.AuthType,
		httpClient: httpClient,
	}
	return client
}

// need add timeout limit
func (this *RpcClient) call(method string, params interface{}) (response rpcResponse, err error) {
	// build http request
	timer := time.NewTimer(2 * time.Second)

	request := rpcRequest{time.Now().UnixNano(), "2.0", method, params}
	payload, err := json.Marshal(request)

	if err != nil {
		return rpcResponse{}, err
	}

	req, err := http.NewRequest("POST", this.address, bytes.NewBuffer(payload))
	if err != nil {
		return rpcResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Add("Accept", "application/json")
	if this.authType == "base" && (this.user != "" || this.password != "") {
		req.SetBasicAuth(this.user, this.password)
	}
	// Timer
	resp, err := this.DoRequest(timer, req)
	if err != nil {
		log.Error(err)
		return rpcResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rpcResponse{}, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return rpcResponse{}, err
	}

	if response.Error != nil {
		return rpcResponse{}, errors.New("rpc call error")
	}

	return response, nil
}

// Do Request and timeout limit
func (this *RpcClient) DoRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	done := make(chan rpcResult, 1)
	go func() {
		resp, err := this.httpClient.Do(req)
		done <- rpcResult{resp, err}
	}()
	select {
	case r := <-done:
		return r.response, r.err
	case <-timer.C:
		return nil, errors.New(req.URL.RawQuery + "Timeout")
	}
}

// get sync status, keep sync
type SyncStatus struct {
}
