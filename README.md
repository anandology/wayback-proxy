# wayback-proxy

HTTP proxy server to proxy content from Internet Archive's wayback machine.


## Usage

Build the proxy server:

    $ go build

Start the proxy server:

    $ export WAYMACK_TIMESTAMP=2015
    $ ./wayback-proxy
    Starting wayback-proxy with timestamp: 2015
    http://localhost:8080/

The timestamp can be specified in even finer details as using YYYYMMDDHHMMSS.

To use the proxy server from command-line:

    $ export http_proxy=http://localhost:8080/
    $ curl http://example.com/
