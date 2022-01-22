package serializer

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

func WriteProtobufToFile(message proto.Message, fileName string) {
	data, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("Error while writing proto to file", err.Error())
	} else {
		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			fmt.Println("Error while writing proto to file", err.Error())
		}
	}
}

func ReadBinaryToProtobuf(fileName string, message proto.Message) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error while reading file", err.Error())
	} else {
		err = proto.Unmarshal(data, message)
		if err != nil {
			fmt.Println("Error while unmarshalling file to proto", err.Error())
		}
	}
}

func WriteProtobufToJson(message proto.Message, fileName string) {
	data, err := ProtoToJSON(message)
	if err != nil {
		fmt.Println("Error while writing proto to json", err.Error())
	} else {
		err = ioutil.WriteFile(fileName, data, 0644)
		if err != nil {
			fmt.Println("Error while writing proto to file", err.Error())
		}
	}
}
