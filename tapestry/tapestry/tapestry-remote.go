package tapestry

import (
	"net/rpc"
	"fmt"
)

/*
	The methods defined in this file parallel the methods defined in tapestry-local.
	These methods take an additional argument, the node on which the method should be invoked.
	Calling any of these methods will invoke the corresponding method on the specified remote node.
*/

// Remote API: ping an address to get tapestry node info
func (tapestry *Tapestry) hello(address string) (rsp Node, err error) {
	err = makeRemoteCall(address, "TapestryRPCServer", "Hello", tapestry.local.node, &rsp)
	return
}

// Helper function to makes a remote call
func makeRemoteNodeCall(remote Node, method string, req interface{}, rsp interface{}) error {
	fmt.Printf("%v(%v)\n", method, req)
	return makeRemoteCall(remote.Address, "TapestryRPCServer", method, req, rsp)
}

// Helper function to makes a remote call
func makeRemoteCall(address string, structtype string, method string, req interface{}, rsp interface{}) error {
	// Dial the server
	client, err := rpc.Dial("tcp", address)
	if err != nil {
		return err
	}

	// Make the request
	fqm := fmt.Sprintf("%v.%v", structtype, method)
	err = client.Call(fqm, req, rsp)

	client.Close()
	if err != nil {
		return err
	}

	return nil
}

// Remote API: makes a remote call to the Register function
func (tapestry *Tapestry) register(remote Node, replica Node, key string) (bool, error) {
	var rsp RegisterResponse
	err := makeRemoteNodeCall(remote, "Register", RegisterRequest{remote, replica, key}, &rsp)
	return rsp.IsRoot, err
}

// Remote API: makes a remote call to the GetNextHop function
func (tapestry *Tapestry) getNextHop(remote Node, id ID) (bool, Node, error) {
	var rsp NextHopResponse
	err := makeRemoteNodeCall(remote, "GetNextHop", NextHopRequest{remote, id}, &rsp)
	return rsp.HasNext, rsp.Next, err
}

// Remote API: makes a remote call to the RemoveBadNodes function
func (tapestry *Tapestry) removeBadNodes(remote Node, toremove []Node) error {
	return makeRemoteNodeCall(remote, "RemoveBadNodes", RemoveBadNodesRequest{remote, toremove}, &Node{})
}