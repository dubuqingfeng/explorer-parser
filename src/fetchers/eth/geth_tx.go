package eth

import (
	"github.com/dubuqingfeng/explorer-parser/src/fetchers"
	"github.com/dubuqingfeng/explorer-parser/src/producer/config"
)

type GethTxFetcher struct {
	NodeConfigs []config.NodeConfig
	fetchers.Fetcher
}

func (fetcher GethTxFetcher) Fetch(title string) (bool, map[string]string) {
	strings := make(map[string]string)
	// async rpc client call
	return true, strings
}
