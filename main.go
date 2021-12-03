package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"os"
	"strings"
)

var files = []string{
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-c5b2ccd1917e73f4d004c2a6eef71c6019ddf470-X_rock.mat`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-045fa2ffa69eaf7ac3987d021e3fc7e48997611c-paper.jld2`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-a63c641bba21050e27b94c58568347a566ef662b-rock.jld2`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-ad90808a6dcb1abde7ef3384230fa15de13e489a-X_scissors.mat`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-d9dabc1decdc317371b688f200d4a2d28cf1d133-Manifest.toml`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-360e5ed7335b8fa4f8aa7a53bee7ad8ad94dce0b-scissors.jld2`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-1ad7bb1a8bfecea37f3e159439083dc43f9b87df-X_paper.mat`,
	`https://storage.googleapis.com/mynerva-testing-static/rps_classificatioin-496e15e2f88a711138a07d416ca6fcd5a41fa4b9-Project.toml`,
}

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	for _, file := range files {
		downloadFile(g, ctx, file)
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func downloadFile(eg *errgroup.Group, ctx context.Context, file string) {
	eg.Go(func() error {
		req, err := http.NewRequestWithContext(ctx, "GET", file, nil)
		if err != nil {
			return err
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		f, err := os.OpenFile(
			file[strings.LastIndex(file, "/")+1:],
			os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(f, res.Body); err != nil {
			return err
		}
		fmt.Printf("download succeeded: %q\n", file)
		return nil
	})
}
