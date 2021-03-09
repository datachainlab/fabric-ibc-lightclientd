package main

import (
	"log"
	"strings"

	ics23 "github.com/confio/ics23/go"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/mem"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/core/03-connection/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/23-commitment/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	pb "github.com/datachainlab/fabric-ibc-lightclientd/types"
	fabrictypes "github.com/datachainlab/fabric-ibc/x/ibc/light-clients/xx-fabric/types"
)

type Lightclient struct {
	ctx   sdk.Context
	cdc   codec.BinaryMarshaler
	store sdk.KVStore
	cs    *fabrictypes.ClientState
}

func mustParseKeyConsensusState(key []byte) clienttypes.Height {
	words := strings.Split(string(key), "/")
	if len(words) != 2 || words[0] != string(host.KeyConsensusStatesPrefix) {
		log.Fatalf("failed to split consensus state key: %v", key)
	}
	return clienttypes.MustParseHeight(words[1])
}

func (lc *Lightclient) State() *pb.State {
	state := &pb.State{
		ClientState:     lc.cs,
		ConsensusStates: make(map[uint64]*fabrictypes.ConsensusState),
	}
	it := lc.store.Iterator(nil, nil)
	defer it.Close()
	for it.Valid() {
		height := mustParseKeyConsensusState(it.Key()).GetVersionHeight()
		consensusState := clienttypes.MustUnmarshalConsensusState(lc.cdc, it.Value()).(*fabrictypes.ConsensusState)
		state.ConsensusStates[height] = consensusState
		it.Next()
	}
	return state
}

func (lc *Lightclient) saveConsensusState(height uint64, consensusState *fabrictypes.ConsensusState) {
	lc.store.Set(
		host.KeyConsensusState(clienttypes.NewHeight(0, height)),
		clienttypes.MustMarshalConsensusState(lc.cdc, consensusState),
	)
}

func NewLightclient(state *pb.State) *Lightclient {
	lc := &Lightclient{
		ctx:   sdk.Context{},
		cdc:   codec.NewProtoCodec(codectypes.NewInterfaceRegistry()),
		store: mem.NewStore(),
		cs:    state.ClientState,
	}

	// save consensus states in store
	for height, consensusState := range state.ConsensusStates {
		lc.saveConsensusState(height, consensusState)
	}

	// create lightclient core
	return lc
}

func (lc *Lightclient) ClientType() string {
	return lc.cs.ClientType()
}

func (lc *Lightclient) GetLatestHeight() clienttypes.Height {
	return lc.cs.GetLatestHeight().(clienttypes.Height)
}

func (lc *Lightclient) IsFrozen() bool {
	return lc.cs.IsFrozen()
}

func (lc *Lightclient) GetFrozenHeight() clienttypes.Height {
	return lc.cs.GetFrozenHeight().(clienttypes.Height)
}

func (lc *Lightclient) Validate() error {
	return lc.cs.Validate()
}

func (lc *Lightclient) GetProofSpecs() []*ics23.ProofSpec {
	return lc.cs.GetProofSpecs()
}

func (lc *Lightclient) CheckHeaderAndUpdateState(header *fabrictypes.Header) error {
	clientState, consensusState, err := lc.cs.CheckHeaderAndUpdateState(lc.ctx, lc.cdc, lc.store, header)
	if err != nil {
		return err
	}
	lc.cs = clientState.(*fabrictypes.ClientState)
	lc.saveConsensusState(header.ChaincodeHeader.Sequence.Value, consensusState.(*fabrictypes.ConsensusState))
	return nil
}

// TODO: define fabrictypes.Misbehaviour
/*
func (lc *Lightclient) CheckMisbehaviourAndUpdateState(misbehaviour *fabrictypes.Misbehaviour) (*fabrictypes.ClientState, error) {
	clientState, err := lc.cs.CheckMisbehaviourAndUpdateState(lc.ctx, lc.cdc, lc.store, misbehaviour)
	return clientState.(*fabrictypes.ClientState), err
}
*/

func (lc *Lightclient) CheckProposedHeaderAndUpdateState(header *fabrictypes.Header) error {
	clientState, consensusState, err := lc.cs.CheckProposedHeaderAndUpdateState(lc.ctx, lc.cdc, lc.store, header)
	if err != nil {
		return err
	}
	lc.cs = clientState.(*fabrictypes.ClientState)
	lc.saveConsensusState(header.ChaincodeHeader.Sequence.Value, consensusState.(*fabrictypes.ConsensusState))
	return nil
}

func (lc *Lightclient) VerifyUpgrade(
	newClient *fabrictypes.ClientState,
	upgradeHeight clienttypes.Height,
	proofUpgrade []byte,
) error {
	return lc.cs.VerifyUpgrade(lc.ctx, lc.cdc, lc.store, newClient, upgradeHeight, proofUpgrade)
}

func (lc *Lightclient) ZeroCustomFields() *fabrictypes.ClientState {
	return lc.cs.ZeroCustomFields().(*fabrictypes.ClientState)
}

func (lc *Lightclient) VerifyClientState(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	counterpartyClientIdentifier string,
	proof []byte,
	clientState *fabrictypes.ClientState,
) error {
	return lc.cs.VerifyClientState(lc.store, lc.cdc, height, prefix, counterpartyClientIdentifier, proof, clientState)
}

func (lc *Lightclient) VerifyClientConsensusState(
	height clienttypes.Height,
	counterpartyClientIdentifier string,
	consensusHeight clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	consensusState *fabrictypes.ConsensusState,
) error {
	return lc.cs.VerifyClientConsensusState(lc.store, lc.cdc, height, counterpartyClientIdentifier, consensusHeight, prefix, proof, consensusState)
}

func (lc *Lightclient) VerifyConnectionState(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	connectionID string,
	connectionEnd connectiontypes.ConnectionEnd,
) error {
	return lc.cs.VerifyConnectionState(lc.store, lc.cdc, height, prefix, proof, connectionID, connectionEnd)
}

func (lc *Lightclient) VerifyChannelState(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	portID,
	channelID string,
	channel channeltypes.Channel,
) error {
	return lc.cs.VerifyChannelState(lc.store, lc.cdc, height, prefix, proof, portID, channelID, channel)
}

func (lc *Lightclient) VerifyPacketCommitment(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	commitmentBytes []byte,
) error {
	return lc.cs.VerifyPacketCommitment(lc.store, lc.cdc, height, prefix, proof, portID, channelID, sequence, commitmentBytes)
}

func (lc *Lightclient) VerifyPacketAcknowledgement(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
	acknowledgement []byte,
) error {
	return lc.cs.VerifyPacketAcknowledgement(lc.store, lc.cdc, height, prefix, proof, portID, channelID, sequence, acknowledgement)
}

func (lc *Lightclient) VerifyPacketReceiptAbsence(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	portID,
	channelID string,
	sequence uint64,
) error {
	return lc.cs.VerifyPacketReceiptAbsence(lc.store, lc.cdc, height, prefix, proof, portID, channelID, sequence)
}

func (lc *Lightclient) VerifyNextSequenceRecv(
	height clienttypes.Height,
	prefix *commitmenttypes.MerklePrefix,
	proof []byte,
	portID,
	channelID string,
	nextSequenceRecv uint64,
) error {
	return lc.cs.VerifyNextSequenceRecv(lc.store, lc.cdc, height, prefix, proof, portID, channelID, nextSequenceRecv)
}
