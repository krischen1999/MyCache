package MyCache

import pb "GeeCache/geecachepb"

// 本地需实现peerpicker
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 远程节点需实现peerGetter
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
