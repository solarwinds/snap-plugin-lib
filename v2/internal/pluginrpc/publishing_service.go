package pluginrpc

const (
	maxPublishChunkSize = 100
)

type publishingService struct {
}

func newPublishingService() PublisherServer {
	return &publishingService{
	}
}

func (ps *publishingService) Publish(stream Publisher_PublishServer) error {
	// todo: implement
	return nil
}
