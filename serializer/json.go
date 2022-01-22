package serializer

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func ProtoToJSON(message proto.Message) ([]byte, error) {
	marshaller := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
		Indent:          "   ",
	}
	return marshaller.Marshal(message)
}
