module github.com/datachainlab/fabric-ibc-lightclientd

go 1.16

require (
	github.com/cosmos/cosmos-sdk v0.40.0-rc3
	github.com/datachainlab/fabric-ibc v0.0.0-20210118090849-c2eaee7a3314
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/tendermint/tm-db v0.6.2
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
