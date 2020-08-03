/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/librato/snap-plugin-lib-go/v2/pluginrpc"
	"google.golang.org/grpc"
)

const (
	defaultGRPCIP          = "127.0.0.1"
	defaultGRPCPort        = 0
	defaultConfig          = "{}"
	defaultFilter          = ""
	defaultTaskID          = ""
	defaultCollectInterval = 5 * time.Second
	defaultPingInterval    = 2 * time.Second
	defaultStreamDuration  = 60 * time.Second

	grpcLoadDelay      = 500 * time.Millisecond
	grpcRequestTimeout = 10 * time.Second
	closeDelay         = 1 * time.Second // wait until Unload() causes end of StreamCollect() to avoid GRPC errors

	filterSeparator = ";"

	maxTaskId = 1024
)

type Options struct {
	PluginIP           string
	CollectorPort      int
	PublisherPort      int
	CollectInterval    time.Duration
	PingInterval       time.Duration
	MaxCollectRequests int
	SendKill           bool
	RequestInfo        bool

	IsStream       bool
	StreamDuration time.Duration

	PluginConfig string
	PluginFilter string
	TaskId       string
}

const (
	unloadMaxRetry   = 5
	unloadRetryDelay = 1 * time.Second

	stoppedByUser = 1
)

type collectChunk struct {
	mts []*pluginrpc.Metric
	//mts      []string
	warnings []string
	err      error
}

///////////////////////////////////////////////////////////////////////////////

func parseCmdLine() *Options {
	opt := &Options{}

	flag.StringVar(&opt.PluginIP,
		"plugin-ip", defaultGRPCIP,
		"IP Address of GRPC Server run by plugin")

	flag.IntVar(&opt.CollectorPort,
		"collector-port", defaultGRPCPort,
		"Port of GRPC Server run by plugin")

	flag.IntVar(&opt.PublisherPort,
		"publisher-port", defaultGRPCPort,
		"Port of GRPC Server run by publisher plugin")

	flag.StringVar(&opt.TaskId,
		"task-id", defaultTaskID,
		"Task identifier used to make GRPC requests ('' means random)")

	flag.StringVar(&opt.PluginConfig,
		"plugin-config", defaultConfig,
		"Plugin configuration (should be valid JSON)")

	flag.StringVar(&opt.PluginFilter,
		"plugin-filter", defaultFilter,
		"Plugin filter (definition which subset of metrics should be gathered), ie. '/example/static/*;/example/global/*'")

	flag.IntVar(&opt.MaxCollectRequests,
		"max-collect-requests", 0,
		"Maximum number of collect requests (default 0 for infinite)")

	flag.DurationVar(&opt.CollectInterval,
		"collect-interval", defaultCollectInterval,
		"Duration between Collect requests")

	flag.DurationVar(&opt.PingInterval,
		"ping-interval", defaultPingInterval,
		"Duration between Ping requests")

	flag.BoolVar(&opt.SendKill,
		"send-kill", false,
		"When set, Kill request will be sent after 'max-collect-requests' collects ")

	flag.BoolVar(&opt.RequestInfo,
		"request-info", false,
		"When set, Info request will be sent after each collect")

	flag.BoolVar(&opt.IsStream,
		"stream", false,
		"When set, expects connecting to streaming plugin")

	flag.DurationVar(&opt.StreamDuration,
		"stream-duration", defaultStreamDuration,
		"Duration of debugging streaming collector, after this time Unload request will be send")

	flag.Parse()

	if opt.TaskId == defaultTaskID {
		opt.TaskId = fmt.Sprintf("task-%d", rand.Intn(maxTaskId))
	}

	return opt
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	rand.Seed(time.Now().UnixNano())

	doneCh := make(chan error)

	opt := parseCmdLine()

	usePublisher := false
	if opt.PublisherPort != defaultGRPCPort {
		usePublisher = true
	}
	// Create connection
	grpcServerCollAddr := fmt.Sprintf("%s:%d", opt.PluginIP, opt.CollectorPort)
	clColl, err := grpc.Dial(grpcServerCollAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Can't start GRPC Server on %s (%v)", grpcServerCollAddr, err)
		os.Exit(1)
	}
	defer func() { _ = clColl.Close() }()

	var clPub *grpc.ClientConn
	if usePublisher {
		grpcServerPubAddr := fmt.Sprintf("%s:%d", opt.PluginIP, opt.PublisherPort)
		clPub, err = grpc.Dial(grpcServerPubAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Can't start GRPC Server on %s (%v)", grpcServerPubAddr, err)
			os.Exit(1)
		}
		defer func() { _ = clPub.Close() }()
	}
	// Load, collect, publish, unload routine
	go func() {

		collClient := pluginrpc.NewCollectorClient(clColl)

		err := doLoadRequest(collClient, opt)
		if err != nil {
			doneCh <- fmt.Errorf("can't send load request to plugin: %v", err)
		}

		var publishClient pluginrpc.PublisherClient
		if usePublisher {
			publishClient = pluginrpc.NewPublisherClient(clPub)
			errPub := doPubLoadRequest(publishClient, opt)
			if errPub != nil {
				doneCh <- fmt.Errorf("can't send load request to plugin: %v", err)
			}
		}
		// Handle ctrl+C
		notifyCh := make(chan os.Signal, 1)
		signal.Notify(notifyCh, os.Interrupt)
		go func() {
			<-notifyCh
			fmt.Printf("!! Ctrl+c pressed !! trying to unload current task\n")

			for i := 0; i < unloadMaxRetry; i++ {
				err := doUnloadRequest(collClient, opt)
				if err != nil {
					fmt.Printf("!! Can't unload plugin (%v), will retry (%d/%d)...\n", err, i+1, unloadMaxRetry)
					time.Sleep(unloadRetryDelay)
					continue
				}

				break
			}
			if usePublisher {
				for i := 0; i < unloadMaxRetry; i++ {
					err := doPubUnloadRequest(publishClient, opt)
					if err != nil {
						fmt.Printf("!! Can't unload plugin (%v), will retry (%d/%d)...\n", err, i+1, unloadMaxRetry)
						time.Sleep(unloadRetryDelay)
						continue
					}

					break
				}
			}
			os.Exit(stoppedByUser)
		}()
		time.Sleep(grpcLoadDelay)

		reqCounter := 0
		for {
			reqCounter++

			if opt.IsStream {
				go func() {
					time.Sleep(opt.StreamDuration)

					err = doUnloadRequest(collClient, opt)
					if err != nil {
						doneCh <- fmt.Errorf("can't send unload request to plugin: %v", err)
					}

					doneCh <- nil
				}()
			}

			var mtsChunks [][]*pluginrpc.Metric

			chunkCh := doCollectRequest(collClient, opt)
			for chunk := range chunkCh {
				if err != nil {
					doneCh <- fmt.Errorf("can't send collect request to plugin: %v", err)
				}

				fmt.Printf("\nReceived %d warning(s)\n", len(chunk.warnings))
				for _, warn := range chunk.warnings {
					fmt.Printf(" %s\n", warn)
				}

				fmt.Printf("\nReceived %d metric(s)\n", len(chunk.mts))
				for _, mt := range chunk.mts {
					fmt.Printf(" %s\n", grpcMetricToString(mt))
				}

				mtsChunks = append(mtsChunks, chunk.mts)
			}
			if usePublisher {
				err := doPublishRequest(publishClient, mtsChunks, opt)
				if err != nil {
					fmt.Printf("Not good")
				}
			}
			if opt.RequestInfo {
				info, err := doInfoRequest(collClient, opt)
				if err != nil {
					doneCh <- fmt.Errorf("can't send info request to plugin: %v", err)
				}

				fmt.Printf("\nReceived info:\n %v\n", string(info))
			}

			if reqCounter == opt.MaxCollectRequests || opt.IsStream {
				break
			}
			time.Sleep(opt.CollectInterval)
		}

		time.Sleep(grpcLoadDelay)

		if !opt.IsStream {
			err = doUnloadRequest(collClient, opt)
			if err != nil {
				doneCh <- fmt.Errorf("can't send unload request to plugin: %v", err)
			}

			doneCh <- nil
		}
	}()

	// ping routine
	contClient := pluginrpc.NewControllerClient(clColl)

	var contPubClient pluginrpc.ControllerClient

	if usePublisher {
		contPubClient = pluginrpc.NewControllerClient(clPub)
	}

	go func() {
		for {
			req := &pluginrpc.PingRequest{}
			_, err := contClient.Ping(context.Background(), req)
			if err != nil {
				doneCh <- fmt.Errorf("can't start: %v", err)
			}
			if usePublisher {
				_, err2 := contPubClient.Ping(context.Background(), req)
				if err2 != nil {
					doneCh <- fmt.Errorf("can't start: %v", err2)
				}
			}
			time.Sleep(opt.PingInterval)
		}
	}()

	doneErr := <-doneCh
	time.Sleep(closeDelay)

	if opt.SendKill {
		err := doKillRequest(contClient)
		if err != nil {
			doneCh <- fmt.Errorf("can't send kill request to plugin: %v", err)
		}
	}

	if doneErr != nil {
		fmt.Printf("Snap-mock exists because of error: %v", doneErr)
	}
}

///////////////////////////////////////////////////////////////////////////////

func doLoadRequest(cc pluginrpc.CollectorClient, opt *Options) error {
	filter := []string{}
	if opt.PluginFilter != defaultFilter {
		filter = strings.Split(opt.PluginFilter, filterSeparator)
	}

	reqLoad := &pluginrpc.LoadCollectorRequest{
		TaskId:          opt.TaskId,
		JsonConfig:      []byte(opt.PluginConfig),
		MetricSelectors: filter,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Load(ctx, reqLoad)

	return err
}

func doPubLoadRequest(cc pluginrpc.PublisherClient, opt *Options) error {
	reqLoad := &pluginrpc.LoadPublisherRequest{
		TaskId:     opt.TaskId,
		JsonConfig: []byte(opt.PluginConfig), // FIXME
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Load(ctx, reqLoad)

	return err
}

func doPubUnloadRequest(cc pluginrpc.PublisherClient, opt *Options) error {
	reqUnload := &pluginrpc.UnloadPublisherRequest{
		TaskId: opt.TaskId,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Unload(ctx, reqUnload)

	return err
}

func doUnloadRequest(cc pluginrpc.CollectorClient, opt *Options) error {
	reqUnload := &pluginrpc.UnloadCollectorRequest{
		TaskId: opt.TaskId,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Unload(ctx, reqUnload)

	return err
}

func doKillRequest(cc pluginrpc.ControllerClient) error {
	reqKill := &pluginrpc.KillRequest{}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	_, err := cc.Kill(ctx, reqKill)
	return err
}

func doInfoRequest(cc pluginrpc.CollectorClient, opt *Options) ([]byte, error) {
	reqInfo := &pluginrpc.InfoRequest{
		TaskId: opt.TaskId,
	}

	ctx, fn := context.WithTimeout(context.Background(), grpcRequestTimeout)
	defer fn()

	resp, err := cc.Info(ctx, reqInfo)
	if err != nil {
		return nil, err
	}
	return resp.Info, nil
}

func doPublishRequest(pc pluginrpc.PublisherClient, mts [][]*pluginrpc.Metric, opt *Options) error {
	stream, _ := pc.Publish(context.Background())

	for _, chunk := range mts {
		reqPubl := &pluginrpc.PublishRequest{
			TaskId:    opt.TaskId,
			MetricSet: chunk,
		}
		stream.Send(reqPubl)
	}

	_, err := stream.CloseAndRecv()
	return err
}

func doCollectRequest(cc pluginrpc.CollectorClient, opt *Options) chan collectChunk {
	var recvMts []*pluginrpc.Metric
	var recvWarns []string

	chunkCh := make(chan collectChunk)

	reqColl := &pluginrpc.CollectRequest{
		TaskId: opt.TaskId,
	}

	go func() {
		ctx := context.Background()

		if !opt.IsStream {
			var fn context.CancelFunc
			ctx, fn = context.WithTimeout(context.Background(), grpcRequestTimeout)
			defer fn()
		}

		defer func() { close(chunkCh) }()

		stream, err := cc.Collect(ctx, reqColl)
		if err != nil {
			chunkCh <- collectChunk{
				mts:      recvMts,
				warnings: recvWarns,
				err:      fmt.Errorf("can't send collect request to plugin: %v", err),
			}
			return
		}

		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				chunkCh <- collectChunk{
					mts:      recvMts,
					warnings: recvWarns,
					err:      fmt.Errorf("error when receiving collect reply from plugin (%v)", err),
				}
				return
			}

			recvMts = append(recvMts, resp.MetricSet...)

			for _, warns := range resp.Warnings {
				recvWarns = append(recvWarns, grpcWarningToString(warns))
			}

			if opt.IsStream {
				chunkCh <- collectChunk{
					mts:      recvMts,
					warnings: recvWarns,
					err:      nil,
				}

				recvMts = nil
				recvWarns = nil
			}
		}

		chunkCh <- collectChunk{
			mts:      recvMts,
			warnings: recvWarns,
			err:      nil,
		}
	}()

	return chunkCh
}

///////////////////////////////////////////////////////////////////////////////

func grpcMetricToString(metric *pluginrpc.Metric) string {
	var nsStr []string
	for _, ns := range metric.Namespace {
		nsStr = append(nsStr, ns.Value)
	}

	return fmt.Sprintf("%s %v [%v]", strings.Join(nsStr, "."), metric.Value, metric.Tags)
}

func grpcWarningToString(warning *pluginrpc.Warning) string {
	return fmt.Sprintf("[%s] %s", time.Unix(warning.Timestamp.Sec, warning.Timestamp.Nsec), warning.Message)
}
