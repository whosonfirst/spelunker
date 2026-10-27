package resolver

import (
	"context"
	"fmt"
	"net/url"

	"github.com/parquet-go/parquet-go"
	sfom_parquet "github.com/sfomuseum/go-parquet"
	"github.com/whosonfirst/go-whosonfirst/v4/findingaid/producer"
)

// type ParquetResolver implements the `Resolver` interface for data stored in a Parquet database
// produced by `github.com/whosonfirst/go-whosonfirst/v4/findingaid/producer.ParquetProducer`
type ParquetResolver struct {
	Resolver
	reader       sfom_parquet.ReadCloserAt
	parquet_file *parquet.File
}

func init() {
	ctx := context.Background()
	RegisterResolver(ctx, "parquet", NewParquetResolver)
}

// NewParquetResolver will return a new `Resolver` instance for resolving repository names
// and IDs stored in a Parquet database.
func NewParquetResolver(ctx context.Context, uri string) (Resolver, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	reader, sz, err := sfom_parquet.OpenURI(u.Path)

	pf, err := parquet.OpenFile(reader, sz)

	if err != nil {
		return nil, err
	}

	f := &ParquetResolver{
		reader:       reader,
		parquet_file: pf,
	}

	return f, nil
}

// GetRepo returns the name of the repository associated with this ID in a Who's On First finding aid.
func (r *ParquetResolver) GetRepo(ctx context.Context, id int64) (string, error) {

	rec, err := sfom_parquet.FindRecordByIdWithParquetFile[*producer.ParquetRecord, int64](r.reader, r.parquet_file, id)

	if err != nil {
		return "", err
	}

	return rec.Repo, nil
}

func (r ParquetResolver) Close() error {

	return r.reader.Close()
}
