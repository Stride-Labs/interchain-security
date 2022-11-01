package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Status int

const (
	// ModuleName defines the CCV provider module name
	ModuleName = "provider"

	// StoreKey is the store key string for IBC transfer
	StoreKey = ModuleName

	// RouterKey is the message route for IBC transfer
	RouterKey = ModuleName

	// QuerierRoute is the querier route for IBC transfer
	QuerierRoute = ModuleName
)

// Iota generated keys/byte prefixes (as a byte), supports 256 possible values
const (

	// PortKey defines the key to store the port ID in store
	PortByteKey byte = iota

	// MaturedUnbondingOpsByteKey is the byte key that stores the list of all unbonding operations ids
	// that have matured from a consumer chain perspective,
	// i.e., no longer waiting on the unbonding period to elapse on any consumer chain
	MaturedUnbondingOpsByteKey

	// ValidatorSetUpdateIdByteKey is the byte key that stores the current validator set update id
	ValidatorSetUpdateIdByteKey

	// SlashGasMeterBytePrefix is the byte prefix for storing the slash gas meter
	SlashGasMeterBytePrefix

	// LastSlashGasReplenishTimeBytePrefix is the byte prefix for storing
	// the last time the slash gas meter was replenished
	LastSlashGasReplenishTimeBytePrefix

	// NumPendingSlashPacketsBytePrefix is the byte prefix for storing the number of pending slash packets in the queue
	NumPendingSlashPacketsBytePrefix

	// ChainToChannelBytePrefix is the byte prefix for storing mapping
	// from chainID to the channel ID that is used to send over validator set changes.
	ChainToChannelBytePrefix

	// ChannelToChainBytePrefix is the byte prefix for storing mapping
	// from the CCV channel ID to the consumer chain ID.
	ChannelToChainBytePrefix

	// ChainToClientBytePrefix is the byte prefix for storing the client ID for a given consumer chainID.
	ChainToClientBytePrefix

	// InitTimeoutTimestampBytePrefix is the byte prefix for storing
	// the init timeout timestamp for a given consumer chainID.
	InitTimeoutTimestampBytePrefix

	// PendingCAPBytePrefix is the byte prefix for storing pending consumer addition proposals before the spawn time occurs.
	// The key includes the BigEndian timestamp to allow for efficient chronological iteration
	PendingCAPBytePrefix

	// PendingCRPBytePrefix is the byte prefix for storing pending consumer removal proposals before the stop time occurs.
	// The key includes the BigEndian timestamp to allow for efficient chronological iteration
	PendingCRPBytePrefix

	// UnbondingOpBytePrefix is the byte prefix that stores a record of all the ids of consumer chains that
	// need to unbond before a given delegation can unbond on this chain.
	UnbondingOpBytePrefix

	// UnbondingOpIndexBytePrefix is byte prefix of the index for looking up which unbonding
	// delegation entries are waiting for a given consumer chain to unbond
	UnbondingOpIndexBytePrefix

	// ValsetUpdateBlockHeightBytePrefix is the byte prefix that will store the mapping from valset update ID to block height
	ValsetUpdateBlockHeightBytePrefix

	// ConsumerGenesisBytePrefix stores consumer genesis state material (consensus state and client state) indexed by consumer chain id
	ConsumerGenesisBytePrefix

	// SlashAcksBytePrefix is the byte prefix that will store consensus address of consumer chain validators successfully slashed on the provider chain
	SlashAcksBytePrefix

	// InitChainHeightBytePrefix is the byte prefix that will store the mapping from a chain id to the corresponding block height on the provider
	// this consumer chain was initialized
	InitChainHeightBytePrefix

	// PendingVSCsBytePrefix is the byte prefix that will store pending ValidatorSetChangePacket data
	PendingVSCsBytePrefix

	// LockUnbondingOnTimeoutBytePrefix is the byte prefix that will store the consumer chain id which unbonding operations are locked on CCV channel timeout
	LockUnbondingOnTimeoutBytePrefix

	// PendingSlashPacketEntryBytePrefix is the byte prefix storing pending slash packet entries
	PendingSlashPacketEntryBytePrefix
)

const (
	// UnbondingOpIndexKey should be of set length: prefix + hashed chain ID + uint64
	UnbondingOpIndexKeySize = 1 + 32 + 8
)

// PortKey returns the key to the port ID in the store
func PortKey() []byte {
	return []byte{PortByteKey}
}

// MaturedUnbondingOpsKey returns the key for storing the list of matured unbonding operations.
func MaturedUnbondingOpsKey() []byte {
	return []byte{MaturedUnbondingOpsByteKey}
}

// ValidatorSetUpdateIdKey is the key that stores the current validator set update id
func ValidatorSetUpdateIdKey() []byte {
	return []byte{ValidatorSetUpdateIdByteKey}
}

// SlashGasMeterKey returns the key under which the slash meter is stored
func SlashGasMeterKey() []byte {
	return []byte{SlashGasMeterBytePrefix}
}

// LastSlashGasReplenishTimeKey returns the key under which the last time the slash gas meter was replenished is stored
func LastSlashGasReplenishTimeKey() []byte {
	return []byte{LastSlashGasReplenishTimeBytePrefix}
}

// ChainToChannelKey returns the key under which the CCV channel ID will be stored for the given consumer chain.
func ChainToChannelKey(chainID string) []byte {
	return append([]byte{ChainToChannelBytePrefix}, []byte(chainID)...)
}

// ChannelToChainKey returns the key under which the consumer chain ID will be stored for the given channelID.
func ChannelToChainKey(channelID string) []byte {
	return append([]byte{ChannelToChainBytePrefix}, []byte(channelID)...)
}

// ChainToClientKey returns the key under which the clientID for the given chainID is stored.
func ChainToClientKey(chainID string) []byte {
	return append([]byte{ChainToClientBytePrefix}, []byte(chainID)...)
}

// InitTimeoutTimestampKey returns the key under which the init timeout timestamp for the given chainID is stored.
func InitTimeoutTimestampKey(chainID string) []byte {
	return append([]byte{InitTimeoutTimestampBytePrefix}, []byte(chainID)...)
}

// PendingCAPKey returns the key under which a pending consumer addition proposal is stored
func PendingCAPKey(timestamp time.Time, chainID string) []byte {
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len([]byte{PendingCAPBytePrefix})

	bz := make([]byte, prefixL+8+timeBzL+len(chainID))
	// copy the prefix
	copy(bz[:prefixL], []byte{PendingCAPBytePrefix})
	// copy the time length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))
	// copy the time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	// copy the chainId
	copy(bz[prefixL+8+timeBzL:], chainID)
	return bz
}

// ParsePendingCAPKey returns the time and chain ID for a pending consumer addition proposal key
// or an error if unparsable
func ParsePendingCAPKey(bz []byte) (time.Time, string, error) {
	expectedPrefix := []byte{PendingCAPBytePrefix}
	prefixL := len(expectedPrefix)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, expectedPrefix) {
		return time.Time{}, "", fmt.Errorf("invalid prefix; expected: %X, got: %X", expectedPrefix, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	timestamp, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, "", err
	}

	chainID := string(bz[prefixL+8+int(timeBzL):])
	return timestamp, chainID, nil
}

// PendingCRPKey returns the key under which pending consumer removal proposals are stored
func PendingCRPKey(timestamp time.Time, chainID string) []byte {
	timeBz := sdk.FormatTimeBytes(timestamp)
	timeBzL := len(timeBz)
	prefixL := len([]byte{PendingCRPBytePrefix})

	bz := make([]byte, prefixL+8+timeBzL+len(chainID))
	// copy the prefix
	copy(bz[:prefixL], []byte{PendingCRPBytePrefix})
	// copy the time length
	copy(bz[prefixL:prefixL+8], sdk.Uint64ToBigEndian(uint64(timeBzL)))
	// copy the time bytes
	copy(bz[prefixL+8:prefixL+8+timeBzL], timeBz)
	// copy the chainId
	copy(bz[prefixL+8+timeBzL:], chainID)
	return bz
}

// ParsePendingCRPKey returns the time and chain ID for a pending consumer removal proposal key or an error if unparseable
func ParsePendingCRPKey(bz []byte) (time.Time, string, error) {
	expectedPrefix := []byte{PendingCRPBytePrefix}
	prefixL := len(expectedPrefix)
	if prefix := bz[:prefixL]; !bytes.Equal(prefix, expectedPrefix) {
		return time.Time{}, "", fmt.Errorf("invalid prefix; expected: %X, got: %X", expectedPrefix, prefix)
	}

	timeBzL := sdk.BigEndianToUint64(bz[prefixL : prefixL+8])
	timestamp, err := sdk.ParseTimeBytes(bz[prefixL+8 : prefixL+8+int(timeBzL)])
	if err != nil {
		return time.Time{}, "", err
	}

	chainID := string(bz[prefixL+8+int(timeBzL):])
	return timestamp, chainID, nil
}

// UnbondingOpIndexKey returns an unbonding op index key
// Note: chainId is hashed to a fixed length sequence of bytes here to prevent
// injection attack between chainIDs.
func UnbondingOpIndexKey(chainID string, valsetUpdateID uint64) []byte {
	return AppendMany([]byte{UnbondingOpIndexBytePrefix}, HashString(chainID),
		sdk.Uint64ToBigEndian(valsetUpdateID))
}

// ParseUnbondingOpIndexKey parses an unbonding op index key for VSC ID
// Removes the prefix + chainID from index key and returns only the key part.
func ParseUnbondingOpIndexKey(key []byte) (vscID []byte, err error) {
	if len(key) != UnbondingOpIndexKeySize {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrLogic, "key provided is incorrect: the key has incorrect length, expected %d, got %d",
			UnbondingOpIndexKeySize, len(key),
		)
	}
	return key[1+32:], nil
}

// UnbondingOpKey returns the key that stores a record of all the ids of consumer chains that
// need to unbond before a given delegation can unbond on this chain
func UnbondingOpKey(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return append([]byte{UnbondingOpBytePrefix}, bz...)
}

// ValsetUpdateBlockHeightKey returns the key that storing the mapping from valset update ID to block height
func ValsetUpdateBlockHeightKey(valsetUpdateId uint64) []byte {
	vuidBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(vuidBytes, valsetUpdateId)
	return append([]byte{ValsetUpdateBlockHeightBytePrefix}, vuidBytes...)
}

// ConsumerGenesisKey returns the key corresponding to consumer genesis state material
// (consensus state and client state) indexed by consumer chain id
func ConsumerGenesisKey(chainID string) []byte {
	return append([]byte{ConsumerGenesisBytePrefix}, []byte(chainID)...)
}

// SlashAcksKey returns the key under which slashing acks are stored for a given chain ID
func SlashAcksKey(chainID string) []byte {
	return append([]byte{SlashAcksBytePrefix}, []byte(chainID)...)
}

// InitChainHeightKey returns the key under which the block height for a given chain ID is stored
func InitChainHeightKey(chainID string) []byte {
	return append([]byte{InitChainHeightBytePrefix}, []byte(chainID)...)
}

// PendingVSCsKey returns the key under which
// pending ValidatorSetChangePacket data is stored for a given chain ID
func PendingVSCsKey(chainID string) []byte {
	return append([]byte{PendingVSCsBytePrefix}, []byte(chainID)...)
}

// LockUnbondingOnTimeoutKey returns the key that will store the consumer chain id which unbonding operations are locked
// on CCV channel timeout
func LockUnbondingOnTimeoutKey(chainID string) []byte {
	return append([]byte{LockUnbondingOnTimeoutBytePrefix}, []byte(chainID)...)
}

// PendingSlashPacketEntryKey returns the key for storing a pending slash packet entry.
//
// Note: It's not expected for a single consumer chain to send a slash packet for the same validator more than once
// in the same block. Hence why this key should ber unique per slash packet. However, if a malicious consumer did send
// duplicate slash packets in the same block, the slash packet entry would simply be overwritten.
func PendingSlashPacketEntryKey(packetEntry SlashPacketEntry) []byte {
	timeBz := sdk.FormatTimeBytes(packetEntry.RecvTime)
	timeBzL := len(timeBz)
	return AppendMany(
		[]byte{PendingSlashPacketEntryBytePrefix},
		sdk.Uint64ToBigEndian(uint64(timeBzL)),
		timeBz,
		HashBytes(packetEntry.ValAddr),
		[]byte(packetEntry.ConsumerChainID),
	)
}

// ParsePendingSlashPacketEntryKey returns the received time and chainID for a pending slash packet entry key
func ParsePendingSlashPacketEntryKey(bz []byte) (time.Time, string) {
	// Prefix is in first byte
	expectedPrefix := []byte{PendingSlashPacketEntryBytePrefix}
	if prefix := bz[:1]; !bytes.Equal(prefix, expectedPrefix) {
		panic(fmt.Sprintf("invalid prefix; expected: %X, got: %X", expectedPrefix, prefix))
	}
	// 8 bytes for uint64 storing timestamp length
	timeBzL := sdk.BigEndianToUint64(bz[1:9])
	recvTime, err := sdk.ParseTimeBytes(bz[9 : 9+timeBzL])
	if err != nil {
		panic(err)
	}
	// ChainID is stored after 32 byte hashed validator address
	chainID := string(bz[9+int(timeBzL)+32:])
	return recvTime, chainID
}

// AppendMany appends a variable number of byte slices together
func AppendMany(byteses ...[]byte) (out []byte) {
	for _, bytes := range byteses {
		out = append(out, bytes...)
	}
	return out
}

// HashString outputs a fixed length 32 byte hash for any string
func HashString(x string) []byte {
	return HashBytes([]byte(x))
}

// HashBytes outputs a fixed length 32 byte hash for any byte slice
func HashBytes(x []byte) []byte {
	hash := sha256.Sum256(x)
	return hash[:]
}
