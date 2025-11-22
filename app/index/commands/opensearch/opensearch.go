package opensearch

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	_ "github.com/whosonfirst/go-whosonfirst-database/opensearch/writer"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/whosonfirst/go-whosonfirst-database/opensearch/client"
	"github.com/whosonfirst/go-whosonfirst-database/opensearch/schema/v2"
	"github.com/whosonfirst/go-whosonfirst-iterwriter/v4"
	iterwriter_app "github.com/whosonfirst/go-whosonfirst-iterwriter/v4/app/iterwriter"
	"github.com/whosonfirst/go-writer/v3"
	"github.com/whosonfirst/spelunker/v2/app/index/commands"
)

type IndexOpenSearchCommand struct {
	commands.Command
}

func init() {
	ctx := context.Background()
	commands.RegisterCommand(ctx, "opensearch", NewIndexOpenSearchCommand)
}

func NewIndexOpenSearchCommand(ctx context.Context, cmd string) (commands.Command, error) {
	c := &IndexOpenSearchCommand{}
	return c, nil
}

func (c *IndexOpenSearchCommand) Run(ctx context.Context, args []string) error {

	fs := DefaultFlagSet()
	fs.Parse(args)

	sources := fs.Args()

	if verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose (debug) logging enabled")
	}

	writer_uri := client_uri

	wr, err := writer.NewWriter(ctx, writer_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new writer, %w", err)
	}

	if create_index {

		u, _ := url.Parse(writer_uri)
		os_index := strings.TrimLeft(u.Path, "/")

		slog.Debug("Create index", "name", os_index)

		mappings_r, err := v2.FS.Open("mappings.spelunker.json")

		if err != nil {
			return fmt.Errorf("Failed to open mappings for reading, %w", err)
		}

		defer mappings_r.Close()

		settings_r, err := v2.FS.Open("settings.spelunker.json")

		if err != nil {
			return fmt.Errorf("Failed to open settings for reading, %w", err)
		}

		defer settings_r.Close()

		os_client, err := client.NewClient(ctx, client_uri)

		if err != nil {
			return fmt.Errorf("Failed to create Opensearch client, %w", err)
		}

		mappings_req := opensearchapi.IndicesCreateReq{
			Index: os_index,
			Body:  mappings_r,
		}

		_, err = os_client.Indices.Create(ctx, mappings_req)

		if err != nil {
			return fmt.Errorf("Failed to create index, %w", err)
		}

		settings_req := opensearchapi.SettingsPutReq{
			Indices: []string{
				os_index,
			},
			Body: settings_r,
		}

		_, err = os_client.Indices.Settings.Put(ctx, settings_req)

		if err != nil {
			return fmt.Errorf("Failed to put settings, %w", err)
		}

	}

	cb_func := iterwriter.DefaultIterwriterCallback(forgiving)

	opts := &iterwriter_app.RunOptions{
		CallbackFunc:  cb_func,
		Writer:        wr,
		IteratorURI:   iterator_uri,
		IteratorPaths: sources,
		Verbose:       verbose,
	}

	/*

				/usr/local/data/whosonfirst/whosonfirst-data-admin-ca
		2025/11/15 17:49:02 INFO Iterator stats elapsed=1m0.001272166s seen=24419 allocated="195 MB" "total allocated"="11 GB" sys="346 MB" numgc=159
		2025/11/15 17:49:30 ERROR Failed to index record path=112/576/680/5/1125766805.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"
		2025/11/15 17:49:30 ERROR Failed to index record path=112/607/179/7/1126071797.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"
		2025/11/15 17:49:30 ERROR Failed to index record path=112/611/017/5/1126110175.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"
		2025/11/15 17:49:30 ERROR Failed to index record path=112/611/366/1/1126113661.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"
		2025/11/15 17:49:30 ERROR Failed to index record path=115/886/315/9/1158863159.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"
		2025/11/15 17:49:30 ERROR Failed to index record path=115/886/830/9/1158868309.geojson type=mapper_exception reason="timed out while waiting for a dynamic mapping update"

	*/

	return iterwriter_app.RunWithOptions(ctx, opts)
}
