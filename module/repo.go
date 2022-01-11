package module

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

type repoListData struct {
	Data struct {
		Repos []struct {
			RepoName string `json:"repoName"`
		} `json:"repos"`
		Total int `json:"total"`
	} `json:"data"`
}

type repoTagData struct {
	Data struct {
		Tags []struct {
			Tag string `json:"tag"`
		} `json:"tags"`
	} `json:"data"`
}

func createAliClient(pathPattern, queryParams string) []byte {
	client, err := sdk.NewClientWithAccessKey("cn-shanghai", Conf.Accesskey, Conf.Accesstoken)

	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https" // https | http
	request.Domain = "cr.cn-shanghai.aliyuncs.com"
	request.Version = "2016-06-07"
	request.Headers["Content-Type"] = "application/json"

	request.PathPattern = pathPattern
	request.QueryParams["PageSize"] = queryParams

	body := `{}`
	request.Content = []byte(body)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	return response.GetHttpContentBytes()
}

func getRepoList(repoList map[string][]string) {
	var repoListData repoListData
	pathPattern := "/repos/hff-online"
	repoData := createAliClient(pathPattern, "")
	_ = json.Unmarshal(repoData, &repoListData)
	for i := 0; i < repoListData.Data.Total; i++ {
		repoList[repoListData.Data.Repos[i].RepoName] = nil
	}
}

func getRepoTag(repoName string) []string {
	var (
		repoTagData repoTagData
		tmpSlice    []string
	)
	pathPattern := fmt.Sprintf("/repos/hff-online/%s/tags", repoName)
	repoData := createAliClient(pathPattern, "4")
	_ = json.Unmarshal(repoData, &repoTagData)
	for i := 0; i < 4; i++ {
		if repoTagData.Data.Tags[i].Tag == "latest" {
			continue
		}

		tmpSlice = append(tmpSlice, repoTagData.Data.Tags[i].Tag)
	}

	return tmpSlice
}

func AssembleRepoInfo() map[string][]string {
	repoList := make(map[string][]string)
	getRepoList(repoList)
	m := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for k := range repoList {
		wg.Add(1)
		tmp := k
		go func() {
			defer wg.Done()
			rs := getRepoTag(tmp)
			m.Lock()
			repoList[tmp] = rs
			m.Unlock()
		}()
	}
	wg.Wait()
	return repoList
}
