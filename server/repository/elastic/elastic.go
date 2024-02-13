package elastic

import (
	"github.com/04Akaps/kafka-go/config"
	"github.com/olivere/elastic/v7"
	"log"
)

type Elastic struct {
	config *config.Config

	Client *elastic.Client
}

type ElasticImpl interface {
}

func NewElastic(cfg *config.Config) (ElasticImpl, error) {
	e := &Elastic{config: cfg}

	var err error

	if e.Client, err = elastic.NewClient(
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
