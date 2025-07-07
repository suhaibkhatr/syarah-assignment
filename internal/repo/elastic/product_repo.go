package elastic

import (
	"context"
	"gift-store/internal/models"
	"gift-store/internal/sink"
	"strconv"
)

type ProductRepo struct {
	sink *sink.ElasticSink
}

func NewProductRepo(sink *sink.ElasticSink) *ProductRepo {
	return &ProductRepo{sink: sink}
}

func (e *ProductRepo) InsertOrUpdate(p models.Product, index string) error {
	_, err := e.sink.Es.Index().
		Index(index).
		Id(strconv.Itoa(p.ID)).
		BodyJson(p).
		Do(context.Background())
	return err
}

func (e *ProductRepo) Delete(id int, index string) error {
	_, err := e.sink.Es.Delete().
		Index(index).
		Id(strconv.Itoa(id)).
		Do(context.Background())
	return err
}
