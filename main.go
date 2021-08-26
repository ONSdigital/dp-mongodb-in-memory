package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strconv"

	"github.com/ONSdigital/dp-mongodb-in-memory/download"
	"github.com/ONSdigital/log.go/v2/log"
)

// Server represents a running MongoDB server.
type Server struct {
	cmd *exec.Cmd
	// watcherCmd *exec.Cmd
	dbDir string
	port  int
}

func init() {
	log.Namespace = "dp-mongodb-in-memory"
}

// Start runs a MongoDB server at a given version using a random free port
// and returns the Server.
func Start(version string) (*Server, error) {
	port, err := getFreePort()
	if err != nil {
		log.Fatal(context.Background(), "Could not find a free port", err)
		return nil, err
	}

	binPath, err := getOrDownloadBinPath(version)
	if err != nil {
		log.Fatal(context.Background(), "Could not find mongodb", err)
		return nil, err
	}

	// Create a db dir. Even the ephemeralForTest engine needs a dbpath.
	dbDir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(context.Background(), "Error creating data directory", err)
		return nil, err
	}

	log.Info(context.Background(), "Starting mongod server", log.Data{"binPath": binPath, "dbDir": dbDir, "port": port})
	cmd := exec.Command(binPath, "--storageEngine", "ephemeralForTest", "--dbpath", dbDir, "--port", strconv.Itoa(port))

	// Run the server
	err = cmd.Start()
	if err != nil {
		remErr := os.RemoveAll(dbDir)
		if remErr != nil {
			log.Error(context.Background(), "Error removing data directory", remErr, log.Data{"dir": dbDir})
		}
		log.Fatal(context.Background(), "Could not start mongodb", err)
		return nil, err
	}

	return &Server{
		cmd: cmd,
		// watcherCmd: watcherCmd,
		dbDir: dbDir,
		port:  port,
	}, nil
}

// Stop kills the mongo server.
func (s *Server) Stop() {
	err := s.cmd.Process.Kill()
	if err != nil {
		log.Error(context.Background(), "Error stopping mongod process", err)
		return
	}

	// err = s.watcherCmd.Process.Kill()
	// if err != nil {
	// 	s.logger.Warnf("error stopping watcher process: %s", err)
	// 	return
	// }

	err = os.RemoveAll(s.dbDir)
	if err != nil {
		log.Error(context.Background(), "Error removing data directory", err, log.Data{"dir": s.dbDir})
		return
	}
}

// Port returns the port the server is listening on.
func (s *Server) Port() int {
	return s.port
}

// URI returns a mongodb:// URI to connect to
func (s *Server) URI() string {
	return fmt.Sprintf("mongodb://localhost:%d", s.port)
}

func getOrDownloadBinPath(version string) (string, error) {
	config, err := download.NewConfig(version)
	if err != nil {
		log.Error(context.Background(), "Failed to create config", err)
		return "", err
	}

	binPath, err := download.GetMongoDB(*config)
	if err != nil {
		return "", err
	}
	return binPath, nil
}

func getFreePort() (int, error) {
	// Based on: https://github.com/phayes/freeport/blob/master/freeport.go
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
