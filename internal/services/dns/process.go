package dns

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/vladopajic/go-actor/actor"
	"link-society.com/flowg/internal/utils/proctree"

	"github.com/miekg/dns"
)

type procHandler struct {
	client *dns.Client
	logger *slog.Logger
	opts   *DnsServiceOptions
	fqdn   string

	// Type Protocol FQDN
	tpfqdn string

	LocalEndpointResolver func() (*url.URL, string, error)
}

const (
	getNodesMaxRetries  = 10
	healthCheckInterval = 5 * time.Second
	healthCheckTimeout  = 1 * time.Second
	shutdownTimeout     = 5 * time.Second
)

func (h *procHandler) Init(ctx actor.Context) proctree.ProcessResult {
	// If no dns server address is provided then stop the dns service as it is not needed
	if h.opts.DnsServerAddress == "" {
		h.logger.InfoContext(ctx, "no dns server address provided")
		return proctree.Continue()
	}

	h.client = new(dns.Client)
	h.client.DialTimeout = 5 * time.Second
	h.client.ReadTimeout = 5 * time.Second
	h.client.WriteTimeout = 5 * time.Second

	// Register node with DNS server
	if err := h.registerNode(ctx); err != nil {
		h.logger.ErrorContext(
			ctx,
			"failed to register node with dns server",
			slog.Any("error", err),
		)
		return proctree.Terminate(err)
	}

	// Set the JoinNode in ClusterJoinNode
	err := h.setJoinNodes(ctx)
	if err != nil {
		/* Log the error but don't terminate the process
		because the first node that starts up in the cluster
		will never find any other nodes */
		h.logger.WarnContext(
			ctx,
			"failed to get service nodes from dns server",
			slog.Any("error", err),
		)
	}

	return proctree.Continue()
}

func (h *procHandler) DoWork(ctx actor.Context) proctree.ProcessResult {
	<-ctx.Done()
	return proctree.Terminate(ctx.Err())
}

func (h *procHandler) Terminate(ctx actor.Context, err error) error {
	if h.opts.DnsServerAddress == "" {
		return err
	}

	h.logger.InfoContext(ctx, "Deregistering service with dns server")

	var recordsToDelete []dns.RR

	m := new(dns.Msg)
	// Set the zone to update
	m.SetUpdate(h.opts.DomainName)

	// Construct the A record for deletion.
	aRec, err := dns.NewRR(fmt.Sprintf("%s A 0.0.0.0", h.fqdn))
	if err != nil {
		return fmt.Errorf("failed to create A record for deletion: %w", err)
	}

	aRec.Header().Rrtype = dns.TypeA
	aRec.Header().Class = dns.ClassANY // ClassANY indicates deletion
	aRec.Header().Ttl = 0              // TTL 0 indicates deletion
	recordsToDelete = append(recordsToDelete, aRec)

	// Construct the SRV record for deletion.
	srvRec, err := dns.NewRR(fmt.Sprintf("%s SRV 0 0 0 .", h.tpfqdn))
	if err != nil {
		return fmt.Errorf("failed to create SRV record for deletion: %w", err)
	}
	srvRec.Header().Rrtype = dns.TypeSRV
	srvRec.Header().Class = dns.ClassANY // ClassANY indicates deletion
	srvRec.Header().Ttl = 0              // TTL 0 indicates deletion
	recordsToDelete = append(recordsToDelete, srvRec)

	m.Insert(recordsToDelete)

	_, err = h.exchange(m)
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"DNS update failed",
			slog.Any("error", err),
		)
		return err
	}

	return err
}

func (h *procHandler) registerNode(ctx actor.Context) error {
	localEndpoint, ip, err := h.LocalEndpointResolver() // do we then really need this here?
	if err != nil {
		return err
	}

	var mgmtPortString string
	_, mgmtPortString, err = net.SplitHostPort(h.opts.MgmtBindAddress)
	if err != nil {
		h.logger.ErrorContext(
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
	m.SetUpdate(fmt.Sprintf("%s.", h.opts.DomainName))

	// Construct the FQDN for the service
	// A trailing dot is needed for FQDN.
	h.fqdn = fmt.Sprintf("%s.%s.", h.opts.ServiceName, h.opts.DomainName)

	// Add an A record for the service
	// This record maps the service name to its IP address
	aRec, err := dns.NewRR(fmt.Sprintf("%s A %s", h.fqdn, ip))
	if err != nil {
		h.logger.ErrorContext(
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
	h.tpfqdn = fmt.Sprintf("_%s._tcp.%s", h.opts.ServiceName, h.opts.DomainName)

	// Add an SRV record for the service
	// SRV records specify the host and port for specific services.
	srvRec, err := dns.NewRR(fmt.Sprintf("%s 85400 IN SRV 10 10 %s %s", h.tpfqdn, mgmtPortString, localEndpoint.Hostname()))
	if err != nil {
		h.logger.ErrorContext(
			ctx,
			"Failed to create SRV record",
			slog.Any("error", err),
		)
		return err
	}
	m.Insert([]dns.RR{srvRec})
	recordsToInsert = append(recordsToInsert, srvRec)

	record1 := fmt.Sprintf("%+v", recordsToInsert[0])
	record2 := fmt.Sprintf("%+v", recordsToInsert[1])

	h.logger.InfoContext(ctx, record1)
	h.logger.InfoContext(ctx, record2)

	_, err = h.exchange(m)
	if err != nil {
		h.logger.ErrorContext(
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
func (h *procHandler) setJoinNodes(ctx actor.Context) error {
	retryCount := 0
	delay := 100 * time.Millisecond

	//serviceName := "_my-service._tcp.example.com." // IMPORTANT: Replace with your actual service name
	// genai suggests that the type is not the protocol

	// i havent add to exclude myself when i get the list of nodes back fomr the dns server

	var scheme string
	if h.opts.MgmtTlsEnabled {
		scheme = "https"
	} else {
		scheme = "https"
	}

	h.logger.InfoContext(ctx, "GOing to ask for other nodes")

	m := new(dns.Msg)
	// Set the question to query for SRV records
	m.SetQuestion(dns.Fqdn(h.tpfqdn), dns.TypeSRV)
	m.RecursionDesired = true

	for retryCount <= getNodesMaxRetries {
		r, err := h.exchange(m)
		if err != nil {
			h.logger.ErrorContext(
				ctx,
				"failed to get nodes from dns server",
				slog.Any("error", err),
			)
			return err
		}

		for _, a := range r.Answer {
			if srv, ok := a.(*dns.SRV); ok {
				h.opts.ClusterJoinNode.JoinNodeEndpoint = &url.URL{
					Scheme: scheme,
					Host:   net.JoinHostPort(srv.Target, strconv.FormatUint(uint64(srv.Port), 10)),
				}
				return nil
			}
		}

		retryCount++
		if retryCount <= getNodesMaxRetries {
			h.logger.InfoContext(ctx, "did not find other nodes, will try again with a delay")
			time.Sleep(delay)
			// Add jitter to the delay
			delay += time.Duration(rand.IntN(int(delay / 4)))
		}
	}

	return fmt.Errorf("failed to find other nodes")
}

func (h *procHandler) exchange(m *dns.Msg) (r *dns.Msg, err error) {
	r, _, err = h.client.Exchange(m, h.opts.DnsServerAddress)

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
