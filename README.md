# dp-mongodb-in-memory

Library that runs an in-memory MongoDB instance for Go unit tests.
It's based on [memongo](https://github.com/benweissmann/memongo).

## How it works

- It detects your operating system and platform to determine the download URL for the right MongoDB binary.

- It will download MongoDB and store it in a [cache location](#cache-location). Any following execution will use the copy from the cache. Therefore internet connection is only required the first time a particular MongoDB version is used.

- It will start a process running the downloaded `mongod` binary.

- It exposes a number of `Start` endpoints to start the server either in standalone mode or replica set mode. It uses the `ephemeralForTest` storage engine in standalone mode, and the `wiredTiger` storage engine in replica set mode. A temporary directory and port may be supplied, or will be generated if not supplied.

- A `Server` object is returned form the above endpoints from which the server URI, port, database directory, and replica set name (if applicable) may be retrieved

- Additionally, a _watcher_ process will start in background ensuring that the mongod process is killed when the current process exits. This guarantees that no process is left behind even if the tests exit uncleanly or you don't call `Stop()`.

### Supported versions

The following Unix systems are supported:

- MacOS
- Ubuntu 16.04 or greater
- Debian 9.2 or greater
- Amazon Linux 2

The supported MongoDB versions are 4.4 and above.

### Cache location

The downloaded mongodb binary will be stored in a local cache: a folder named `dp-mongodb-in-memory` living on the machine base cache directory. That is `$XDG_CACHE_HOME` if such environment variable is set or `~/.cache` (Linux) and `~/Library/Caches` (MacOS) if not.

## Installation

To install this package run:

```bash
go get github.com/ONSdigital/dp-mongodb-in-memory
```

## Usage

Call:
    `Start(ctx, version)`, `StartWithReplicaSet(ctx, version, replicaSetName)`, or `StartWithOptions(ctx, version, ...options)` where version is the MongoDB version you want to use. You can then use `URI()` to connect a client to it.
Call `Stop()` when you are done with the server.

```go
package example

import (
    "context"
    "testing"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    mim "github.com/ONSdigital/dp-mongodb-in-memory"
)

func TestExample(t *testing.T) {
    testCtx := context.Background()

    server, err := mim.Start(testCtx, "5.0.2")
    // OR
    server, err = mim.StartWithReplicaSet(testCtx, "5.0.2", "my-replica-set")
    // OR
    server, err = mim.StartWithOptions(testCtx, "5.0.2", mim.WithPort(27017), mim.WithDatabaseDir("/var/tmp/my-temp-dir"))

    if err != nil {
        // Deal with error
    }
    defer server.Stop(testCtx)

    client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
    if err != nil {
        // Deal with error
    }

    //Use client as needed
    err = client.Ping(testCtx, nil)
    if err != nil {
        // Deal with error
    }
}

```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Released under MIT license, see [LICENSE](LICENSE.md) for details.
