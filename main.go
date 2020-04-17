package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func makeFileDescriptor() pref.FileDescriptor {
	// make FileDescriptorProto
	pb := &descriptorpb.FileDescriptorProto{
		Syntax:  proto.String("proto3"),
		Name:    proto.String("example.proto"),
		Package: proto.String("example"),
		MessageType: []*descriptorpb.DescriptorProto{
			// define Foo message
			&descriptorpb.DescriptorProto{
				Name: proto.String("Foo"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("id"),
						JsonName: proto.String("id"),
						Number:   proto.Int32(1),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.Int32Kind).Enum(),
					},
					{
						Name:     proto.String("title"),
						JsonName: proto.String("title"),
						Number:   proto.Int32(2),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.StringKind).Enum(),
					},
				},
			},

			// define Bar message
			&descriptorpb.DescriptorProto{
				Name: proto.String("Bar"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("bar_map"),
						JsonName: proto.String("bar_map"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum(),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.MessageKind).Enum(),
						TypeName: proto.String(".example.Bar.BarMapEntry"),
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name: proto.String("BarMapEntry"),
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:     proto.String("key"),
								JsonName: proto.String("key"),
								Number:   proto.Int32(1),
								Label:    descriptorpb.FieldDescriptorProto_Label(pref.Optional).Enum(),
								Type:     descriptorpb.FieldDescriptorProto_Type(pref.StringKind).Enum(),
							}, {
								Name:     proto.String("value"),
								JsonName: proto.String("value"),
								Number:   proto.Int32(2),
								Label:    descriptorpb.FieldDescriptorProto_Label(pref.Optional).Enum(),
								Type:     descriptorpb.FieldDescriptorProto_Type(pref.MessageKind).Enum(),
								TypeName: proto.String(".example.Foo"),
							},
						},
						Options: &descriptorpb.MessageOptions{
							MapEntry: proto.Bool(true),
						},
					},
				},
			},

			// define Baz message
			&descriptorpb.DescriptorProto{
				Name: proto.String("Baz"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:     proto.String("baz_list"),
						JsonName: proto.String("baz_list"),
						Number:   proto.Int32(1),
						Label:    descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum(),
						Type:     descriptorpb.FieldDescriptorProto_Type(pref.MessageKind).Enum(),
						TypeName: proto.String(".example.Foo"),
					},
				},
			},
		},
	}

	// get FileDescriptor
	fd, err := protodesc.NewFile(pb, nil)
	if err != nil {
		panic(err)
	}
	return fd
}

func makeFooMsg(fd pref.FileDescriptor) *dynamicpb.Message {
	fooMessageDescriptor := fd.Messages().ByName("Foo")
	msg := dynamicpb.NewMessage(fooMessageDescriptor)
	msg.Set(fooMessageDescriptor.Fields().ByName("id"), pref.ValueOfInt32(42))
	msg.Set(fooMessageDescriptor.Fields().ByNumber(2), pref.ValueOfString("aloha"))
	return msg
}

func makeBarMsg(fd pref.FileDescriptor) *dynamicpb.Message {
	barMessageDescriptor := fd.Messages().ByName("Bar")
	msg := dynamicpb.NewMessage(barMessageDescriptor)
	mf := barMessageDescriptor.Fields().ByName("bar_map")
	mp := msg.NewField(mf)

	fooMsg := makeFooMsg(fd)

	mp.Map().Set(pref.MapKey(pref.ValueOfString("key1")), pref.ValueOfMessage(fooMsg))
	mp.Map().Set(pref.MapKey(pref.ValueOfString("key2")), pref.ValueOfMessage(fooMsg))
	msg.Set(mf, mp)
	return msg
}

func makeBazMsg(fd pref.FileDescriptor) *dynamicpb.Message {
	bazMessageDescriptor := fd.Messages().ByName("Baz")
	msg := dynamicpb.NewMessage(bazMessageDescriptor)
	lf := bazMessageDescriptor.Fields().ByName("baz_list")
	fooMsg := makeFooMsg(fd)
	lst := msg.NewField(lf).List()
	lst.Append(pref.ValueOf(fooMsg))
	lst.Append(pref.ValueOf(fooMsg))
	lst.Append(pref.ValueOf(fooMsg))
	msg.Set(lf, pref.ValueOf(lst))
	return msg
}

func useFooMsg(fd pref.FileDescriptor, data []byte) {
	fooMessageDescriptor := fd.Messages().ByName("Foo")
	msg := dynamicpb.NewMessage(fooMessageDescriptor)
	if err := proto.Unmarshal(data, msg); err != nil {
		panic(err)
	}

	// iterate over all fields
	msg.Range(func(descriptor pref.FieldDescriptor, value pref.Value) bool {
		fmt.Printf("field: %v value: %v \n", descriptor.Name(), value)
		return true
	})

	// get single field's value
	v := msg.Get(fooMessageDescriptor.Fields().ByName("id"))
	fmt.Printf("get %v \n", v)
}

func useBarMsg(fd pref.FileDescriptor, data []byte) {
	barMessageDescriptor := fd.Messages().ByName("Bar")
	msg := dynamicpb.NewMessage(barMessageDescriptor)
	if err := proto.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	mp := msg.Get(barMessageDescriptor.Fields().ByName("bar_map")).Map()

	// iterate over map field
	mp.Range(func(key pref.MapKey, value pref.Value) bool {
		fmt.Printf("key: %v value: %v  \n", key.String(), value.Message())
		return true
	})
}

func useBazMsg(fd pref.FileDescriptor, data []byte) {
	bazMessageDescriptor := fd.Messages().ByName("Baz")
	msg := dynamicpb.NewMessage(bazMessageDescriptor)
	if err := proto.Unmarshal(data, msg); err != nil {
		panic(err)
	}
	lf := bazMessageDescriptor.Fields().ByName("baz_list")
	lst := msg.Get(lf).List()
	length := lst.Len()
	for i := 0; i < length; i++ {
		ele := lst.Get(i)
		fmt.Printf("index: %v value: %v  \n", i, ele.Message())
	}
}

func marshalMsg(msg *dynamicpb.Message) []byte {
	var (
		data []byte
		err  error
	)
	if data, err = proto.Marshal(msg); err != nil {
		panic(err)
	}
	return data
}

func main() {
	fd := makeFileDescriptor()
	var (
		msg  *dynamicpb.Message
		data []byte
	)

	// foo
	fmt.Println("example of Foo ---")
	msg = makeFooMsg(fd)
	data = marshalMsg(msg)
	useFooMsg(fd, data)

	// bar
	fmt.Println("example of Bar ---")
	msg = makeBarMsg(fd)
	data = marshalMsg(msg)
	useBarMsg(fd, data)

	// baz
	fmt.Println("example of Baz ---")
	msg = makeBazMsg(fd)
	data = marshalMsg(msg)
	useBazMsg(fd, data)
}
