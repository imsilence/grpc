package file

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/resolver"

	"github.com/fsnotify/fsnotify"
)

// 创建Resolver
func NewBuilderWithScheme(scheme string) *Resolver {
	return &Resolver{
		ResolveNowCallback: func(opts resolver.ResolveNowOptions) {},
		scheme:             scheme,
	}
}

// 定义Resolver
type Resolver struct {
	ResolveNowCallback func(resolver.ResolveNowOptions)
	scheme             string
	CC                 resolver.ClientConn
	bootstrapState     *resolver.State
}

// 初始化state
func (r *Resolver) InitialState(s resolver.State) {
	log.Print("InitialState")
	r.bootstrapState = &s
}

// build resolver
func (r *Resolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	log.Print("Build")

	r.CC = cc
	if r.bootstrapState != nil {
		r.UpdateState(*r.bootstrapState)
	}
	go r.watcher()
	return r, nil
}

// 监听变化
func (r *Resolver) watcher() {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && filepath.Base(event.Name) == "server.names" {
					addrs := make([]resolver.Address, 0, 10)
					if file, err := os.Open(event.Name); err == nil {
						defer file.Close()
						scanner := bufio.NewScanner(file)
						for scanner.Scan() {
							line := strings.TrimSpace(scanner.Text())
							if line == "" || strings.HasPrefix(line, "#") {
								continue
							}
							addrs = append(addrs, resolver.Address{Addr: line})
						}
					}
					r.UpdateState(resolver.State{Addresses: addrs})

				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("etc/")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// 获取协议
func (r *Resolver) Scheme() string {
	return r.scheme
}

//
func (r *Resolver) ResolveNow(o resolver.ResolveNowOptions) {
	log.Print("ResolveNow")
	r.ResolveNowCallback(o)
}

// 关闭
func (*Resolver) Close() {}

// 更新state
func (r *Resolver) UpdateState(s resolver.State) {
	log.Print("UpdateState")
	r.CC.UpdateState(s)
}
