package main

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStart(t *testing.T) {
	versions := []string{"4.4.8", "5.0.2"}

	for _, version := range versions {
		Convey("Given the version "+version, t, func() {
			Convey("When the Start method is called", func() {
				server, err := Start(version)
				Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
					So(server, ShouldNotBeNil)
				})
				Convey("And the server accepts connections", func() {
					client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(server.URI()))
					So(err, ShouldBeNil)
					So(client, ShouldNotBeNil)
					So(client.Ping(context.Background(), nil), ShouldBeNil)
				})
			})
		})
	}
}
