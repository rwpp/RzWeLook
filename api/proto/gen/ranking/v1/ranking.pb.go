// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: ranking/v1/ranking.proto

package rankingv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Author struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id   int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"` // 添加其他作者相关字段
}

func (x *Author) Reset() {
	*x = Author{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Author) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Author) ProtoMessage() {}

func (x *Author) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Author.ProtoReflect.Descriptor instead.
func (*Author) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{0}
}

func (x *Author) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Author) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Article struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title   string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Status  int32                  `protobuf:"varint,3,opt,name=status,proto3" json:"status,omitempty"`
	Content string                 `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	Author  *Author                `protobuf:"bytes,5,opt,name=author,proto3" json:"author,omitempty"`
	Ctime   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=ctime,proto3" json:"ctime,omitempty"`
	Utime   *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=utime,proto3" json:"utime,omitempty"`
}

func (x *Article) Reset() {
	*x = Article{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Article) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Article) ProtoMessage() {}

func (x *Article) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Article.ProtoReflect.Descriptor instead.
func (*Article) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{1}
}

func (x *Article) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Article) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Article) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *Article) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Article) GetAuthor() *Author {
	if x != nil {
		return x.Author
	}
	return nil
}

func (x *Article) GetCtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Ctime
	}
	return nil
}

func (x *Article) GetUtime() *timestamppb.Timestamp {
	if x != nil {
		return x.Utime
	}
	return nil
}

type RankTopNRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RankTopNRequest) Reset() {
	*x = RankTopNRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RankTopNRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RankTopNRequest) ProtoMessage() {}

func (x *RankTopNRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RankTopNRequest.ProtoReflect.Descriptor instead.
func (*RankTopNRequest) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{2}
}

type RankTopNResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RankTopNResponse) Reset() {
	*x = RankTopNResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RankTopNResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RankTopNResponse) ProtoMessage() {}

func (x *RankTopNResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RankTopNResponse.ProtoReflect.Descriptor instead.
func (*RankTopNResponse) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{3}
}

type TopNRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *TopNRequest) Reset() {
	*x = TopNRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopNRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopNRequest) ProtoMessage() {}

func (x *TopNRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopNRequest.ProtoReflect.Descriptor instead.
func (*TopNRequest) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{4}
}

type TopNResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Articles []*Article `protobuf:"bytes,1,rep,name=articles,proto3" json:"articles,omitempty"`
}

func (x *TopNResponse) Reset() {
	*x = TopNResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ranking_v1_ranking_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TopNResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TopNResponse) ProtoMessage() {}

func (x *TopNResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ranking_v1_ranking_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TopNResponse.ProtoReflect.Descriptor instead.
func (*TopNResponse) Descriptor() ([]byte, []int) {
	return file_ranking_v1_ranking_proto_rawDescGZIP(), []int{5}
}

func (x *TopNResponse) GetArticles() []*Article {
	if x != nil {
		return x.Articles
	}
	return nil
}

var File_ranking_v1_ranking_proto protoreflect.FileDescriptor

var file_ranking_v1_ranking_proto_rawDesc = []byte{
	0x0a, 0x18, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x61, 0x6e,
	0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x72, 0x61, 0x6e, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2c, 0x0a, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f,
	0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xf1, 0x01, 0x0a, 0x07, 0x41, 0x72, 0x74, 0x69, 0x63, 0x6c,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x2a, 0x0a, 0x06, 0x61, 0x75, 0x74,
	0x68, 0x6f, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x72, 0x61, 0x6e, 0x6b,
	0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x52, 0x06, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x30, 0x0a, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x05, 0x63, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x05, 0x75, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x05, 0x75, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x11, 0x0a, 0x0f, 0x52, 0x61, 0x6e,
	0x6b, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x12, 0x0a, 0x10,
	0x52, 0x61, 0x6e, 0x6b, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x0d, 0x0a, 0x0b, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x3f, 0x0a, 0x0c, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2f, 0x0a, 0x08, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x13, 0x2e, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x52, 0x08, 0x61, 0x72, 0x74, 0x69, 0x63, 0x6c, 0x65, 0x73,
	0x32, 0x96, 0x01, 0x0a, 0x0e, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x08, 0x52, 0x61, 0x6e, 0x6b, 0x54, 0x6f, 0x70, 0x4e, 0x12,
	0x1b, 0x2e, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e,
	0x6b, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x72,
	0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x6e, 0x6b, 0x54, 0x6f,
	0x70, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x04,
	0x54, 0x6f, 0x70, 0x4e, 0x12, 0x17, 0x2e, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76,
	0x31, 0x2e, 0x54, 0x6f, 0x70, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e,
	0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x70, 0x4e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0xa4, 0x01, 0x0a, 0x0e, 0x63, 0x6f,
	0x6d, 0x2e, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x52, 0x61,
	0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3b, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x77, 0x70, 0x70, 0x2f, 0x52, 0x7a,
	0x57, 0x65, 0x4c, 0x6f, 0x6f, 0x6b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x3b,
	0x72, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x52, 0x58, 0x58, 0xaa,
	0x02, 0x0a, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x0a, 0x52,
	0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x16, 0x52, 0x61, 0x6e, 0x6b,
	0x69, 0x6e, 0x67, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x0b, 0x52, 0x61, 0x6e, 0x6b, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ranking_v1_ranking_proto_rawDescOnce sync.Once
	file_ranking_v1_ranking_proto_rawDescData = file_ranking_v1_ranking_proto_rawDesc
)

func file_ranking_v1_ranking_proto_rawDescGZIP() []byte {
	file_ranking_v1_ranking_proto_rawDescOnce.Do(func() {
		file_ranking_v1_ranking_proto_rawDescData = protoimpl.X.CompressGZIP(file_ranking_v1_ranking_proto_rawDescData)
	})
	return file_ranking_v1_ranking_proto_rawDescData
}

var file_ranking_v1_ranking_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_ranking_v1_ranking_proto_goTypes = []interface{}{
	(*Author)(nil),                // 0: ranking.v1.Author
	(*Article)(nil),               // 1: ranking.v1.Article
	(*RankTopNRequest)(nil),       // 2: ranking.v1.RankTopNRequest
	(*RankTopNResponse)(nil),      // 3: ranking.v1.RankTopNResponse
	(*TopNRequest)(nil),           // 4: ranking.v1.TopNRequest
	(*TopNResponse)(nil),          // 5: ranking.v1.TopNResponse
	(*timestamppb.Timestamp)(nil), // 6: google.protobuf.Timestamp
}
var file_ranking_v1_ranking_proto_depIdxs = []int32{
	0, // 0: ranking.v1.Article.author:type_name -> ranking.v1.Author
	6, // 1: ranking.v1.Article.ctime:type_name -> google.protobuf.Timestamp
	6, // 2: ranking.v1.Article.utime:type_name -> google.protobuf.Timestamp
	1, // 3: ranking.v1.TopNResponse.articles:type_name -> ranking.v1.Article
	2, // 4: ranking.v1.RankingService.RankTopN:input_type -> ranking.v1.RankTopNRequest
	4, // 5: ranking.v1.RankingService.TopN:input_type -> ranking.v1.TopNRequest
	3, // 6: ranking.v1.RankingService.RankTopN:output_type -> ranking.v1.RankTopNResponse
	5, // 7: ranking.v1.RankingService.TopN:output_type -> ranking.v1.TopNResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_ranking_v1_ranking_proto_init() }
func file_ranking_v1_ranking_proto_init() {
	if File_ranking_v1_ranking_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ranking_v1_ranking_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Author); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ranking_v1_ranking_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Article); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ranking_v1_ranking_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RankTopNRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ranking_v1_ranking_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RankTopNResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ranking_v1_ranking_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopNRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_ranking_v1_ranking_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TopNResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_ranking_v1_ranking_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ranking_v1_ranking_proto_goTypes,
		DependencyIndexes: file_ranking_v1_ranking_proto_depIdxs,
		MessageInfos:      file_ranking_v1_ranking_proto_msgTypes,
	}.Build()
	File_ranking_v1_ranking_proto = out.File
	file_ranking_v1_ranking_proto_rawDesc = nil
	file_ranking_v1_ranking_proto_goTypes = nil
	file_ranking_v1_ranking_proto_depIdxs = nil
}
