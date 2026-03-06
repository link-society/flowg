package cluster

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"net/url"

	"codeberg.org/miekg/dns"
)

type DnsClusterFormationStrategy struct {
	DnsServer string
	DnsDomain string
	DnsScript string
	NodeID    string
}

var _ ClusterFormationStrategy = (*DnsClusterFormationStrategy)(nil)

const FLOWG_CLUSTER = "flowg-cluster="

func (s *DnsClusterFormationStrategy) Join(ctx context.Context, resolver LocalEndpointResolverCallback) ([]*ClusterJoinNode, error) {
	msg := dns.NewMsg(s.DnsDomain, dns.TypeTXT)
	client := new(dns.Client)

	r, _, err := client.Exchange(ctx, msg, "udp", s.DnsServer)
	if err != nil {
		return nil, err
	}

	var nodes []*ClusterJoinNode

	localEndpoint, err := resolver()
	if err != nil {
		return nil, err
	}

	alreadyAdded := false
	for _, answer := range r.Answer {
		if data, ok := answer.(*dns.TXT); ok {
			txt := data.Txt[0]
			if !strings.HasPrefix(txt, FLOWG_CLUSTER) {
				continue
			}

			strs := strings.Split(txt[len(FLOWG_CLUSTER):], ";")
			if len(strs) != 2 {
				return nil, errors.New("incorrect dns cluster format")
			}

			endpoint, err := url.Parse(strings.TrimSpace(strs[1]))
			if err != nil {
				return nil, err
			}

			node := &ClusterJoinNode{
				JoinNodeID:       strings.TrimSpace(strs[0]),
				JoinNodeEndpoint: endpoint,
			}

			if node.JoinNodeID == s.NodeID {
				alreadyAdded = true
				continue
			}

			fmt.Printf("TXT: %s\n", txt)

			nodes = append(nodes, node)
		}
	}

	if len(s.DnsScript) > 0 && !alreadyAdded {
		err := exec.Command(s.DnsScript, "set", "TXT", s.DnsDomain, fmt.Sprintf("%s%s;%s", FLOWG_CLUSTER, s.NodeID, localEndpoint)).Run()
		if err != nil {
			return nil, err
		}
	}

	return nodes, nil
}

func (s *DnsClusterFormationStrategy) Leave(ctx context.Context) error {
	if len(s.DnsScript) > 0 {
		return exec.Command(s.DnsScript, "del", "TXT", s.DnsDomain, fmt.Sprintf("%s%s", FLOWG_CLUSTER, s.NodeID)).Run()
	}

	return nil
}
