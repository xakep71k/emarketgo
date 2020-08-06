package db

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		fmt.Printf("new client: %v\n", err)
	}

	ctx := DefaultContext()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Printf("connect: %v\n", err)
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Printf("mongo unreachable: %v\n", err)
		return nil, err
	}

	return client, nil
}

func DefaultContext() context.Context {
	ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx
}

type DockerMongo struct {
	dataPath string
}

func NewDockerMongo(dataPath string) *DockerMongo {
	return &DockerMongo{
		dataPath: dataPath,
	}
}

func (d *DockerMongo) Start() error {
	d.Stop()
	return runProcess(
		30*time.Second,
		"docker",
		"run",
		"--rm",
		"--name",
		"emarket-mongo",
		"--network",
		"host",
		"-v",
		fmt.Sprintf("%v:/data/db", d.dataPath),
		"-d",
		"mongo",
	)
}

func (d *DockerMongo) Stop() {
	_ = runProcess(
		30*time.Second,
		"docker",
		"stop",
		"emarket-mongo",
	)
	return
}

func runProcess(timeout time.Duration, command string, arg ...string) error {
	cmd := exec.Command(command, arg...)

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Printf("%v %v\n", command, strings.Join(arg, " "))

	timer := time.AfterFunc(30*time.Second, func() {
		cmd.Process.Kill()
	})

	if err := cmd.Wait(); err != nil {
		return err
	}

	timer.Stop()

	return nil
}
