#Protobuf

Protocol Buffers are language-neutral, platform-neutral extensible mechanisms for serializing structured data. Protobuf code can be converted and used with different languages.

install protobuf and protobuf go runtime

- brew install protobuf
- go install google.golang.org/protobuf/cmd/protoc-gen-go


compile protobuf code to go struct

protoc api/v1/*.proto \
--go_out=. \
--go_opt=paths=source_relative \
--proto_path=.