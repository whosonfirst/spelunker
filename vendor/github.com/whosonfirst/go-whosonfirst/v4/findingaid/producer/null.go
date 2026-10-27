package producer

import (
	"context"
)

type NullProducer struct {
	Producer
}

func init() {
	ctx := context.Background()
	RegisterProducer(ctx, "null", NewNullProducer)
}

func NewNullProducer(ctx context.Context, uri string) (Producer, error) {
	p := &NullProducer{}
	return p, nil
}

func (p *NullProducer) PopulateWithIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) error {
	return nil
}

func (p *NullProducer) Close(ctx context.Context) error {
	return nil
}
