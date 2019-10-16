package elasticClient

import (
	conf "../Config"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"
)

// Index and Mapping in ElasticSearch:
/*

private_im_user:

{
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1
		},
	"mappings": {
		"properties": {
			"id": {
				"type": "long",
				"index":false
			},
			"name": {
				 "type": "text",
				 "analyzer": "ik_max_word",
				 "search_analyzer": "ik_max_word"
				},
			"email": {"type": "keyword"},
			"avatar": {
				"type":"text",
				"index":false
				},
			"gender": {
				"type":"short",
				"index":false
			},
			"is_delete": {
				"type":"boolean",
			}
		}
	}
}
*/

/*

private_im_group_chat:

{
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1
		},
	"mappings": {
		"properties": {
			"id": {
				"type": "long",
				"index":false
			},
			"name": {
				 "type": "text",
				 "analyzer": "ik_max_word",
				 "search_analyzer": "ik_max_word"
			},
			"avatar": {
				"type":"text",
				"index":false
			},
			"is_delete": {
				"type":"boolean"
			},
			"manager_name": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_max_word"
			},
			"manager_avatar: {
				"type":"text",
				"index":false
			}
		}
	}
}
*/

/*

private_im_subscription:

{
	"settings": {
		"number_of_shards": 3,
		"number_of_replicas": 1
		},
	"mappings": {
		"properties": {
			"id": {
				"type": "long",
				"index":false
			},
			"name": {
				 "type": "text",
				 "analyzer": "ik_max_word",
				 "search_analyzer": "ik_max_word"
			},
			"intro": {
				 "type": "text",
				 "analyzer": "ik_max_word",
				 "search_analyzer": "ik_max_word"
			},
			"avatar": {
				"type":"text",
				"index":false
			},
			"is_delete": {
				"type":"boolean"
			},
			"manager_name": {
				"type": "text",
				"analyzer": "ik_max_word",
				"search_analyzer": "ik_max_word"
			},
			"manager_avatar": {
				"type":"text",
				"index":false
			}
		}
	}
}

*/
const (
	HTTPControlPrefix     = "http://"
	UserIndexName         = "private_im_user"
	GroupChatIndexName    = "private_im_group_chat"
	SubscriptionIndexName = "private_im_subscription"
)

func makeAndSendRequest(method, url string, body io.Reader) ([]byte, error) {
	request, _ := http.NewRequest(method, url, body)
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if nil != err {
		return nil, err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return nil, errors.New(fmt.Sprintf("fail:\n%s", data))
	}
	_ = request.Body.Close()
	_ = resp.Body.Close()
	return data, nil
}

func ChangeReplicaNumber(indexName string, replica int) error {
	url := HTTPControlPrefix + conf.ElasticsearchServerAddress + "/" + indexName + "/_settings"
	jsonBodyStr := fmt.Sprintf(`{"number_of_replicas": %d}"`, replica)
	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("PUT", url, body)
	return err
}

func UserIndexDocSave(id int64, name, email, avatar string, gender int32) error {
	tempDataFormatStr := `{
	"id": %d,
	"name": "%s",
	"email":"%s",
	"avatar":"%s",
	"gender": %d,
	"is_delete": false
}`

	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, id, name, email, avatar, gender)
	url := fmt.Sprintf("%s%s/%s/_doc/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, UserIndexName, id)
	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("PUT", url, body)
	return err

}

// Can search by name or email
func UserIndexDocSearch(target string, page, perPage int) ([]byte, error) {
	tempDataFormatStr := `{	
	"from": %d,
	"size": %d,
	"query": { 
		"bool": {
    		"should": [
        		{ "match": { "name": "%s" } },
				{ "bool" : { "must": { "term": { "email": "%s" } } } }
			],
    		"filter": [ 
        		{ "term":  { "is_delete": false } }
    		]
    	}
	}
}
`
	fromValue := (page - 1) * perPage
	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, fromValue, perPage, target, target)
	url := fmt.Sprintf("%s%s/%s/_search", HTTPControlPrefix, conf.ElasticsearchServerAddress, UserIndexName)
	body := bytes.NewBufferString(jsonBodyStr)
	data, err := makeAndSendRequest("GET", url, body)
	if nil != err {
		return nil, err
	}
	return data, nil

}

// Update the name, email, avatar or gender in document of private_im_user index.
func UserIndexDocUpdate(id int64, field string, value interface{}) error {
	var jsonBodyStr string
	switch field {
	case "name", "email", "avatar":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%s"}}`, field, value.(string))
	case "gender":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%d"}}`, field, value.(int))
	case "is_delete":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%t"}}`, field, value.(bool))
	}
	url := fmt.Sprintf("%s%s/%s/_update/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, UserIndexName, id)

	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("POST", url, body)
	return err
}

// -----------------------------------------------------------------------------

func GroupChatIndexDocSave(id int64, name, avatar, managerName, managerAvatar string) error {
	tempDataFormatStr := `{
	"id": %d,
	"name": "%s",
	"avatar":"%s",
	"is_delete": false,
	"manager_name": "%s",
	"manager_avatar": "%s"
}`
	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, id, name, avatar, managerName, managerAvatar)
	url := fmt.Sprintf("%s%s/%s/_doc/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, GroupChatIndexName, id)
	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("PUT", url, body)
	return err
}

func GroupChatIndexDocSearch(target string, page, perPage int) ([]byte, error) {
	tempDataFormatStr := `
{	
	"from": %d,
	"size": %d,
	"query": { 
		"bool": { 
    		"should": [
        		{ "match": { "name": "%s" } },
        		{ "match": { "manager_name": "%s" } }
        	],
    		"filter": [ 
        		{ "term":  { "is_delete": false } }
    		]
    	}
	}
}`
	fromValue := (page - 1) * perPage
	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, fromValue, perPage, target, target)
	url := fmt.Sprintf("%s%s/%s/_search", HTTPControlPrefix, conf.ElasticsearchServerAddress, GroupChatIndexName)
	body := bytes.NewBufferString(jsonBodyStr)
	data, err := makeAndSendRequest("GET", url, body)
	if nil != err {
		return nil, err
	}
	return data, nil
}

// update the name, avatar, manager_name, manager_avatar, is_delete in document of private_im_group_chat
func GroupChatIndexDocUpdate(id int64, field string, value interface{}) error {
	var jsonBodyStr string
	switch field {
	case "name", "avatar", "manager_name", "manager_avatar":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%s"}}`, field, value.(string))
	case "is_delete":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%t"}}`, field, value.(bool))
	}
	url := fmt.Sprintf("%s%s/%s/_update/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, GroupChatIndexName, id)

	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("POST", url, body)
	return err
}

// -----------------------------------------------------------------------------

func SubscriptionIndexSave(id int64, name, intro, avatar, managerName, managerAvatar string) error {
	tempDataFormatStr := `{
	"id": %d,
	"name": "%s",
	"intro": "%s", 
	"avatar":"%s",
	"is_delete": false,
	"manager_name": "%s",
	"manager_avatar": "%s"
}`
	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, id, name, intro, avatar, managerName, managerAvatar)
	url := fmt.Sprintf("%s%s/%s/_doc/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, SubscriptionIndexName, id)
	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("PUT", url, body)
	return err
}

func SubscriptionIndexSearch(target string, page, perPage int) ([]byte, error) {
	tempDataFormatStr := `
{	
	"from": %d,
	"size": %d,
	"query": { 
		"bool": { 
    		"should": [
        		{ "match": { "name": "%s" } },
				{ "match": { "intro": "%s" } },
        		{ "match": { "manager_name": "%s" } }
        	],
    		"filter": [ 
        		{ "term":  { "is_delete": false } }
    		]
    	}
	}
}`
	fromValue := (page - 1) * perPage
	jsonBodyStr := fmt.Sprintf(tempDataFormatStr, fromValue, perPage, target, target, target)
	url := fmt.Sprintf("%s%s/%s/_search", HTTPControlPrefix, conf.ElasticsearchServerAddress, SubscriptionIndexName)
	body := bytes.NewBufferString(jsonBodyStr)
	data, err := makeAndSendRequest("GET", url, body)
	if nil != err {
		return nil, err
	}
	return data, nil
}

func SubscriptionIndexUpdate(id int64, field string, value interface{}) error {
	var jsonBodyStr string
	switch field {
	case "name", "avatar", "intro", "manager_name", "manager_avatar":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%s"}}`, field, value.(string))
	case "is_delete":
		jsonBodyStr = fmt.Sprintf(`{"doc": {"%s": "%t"}}`, field, value.(bool))
	}
	url := fmt.Sprintf("%s%s/%s/_update/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, SubscriptionIndexName, id)

	body := bytes.NewBufferString(jsonBodyStr)
	_, err := makeAndSendRequest("POST", url, body)
	return err
}

// -----------------------------------------------------------------------------

func IndexDocDelete(indexName string, docId int64) error {
	url := fmt.Sprintf("%s%s/%s/_doc/%d", HTTPControlPrefix, conf.ElasticsearchServerAddress, indexName, docId)
	body := bytes.NewBufferString("")
	_, err := makeAndSendRequest("DELETE", url, body)
	return err
}
