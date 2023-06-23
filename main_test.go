package mim

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStart_All(t *testing.T) {
	versions := []string{"4.4.8", "5.0.2"}
	testCtx := context.Background()

	for _, version := range versions {
		Convey("Given the version "+version, t, func() {

			Convey("When the Start method is called", func() {
				server, err := Start(testCtx, version)
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run as a standalone server with default options", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldNotBeEmpty)
						So(server.ReplicaSet(), ShouldBeEmpty)
						So(server.DBdir(), ShouldNotBeBlank)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, strconv.Itoa(server.Port()))
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, server.DBdir())
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "ephemeralForTest")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})
						})
					})
				})
			})

			Convey("When the StartWithReplicaSet method is called with an empty replica set option", func() {
				server, err := StartWithReplicaSet(testCtx, version, "")
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run as a replica set with the default name ('rs0')", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldNotBeEmpty)
						So(server.ReplicaSet(), ShouldEqual, "rs0")
						So(server.DBdir(), ShouldNotBeBlank)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, strconv.Itoa(server.Port()))
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, server.DBdir())
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "wiredTiger")
						So(server.cmd.Args[9], ShouldEqual, "--replSet")
						So(server.cmd.Args[10], ShouldEqual, "rs0")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()).SetReplicaSet(server.ReplicaSet()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})

						})
					})
				})
			})

			Convey("When the StartWithOptions method is called with no options", func() {
				server, err := StartWithOptions(testCtx, version)
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run as a standalone server with default options", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldNotBeEmpty)
						So(server.ReplicaSet(), ShouldBeEmpty)
						So(server.DBdir(), ShouldNotBeBlank)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, strconv.Itoa(server.Port()))
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, server.DBdir())
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "ephemeralForTest")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})
						})
					})
				})
			})

			Convey("When the StartWithOptions method is called with a replica set option", func() {
				server, err := StartWithOptions(testCtx, version, WithReplicaSet("my-replica-set"))
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run as a replica set with the given name", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldNotBeEmpty)
						So(server.ReplicaSet(), ShouldEqual, "my-replica-set")
						So(server.DBdir(), ShouldNotBeBlank)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, strconv.Itoa(server.Port()))
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, server.DBdir())
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "wiredTiger")
						So(server.cmd.Args[9], ShouldEqual, "--replSet")
						So(server.cmd.Args[10], ShouldEqual, "my-replica-set")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()).SetReplicaSet(server.ReplicaSet()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})

						})
					})
				})
			})

			Convey("When the StartWithOptions method is called with an empty replica set option", func() {
				server, err := StartWithOptions(testCtx, version, WithReplicaSet(""))
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run as a standalone server - as if no option had been provided", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldNotBeEmpty)
						So(server.ReplicaSet(), ShouldBeEmpty)
						So(server.DBdir(), ShouldNotBeBlank)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, strconv.Itoa(server.Port()))
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, server.DBdir())
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "ephemeralForTest")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()).SetReplicaSet(server.ReplicaSet()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})

						})
					})
				})
			})

			Convey("When the StartWithOptions method is called with the port and database directory options", func() {
				tempDir, err := os.MkdirTemp("", "")
				if err != nil {
					t.Fatalf("Error creating data directory: %v", err)
				}
				server, err := StartWithOptions(testCtx, version, WithPort(27017), WithDatabaseDir(tempDir))
				defer server.Stop(testCtx)

				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					Convey("And the mongod process has run with the supplied options", func() {
						So(server, ShouldNotBeNil)
						So(server.Port(), ShouldEqual, 27017)
						So(server.ReplicaSet(), ShouldBeEmpty)
						So(server.DBdir(), ShouldEqual, tempDir)
						So(server.cmd, ShouldNotBeNil)
						So(server.cmd.Args[0], ShouldEndWith, "mongod")
						So(server.cmd.Args[1], ShouldEqual, "--bind_ip")
						So(server.cmd.Args[2], ShouldEqual, "localhost")
						So(server.cmd.Args[3], ShouldEqual, "--port")
						So(server.cmd.Args[4], ShouldEqual, "27017")
						So(server.cmd.Args[5], ShouldEqual, "--dbpath")
						So(server.cmd.Args[6], ShouldEqual, tempDir)
						So(server.cmd.Args[7], ShouldEqual, "--storageEngine")
						So(server.cmd.Args[8], ShouldEqual, "ephemeralForTest")

						Convey("And the watcher process has run", func() {
							expectedScript := fmt.Sprintf("while ps -o pid= -p %d; "+
								"do sleep 1; "+
								"done; "+
								"kill -9 %d",
								os.Getpid(), server.cmd.Process.Pid)

							So(server, ShouldNotBeNil)
							So(server.watcherCmd, ShouldNotBeNil)
							So(server.watcherCmd.Args[0], ShouldEndWith, "/bin/sh")
							So(server.watcherCmd.Args[1], ShouldEqual, "-c")
							So(server.watcherCmd.Args[2], ShouldEqual, expectedScript)

							Convey("And the server accepts connections", func() {
								client, err := mongo.Connect(testCtx, options.Client().ApplyURI(server.URI()))
								So(err, ShouldBeNil)
								So(client, ShouldNotBeNil)
								So(client.Ping(testCtx, nil), ShouldBeNil)
							})
						})
					})
				})
			})
		})
	}
}

func TestGetFreeMongoPort(t *testing.T) {
	Convey("When getFreeMongoPort() is called n times, where n > 1", t, func() {
		n := 10

		Convey("A free, usable port should be returned every time", func() {
			for i := 0; i < n; i++ {
				port, err := getFreeMongoPort()
				So(err, ShouldBeNil)
				So(port, ShouldNotEqual, 0)

				l, e := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
				So(e, ShouldBeNil)
				So(l, ShouldNotBeNil)
				_ = l.Close()
			}
		})
	})
}
