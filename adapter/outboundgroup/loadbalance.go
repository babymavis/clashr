package outboundgroup

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"

	"github.com/Dreamacro/clash/adapter/outbound"
	"github.com/Dreamacro/clash/common/murmur3"
	"github.com/Dreamacro/clash/common/singledo"
	"github.com/Dreamacro/clash/component/dialer"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/constant/provider"

	"golang.org/x/net/publicsuffix"
)

type strategyFn = func(proxies []C.Proxy, metadata *C.Metadata) C.Proxy

type LoadBalance struct {
	*outbound.Base
	disableUDP bool
	filter     string
	single     *singledo.Single
	providers  []provider.ProxyProvider
	strategyFn strategyFn
}

var errStrategy = errors.New("unsupported strategy")

func parseStrategy(config map[string]interface{}) string {
	if elm, ok := config["strategy"]; ok {
		if strategy, ok := elm.(string); ok {
			return strategy
		}
	}
	return "random"
}

func getKey(metadata *C.Metadata) string {
	if metadata == nil {
		return ""
	}

	if metadata.Host != "" {
		// ip host
		if ip := net.ParseIP(metadata.Host); ip != nil {
			return metadata.Host
		}

		if etld, err := publicsuffix.EffectiveTLDPlusOne(metadata.Host); err == nil {
			return etld
		}
	}

	if metadata.DstIP == nil {
		return ""
	}

	return metadata.DstIP.String()
}

func jumpHash(key uint64, buckets int32) int32 {
	var b, j int64

	for j < int64(buckets) {
		b = j
		key = key*2862933555777941757 + 1
		j = int64(float64(b+1) * (float64(int64(1)<<31) / float64((key>>33)+1)))
	}

	return int32(b)
}

// DialContext implements C.ProxyAdapter
func (lb *LoadBalance) DialContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (c C.Conn, err error) {
	defer func() {
		if err == nil {
			c.AppendToChains(lb)
		}
	}()

	proxy := lb.Unwrap(metadata)

	c, err = proxy.DialContext(ctx, metadata, lb.Base.DialOptions(opts...)...)
	return
}

// ListenPacketContext implements C.ProxyAdapter
func (lb *LoadBalance) ListenPacketContext(ctx context.Context, metadata *C.Metadata, opts ...dialer.Option) (pc C.PacketConn, err error) {
	defer func() {
		if err == nil {
			pc.AppendToChains(lb)
		}
	}()

	proxy := lb.Unwrap(metadata)
	return proxy.ListenPacketContext(ctx, metadata, lb.Base.DialOptions(opts...)...)
}

// SupportUDP implements C.ProxyAdapter
func (lb *LoadBalance) SupportUDP() bool {
	return !lb.disableUDP
}

func strategyRandom() strategyFn {
	return func(proxies []C.Proxy, metadata *C.Metadata) C.Proxy {
		aliveProxies := make([]C.Proxy, 0, len(proxies))
		for _, proxy := range proxies {
			if proxy.Alive() {
				aliveProxies = append(aliveProxies, proxy)
			}
		}
		aliveNum := int64(len(aliveProxies))
		if aliveNum == 0 {
			return proxies[0]
		}
		idx, err := rand.Int(rand.Reader, big.NewInt(aliveNum))
		if err != nil {
			return aliveProxies[0]
		}

		return aliveProxies[idx.Int64()]
	}
}

func strategyRoundRobin() strategyFn {
	idx := 0
	return func(proxies []C.Proxy, metadata *C.Metadata) C.Proxy {
		length := len(proxies)
		for i := 0; i < length; i++ {
			idx = (idx + 1) % length
			proxy := proxies[idx]
			if proxy.Alive() {
				return proxy
			}
		}

		return proxies[0]
	}
}

func strategyConsistentHashing() strategyFn {
	maxRetry := 5
	return func(proxies []C.Proxy, metadata *C.Metadata) C.Proxy {
		key := uint64(murmur3.Sum32([]byte(getKey(metadata))))
		buckets := int32(len(proxies))
		for i := 0; i < maxRetry; i, key = i+1, key+1 {
			idx := jumpHash(key, buckets)
			proxy := proxies[idx]
			if proxy.Alive() {
				return proxy
			}
		}

		return proxies[0]
	}
}

// Unwrap implements C.ProxyAdapter
func (lb *LoadBalance) Unwrap(metadata *C.Metadata) C.Proxy {
	proxies := lb.proxies(true)
	return lb.strategyFn(proxies, metadata)
}

func (lb *LoadBalance) proxies(touch bool) []C.Proxy {
	elm, _, _ := lb.single.Do(func() (interface{}, error) {
		return getProvidersProxies(lb.providers, touch, lb.filter), nil
	})

	return elm.([]C.Proxy)
}

// MarshalJSON implements C.ProxyAdapter
func (lb *LoadBalance) MarshalJSON() ([]byte, error) {
	var all []string
	for _, proxy := range lb.proxies(false) {
		all = append(all, proxy.Name())
	}
	return json.Marshal(map[string]interface{}{
		"type": lb.Type().String(),
		"all":  all,
	})
}

func NewLoadBalance(option *GroupCommonOption, providers []provider.ProxyProvider, strategy string) (lb *LoadBalance, err error) {
	var strategyFn strategyFn
	switch strategy {
	case "random":
		strategyFn = strategyRandom()
	case "consistent-hashing":
		strategyFn = strategyConsistentHashing()
	case "round-robin":
		strategyFn = strategyRoundRobin()
	default:
		return nil, fmt.Errorf("%w: %s", errStrategy, strategy)
	}
	return &LoadBalance{
		Base: outbound.NewBase(outbound.BaseOption{
			Name:        option.Name,
			Type:        C.LoadBalance,
			Interface:   option.Interface,
			RoutingMark: option.RoutingMark,
		}),
		single:     singledo.NewSingle(defaultGetProxiesDuration),
		providers:  providers,
		strategyFn: strategyFn,
		disableUDP: option.DisableUDP,
		filter:     option.Filter,
	}, nil
}
