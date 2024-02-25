package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/04Akaps/kafka-go/config"
	"github.com/04Akaps/kafka-go/server/types"
	"github.com/04Akaps/kafka-go/server/util"
	"github.com/olivere/elastic/v7"
	"log"
)

type Elastic struct {
	config     *config.Config
	elasticLog <-chan interface{}

	client *elastic.Client
}

type ElasticImpl interface {
	SearchData(index string, query elastic.Query, sort types.Sort, paging types.Paging) ([]interface{}, error)
}

func NewElastic(cfg *config.Config, elasticLog <-chan interface{}) (ElasticImpl, error) {
	e := &Elastic{config: cfg, elasticLog: elasticLog}

	var err error

	if e.client, err = elastic.NewClient(
		elastic.SetBasicAuth(
			cfg.Elastic.User,
			cfg.Elastic.Password,
		),
		elastic.SetURL(cfg.Elastic.URI),
		elastic.SetSniff(false),
	); err != nil {
		return nil, err
	} else {
		go e.subscribeLog()

		log.Println("Success To Connect Elastic Search")
		return e, nil
	}
}

func (e *Elastic) SearchData(index string, query elastic.Query, sort types.Sort, paging types.Paging) ([]interface{}, error) {
	client := e.client

	baseQuery := client.Search(index).Query(query)

	if sort.Field != "" {
		// default asc
		baseQuery.Sort(sort.Field, true)
	}

	if paging.PageSize != 0 {
		baseQuery.From(int(paging.Page))
		baseQuery.Size(int(paging.PageSize))
	}

	baseQuery.Pretty(true)

	if err := checkIndexExisted(client, index); err != nil {
		return nil, err
	} else if result, err := baseQuery.Do(util.Context()); err != nil {
		return nil, err
	} else {
		var res []interface{}

		for _, hit := range result.Hits.Hits {
			var model interface{}

			if err := json.Unmarshal(hit.Source, &model); err != nil {
				return nil, err
			} else {
				res = append(res, model)
			}
		}

		return res, nil
	}
}

func (e *Elastic) subscribeLog() {
	for {
		select {
		case event := <-e.elasticLog:

			var decode types.KafkaEvent

			if b, err := json.Marshal(event); err != nil {
				log.Println("Failed To Marshal", err.Error())
			} else if err = json.Unmarshal(b, &decode); err != nil {
				log.Println("Failed To UnMarshal", err.Error())
			} else {
				if err = checkIndexExisted(e.client, decode.Index); err == nil {
					go e.createData(decode.Index, decode, 0)
				}
			}

		}
	}
}

func (e *Elastic) createData(index string, data interface{}, retryCount int64) {
	if retryCount > 3 {
		log.Println("Reached maximum retry limit, giving up.", index, data)
		return
	} else if d, ok := data.(types.KafkaEvent); ok {
		client := e.client
		if err := checkIndexExisted(client, index); err != nil {
			log.Println("Failed To Create Index Try Again")
			e.createData(index, data, retryCount+1)
		} else if _, err = client.Index().Index(index).Id(d.ElasticId).BodyJson(d.Data).Do(util.Context()); err != nil {
			log.Println("Failed To Create New Elastic Data")
		} else {
			log.Println("Success To Create New Elastic Data")
		}
	}
}

func checkIndexExisted(client *elastic.Client, index string) (err error) {
	ctx := util.Context()
	indices := []string{index}

	existService := elastic.NewIndicesExistsService(client)
	existService.Index(indices)

	exist, err := existService.Do(ctx)

	if err != nil {
		message := fmt.Sprintf("NewIndicesExistsService.Do() %s", err.Error())
		return errors.New(message)
	} else if !exist {
		fmt.Println("nOh no! The index", index, "doesn't exist.")
		fmt.Println("Create the index, and then run the Go script again")
		if _, err = client.CreateIndex(index).Do(ctx); err != nil {
			return err
		} else {
			return nil
		}
	} else if exist {
		return nil
	} else {
		return nil
	}
}
