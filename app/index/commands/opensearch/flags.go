package opensearch

import (
	"flag"
	"fmt"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

var client_uri string
var iterator_uri string
var forgiving bool
var verbose bool

var create_index bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("index")

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)

	fs.BoolVar(&forgiving, "forgiving", true, "Be \"forgiving\" of failed writes, logging the issue(s) but not triggering errors")
	fs.BoolVar(&create_index, "create-index", false, "Create a new OpenSearch index before indexing records.")

	fs.StringVar(&client_uri, "client-uri", "", "A valid whosonfirst/go-whosonfirst-database/opensearch/client URI in the form of \"opensearch://{OPENSEARCH_HOST}:{OPENSEARCH_PORT}/{OPENSEARCH_INDEX}?{QUERY_PARAMETERS}\".")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")
	return fs
}
