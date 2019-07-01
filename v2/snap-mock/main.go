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
	defaultGRPCIP          = "127.0.0.1"
	defaultGRPCPort        = 0
	defaultConfig          = "{}"
	defaultTaskID          = 1
	defaultCollectInterval = 5 * time.Second
	defaultPingInterval    = 2 * time.Second

	grpcLoadDelay      = 500 * time.Millisecond
	grpcRequestTimeout = 10 * time.Second
)

type Options struct {
	PluginIP           string
	PluginPort         int
	CollectInterval    time.Duration
	PingInterval       time.Duration
	MaxCollectRequests int

	PluginConfig string
}

///////////////////////////////////////////////////////////////////////////////

func parseCmdLine() *Options {
	opt := &Options{}

	flag.StringVar(&opt.PluginIP,
		"plugin-ip", defaultGRPCIP,
		"IP Address of GRPC Server run by plugin")

	flag.IntVar(&opt.PluginPort,
		"plugin-port", defaultGRPCPort,
		"Port of GRPC Server run by plugin")

	flag.StringVar(&opt.PluginConfig,
		"plugin-config", defaultConfig,
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

///////////////////////////////////////////////////////////////////////////////

func main() {
	doneCh := make(chan error)

	opt := parseCmdLine()

	// Create connection
	grpcServerAddr := fmt.Sprintf("%s:%d", opt.PluginIP, opt.PluginPort)
	cl, err := grpc.Dial(grpcServerAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Can't to GRPC Server on %s", grpcServerAddr)
		os.Exit(1)
	}
	defer func() { _ = cl.Close() }()

	// Load, collect, unload routine
	go func() {
		cc := pluginrpc.NewCollectorClient(cl)

		err := doLoadRequest(cc, opt)
		if err != nil {
			doneCh <- fmt.Errorf("can't send load request to plugin: %v", err)
		}
		time.Sleep(grpcLoadDelay)

		reqCounter := 0
		for {
			reqCounter++
			recvMts, err := doCollectRequest(cc, opt)
			if err != nil {
				doneCh <- fmt.Errorf("can't send collect request to plugin: %v", err)
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

		time.Sleep(grpcLoadDelay)

		err = doUnloadRequest(cc, opt)
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

///////////////////////////////////////////////////////////////////////////////

func doLoadRequest(cc pluginrpc.CollectorClient, opt *Options) error {
	reqLoad := &pluginrpc.LoadRequest{
		TaskId:          defaultTaskID,
		JsonConfig:      []byte(opt.PluginConfig),
		MetricSelectors: nil,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Load(ctx, reqLoad)

	return err
}

func doUnloadRequest(cc pluginrpc.CollectorClient, _ *Options) error {
	reqUnload := &pluginrpc.UnloadRequest{
		TaskId: defaultTaskID,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Unload(ctx, reqUnload)

	return err
}

func doCollectRequest(cc pluginrpc.CollectorClient, _ *Options) ([]string, error) {
	var recvMts []string

	reqColl := &pluginrpc.CollectRequest{
		TaskId: defaultTaskID,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	stream, err := cc.Collect(ctx, reqColl)
	if err != nil {
		return recvMts, fmt.Errorf("can't send collect request to plugin: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return recvMts, fmt.Errorf("error when receiving collect reply from plugin (%v)", err)
		}

		for _, mt := range resp.MetricSet {
			recvMts = append(recvMts, grpcMetricToString(mt))
		}
	}

	return recvMts, nil
}

///////////////////////////////////////////////////////////////////////////////

func grpcMetricToString(metric *pluginrpc.Metric) string {
	var nsStr []string
	for _, ns := range metric.Namespace {
		nsStr = append(nsStr, ns.Value)
	}

	return fmt.Sprintf("%s %v [%v]", strings.Join(nsStr, "."), metric.Value, metric.Tags)
}
