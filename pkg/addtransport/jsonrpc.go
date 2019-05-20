package addtransport

import (
	"context"
	"encoding/json"
	"fmt"
	"goexamples/pkg/addendpoint"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport/http/jsonrpc"
)

// NewJSONRPCHandler returns a JSON RPC Server/Handler that can be passed to http.Handle().
func NewJSONRPCHandler(endpoints addendpoint.Set, logger log.Logger) *jsonrpc.Server {
	handler := jsonrpc.NewServer(
		makeEndpointCodecMap(endpoints),
		jsonrpc.ServerErrorLogger(logger),
	)
	return handler
}

func makeEndpointCodecMap(endpoints addendpoint.Set) jsonrpc.EndpointCodecMap {
	return jsonrpc.EndpointCodecMap{
		"sum": jsonrpc.EndpointCodec{
			Endpoint: endpoints.SumEndpoint,
			Decode:   decodeSumRequest,
			Encode:   encodeSumResponse,
		},
		"concat": jsonrpc.EndpointCodec{
			Endpoint: endpoints.ConcatEndpoint,
			Decode:   decodeConcatRequest,
			Encode:   encodeConcatResponse,
		},
	}
}

func decodeSumRequest(_ context.Context, msg json.RawMessage) (interface{}, error) {
	var req addendpoint.SumRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    -3200,
			Message: fmt.Sprintf("couldn't unmarshal body to sum request: %s", err),
		}
	}
	return req, nil
}

func encodeSumResponse(_ context.Context, obj interface{}) (json.RawMessage, error) {
	res, ok := obj.(addendpoint.SumResponse)
	if !ok {
		return nil, &jsonrpc.Error{
			Code:    -32000,
			Message: fmt.Sprintf("Asserting result to *SumResponse failed. Got %T, %+v", obj, obj),
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func decodeConcatRequest(_ context.Context, msg json.RawMessage) (interface{}, error) {
	var req addendpoint.ConcatRequest
	err := json.Unmarshal(msg, &req)
	if err != nil {
		return nil, &jsonrpc.Error{
			Code:    -32000,
			Message: fmt.Sprintf("couldn't unmarshal body to concat request: %s", err),
		}
	}
	return req, nil
}

func encodeConcatResponse(_ context.Context, obj interface{}) (json.RawMessage, error) {
	res, ok := obj.(addendpoint.ConcatResponse)
	if !ok {
		return nil, &jsonrpc.Error{
			Code:    -32000,
			Message: fmt.Sprintf("Asserting result to *ConcatResponse failed. Got %T, %+v", obj, obj),
		}
	}
	b, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("couldn't marshal response: %s", err)
	}
	return b, nil
}
