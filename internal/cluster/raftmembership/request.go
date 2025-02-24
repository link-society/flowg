package raftmembership

import (
	"time"

	"github.com/hashicorp/raft"
)

type ChangeRequest struct {
	id   raft.ServerID
	kind changeRequestKind
	done chan error
}

type changeRequestKind interface {
	Handle(raft *raft.Raft, id raft.ServerID) error
}

type statusRequest struct {
	replyTo chan<- []ServerInfo
}

type joinRequest struct {
	addr    raft.ServerAddress
	timeout time.Duration
}

type leaveRequest struct {
	timeout time.Duration
}

func NewStatusRequest(replyTo chan<- []ServerInfo) *ChangeRequest {
	return &ChangeRequest{
		kind: &statusRequest{replyTo: replyTo},
		done: make(chan error, 1),
	}
}

func NewJoinRequest(id raft.ServerID, addr raft.ServerAddress, timeout time.Duration) *ChangeRequest {
	return &ChangeRequest{
		id:   id,
		kind: &joinRequest{addr: addr, timeout: timeout},
		done: make(chan error, 1),
	}
}

func NewLeaveRequest(id raft.ServerID, timeout time.Duration) *ChangeRequest {
	return &ChangeRequest{
		id:   id,
		kind: &leaveRequest{timeout: timeout},
		done: make(chan error, 1),
	}
}

func (r *ChangeRequest) notifyDone(err error) {
	r.done <- err
	close(r.done)
}

func (r *statusRequest) Handle(ra *raft.Raft, id raft.ServerID) error {
	defer close(r.replyTo)

	if ra.State() != raft.Leader {
		return raft.ErrNotLeader
	}

	config := ra.GetConfiguration().Configuration()

	servers := make([]ServerInfo, 0, len(config.Servers))
	for _, server := range config.Servers {
		servers = append(servers, ServerInfo{
			ID:      string(server.ID),
			Address: string(server.Address),
		})
	}

	r.replyTo <- servers

	return nil
}

func (r *joinRequest) Handle(ra *raft.Raft, id raft.ServerID) error {
	if ra.State() != raft.Leader {
		return raft.ErrNotLeader
	}

	return ra.AddVoter(id, r.addr, 0, r.timeout).Error()
}

func (r *leaveRequest) Handle(ra *raft.Raft, id raft.ServerID) error {
	if ra.State() != raft.Leader {
		return raft.ErrNotLeader
	}

	return ra.RemoveServer(id, 0, r.timeout).Error()
}
