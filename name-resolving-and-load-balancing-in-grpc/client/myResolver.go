package main

import (

	"google.golang.org/grpc/resolver"
)

const (
	myScheme   = "makito"
	myEndPoint = "resolver.makito.icu"
)

var addrs = []string{"127.0.0.1:8973", "127.0.0.1:8974", "127.0.0.1:8975"}

// 实现Resolver接口
type myResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *myResolver) ResolveNow(resolver.ResolveNowOptions) {
	addrStrs := r.addrsStore[myEndPoint]
	addrList := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrList[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrList})
}
func (r *myResolver) Close() {}

// myResolverBuilder 需实现 Builder 接口
type myResolverBuilder struct{}

func (*myResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &myResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			myEndPoint: addrs,
		},
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (*myResolverBuilder) Scheme() string { return myScheme }
