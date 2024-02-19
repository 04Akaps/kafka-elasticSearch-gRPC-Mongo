package elastic

import (
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
	elasticLog chan interface{}

	client *elastic.Client
}

type ElasticImpl interface {
}

func NewElastic(cfg *config.Config, elasticLog chan interface{}) (ElasticImpl, error) {
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

		log.Println("Success To Connect Elastic Search")
		return e, nil
	}
}

func (e *Elastic) CreateData(index string, data interface{}) {
	if d, ok := data.(types.KafkaEvent); ok {
		client := e.client
		if err := checkIndexExisted(client, index); err != nil {
			log.Println("Failed To Create Index Try Again")
			e.CreateData(index, data)
		} else if _, err = client.Index().Index(index).Id(d.ElasticId).BodyJson(d.Data).Do(util.Context()); err != nil {
			log.Println("Failed To Create New Elastic Data")
		} else {
			log.Println("Success To Create New Elastic Data")
		}
	} else {

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
