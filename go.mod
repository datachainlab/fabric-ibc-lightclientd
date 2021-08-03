module github.com/datachainlab/fabric-ibc-lightclientd

go 1.16

require (
	github.com/confio/ics23/go v0.6.6
	github.com/cosmos/cosmos-sdk v0.43.0-beta1
	github.com/cosmos/ibc-go v1.0.0-beta1
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger-labs/yui-corda-ibc/go v0.0.0
	github.com/hyperledger-labs/yui-fabric-ibc v0.2.0
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
	github.com/hyperledger-labs/yui-corda-ibc/go => github.com/siburu/yui-corda-ibc/go v0.0.0-20210803083131-26a2c06b0680
)
