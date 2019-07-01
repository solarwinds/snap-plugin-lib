package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/internal/pluginrpc"
	"google.golang.org/grpc"
)

const (
	defaultLoadDelay = 500 * time.Millisecond

	defaultCollectInterval = 5 * time.Second
	defaultPingInterval    = 2 * time.Second

	defaultGRPCTimeout = 10 * time.Second

	defaultTaskID = 1
)

type Options struct {
	PluginIP           string
	PluginPort         int
	CollectInterval    time.Duration
	PingInterval       time.Duration
	MaxCollectRequests int

	PluginConfig string
}

func main() {
	doneCh := make(chan error)

	opt := parseCmdLine()

	grpcServerAddr := fmt.Sprintf("%s:%d", opt.PluginIP, opt.PluginPort)
	cl, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Can't to GRPC Server on %s", grpcServerAddr)
		os.Exit(1)
	}
	defer func() { _ = cl.Close() }()

	// Load, collect, unload
	go func() {
		cc := pluginrpc.NewCollectorClient(cl)

		reqLoad := &pluginrpc.LoadRequest{
			TaskId:          defaultTaskID,
			JsonConfig:      []byte(opt.PluginConfig),
			MetricSelectors: nil,
		}

		ctx, _ := context.WithTimeout(context.Background(), defaultGRPCTimeout)
		_, err := cc.Load(ctx, reqLoad)
		if err != nil {
			doneCh <- fmt.Errorf("can't send load request to plugin: %v", err)
		}
		time.Sleep(defaultLoadDelay)

		reqCounter := 0
		for {
			reqCounter++
			reqColl := &pluginrpc.CollectRequest{
				TaskId: defaultTaskID,
			}

			ctx, _ := context.WithTimeout(context.Background(), defaultGRPCTimeout)
			stream, err := cc.Collect(ctx, reqColl)
			if err != nil {
				doneCh <- fmt.Errorf("can't send collect request to plugin: %v", err)
			}

			var recvMts []string
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					doneCh <- fmt.Errorf("error when receiving collect reply from plugin (%v)", err)
				}

				for _, mt := range resp.MetricSet {
					recvMts = append(recvMts, grpcMetricToString(mt))
				}
			}

			fmt.Printf("Recived %d metric(s)\n", len(recvMts))
			for _, mt := range recvMts {
				fmt.Printf(" %s\n", mt)
			}
			fmt.Printf("\n")

			if reqCounter == opt.MaxCollectRequests {
				break
			}
			time.Sleep(opt.CollectInterval)
		}

		time.Sleep(defaultLoadDelay)
		reqUnload := &pluginrpc.UnloadRequest{
			TaskId: defaultTaskID,
		}

		ctx, _ = context.WithTimeout(context.Background(), defaultGRPCTimeout)
		_, err = cc.Unload(ctx, reqUnload)
		if err != nil {
			doneCh <- fmt.Errorf("can't send unload request to plugin: %v", err)
		}

		doneCh <- nil
	}()

	// ping routine
	go func() {
		for {
			cc := pluginrpc.NewControllerClient(cl)
			req := &pluginrpc.PingRequest{}
			_, err := cc.Ping(context.Background(), req)
			if err != nil {
				doneCh <- fmt.Errorf("can't start ")
			}
			time.Sleep(opt.PingInterval)
		}
	}()

	doneErr := <-doneCh
	if doneErr != nil {
		fmt.Printf("Snap-mock exists because of error: %v", doneErr)
	}
}

func parseCmdLine() *Options {
	opt := &Options{}

	flag.StringVar(&opt.PluginIP,
		"plugin-ip", "127.0.0.1",
		"IP Address of GRPC Server run by plugin")

	flag.IntVar(&opt.PluginPort,
		"plugin-port", 0,
		"Port of GRPC Server run by plugin")

	flag.StringVar(&opt.PluginConfig,
		"plugin-config", "{}",
		"Plugin configuration (should be valid JSON)")

	flag.IntVar(&opt.MaxCollectRequests,
		"max-collect-requests", 0,
		"Maximum number of collect requests (default 0 for infinite)")

	flag.DurationVar(&opt.CollectInterval,
		"collect-interval", defaultCollectInterval,
		"Duration between Collect requests")

	flag.DurationVar(&opt.PingInterval,
		"ping-interval", defaultPingInterval,
		"Duration between Ping requests")

	flag.Parse()

	return opt
}

func grpcMetricToString(metric *pluginrpc.Metric) string {
	var nsStr []string
	for _, ns := range metric.Namespace {
		nsStr = append(nsStr, ns.Value)
	}

	return fmt.Sprintf("%s %v [%v]", strings.Join(nsStr, "."), metric.Value, metric.Tags)
}
