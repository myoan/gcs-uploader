package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func upload(ctx context.Context, bh *storage.BucketHandle, srcPath, dstPath string) error {
	fmt.Printf("upload: %s\n", dstPath)
	obj := bh.Object(dstPath)
	writer := obj.NewWriter(ctx)

	f, err := os.Open(srcPath)

	if err != nil {
		return err
	}
	if _, err := io.Copy(writer, f); err != nil {
		return err
	}
	if err = writer.Close(); err != nil {
		return err
	}

	return nil
}

func walk(root, dir string) []string {
	files, err := ioutil.ReadDir(filepath.Join(root, dir))
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, walk(root, filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}

func readHash(path string) (string, error) {
	fp, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer fp.Close()

	buf := make([]byte, 64)
	hash := ""
	for {
		n, err := fp.Read(buf)
		if n == 0 {
			break
		}
		if err != nil {
			panic(err)
		}
		hash += string(buf[:n])
	}
	return hash, nil
}

func main() {
	var (
		bn   = flag.String("b", "", "backetname")
		cr   = flag.String("cred", "", "credential path")
		in   = flag.String("in", "", "input dir path")
		out  = flag.String("out", "", "output dir path")
		conc = flag.Int("out", 4, "upload cuncurrency")
	)
	flag.Parse()

	if len(*bn) == 0 {
		panic("err: undefined bucket name")
	}

	if len(*cr) == 0 {
		panic("err: undefined credential path")
	}

	if len(*in) == 0 {
		panic("err: undefined input path")
	}

	if len(*out) == 0 {
		panic("err: undefined output path")
	}

	ctx := context.Background()

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(*cr))
	if err != nil {
		panic("err: failed to create gcs client")
	}

	b := client.Bucket(*bn)
	if _, err = b.Attrs(ctx); err != nil {
		panic("bucket not found")
	}

	fmt.Printf("---------------------------------------------\n\n")
	fmt.Printf("bucket:     %s\n", *bn)
	fmt.Printf("credential: %s\n", *cr)
	fmt.Printf("input:      %s\n", *in)
	fmt.Printf("output:     %s\n", *out)
	fmt.Printf("---------------------------------------------\n\n")

	limit := make(chan struct{}, *conc)
	var wg sync.WaitGroup
	for _, f := range walk(*in, "") {
		wg.Add(1)
		go func(ctx context.Context, bh *storage.BucketHandle, src, dst string) {
			limit <- struct{}{}
			defer wg.Done()
			upload(ctx, bh, src, dst)
			<-limit
		}(ctx, b, f, filepath.Join(*out, f))
	}
	wg.Wait()
}
