// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: model.proto

package pb

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/gogo/protobuf/gogoproto"
	time "time"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func (this *Message) Validate() error {
	return nil
}
func (this *User) Validate() error {
	if !(this.Id > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("Id", fmt.Errorf(`ID must a positive integer`))
	}
	if !(len(this.Username) > 6) {
		return github_com_mwitkow_go_proto_validators.FieldError("Username", fmt.Errorf(`username should between 6-12`))
	}
	if !(len(this.Username) < 12) {
		return github_com_mwitkow_go_proto_validators.FieldError("Username", fmt.Errorf(`username should between 6-12`))
	}
	if !(len(this.Password) > 6) {
		return github_com_mwitkow_go_proto_validators.FieldError("Password", fmt.Errorf(`password should between 6-12`))
	}
	if !(len(this.Password) < 12) {
		return github_com_mwitkow_go_proto_validators.FieldError("Password", fmt.Errorf(`password should between 6-12`))
	}
	if this.CreateDate != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.CreateDate); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("CreateDate", err)
		}
	}
	return nil
}
