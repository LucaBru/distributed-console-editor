package raft

import (
	"context"

	"github.com/mb0/lab/ot"

	pb "dconsole-editor/raft/protos"
)

type entry struct {
	term int
	// should be optional, need to send heartbeat
	cmd ot.Op
}

type NodeModule struct {
	Log []entry
	pb.UnimplementedRaftServer
}

func (s *NodeModule) RequestVote(_ context.Context, in *pb.VoteRequest) (*pb.VoteReply, error) {
	//build the reply starting from the request
	return &pb.VoteReply{Term: 10, Ok: false}, nil
}

func (s *NodeModule) AppendEntry(_ context.Context, in *pb.AppendEntryRequest) (*pb.AppendEntryReply, error) {
	//build the reply starting from the request
	return &pb.AppendEntryReply{Term: 10, Ok: false}, nil
}
