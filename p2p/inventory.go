package p2p

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/erxonxi/coin/blockchain"

	pb "github.com/erxonxi/coin/p2p/p2p"
	proto "github.com/gogo/protobuf/proto"
	uuid "github.com/google/uuid"
)

const inventoryRequest = "/inventory/request/0.0.1"
const inventoryResponse = "/inventory/response/0.0.1"

type InventoryProtocol struct {
	node      *Node                           // local host
	requests  map[string]*pb.InventoryRequest // used to access request data from response handlers
	inventory chan []byte
	done      chan bool
}

func NewInvenotryProtocol(node *Node, done chan bool) *InventoryProtocol {
	p := &InventoryProtocol{node: node, requests: make(map[string]*pb.InventoryRequest), done: done}
	node.SetStreamHandler(inventoryRequest, p.onInventoryRequest)
	node.SetStreamHandler(inventoryResponse, p.onInventoryResponse)
	return p
}

// remote peer requests handler
func (p *InventoryProtocol) onInventoryRequest(s network.Stream) {
	data := &pb.InventoryRequest{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	log.Println("ERROR 1")
	err = proto.Unmarshal(buf, data)
	log.Println("ERROR 2")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("ERROR 3")

	log.Printf("Received inventory request from %s.\n", s.Conn().RemotePeer())

	valid := p.node.authenticateMessage(data, data.MessageData)

	if !valid {
		log.Println("Failed to authenticate message")
		return
	}

	// generate response message
	resp := &pb.InventoryResponse{
		MessageData: p.node.NewMessageData(data.MessageData.Id, false),
		Inventory:   data.Inventory,
	}

	// sign the data
	signature, err := p.node.signProtoMessage(resp)
	if err != nil {
		log.Println("failed to sign response")
		return
	}

	// add the signature to the message
	resp.MessageData.Sign = signature

	// send the response
	ok := p.node.sendProtoMessage(s.Conn().RemotePeer(), inventoryResponse, resp)

	if ok {
		log.Println("Inventory Response OK")
	}

	p.done <- true
}

// remote ping response handler
func (p *InventoryProtocol) onInventoryResponse(s network.Stream) {
	data := &pb.InventoryResponse{}
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	// unmarshal it
	log.Println("Response A")
	err = proto.Unmarshal(buf, data)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Response B")

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

	log.Printf("%s:\nReceived ping response from %s.\nMessage id:%s.\n", s.Conn().LocalPeer(), s.Conn().RemotePeer(), data.MessageData.Id)
	p.done <- true
}

func (p *InventoryProtocol) sendInventory(hostId peer.ID, inventory []byte) bool {
	req := &pb.InventoryRequest{
		MessageData: p.node.NewMessageData(uuid.New().String(), false),
		Inventory:   inventory,
	}

	signature, err := p.node.signProtoMessage(req)
	if err != nil {
		log.Println("failed to sign pb data")
		log.Println(err)
		return false
	}

	req.MessageData.Sign = signature
	ok := p.node.sendProtoMessage(hostId, inventoryRequest, req)
	if !ok {
		return false
	}

	p.requests[req.MessageData.Id] = req
	log.Printf("Sending inventory from %s to %s", p.node.ID(), hostId)
	return true
}

func (p *InventoryProtocol) SendBlock(hostId peer.ID, block *blockchain.Block) bool {
	return p.sendInventory(hostId, GobEncode(block.Serialize()))
}

func (p *InventoryProtocol) SendTransaction(hostId peer.ID, tx *blockchain.Transaction) bool {
	return p.sendInventory(hostId, GobEncode(tx.Serialize()))
}

func GobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
