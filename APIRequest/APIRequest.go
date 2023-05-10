package APIRequest

// API request is a dummy request that hits the load balancer.
// load balancer is supposed to change it destinationIP to the API server's
// IP where this request should go
type APIRequest struct {
	Name          string
	SourceIP      string
	DestinationIP string
}

func NewAPIRequest(srcIP, destIP, name string) *APIRequest {
	return &APIRequest{
		SourceIP:      srcIP,
		DestinationIP: destIP,
		Name:          name,
	}
}
