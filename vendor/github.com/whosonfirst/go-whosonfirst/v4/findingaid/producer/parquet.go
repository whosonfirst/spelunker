package producer

import (
	"context"
	"fmt"
	"net/url"

	"github.com/sfomuseum/go-parquet"
	"github.com/whosonfirst/go-whosonfirst/v4/feature/alt"
	"github.com/whosonfirst/go-whosonfirst/v4/feature/properties"
	"github.com/whosonfirst/go-whosonfirst/v4/iterate"
)

type ParquetRecord struct {
	Id   int64  `parquet:"id"`
	Repo string `parquet:"repo"`
}

type ParquetProducer struct {
	Producer
	writer *parquet.ParquetWriter[*ParquetRecord]
}

func init() {
	ctx := context.Background()
	RegisterProducer(ctx, "parquet", NewParquetProducer)
}

func NewParquetProducer(ctx context.Context, uri string) (Producer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	wr, err := parquet.NewWriter[*ParquetRecord](ctx, u.Path)

	if err != nil {
		return nil, fmt.Errorf("Failed to create parquet writer, %w", err)
	}

	p := &ParquetProducer{
		writer: wr,
	}

	return p, nil
}

func (p *ParquetProducer) PopulateWithIterator(ctx context.Context, iterator_uri string, iterator_sources ...string) error {

	iter, err := iterate.NewIterator(ctx, iterator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create iterator '%s', %w", iterator_uri, err)
	}

	for rec, err := range iter.Iterate(ctx, iterator_sources...) {

		if err != nil {
			return err
		}

		body, err := rec.ReadAllAndClose()

		if err != nil {
			return fmt.Errorf("Failed to read %s, %w", rec.Path, err)
		}

		if alt.IsAlt(body) {
			continue
		}

		id, err := properties.Id(body)

		if err != nil {
			return fmt.Errorf("Failed to derive ID, %w", err)
		}

		repo, err := properties.Repo(body)

		if err != nil {
			return fmt.Errorf("Failed to derive repo, %w", err)
		}

		fa_rec := &ParquetRecord{
			Id:   id,
			Repo: repo,
		}

		p.writer.WriteRow(fa_rec)
	}

	return nil
}

func (p *ParquetProducer) Close(ctx context.Context) error {
	return p.writer.Close()
}
