package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/images/archive"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/platforms"
	"github.com/opencontainers/go-digest"
)

func imageName(s digest.Digest) string {
	return strings.TrimPrefix(string(s), "sha256:")
}

type imageSyncer struct {
	ctx      context.Context
	client   *containerd.Client
	cacheDir string
	archived map[string]struct{}

	mu       sync.Mutex
	imported map[string]struct{}

	ns string
}

func newImageSyncer(dir, ns string) (*imageSyncer, error) {
	sock := "/run/containerd/containerd.sock"
	if x := os.Getenv("CONTAINERD_SOCK"); x != "" {
		sock = x
	}
	client, err := containerd.New(sock)
	if err != nil {
		return nil, err
	}
	is := &imageSyncer{
		ns:       ns,
		ctx:      namespaces.WithNamespace(context.Background(), ns),
		cacheDir: dir,
		client:   client,
		archived: make(map[string]struct{}),
		imported: make(map[string]struct{}),
	}
	files, err := ioutil.ReadDir(is.cacheDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		is.archived[f.Name()] = struct{}{}
	}
	images, err := is.client.ListImages(is.ctx)
	if err != nil {
		return nil, err
	}
	is.mu.Lock()
	for _, image := range images {
		is.imported[imageName(image.Target().Digest)] = struct{}{}
	}
	is.mu.Unlock()
	return is, nil
}

func (is *imageSyncer) Sync() error {
	if err := is.exportAll(); err != nil {
		return err
	}
	return is.importAll()
}

func (is *imageSyncer) exportImage(image containerd.Image, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	opts := []archive.ExportOpt{
		archive.WithPlatform(platforms.Default()),
		archive.WithImage(is.client.ImageService(), image.Name()),
	}
	return is.client.Export(is.ctx, f, opts...)
}

func (is *imageSyncer) exportAll() error {
	images, err := is.client.ListImages(is.ctx)
	if err != nil {
		return err
	}
	for _, image := range images {
		name := imageName(image.Target().Digest)
		if _, ok := is.archived[name]; ok {
			continue
		}
		if err := is.exportImage(image, filepath.Join(is.cacheDir, name)); err != nil {
			log.Print(err)
			continue
		}
		log.Printf("export: %s\n", name)
		is.mu.Lock()
		is.archived[name] = struct{}{}
		is.imported[name] = struct{}{}
		is.mu.Unlock()
	}
	return nil
}

func (is *imageSyncer) importAll() error {
	files, err := ioutil.ReadDir(is.cacheDir)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, f := range files {
		wg.Add(1)
		go func(f os.FileInfo) {
			defer wg.Done()
			name := f.Name()
			is.mu.Lock()
			_, ok := is.imported[name]
			is.mu.Unlock()
			if ok {
				return
			}
			//log.Printf("in = %+v\n", filepath.Join(is.cacheDir, name))
			if err := is.importImage(filepath.Join(is.cacheDir, name)); err != nil {
				log.Printf("error importing %q: %v", name, err)
				return
			}
			log.Printf("import: %s\n", name)
			is.imported[name] = struct{}{}
		}(f)
	}
	wg.Wait()
	return nil
}

func (is *imageSyncer) importImage(path string) error {
	opts := []containerd.ImportOpt{
		//containerd.WithImageRefTranslator(archive.AddRefPrefix("foo/bar")),
	}
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	ctx := namespaces.WithNamespace(context.Background(), is.ns)
	_, err = is.client.Import(ctx, f, opts...)
	if err != nil {
		return err
	}
	//log.Printf("out = %+v\n", path)
	//for _, imgrec := range imgrecs {
	//fmt.Printf("imgrec.Name = %+v\n", imgrec.Name)
	//}
	return nil
}
