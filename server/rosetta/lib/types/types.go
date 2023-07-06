package types

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// SpecVersion defines the specification of rosetta
const SpecVersion = ""

// NetworkInformationProvider defines the interface used to provide information regarding
// the network and the version of the cosmos sdk used
type NetworkInformationProvider interface {
	// SupportedOperations lists the operations supported by the implementation
	SupportedOperations() []string
	// OperationStatuses returns the list of statuses supported by the implementation
	OperationStatuses() []*types.OperationStatus
	// Version returns the version of the node
	Version() string
}

// Client defines the API the client implementation should provide.
type Client interface {
	// Bootstrap Needed if the client needs to perform some action before connecting.
	Bootstrap() error
	// Ready checks if the servicer constraints for queries are satisfied
	// for example the node might still not be ready, it's useful in process
	// when the rosetta instance might come up before the node itself
	// the servicer must return nil if the node is ready
	Ready() error

	// Data API

	AccountInfo(ctx context.Context, addr string, height *int64) (*SignerData, error)
	// Balances fetches the balance of the given address
	// if height is not nil, then the balance will be displayed
	// at the provided height, otherwise last block balance will be returned
	Balances(ctx context.Context, addr string, height *int64) ([]*types.Amount, error)
	// BlockByHash gets a block and its transaction at the provided height
	BlockByHash(ctx context.Context, hash string) (BlockResponse, error)
	// BlockByHeight gets a block given its height, if height is nil then last block is returned
	BlockByHeight(ctx context.Context, height *int64) (BlockResponse, error)
	// BlockTransactionsByHash gets the block, parent block and transactions
	// given the block hash.
	BlockTransactionsByHash(ctx context.Context, hash string) (BlockTransactionsResponse, error)
	// BlockTransactionsByHeight gets the block, parent block and transactions
	// given the block hash.
	BlockTransactionsByHeight(ctx context.Context, height *int64) (BlockTransactionsResponse, error)
	// GetTx gets a transaction given its hash
	GetTx(ctx context.Context, hash string) (*types.Transaction, error)
	// GetUnconfirmedTx gets an unconfirmed Tx given its hash
	// NOTE(fdymylja): NOT IMPLEMENTED YET!
	GetUnconfirmedTx(ctx context.Context, hash string) (*types.Transaction, error)
	// Mempool returns the list of the current non confirmed transactions
	Mempool(ctx context.Context) ([]*types.TransactionIdentifier, error)
	// Peers gets the peers currently connected to the node
	Peers(ctx context.Context) ([]*types.Peer, error)
	// Status returns the node status, such as sync data, version etc
	Status(ctx context.Context) (*types.SyncStatus, error)
	// LastBlockHeight returns last committed block height
	LastBlockHeight(ctx context.Context) (int64, error)
	// Rewards fetches the pending rewards of the given delegator address
	Rewards(ctx context.Context, delegator string, validator string, height *int64) ([]*types.Amount, error)
	// UnbondingDelegations fetches the unbonging delegations of the given delegator address
	UnbondingDelegations(ctx context.Context, delegator string, height *int64) ([]*types.Amount, error)
	// Delegations fetches the delegations of the given delegator address
	Delegations(ctx context.Context, delegator string, height *int64) ([]*types.Amount, error)

	// Construction API

	// PostTx posts txBytes to the node and returns the transaction identifier plus metadata related
	// to the transaction itself.
	PostTx(txBytes []byte) (res *types.TransactionIdentifier, meta map[string]interface{}, err error)
	// ConstructionMetadataFromOptions builds metadata map from an option map
	ConstructionMetadataFromOptions(ctx context.Context, options map[string]interface{}) (meta map[string]interface{}, err error)
	OfflineClient
}

// OfflineClient defines the functionalities supported without having access to the node
type OfflineClient interface {
	NetworkInformationProvider
	// SignedTx returns the signed transaction given the tx bytes (msgs) plus the signatures
	SignedTx(ctx context.Context, txBytes []byte, sigs []*types.Signature) (signedTxBytes []byte, err error)
	// TxOperationsAndSignersAccountIdentifiers returns the operations related to a transaction and the account
	// identifiers if the transaction is signed
	TxOperationsAndSignersAccountIdentifiers(signed bool, hexBytes []byte) (ops []*types.Operation, signers []*types.AccountIdentifier, err error)
	// ConstructionPayload returns the construction payload given the request
	ConstructionPayload(ctx context.Context, req *types.ConstructionPayloadsRequest) (resp *types.ConstructionPayloadsResponse, err error)
	// PreprocessOperationsToOptions returns the options given the preprocess operations
	PreprocessOperationsToOptions(ctx context.Context, req *types.ConstructionPreprocessRequest) (resp *types.ConstructionPreprocessResponse, err error)
	// AccountIdentifierFromPublicKey returns the account identifier given the public key
	AccountIdentifierFromPublicKey(pubKey *types.PublicKey) (*types.AccountIdentifier, error)
}

type BlockTransactionsResponse struct {
	BlockResponse
	Transactions []*types.Transaction
}

type BlockResponse struct {
	Block                *types.BlockIdentifier
	ParentBlock          *types.BlockIdentifier
	MillisecondTimestamp int64
	TxCount              int64
}

// API defines the exposed APIs
// if the service is online
type API interface {
	DataAPI
	ConstructionAPI
}

// DataAPI defines the full data API implementation
type DataAPI interface {
	server.NetworkAPIServicer
	server.AccountAPIServicer
	server.BlockAPIServicer
	server.MempoolAPIServicer
}

var _ server.ConstructionAPIServicer = ConstructionAPI(nil)

// ConstructionAPI defines the full construction API with
// the online and offline endpoints
type ConstructionAPI interface {
	ConstructionOnlineAPI
	ConstructionOfflineAPI
}

// ConstructionOnlineAPI defines the construction methods
// allowed in an online implementation
type ConstructionOnlineAPI interface {
	ConstructionMetadata(
		context.Context,
		*types.ConstructionMetadataRequest,
	) (*types.ConstructionMetadataResponse, *types.Error)
	ConstructionSubmit(
		context.Context,
		*types.ConstructionSubmitRequest,
	) (*types.TransactionIdentifierResponse, *types.Error)
}

// ConstructionOfflineAPI defines the construction methods
// allowed
type ConstructionOfflineAPI interface {
	ConstructionCombine(
		context.Context,
		*types.ConstructionCombineRequest,
	) (*types.ConstructionCombineResponse, *types.Error)
	ConstructionDerive(
		context.Context,
		*types.ConstructionDeriveRequest,
	) (*types.ConstructionDeriveResponse, *types.Error)
	ConstructionHash(
		context.Context,
		*types.ConstructionHashRequest,
	) (*types.TransactionIdentifierResponse, *types.Error)
	ConstructionParse(
		context.Context,
		*types.ConstructionParseRequest,
	) (*types.ConstructionParseResponse, *types.Error)
	ConstructionPayloads(
		context.Context,
		*types.ConstructionPayloadsRequest,
	) (*types.ConstructionPayloadsResponse, *types.Error)
	ConstructionPreprocess(
		context.Context,
		*types.ConstructionPreprocessRequest,
	) (*types.ConstructionPreprocessResponse, *types.Error)
}

// SignerData contains information on the signers when the request
// is being created, used to populate the account information
type SignerData struct {
	AccountNumber uint64 `json:"account_number"`
	Sequence      uint64 `json:"sequence"`
}

// BalanceType used to query different account balance
type BalanceType int

const (
	Unrecognized BalanceType = iota
	AvailableBalance
	PendingRewards
	UnbondingBalance
	DelegatedBalance
)

func (i BalanceType) String() string {
	switch i {
	case AvailableBalance:
		return "available_balance"
	case PendingRewards:
		return "pending_rewards"
	case UnbondingBalance:
		return "unbonding_balance"
	case DelegatedBalance:
		return "delegated_balance"
	default:
		return "unrecognized_balance_type"
	}
}

func ParseBalanceType(s string) (BalanceType, error) {
	switch strings.ToLower(s) {
	case strings.ToLower(AvailableBalance.String()):
		return AvailableBalance, nil
	case strings.ToLower(PendingRewards.String()):
		return PendingRewards, nil
	case strings.ToLower(UnbondingBalance.String()):
		return UnbondingBalance, nil
	case strings.ToLower(DelegatedBalance.String()):
		return DelegatedBalance, nil
	default:
		return -1, fmt.Errorf("unrecognized balance type option")
	}
}

// UnbondingDelegationMetaData attaches the validator address and completion time to the unbonding delegation
func UnbondingDelegationMetaData(validator string, completionTime time.Time) map[string]interface{} {
	return map[string]interface{}{
		"balance_type":      UnbondingBalance.String(),
		"completion_time":   completionTime.String(),
		"validator_address": validator,
	}
}

// BalanceMetaData attaches the validator address to balance response
func BalanceMetaData(balanceType BalanceType, validator string) map[string]interface{} {
	return map[string]interface{}{
		"balance_type":      balanceType.String(),
		"validator_address": validator,
	}
}
