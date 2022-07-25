package main

import (
	"context"
	"main/src"
)

func main() {
	container := src.BuildContainer()
	err := container.Invoke(func(application *src.Application) {
		ctx := context.Background()
		err := application.Init(ctx)
		if err != nil {
			panic(err)
		}

		err = application.Run(ctx)
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(err)
	}
}
