package cluster

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"
	"time"

	"github.com/hashicorp/go-sockaddr"
	"github.com/miekg/dns"
	"github.com/vladopajic/go-actor/actor"
)

type DnsClusterFormationStrategy struct {
	NodeID           string
	ServiceName      string
	DnsServerAddress string
	DnsDomainName    string
	ServiceTls       bool
	ServiceAddress   string

	fqdn string

	// Type Protocol FQDN
	tpfqdn string

	logger *slog.Logger
	client *dns.Client
}

const (
	getNodesMaxRetries  = 10
	healthCheckInterval = 5 * time.Second
	healthCheckTimeout  = 1 * time.Second
	shutdownTimeout     = 5 * time.Second
)

// Join implements ClusterFormationStrategy.
func (d *DnsClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {

	d.logger = slog.Default().With(slog.String("channel", "cluster.dns"))

	d.logger.InfoContext(ctx, "Initiating registering with a dns server")

	d.client = new(dns.Client)
	d.client.DialTimeout = 5 * time.Second
	d.client.ReadTimeout = 5 * time.Second
	d.client.WriteTimeout = 5 * time.Second

	if err := d.registerNode(ctx); err != nil {
		d.logger.ErrorContext(
			ctx,
			"failed to register node with dns server",
			slog.Any("error", err),
		)
	}

	// Set the JoinNode in ClusterJoinNode
	err := d.setJoinNodes()
	if err != nil {
		/* Log the error but don't terminate the process
		because the first node that starts up in the cluster
		will never find any other nodes */
		d.logger.WarnContext(
			ctx,
			"failed to get service nodes from dns server",
			slog.Any("error", err),
		)
	}

	return nil, nil

}

// Leave implements ClusterFormationStrategy.
func (d *DnsClusterFormationStrategy) Leave(ctx context.Context) error {
	//panic("unimplemented")
	return nil
}

func (d *DnsClusterFormationStrategy) registerNode(ctx actor.Context) error {
	localEndpoint, ip, err := d.localEndpointResolver() // do we then really need this here?
	if err != nil {
		return err
	}

	var mgmtPortString string
	_, mgmtPortString, err = net.SplitHostPort(d.ServiceAddress)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"failed to split manegemnt bind address",
			slog.Any("error", err),
		)
		return err
	}

	var recordsToInsert []dns.RR

	// Create a new DNS update message
	m := new(dns.Msg)

	// Set the zone to update
	m.SetUpdate(fmt.Sprintf("%s.", d.DnsDomainName))

	// Construct the FQDN for the service
	// A trailing dot is needed for FQDN.
	d.fqdn = fmt.Sprintf("%s.%s.", d.ServiceName, d.DnsDomainName)
	d.logger.InfoContext(ctx, fmt.Sprintf("Fully qualified domain name fqdn: %s", d.fqdn))

	aRecString := fmt.Sprintf("%s A %s", d.fqdn, ip)
	d.logger.InfoContext(ctx, fmt.Sprintf("A Record: %s", aRecString))

	// Add an A record for the service
	// This record maps the service name to its IP address
	aRec, err := dns.NewRR(aRecString)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"Failed to create A record",
			slog.Any("error", err),
		)
		return err
	}
	m.Insert([]dns.RR{aRec}) // Use Insert to add records
	recordsToInsert = append(recordsToInsert, aRec)
	// do we really need A records? Isnt target from teh SRV record enough?

	// target has to not just be the domain name, but should also have the node host
	// because otherwise every node will have the same domain so you cant use the domian name to find other nodes
	// h.tpfqdn = fmt.Sprintf("_%s._tcp.%s", h.opts.ServiceName, h.opts.DomainName)
	d.tpfqdn = fmt.Sprintf("_%s._tcp.%s", d.ServiceName, d.DnsDomainName)
	d.logger.InfoContext(ctx, fmt.Sprintf("TPFQDN: %s", d.tpfqdn))

	// Add an SRV record for the service
	// SRV records specify the host and port for specific services.
	srvRec, err := dns.NewRR(fmt.Sprintf("%s 85400 IN SRV 10 10 %s %s", d.tpfqdn, mgmtPortString, localEndpoint.Hostname()))
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"Failed to create SRV record",
			slog.Any("error", err),
		)
		return err
	}
	m.Insert([]dns.RR{srvRec})
	// Trying only to insrt the 1st reord and skipping the 2nd one to see if i get the same domain name must be fully qualified error
	// recordsToInsert = append(recordsToInsert, srvRec)

	record1 := fmt.Sprintf("%+v", recordsToInsert[0])
	// record2 := fmt.Sprintf("%+v", recordsToInsert[1])

	d.logger.InfoContext(ctx, record1)
	// h.logger.InfoContext(ctx, record2)

	_, err = d.exchange(m)
	if err != nil {
		d.logger.ErrorContext(
			ctx,
			"DNS update failed",
			slog.Any("error", err),
		)
		return err
	}

	return nil
}

// setJoinNodes() retries with exponential backoff with jitter to fetch other nodes in the cluster using consul
// and sets one node as a JoinNode in the ClusterJoinNode
// ClusterJoinNode is shared between ConsulService and ManagementServer
func (d *DnsClusterFormationStrategy) setJoinNodes() error {

	// Registering is not working therefore commenting out below code.

	/* retryCount := 0
	delay := 100 * time.Millisecond

	//serviceName := "_my-service._tcp.example.com." // IMPORTANT: Replace with your actual service name
	// genai suggests that the type is not the protocol

	// i havent add to exclude myself when i get the list of nodes back fomr the dns server

	var scheme string
	if d.ServiceTls {
		scheme = "https"
	} else {
		scheme = "http"
	}

	d.logger.InfoContext(ctx, "Going to ask for other nodes")

	m := new(dns.Msg)
	// Set the question to query for SRV records
	m.SetQuestion(dns.Fqdn(d.tpfqdn), dns.TypeSRV)
	m.RecursionDesired = true

	for retryCount <= getNodesMaxRetries {
		r, err := d.exchange(m)
		if err != nil {
			d.logger.ErrorContext(
				ctx,
				"failed to get nodes from dns server",
				slog.Any("error", err),
			)
			return err
		}

		for _, a := range r.Answer {
			if srv, ok := a.(*dns.SRV); ok {
				ClusterJoinNode.JoinNodeEndpoint = &url.URL{
					Scheme: scheme,
					Host:   net.JoinHostPort(srv.Target, strconv.FormatUint(uint64(srv.Port), 10)),
				}
				return nil
			}
		}

		retryCount++
		if retryCount <= getNodesMaxRetries {
			d.logger.InfoContext(ctx, "did not find other nodes, will try again with a delay")
			time.Sleep(delay)
			// Add jitter to the delay
			delay += time.Duration(rand.IntN(int(delay / 4)))
		}
	}

	return fmt.Errorf("failed to find other nodes") */

	return nil
}

func (d *DnsClusterFormationStrategy) exchange(m *dns.Msg) (r *dns.Msg, err error) {
	r, _, err = d.client.Exchange(m, d.DnsServerAddress)

	if err != nil {
		return nil, err // fmt.Errorf("DNS query failed")
	}

	if r == nil {
		return nil, fmt.Errorf("no response received from DNS server")
	}

	if r.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("DNS update failed with RCODE: %s (%d)", dns.RcodeToString[r.Rcode], r.Rcode)
	}

	return r, nil
}

func (d *DnsClusterFormationStrategy) localEndpointResolver() (*url.URL, string, error) {
	host, port, err := net.SplitHostPort(d.ServiceAddress)
	if err != nil {
		return nil, "", fmt.Errorf("failed to bind address: %w", err)
	}

	ip, err := sockaddr.GetPrivateIP()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get private IP: %w", err)
	}
	if ip == "" {
		return nil, "", fmt.Errorf("no private IP found")
	}

	if host == "0.0.0.0" || host == "::" {
		host = ip
	}

	if len(host) == 0 {
		host = "localhost"
	}

	var localEndpoint url.URL
	if d.ServiceTls {
		localEndpoint = url.URL{
			Scheme: "https",
			Host:   net.JoinHostPort(host, port),
		}
	} else {
		localEndpoint = url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(host, port),
		}
	}

	return &localEndpoint, ip, nil
}
