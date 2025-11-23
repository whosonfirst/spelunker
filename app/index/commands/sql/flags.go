package sql

import (
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-iterate/v3"
)

var iterator_uri string

var db_uri string

var optimize bool

var strict_alt_files bool
var index_alt multi.MultiString

var index_relations bool
var relations_uri string

var procs int
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("index")

	valid_schemes := strings.Join(iterate.IteratorSchemes(), ",")
	iterator_desc := fmt.Sprintf("A valid whosonfirst/go-whosonfirst-iterate/v3.Iterator URI. Supported iterator URI schemes are: %s", valid_schemes)

	fs.StringVar(&iterator_uri, "iterator-uri", "repo://", iterator_desc)

	fs.StringVar(&db_uri, "database-uri", "", "A URI in the form of 'sql://{DATABASE_SQL_ENGINE}?dsn={DATABASE_SQL_DSN}'. For example: sql://sqlite3?dsn=test.db")

	fs.BoolVar(&optimize, "optimize", true, "Attempt to optimize the database before closing connection")

	fs.BoolVar(&strict_alt_files, "strict-alt-files", true, "Be strict when indexing alt geometries")

	fs.IntVar(&procs, "processes", (runtime.NumCPU() * 2), "The number of concurrent processes to index data with")

	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging")
	return fs
}
