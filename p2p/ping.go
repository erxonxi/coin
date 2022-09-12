package p2p

import (
	"fmt"
	"io"
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	pb "github.com/erxonxi/coin/p2p/p2p"
	proto "github.com/gogo/protobuf/proto"
	uuid "github.com/google/uuid"
)

// pattern: /protocol-name/request-or-response-message/version
const pingRequest = "/ping/request/0.0.1"
const pingResponse = "/ping/response/0.0.1"

// PingProtocol type
type PingProtocol struct {
	node     *Node                      // local host
	requests map[string]*pb.PingRequest // used to access request data from response handlers
	done     chan bool                  // only for demo purposes to stop main from terminating
}

func NewPingProtocol(node *Node, done chan bool) *PingProtocol {
	p := &PingProtocol{node: node, requests: make(map[string]*pb.PingRequest), done: done}
	node.SetStreamHandler(pingRequest, p.onPingRequest)
	node.SetStreamHandler(pingResponse, p.onPingResponse)
	return p
}

// remote peer requests handler
func (p *PingProtocol) onPingRequest(s network.Stream) {

	// get request data
	data := &pb.PingRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("%s:\nReceived ping request from %s.\nMessage:\n%s\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.Message)

	valid := p.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// generate response message
	log.Printf("%s:\nSending ping response to %s.\nMessage id: %s...\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)

	resp := &pb.PingResponse{MessageData: p.node.NewMessageData(data.MessageData.Id, false),
		Message: fmt.Sprintf("Ping response from %s", p.node.ID())}

	// sign the data
	signature, err := p.node.signProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	// send the response
	ok := p.node.sendProtoMessage(s.Conn().RemotePeer(), pingResponse, resp)

	if ok {
		log.Printf("%s:\nPing response to %s sent\n", s.Conn().LocalPeer().String(), s.Conn().RemotePeer().String())
	}
	p.done <- true
}

// remote ping response handler
func (p *PingProtocol) onPingResponse(s network.Stream) {
	data := &pb.PingResponse{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}

	valid := p.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// locate request data and remove it if found
	_, ok := p.requests[data.MessageData.Id]
	if ok {
		// remove request from map as we have processed it here
		delete(p.requests, data.MessageData.Id)
	} else {
		log.Println("Failed to locate request data boject for response")
		return
	}

	log.Printf("%s:\nReceived ping response from %s.\nMessage id:%s.\nMessage:\n%s.\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id, data.Message)
	p.done <- true
}

func (p *PingProtocol) Ping(hostId peer.ID) bool {
	log.Printf("%s: Sending ping to: %s....", p.node.ID(), hostId)

	// create message data
	req := &pb.PingRequest{MessageData: p.node.NewMessageData(uuid.New().String(), false),
		Message: fmt.Sprintf("Ping from %s", p.node.ID())}

	// sign the data
	signature, err := p.node.signProtoMessage(req)
	if err != nil {
		log.Println("failed to sign pb data")
		return false
	}

	// add the signature to the message
	req.MessageData.Sign = signature

	ok := p.node.sendProtoMessage(hostId, pingRequest, req)
	if !ok {
		return false
	}

	// store ref request so response handler has access to it
	p.requests[req.MessageData.Id] = req
	log.Printf("%s: Ping to: %s was sent. Message Id: %s, Message: %s", p.node.ID(), hostId, req.MessageData.Id, req.Message)
	return true
}
