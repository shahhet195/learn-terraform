package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

var (
	flagDir    = flag.String("dir", ".", "Directory to scan recursively")
	flagInject = flag.String("inject", "README.md", "File to inject module index")
)

const (
	MarkerBegin = "<!-- BEGIN_MODULE_INDEX -->"
	MarkerEnd   = "<!-- END_MODULE_INDEX -->"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigs
		cancelFn()
	}()

	dirs := make(map[string]struct{})
	err := filepath.Walk(*flagDir,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			// Skip everything in .terraform directory
			pathParts := strings.Split(path, "/")
			for _, pathPart := range pathParts {
				if pathPart == ".terraform" {
					return nil
				}
			}
			if filepath.Base(path) == "module.yml" {
				dirs[filepath.Dir(path)] = struct{}{}
			}

			if err != nil {
				return err
			}
			return nil
		})
	failOnError(err)

	modErrors := make(map[string]error)
	var modules []ModuleYml
	for dir, _ := range dirs {
		select {
		case <-ctx.Done():
			modErrors[dir] = errors.New("Context canceled")
			continue
		default:
		}

		var m ModuleYml
		b, err := ioutil.ReadFile(filepath.Join(dir, "module.yml"))
		if err != nil {
			modErrors[dir] = err
			continue
		}

		err = yaml.Unmarshal(b, &m)
		if err != nil {
			modErrors[dir] = err
			continue
		}

		m.AbsDirectory = dir
		relPath, err := filepath.Rel(*flagDir, dir)
		failOnError(err)
		m.RelDirectory = relPath
		parts := strings.Split(relPath, "/")
		m.Parts = parts
		modules = append(modules, m)
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i].RelDirectory < modules[j].RelDirectory
	})

	if len(modErrors) > 0 {
		log.Println("Failed:")
		for dir, err := range modErrors {
			log.Println("-----------------", dir, "---------------------------")
			log.Println(err)
		}
		os.Exit(1)
	}

	byFirstLevel := make(map[string][]ModuleYml)
	for _, m := range modules {
		key := strings.Join(m.Parts, "/") + "test"
		if len(m.Parts) > 2 {
			key = strings.Join(m.Parts[0:2], "/")
		}

		byFirstLevel[key] = append(byFirstLevel[key], m)
	}

	var templateError error
	tmplB, err := ioutil.ReadFile("template.tpl")
	failOnError(err)
	t := template.Must(template.New("main").Funcs(sprig.TxtFuncMap()).Parse(string(tmplB)))
	var tbuf bytes.Buffer
	err = t.Execute(&tbuf, struct {
		Modules      []ModuleYml
		ByFirstLevel map[string][]ModuleYml
	}{
		Modules:      modules,
		ByFirstLevel: byFirstLevel,
	})
	failOnError(err)
	failOnError(templateError)

	// TODO(blang): Cleanup this mess
	tmplResult := MarkerBegin + "\n"
	tmplResult += string(tbuf.Bytes())
	tmplResult += MarkerEnd

	injectFile := filepath.Join(*flagDir, *flagInject)
	readmeB, err := ioutil.ReadFile(injectFile)
	failOnError(err)
	markerRegexp := regexp.MustCompile(`(?s)` + MarkerBegin + `(.*)` + MarkerEnd)
	content := markerRegexp.ReplaceAllString(string(readmeB), string(tmplResult))
	stat, err := os.Stat(injectFile)
	failOnError(err)
	ioutil.WriteFile(injectFile, []byte(content), stat.Mode().Perm())
	log.Println("Successfully written ", injectFile)
}

type ModuleYml struct {
	ShortDescription string   `json:"short_description"`
	Parts            []string `json:"-"`
	RelDirectory     string   `json:"-"`
	AbsDirectory     string   `json:"-"`
	err              error    `json:"-"`
}
