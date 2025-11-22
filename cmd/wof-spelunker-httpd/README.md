# wof-spelunker-httpd

Start the Spelunker web application.

```
$> ./bin/wof-spelunker-httpd -h
Start the Spelunker web application.
Usage:
	./bin/wof-spelunker-httpd [options]
Valid options are:
  -authenticator-uri string
    	A valid aaronland/go-http/v3/auth.Authenticator URI. This is future-facing work and can be ignored for now. (default "null://")
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -protomaps-max-data-zoom int
    	The maximum zoom (tile) level for data in a PMTiles database
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
  -root-url string
    	The root URL for all public-facing URLs and links. If empty then the value of the -server-uri flag will be used.
  -server-uri string
    	A valid `aaronland/go-http/v3/server.Server URI. (default "http://localhost:8080")
  -spelunker-uri string
    	A URI in the form of '{SPELUNKER_SCHEME}://{IMPLEMENTATION_DETAILS}' referencing the underlying Spelunker database. For example: sql://sqlite3?dsn=spelunker.db (default "null://")
```

## Things the Spelunker web application does NOT do

* It does not provide the ability to edit records. Most of the pieces to do that are available at this point but they have not been wired in to the Spelunker at this point.
* It does not define any methods for querying spatial data. This _might_ be possible in a version "3" of the Spelunker but it's too soon to say righr now.
* It does not provide any kind of authentication or authorization mechanism. Yet.

## Building

The `wof-spelunker-httpd` depends on Go language build tags. The default `cli` Makefile target to compile command line tools build the `wof-spelunker-httpd` tool with support for all the database implementations included in this package. For example:

```
$> cd spelunker
$> make cli
go build -mod vendor -tags="sqlite3,icu,json1,fts5,opensearch" -ldflags="-s -w" -o bin/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go
```

If you only want to build the `wof-spelunker-httpd` tool with support for SQLite-backed database you can run the `cli-sqlite` Makefile target:

```
$> make cli-sqlite
go build -mod vendor -tags="sqlite3,icu,json1,fts5" -ldflags="-s -w" -o bin/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go
```

_Note that the default SQLite-backed implementation depends on being able to compile the [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) package._

If you only want to build the `wof-spelunker-httpd` tool with support for an OpenSearch-backed database you can run the `cli-opensearch` Makefile target:

```
$> make cli-opensearch
go build -mod vendor -tags="opensearch" -ldflags="-s -w" -o bin/wof-spelunker-httpd cmd/wof-spelunker-httpd/main.go
```

### Build tags

| Target | Tags | Notes |
| --- | --- | --- |
| MySQL | `mysql` | Support for MySQL should probably still be considered "alpha" at best. |
| Postgres | `postgres` | Support for Postgres should probably still be considered "alpha" at best. |
| SQLite | `sqlite3,icu,json1,fts5` | |
| OpenSearch | `opensearch` | |

## Indexing

For details on indexing records in to the Spelunker web application please consult the documentation for the [wof-spelunker-index](../wof-spelunker-index) tool.

## Maps

Map configuration for the `wof-spelunker-httpd` application is controlled by these flags:

```
  -map-provider string
    	Valid options are: leaflet, protomaps (default "leaflet")
  -map-tile-uri string
    	A valid Leaflet tile layer URI. See documentation for special-case (interpolated tile) URIs. (default "https://tile.openstreetmap.org/{z}/{x}/{y}.png")
  -protomaps-max-data-zoom int
    	The maximum zoom (tile) level for data in a PMTiles database (default 15)
  -protomaps-theme string
    	A valid Protomaps theme label. (default "white")
```

_Under the hood the `wof-spelunker-httpd` application is using the [aaronland/go-http-maps](https://github.com/aaronland/go-http-maps) package to manage map configuration._

### Raster tiles

The default map for the `wof-spelunker-httpd` application uses the [LeafletJS](https://leafletjs.com/) package to display raster tiles provided by the [OpenStreetMap](https://) project.

You can use alternate raster tiles by specifying their "ZXY" tile URL in the `-map-file-uri` flag.

### Protomaps tiles

You can use also use map data encoded in [PMTiles](https://docs.protomaps.com/pmtiles/) database for rendering base maps by setting the `-map-provider` flag to "protomaps".

The value of the `-map-tile-uri` should be one of the following:

* `api://{PROTOMAPS_API_KEY}` to load PMTiles data from the [Protomaps API](https://protomaps.com/api).
* `file:///path/to/local/pmtiles.db` to load PMTiles data from a local PMTile database file.

## Examples

### database/sql

New `database/sql`-backed Spelunker instances are created by passing a URI to the `NewSpelunker` method in the form of:

```
sql://{DATABASE_ENGINE}?dsn={DATABASE_ENGINE_DSN}
```

Where `{DATABASE_ENGINE}` is a registered (as in "imported") Go language [database/sql](https://pkg.go.dev/database/sql) driver name and `{DATABASE_ENGINE_DSN}` is a driver-specific DSN string for connecting to that database.

See [sql/README.md](../../sql/README.md) for details.

#### SQLite

For example to start the `wof-spelunker-httpd` application using data stored in a local SQLite database:

```
./bin/wof-spelunker-httpd \
	-spelunker-uri 'sql://sqlite3?dsn=/usr/local/data/sfom.db'
```

### OpenSearch

New OpenSearch-backed Spelunker instances are created by passing a URI to the `NewSpelunker` method in the form of:

```
opensearch://?{QUERY_PARAMETERS}
```

Where {QUERY_PARAMETERS} may be one or more of the following:
* `client-uri={STRING}. A URI in the form of "opensearch://?client-uri={GO_WHOSONFIRST_DATABASE_OPENSEARCH_CLIENT_URI}" for connecting to OpenSearch.
* `reader-uri={STRING}. A valid "whosonfirst/go-reader/v2.Reader" URI used to read raw "source" Who's On First documents (because documents are indexed in a truncated form in OpenSearch).
* `cache-uri={STRING}. A valid "whosonfirst/go-cache.Cache" URI used to cache data retrieved from a "reader-uri" source.

For example:

```
./bin/wof-spelunker-httpd \
	-spelunker-uri 'opensearch://?client-uri=https%3A%2F%2Flocalhost%3A9200%2Fspelunker%3Fusername%3Dadmin%26password%3Ddkjfhsjdkfkjdjhksfhskd98475kjHkzjxckj%26insecure%3Dtrue%26require-tls%3Dtrue&cache-uri=ristretto%3A%2F%2F&reader-uri=https%3A%2F%2Fdata.whosonfirst.org'
```

_The need to URL-escape the `client-uri`, `-reader-uri` and `-cache-uri` parameters is not great but that's how things work today._

#### client-uri

The value of the `client-uri` query parameter is a URL-escaped URI for instantiating a [opensearchapi.Client](https://pkg.go.dev/github.com/opensearch-project/opensearch-go/v4/opensearchapi#Client) instance using the [whosonfirst/go-whosonfirst-database/opensearch/client](https://github.com/whosonfirst/go-whosonfirst-database/tree/main/opensearch/client) package. The `?client-uri` paramater is expected to take the form of:

```
opensearch://{OPENSEARCH_HOST}:{OPENSEARCH_PORT}/{OPENSEARCH_INDEX}?{QUERY_PARAMETERS}
```

Where {QUERY_PARAMETERS} may be one or more of the following:
* `debug={BOOLEAN}`. A boolean value to configure the underlying OpenSearch client to write request and response bodies to STDOUT.
* `insecure={BOOLEAN}`. A boolean value to disable TLS "InsecureSkipVerify" checks (for custom certificate authorities and the like).
* `require-tls={BOOLEAN}`. A boolean value to ensure that all connections are made over HTTPS even if the OpenSearch port is not 443.
* `username={STRING}`. The OpenSearch username for authenticated connections.
* `password={STRING}`. The OpenSearch password for authenticated connections.
* `aws-credentials-uri={STRING}`. A a valid `aaronland/go-aws-auth` URI used to create a Golang AWS authentication config used to sign requests to an AWS-hosted OpenSearch instance.

See [opensearch/README.md](../../opensearch/README.md) for details.

#### reader-uri

The OpenSearch Spelunker implementation does not index Who's On First records in their entirety. Specifically geometries are excluded from the OpenSearch index with the goal of minimizing the overall size of the index. While that decision may change in future releases it is the way things work today.

In order to account for the need to retrieve and display complete Who's On First GeoJSON Feature records the OpenSearch Spelunker implementation relies on the [whosonfirst/go-reader/v2.Reader](https://github.com/whosonfirst/go-reader) interface to do so. The interface is very simple:

```
// Reader is an interface for reading data from multiple sources or targets.
type Reader interface {
	// Reader returns a `io.ReadSeekCloser` instance for a URI resolved by the instance implementing the `Reader` interface.
	Read(context.Context, string) (io.ReadSeekCloser, error)
	// Exists returns a boolean value indicating whether a URI already exists.
	Exists(context.Context, string) (bool, error)
	// The absolute path for the file is determined by the instance implementing the `Reader` interface.
	ReaderURI(context.Context, string) string
}
```

And there are number of implementations allowing you to read records from files stored locally on disk, to files stored in GitHub, to files hosted from the `data.whosonfirst.org` website. For example:

To read files from a Who's On First data repository on disk:

```
?reader-uri=repo:///path/to-whosonfirst-data-reposiroty
```

To read files from a Who's On First data repository hosted on the `data.whosonfirst.org` website:

```
?reader-uri=https://data.whosonfirst.org/geojson
```

To read files from a Who's On First data repository hosted on GitHub:

```
?reader-uri=URLESCAPE(github://whosonfirst-data/whosonfirst-data-admin-us?branch=master&prefix=data)
```

_See [whosonfirst/go-reader-github](https://github.com/whosonfirst/go-reader-github) for details._

To read files from a Who's On First data repository hosted on GitHub deriving the name of the repository dynamically using a "finding aid":

```
?reader-uri=URLESCAPE(findingaid://https/data.whosonfirst.org/findingaid?template=https://raw.githubusercontent.com/whosonfirst-data/{repo}/master/data/)
```

_See [whosonfirst/go-reader-findingaid](https://github.com/whosonfirst/go-reader-findingaid) for details._

#### cache-uri

In order to minimize the number of requests to an external "reader" source the OpenSearch Spelunker implementation caches data using an implementation of the [whosonfirst-/go-cache](https://github.com/whosonfirst/go-cache) interface. While there are a number of implementations if you just want a reliable in-memory the easiest thing is to use the [whosonfirst/go-cache-ristretto](https://github.com/whosonfirst/go-cache-ristretto) implementation which is bundled with this package. For example:

```
?cache-uri=ristretto://
```

If you want to disable caching entirely, use the "null" cache implementation. For example:

```
?cache-uri=null://
```

## Endpoints

### Endpoints for humans

#### /

The root URL for the Spelunker. For example `http://localhost:8080/`.

#### /about

The URL for the Spelunker "about" page. For example `http://localhost:8080/about/`.

#### /concordances

The URL for the page listing all the concordances in a Spelunker index. For example `http://localhost:8080/concordances/`.

#### /concordances/{namespace}

The URL for the page listing all the records matching a specific concordance namespace in a Spelunker index. For example `http://localhost:8080/concordances/gp`.

#### /concordances/{namespace}:{predicate}

The URL for the page listing all the records matching a specific concordance namespace and predicate in a Spelunker index. For example `http://localhost:8080/concordances/foo:bar`.

#### /concordances/{namespace}:{predicate}={value}

The URL for the page listing all the records matching a specific concordance in a Spelunker index. For example `http://localhost:8080/concordances/foo:bar=baz`.

#### /id/{id}

The URL for the page to display a specific Who's On First record. For example `http://localhost:8080/id/1234`.

#### /id/{id}/descendants

The URL for the page to display all of descendants of a specific Who's On First record. For example `http://localhost:8080/id/1234/descendants`.`

#### nullisland

The URL for the page to display all of the records in a Spelunker index that are "visiting" Null Island (have lat,lon coordinates of "0.0,0.0"). For example `http://localhost:8080/nullisland`.`

#### /placetypes

The URL for the page listing all the placetypes in a Spelunker index. For example `http://localhost:8080/placetypes/`.

#### /placetypes/{placetype}

The URL for the page listing all the records matching a specific placetype in a Spelunker index. For example `http://localhost:8080/placetypes/locality/`.

#### /recent/{duration}

The URL for the page listing all the records in a Spelunker index that have been updated within a given time period. For example `http://localhost:8080/recent/P14D`.

#### /search

The URL for the page for search records in a Spelunker index. For example `http://localhost:8080/search?q=Montreal`.

### Endpoints for machines

#### /concordances/{namespace}/facets

#### /concordances/{namespace}:{predicate}/facets

#### /concordances/{namespace}:{predicate}={value}/facets

#### /findingaid

#### /id/{id}/descendants/facets

#### /id/{id}/geojson

#### /id/{id}/geojsonld

#### /id/{id}/navplace

#### /id/{id}/select

#### /id/{id}/spr

#### /id/{id}/svg

#### /id/{id}/wkt

#### /nullisland/facets

#### /opensearch

#### /placetypes/{placetype}/facets

#### /recent/{duration}/facets

#### /search/facets
