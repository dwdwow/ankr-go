package ankr

func TruePtr() *bool {
	return &[]bool{true}[0]
}

func FalsePtr() *bool {
	return &[]bool{false}[0]
}

// Chain represents supported blockchain networks
type Chain string

// Mainnet chains
const (
	ChainArbitrum     Chain = "arbitrum"
	ChainAvalanche    Chain = "avalanche"
	ChainBase         Chain = "base"
	ChainBSC          Chain = "bsc"
	ChainEthereum     Chain = "eth"
	ChainFantom       Chain = "fantom"
	ChainFlare        Chain = "flare"
	ChainGnosis       Chain = "gnosis"
	ChainLinea        Chain = "linea"
	ChainOptimism     Chain = "optimism"
	ChainPolygon      Chain = "polygon"
	ChainPolygonZkEVM Chain = "polygon_zkevm"
	ChainScroll       Chain = "scroll"
	ChainStellar      Chain = "stellar"
	ChainStory        Chain = "story_mainnet"
	ChainSyscoin      Chain = "syscoin"
	ChainTelos        Chain = "telos"
	ChainXai          Chain = "xai"
	ChainXLayer       Chain = "xlayer"
)

// Testnet chains
const (
	ChainAvalancheFuji   Chain = "avalanche_fuji"
	ChainBaseSepolia     Chain = "base_sepolia"
	ChainEthereumHolesky Chain = "eth_holesky"
	ChainEthereumSepolia Chain = "eth_sepolia"
	ChainOptimismTestnet Chain = "optimism_testnet"
	ChainPolygonAmoy     Chain = "polygon_amoy"
	ChainStoryTestnet    Chain = "story_aeneid_testnet"
)

const JSONRPC = "2.0"

type RPCReqBody struct {
	ID      int64  `json:"id" bson:"id"`
	JSONRPC string `json:"jsonrpc" bson:"jsonrpc"`
	Method  string `json:"method" bson:"method"`
	Params  any    `json:"params" bson:"params"`
}

type RPCRespBody[Result any] struct {
	JSONRPC string        `json:"jsonrpc" bson:"jsonrpc"`
	ID      int64         `json:"id" bson:"id"`
	Result  Result        `json:"result" bson:"result"`
	Error   *RPCRespError `json:"error" bson:"error"`
}

type RPCRespError struct {
	Code    int    `json:"code" bson:"code"`
	Message string `json:"message" bson:"message"`
	Data    any    `json:"data" bson:"data"`
}

// GetNFTsByOwnerReq represents the request parameters for ankr_getNFTsByOwner
type GetNFTsByOwnerReq struct {
	// WalletAddress is the account address to query for NFTs; supports ENS
	WalletAddress string `json:"walletAddress,omitempty" bson:"walletAddress,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// PageSize is the number of page results (default=10, max=50)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"50"`

	// PageToken is provided at the end of response and used to fetch next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`

	// Filter filters the request by smart contract address and optional NFT ID
	// Format: map[string][]string
	// Example: {"0xd8682bfa6918b0174f287b888e765b9a1b4dc9c3": []} - all NFTs from address
	// Example: {"0xd8682bfa6918b0174f287b888e765b9a1b4dc9c3": ["8937"]} - specific NFT
	Filter map[string][]string `json:"filter,omitempty" bson:"filter,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTsByOwnerReq) setPageToken(token string) {
	r.PageToken = token
}

// NFTTrait represents a trait/attribute of an NFT
type NFTTrait struct {
	// TraitType is the trait's descriptive name
	TraitType string `json:"trait_type" bson:"trait_type"`

	// Value is the value description
	Value string `json:"value" bson:"value"`
}

// NFT represents an NFT asset
type NFT struct {
	// Blockchain is one of the supported chains
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// CollectionName is the collection name the NFT belongs to
	CollectionName string `json:"collectionName" bson:"collectionName"`

	// ContractAddress is the NFT collection's EVM-compatible contract address
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// ContractType is the type of contract - either ERC721 or ERC1155
	ContractType string `json:"contractType" bson:"contractType"`

	// Name is the name of the NFT asset
	Name string `json:"name" bson:"name"`

	// TokenID is the ID of the NFT asset
	TokenID string `json:"tokenId" bson:"tokenId"`

	// ImageURL is a URL that points to the actual digital file, usually an IPFS link
	ImageURL string `json:"imageUrl" bson:"imageUrl"`

	// Symbol is the symbol of the NFT asset
	Symbol string `json:"symbol" bson:"symbol"`

	// Traits are the attributes of the NFT asset
	Traits []NFTTrait `json:"traits" bson:"traits"`
}

// GetNFTsByOwnerResp represents the response for ankr_getNFTsByOwner
type GetNFTsByOwnerResp struct {
	// Assets is the list of NFT assets
	Assets []NFT `json:"assets" bson:"assets"`

	// NextPageToken is provided at the end of response for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTsByOwnerResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetNFTMetadataReq represents the request parameters for ankr_getNFTMetadata
type GetNFTMetadataReq struct {
	// Blockchain is the supported chain for the NFT
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// ContractAddress is the address of the NFT contract the metadata belongs to; supports ENS
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`

	// ForceFetch determines the source of NFT metadata
	// true: fetch from the contract
	// false: fetch from the database
	ForceFetch bool `json:"forceFetch,omitempty" bson:"forceFetch,omitempty" default:"false"`

	// SkipSyncCheck if set to true, the info will be returned regardless of indexer health
	SkipSyncCheck bool `json:"skipSyncCheck,omitempty" bson:"skipSyncCheck,omitempty" default:"false"`

	// TokenID is the token ID of the NFT the metadata belongs to
	// Created by the contract when minting the NFT
	TokenID string `json:"tokenId,omitempty" bson:"tokenId,omitempty"`
}

// NFTMetadataAttributes represents additional information on the NFT
type NFTMetadataAttributes struct {
	// ContractType is the contract type of the NFT (e.g., ERC721, ERC1155)
	ContractType string `json:"contractType" bson:"contractType"`

	// TokenURL is a URL that points to the place storing an NFT's metadata
	TokenURL string `json:"tokenUrl" bson:"tokenUrl"`

	// ImageURL is a URL that points to the actual digital file, usually an IPFS link
	ImageURL string `json:"imageUrl" bson:"imageUrl"`

	// Name is the name of the token
	Name string `json:"name" bson:"name"`

	// Description is the description of the NFT
	Description string `json:"description" bson:"description"`

	// Traits is an array of pre-defined NFT traits
	Traits []NFTTrait `json:"traits" bson:"traits"`
}

// NFTMetadata represents the metadata of an NFT
type NFTMetadata struct {
	// Blockchain is one of the supported chains
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address of the NFT Collection; supports ENS
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// ContractType is the contract type of the NFT (e.g., ERC721, ERC1155)
	ContractType string `json:"contractType" bson:"contractType"`

	// TokenID is the token ID of the NFT the metadata belongs to
	TokenID string `json:"tokenId" bson:"tokenId"`

	// Attributes contains additional information on the NFT
	Attributes NFTMetadataAttributes `json:"attributes" bson:"attributes"`
}

// GetNFTMetadataResp represents the response for ankr_getNFTMetadata
type GetNFTMetadataResp struct {
	// Metadata contains the NFT metadata and attributes
	Metadata NFTMetadata `json:"metadata" bson:"metadata"`
}

// GetNFTHoldersReq represents the request parameters for ankr_getNFTHolders
type GetNFTHoldersReq struct {
	// Blockchain is the supported blockchain for the NFT
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// ContractAddress is the contract address of the NFT collection; supports ENS
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000, default: 1000)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"1000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTHoldersReq) setPageToken(token string) {
	r.PageToken = token
}

// GetNFTHoldersResp represents the response for ankr_getNFTHolders
type GetNFTHoldersResp struct {
	// Holders is a list of wallet addresses that hold the NFT
	Holders []string `json:"holders" bson:"holders"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTHoldersResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetNFTTransfersReq represents the request parameters for ankr_getNftTransfers
type GetNFTTransfersReq struct {
	// Address is an address (or list of addresses) to search for transactions
	Address []string `json:"address,omitempty" bson:"address,omitempty"`

	// Blockchain is the supported blockchains to search in
	Blockchain []Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" bson:"descOrder,omitempty" default:"true"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	FromBlock int64 `json:"fromBlock,omitempty" bson:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	ToBlock int64 `json:"toBlock,omitempty" bson:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty" bson:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty" bson:"toTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get (max: 10000, default: 100)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"100"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTTransfersReq) setPageToken(token string) {
	r.PageToken = token
}

// NFTTransfer represents an NFT transfer transaction
type NFTTransfer struct {
	// BlockHeight is the block number where the transfer occurred
	BlockHeight int64 `json:"blockHeight" bson:"blockHeight"`

	// Blockchain is the blockchain where the transfer occurred
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// CollectionName is the name of the NFT collection
	CollectionName string `json:"collectionName" bson:"collectionName"`

	// CollectionSymbol is the symbol of the NFT collection
	CollectionSymbol string `json:"collectionSymbol" bson:"collectionSymbol"`

	// ContractAddress is the contract address of the NFT
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// FromAddress is the address that sent the NFT
	FromAddress string `json:"fromAddress" bson:"fromAddress"`

	// ImageURL is the URL of the NFT image
	ImageURL string `json:"imageUrl" bson:"imageUrl"`

	// Name is the name of the NFT
	Name string `json:"name" bson:"name"`

	// Timestamp is the timestamp when the transfer occurred
	Timestamp int64 `json:"timestamp" bson:"timestamp"`

	// ToAddress is the address that received the NFT
	ToAddress string `json:"toAddress" bson:"toAddress"`

	// TokenID is the token ID of the NFT
	TokenID string `json:"tokenId" bson:"tokenId"`

	// TransactionHash is the hash of the transaction
	TransactionHash string `json:"transactionHash" bson:"transactionHash"`

	// Type is the type of the NFT contract (e.g., ERC721, ERC1155)
	Type string `json:"type" bson:"type"`

	// Value is the amount/value transferred (for ERC1155 tokens)
	Value string `json:"value" bson:"value"`
}

// GetNFTTransfersResp represents the response for ankr_getNftTransfers
type GetNFTTransfersResp struct {
	// Transfers is a list of NFT transfer transactions
	Transfers []NFTTransfer `json:"transfers" bson:"transfers"`

	// NextPageToken is provided at the end of response for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTTransfersResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetAccountBalanceReq represents the request parameters for ankr_getAccountBalance
type GetAccountBalanceReq struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// NativeFirst sets sorting order - native network token first (true) or not (false)
	NativeFirst *bool `json:"nativeFirst,omitempty" bson:"nativeFirst,omitempty" default:"true"`

	// OnlyWhitelisted shows only tokens listed on CoinGecko (true) or all tokens (false)
	OnlyWhitelisted *bool `json:"onlyWhitelisted,omitempty" bson:"onlyWhitelisted,omitempty" default:"true"`

	// PageSize is the number of results you'd like to get (max: all; default: all)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`

	// WalletAddress is the account address to query for balance; supports ENS
	WalletAddress string `json:"walletAddress,omitempty" bson:"walletAddress,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetAccountBalanceReq) setPageToken(token string) {
	r.PageToken = token
}

// TokenAsset represents a token asset in account balance
type TokenAsset struct {
	// Balance is the token balance
	Balance string `json:"balance" bson:"balance"`

	// BalanceRawInteger is the raw integer balance
	BalanceRawInteger string `json:"balanceRawInteger" bson:"balanceRawInteger"`

	// BalanceUsd is the USD value of the balance
	BalanceUsd string `json:"balanceUsd" bson:"balanceUsd"`

	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the token contract address
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// HolderAddress is the address holding the token
	HolderAddress string `json:"holderAddress" bson:"holderAddress"`

	// Thumbnail is the token thumbnail image URL
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals" bson:"tokenDecimals"`

	// TokenName is the name of the token
	TokenName string `json:"tokenName" bson:"tokenName"`

	// TokenPrice is the USD price of the token
	TokenPrice string `json:"tokenPrice" bson:"tokenPrice"`

	// TokenSymbol is the symbol of the token
	TokenSymbol string `json:"tokenSymbol" bson:"tokenSymbol"`

	// TokenType is the type of the token
	TokenType string `json:"tokenType" bson:"tokenType"`
}

// GetAccountBalanceResp represents the response for ankr_getAccountBalance
type GetAccountBalanceResp struct {
	// Assets is the list of token assets
	Assets []TokenAsset `json:"assets" bson:"assets"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`

	// TotalBalanceUsd is the total USD value of all assets
	TotalBalanceUsd string `json:"totalBalanceUsd" bson:"totalBalanceUsd"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetAccountBalanceResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetCurrenciesReq represents the request parameters for ankr_getCurrencies
type GetCurrenciesReq struct {
	// Blockchain is the supported chain to get currencies for
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`
}

// Currency represents a currency on a blockchain
type Currency struct {
	// Address is the contract address of the currency
	Address string `json:"address" bson:"address"`

	// Blockchain is the blockchain where the currency exists
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// Decimals is the number of decimals for the currency
	Decimals int32 `json:"decimals" bson:"decimals"`

	// Name is the name of the currency
	Name string `json:"name" bson:"name"`

	// Symbol is the symbol of the currency
	Symbol string `json:"symbol" bson:"symbol"`

	// Thumbnail is the currency thumbnail image URL
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`
}

// GetCurrenciesResp represents the response for ankr_getCurrencies
type GetCurrenciesResp struct {
	// Currencies is the list of currencies on the blockchain
	Currencies []Currency `json:"currencies" bson:"currencies"`
}

// GetTokenPriceReq represents the request parameters for ankr_getTokenPrice
type GetTokenPriceReq struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	// If not provided, returns the native coin price of the blockchain specified
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`
}

// GetTokenPriceResp represents the response for ankr_getTokenPrice
type GetTokenPriceResp struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// UsdPrice is the USD price of the token
	UsdPrice string `json:"usdPrice" bson:"usdPrice"`
}

// GetTokenHoldersReq represents the request parameters for ankr_getTokenHolders
type GetTokenHoldersReq struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenHoldersReq) setPageToken(token string) {
	r.PageToken = token
}

// TokenHolder represents a token holder
type TokenHolder struct {
	// Balance is the token balance held by the address
	Balance string `json:"balance" bson:"balance"`

	// BalanceRawInteger is the raw integer balance
	BalanceRawInteger string `json:"balanceRawInteger" bson:"balanceRawInteger"`

	// HolderAddress is the address holding the tokens
	HolderAddress string `json:"holderAddress" bson:"holderAddress"`
}

// GetTokenHoldersResp represents the response for ankr_getTokenHolders
type GetTokenHoldersResp struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// Holders is the list of token holders
	Holders []TokenHolder `json:"holders" bson:"holders"`

	// HoldersCount is the total number of holders
	HoldersCount int64 `json:"holdersCount" bson:"holdersCount"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals" bson:"tokenDecimals"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenHoldersResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetTokenHoldersCountReq represents the request parameters for ankr_getTokenHoldersCount
type GetTokenHoldersCountReq struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenHoldersCountReq) setPageToken(token string) {
	r.PageToken = token
}

// HolderCountHistory represents the holder count history
type HolderCountHistory struct {
	// HolderCount is the number of holders at this point in time
	HolderCount int64 `json:"holderCount" bson:"holderCount"`

	// LastUpdatedAt is the timestamp when this data was last updated
	LastUpdatedAt string `json:"lastUpdatedAt" bson:"lastUpdatedAt"`

	// TotalAmount is the total amount of tokens
	TotalAmount string `json:"totalAmount" bson:"totalAmount"`

	// TotalAmountRawInteger is the raw integer total amount
	TotalAmountRawInteger string `json:"totalAmountRawInteger" bson:"totalAmountRawInteger"`
}

// GetTokenHoldersCountResp represents the response for ankr_getTokenHoldersCount
type GetTokenHoldersCountResp struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// HolderCountHistory is the history of holder counts over time
	HolderCountHistory []HolderCountHistory `json:"holderCountHistory" bson:"holderCountHistory"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals" bson:"tokenDecimals"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenHoldersCountResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetTokenTransfersReq represents the request parameters for ankr_getTokenTransfers
type GetTokenTransfersReq struct {
	// Address is an address or list of addresses to search for token transfers
	Address []string `json:"address,omitempty" bson:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" bson:"descOrder,omitempty" default:"true"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty" bson:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty" bson:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the UNIX timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty" bson:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the UNIX timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty" bson:"toTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided in the response
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenTransfersReq) setPageToken(token string) {
	r.PageToken = token
}

// TokenTransfer represents a token transfer transaction
type TokenTransfer struct {
	// BlockHeight is the block number where the transfer occurred
	BlockHeight int64 `json:"blockHeight" bson:"blockHeight"`

	// Blockchain is the blockchain where the transfer occurred
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress" bson:"contractAddress"`

	// FromAddress is the address that sent the tokens
	FromAddress string `json:"fromAddress" bson:"fromAddress"`

	// Thumbnail is the token thumbnail image URL
	Thumbnail string `json:"thumbnail" bson:"thumbnail"`

	// Timestamp is the timestamp when the transfer occurred
	Timestamp int64 `json:"timestamp" bson:"timestamp"`

	// ToAddress is the address that received the tokens
	ToAddress string `json:"toAddress" bson:"toAddress"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals" bson:"tokenDecimals"`

	// TokenName is the name of the token
	TokenName string `json:"tokenName" bson:"tokenName"`

	// TokenSymbol is the symbol of the token
	TokenSymbol string `json:"tokenSymbol" bson:"tokenSymbol"`

	// TransactionHash is the hash of the transaction
	TransactionHash string `json:"transactionHash" bson:"transactionHash"`

	// Value is the amount of tokens transferred
	Value string `json:"value" bson:"value"`

	// ValueRawInteger is the raw integer value
	ValueRawInteger string `json:"valueRawInteger" bson:"valueRawInteger"`
}

// GetTokenTransfersResp represents the response for ankr_getTokenTransfers
type GetTokenTransfersResp struct {
	// Transfers is a list of token transfer transactions
	Transfers []TokenTransfer `json:"transfers" bson:"transfers"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenTransfersResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetBlockchainStatsReq represents the request parameters for ankr_getBlockchainStats
type GetBlockchainStatsReq struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`
}

// BlockchainStat represents blockchain statistics
type BlockchainStat struct {
	// Blockchain is the blockchain identifier
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// TotalTransactionsCount is the total number of transactions
	TotalTransactionsCount int64 `json:"totalTransactionsCount" bson:"totalTransactionsCount"`

	// TotalEventsCount is the total number of events
	TotalEventsCount int64 `json:"totalEventsCount" bson:"totalEventsCount"`

	// LatestBlockNumber is the latest block number
	LatestBlockNumber int64 `json:"latestBlockNumber" bson:"latestBlockNumber"`

	// BlockTimeMs is the average block time in milliseconds
	BlockTimeMs int64 `json:"blockTimeMs" bson:"blockTimeMs"`

	// NativeCoinUsdPrice is the USD price of the native coin
	NativeCoinUsdPrice string `json:"nativeCoinUsdPrice" bson:"nativeCoinUsdPrice"`
}

// GetBlockchainStatsResp represents the response for ankr_getBlockchainStats
type GetBlockchainStatsResp struct {
	// Stats is the list of blockchain statistics
	Stats []BlockchainStat `json:"stats" bson:"stats"`
}

// GetBlocksReq represents the request parameters for ankr_getBlocks
type GetBlocksReq struct {
	// Blockchain is the supported chain for the blocks
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" bson:"decodeLogs,omitempty" default:"false"`

	// DecodeTxData sets to true to decode transaction data, or to false if not interested in it
	DecodeTxData bool `json:"decodeTxData,omitempty" bson:"decodeTxData,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" bson:"descOrder,omitempty" default:"true"`

	// FromBlock is the first block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty" bson:"fromBlock,omitempty"`

	// ToBlock is the last block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty" bson:"toBlock,omitempty"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	// Note that logs are stored inside transactions, so make sure includeTxs is also set to true
	IncludeLogs bool `json:"includeLogs,omitempty" bson:"includeLogs,omitempty" default:"false"`

	// IncludeTxs sets to true to include transactions, or to false to exclude them
	IncludeTxs *bool `json:"includeTxs,omitempty" bson:"includeTxs,omitempty" default:"true"`
}

// EthBlockDetails represents Ethereum block details
type EthBlockDetails struct {
	// Difficulty is the block difficulty
	Difficulty string `json:"difficulty" bson:"difficulty"`

	// ExtraData is the extra data field
	ExtraData string `json:"extraData" bson:"extraData"`

	// GasLimit is the gas limit for the block
	GasLimit int64 `json:"gasLimit" bson:"gasLimit"`

	// GasUsed is the gas used in the block
	GasUsed int64 `json:"gasUsed" bson:"gasUsed"`

	// Miner is the miner address
	Miner string `json:"miner" bson:"miner"`

	// Nonce is the block nonce
	Nonce string `json:"nonce" bson:"nonce"`

	// Sha3Uncles is the SHA3 of uncles
	Sha3Uncles string `json:"sha3Uncles" bson:"sha3Uncles"`

	// Size is the block size
	Size string `json:"size" bson:"size"`

	// StateRoot is the state root hash
	StateRoot string `json:"stateRoot" bson:"stateRoot"`

	// TotalDifficulty is the total difficulty
	TotalDifficulty string `json:"totalDifficulty" bson:"totalDifficulty"`
}

// BlockDetails represents block details
type BlockDetails struct {
	// EthBlock contains Ethereum-specific block details
	EthBlock EthBlockDetails `json:"ethBlock" bson:"ethBlock"`
}

// Block represents a blockchain block
type Block struct {
	// BlockHash is the hash of the block
	BlockHash string `json:"hash" bson:"hash"`

	// Number is the height of the block
	Number string `json:"number" bson:"number"`

	// BlockchainLogo is the logo URL of the blockchain
	BlockchainLogo string `json:"blockchainLogo,omitempty" bson:"blockchainLogo,omitempty"`

	// BlockchainName is the name of the blockchain
	BlockchainName string `json:"blockchain" bson:"blockchain"`

	// Details contains blockchain-specific block details
	Details BlockDetails `json:"details" bson:"details"`

	// LogsBloom is the bloom filter for the logs of the block
	LogsBloom string `json:"logsBloom" bson:"logsBloom"`

	// MixHash is the mixed hash
	MixHash string `json:"mixHash" bson:"mixHash"`

	// Nonce is the block nonce
	Nonce string `json:"nonce" bson:"nonce"`

	// ParentHash is the hash of the parent block
	ParentHash string `json:"parentHash" bson:"parentHash"`

	// ReceiptsRoot is the root hash of the receipts trie
	ReceiptsRoot string `json:"receiptsRoot" bson:"receiptsRoot"`

	// Sha3Uncles is the SHA3 of uncles
	Sha3Uncles string `json:"sha3Uncles" bson:"sha3Uncles"`

	// StateRoot is the state root hash
	StateRoot string `json:"stateRoot" bson:"stateRoot"`

	// Miner is the miner address
	Miner string `json:"miner" bson:"miner"`

	// Difficulty is the block difficulty
	Difficulty string `json:"difficulty" bson:"difficulty"`

	// ExtraData is the extra data field
	ExtraData string `json:"extraData" bson:"extraData"`

	// Size is the block size
	Size string `json:"size" bson:"size"`

	// GasLimit is the gas limit for the block
	GasLimit string `json:"gasLimit" bson:"gasLimit"`

	// GasUsed is the gas used in the block
	GasUsed string `json:"gasUsed" bson:"gasUsed"`

	// Timestamp is the timestamp of the block
	Timestamp string `json:"timestamp" bson:"timestamp"`

	// TransactionsRoot is the root hash of the transaction trie
	TransactionsRoot string `json:"transactionsRoot" bson:"transactionsRoot"`

	// TotalDifficulty is the total difficulty
	TotalDifficulty string `json:"totalDifficulty" bson:"totalDifficulty"`

	// TransactionsCount is the number of transactions in the block
	TransactionsCount int32 `json:"transactionsCount" bson:"transactionsCount"`

	// Transactions contains the transactions in the block
	Transactions []Tx `json:"transactions" bson:"transactions"`

	// Uncles contains the uncles block hashes
	Uncles []any `json:"uncles" bson:"uncles"`
}

// GetBlocksResp represents the response for ankr_getBlocks
type GetBlocksResp struct {
	// Blocks is the list of blocks
	Blocks []Block `json:"blocks" bson:"blocks"`
}

// GetLogsReq represents the request parameters for ankr_getLogs
type GetLogsReq struct {
	// Address is a contract address or list of addresses from which the logs originate
	// Supported value formats: hex or array of hexes
	Address []string `json:"address,omitempty" bson:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" bson:"decodeLogs,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" bson:"descOrder,omitempty" default:"true"`

	// FromBlock is the first block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty" bson:"fromBlock,omitempty"`

	// FromTimestamp is the first timestamp of the range
	FromTimestamp int64 `json:"fromTimestamp,omitempty" bson:"fromTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get (max: 10000; default: 1000)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`

	// ToBlock is the last block included in the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty" bson:"toBlock,omitempty"`

	// ToTimestamp is the last timestamp of the range
	ToTimestamp int64 `json:"toTimestamp,omitempty" bson:"toTimestamp,omitempty"`

	// Topics is the data the log contains
	Topics [][]string `json:"topics,omitempty" bson:"topics,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetLogsReq) setPageToken(token string) {
	r.PageToken = token
}

// EventInput represents an event input parameter
type EventInput struct {
	// Indexed indicates if the input is indexed
	Indexed bool `json:"indexed" bson:"indexed"`

	// Name is the name of the input
	Name string `json:"name" bson:"name"`

	// Size is the size of the input
	Size int32 `json:"size" bson:"size"`

	// Type is the type of the input
	Type string `json:"type" bson:"type"`

	// ValueDecoded is the decoded value of the input
	ValueDecoded string `json:"valueDecoded" bson:"valueDecoded"`
}

// Event represents a decoded event
type Event struct {
	// Anonymous indicates if the event is anonymous
	Anonymous bool `json:"anonymous" bson:"anonymous"`

	// ID is the event ID
	ID string `json:"id" bson:"id"`

	// Inputs is the list of event inputs
	Inputs []EventInput `json:"inputs" bson:"inputs"`

	// Name is the name of the event
	Name string `json:"name" bson:"name"`

	// Signature is the event signature
	Signature string `json:"signature" bson:"signature"`

	// String is the string representation
	String string `json:"string" bson:"string"`

	// Verified indicates if the event is verified
	Verified bool `json:"verified" bson:"verified"`
}

// Log represents a blockchain log entry
type Log struct {
	// Address is the contract address that emitted the log
	Address string `json:"address" bson:"address"`

	// BlockHash is the hash of the block containing the log
	BlockHash string `json:"blockHash" bson:"blockHash"`

	// BlockNumber is the number of the block containing the log
	BlockNumber string `json:"blockNumber" bson:"blockNumber"`

	// Data is the log data
	Data string `json:"data" bson:"data"`

	// Event contains decoded event information
	Event Event `json:"event,omitempty" bson:"event,omitempty"`

	// LogIndex is the index of the log in the block
	LogIndex string `json:"logIndex" bson:"logIndex"`

	// Removed indicates if the log was removed
	Removed bool `json:"removed" bson:"removed"`

	// Topics is the list of log topics
	Topics []string `json:"topics" bson:"topics"`

	// TransactionHash is the hash of the transaction containing the log
	TransactionHash string `json:"transactionHash" bson:"transactionHash"`

	// TransactionIndex is the index of the transaction in the block
	TransactionIndex string `json:"transactionIndex" bson:"transactionIndex"`
}

// GetLogsResp represents the response for ankr_getLogs
type GetLogsResp struct {
	// Logs is the list of log entries
	Logs []Log `json:"logs" bson:"logs"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetLogsResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetTxsByHashReq represents the request parameters for ankr_getTransactionsByHash
type GetTxsByHashReq struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// TransactionHash is the hash of the transaction you'd like to request details for
	TransactionHash string `json:"transactionHash,omitempty" bson:"transactionHash,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" bson:"decodeLogs,omitempty" default:"false"`

	// DecodeTxData sets to true to decode transaction data, or to false if not interested in it
	DecodeTxData bool `json:"decodeTxData,omitempty" bson:"decodeTxData,omitempty" default:"false"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	IncludeLogs bool `json:"includeLogs,omitempty" bson:"includeLogs,omitempty" default:"false"`
}

// MethodInput represents a method input parameter
type MethodInput struct {
	// Name is the name of the input
	Name string `json:"name" bson:"name"`

	// Size is the size of the input
	Size int32 `json:"size" bson:"size"`

	// Type is the type of the input
	Type string `json:"type" bson:"type"`

	// ValueDecoded is the decoded value of the input
	ValueDecoded string `json:"valueDecoded" bson:"valueDecoded"`
}

// Method represents a decoded method
type Method struct {
	// ID is the method ID
	ID string `json:"id" bson:"id"`

	// Inputs is the list of method inputs
	Inputs []MethodInput `json:"inputs" bson:"inputs"`

	// Name is the name of the method
	Name string `json:"name" bson:"name"`

	// Signature is the method signature
	Signature string `json:"signature" bson:"signature"`

	// String is the string representation
	String string `json:"string" bson:"string"`

	// Verified indicates if the method is verified
	Verified bool `json:"verified" bson:"verified"`
}

// Tx represents a blockchain transaction
type Tx struct {
	// BlockHash is the hash of the block containing the transaction
	BlockHash string `json:"blockHash" bson:"blockHash"`

	// BlockNumber is the number of the block containing the transaction
	BlockNumber string `json:"blockNumber" bson:"blockNumber"`

	// Blockchain is the blockchain where the transaction occurred
	Blockchain string `json:"blockchain" bson:"blockchain"`

	// ContractAddress is the contract address (for contract creation transactions)
	ContractAddress string `json:"contractAddress,omitempty" bson:"contractAddress,omitempty"`

	// CumulativeGasUsed is the cumulative gas used
	CumulativeGasUsed string `json:"cumulativeGasUsed" bson:"cumulativeGasUsed"`

	// From is the sender address
	From string `json:"from" bson:"from"`

	// Gas is the gas limit
	Gas string `json:"gas" bson:"gas"`

	// GasPrice is the gas price
	GasPrice string `json:"gasPrice" bson:"gasPrice"`

	// GasUsed is the gas used
	GasUsed string `json:"gasUsed" bson:"gasUsed"`

	// Hash is the transaction hash
	Hash string `json:"hash" bson:"hash"`

	// Input is the transaction input data
	Input string `json:"input" bson:"input"`

	// Logs is the list of logs emitted by the transaction
	Logs []Log `json:"logs,omitempty" bson:"logs,omitempty"`

	// LogsBloom is the logs bloom filter
	LogsBloom string `json:"logsBloom" bson:"logsBloom"`

	// Method contains decoded method information
	Method Method `json:"method" bson:"method"`

	// Nonce is the transaction nonce
	Nonce string `json:"nonce" bson:"nonce"`

	// R is the R signature component
	R string `json:"r" bson:"r"`

	// S is the S signature component
	S string `json:"s" bson:"s"`

	// Status is the transaction status
	Status string `json:"status" bson:"status"`

	// Timestamp is the timestamp of the transaction
	Timestamp string `json:"timestamp" bson:"timestamp"`

	// To is the recipient address
	To string `json:"to" bson:"to"`

	// TransactionHash is the transaction hash
	TransactionHash string `json:"transactionHash" bson:"transactionHash"`

	// TransactionIndex is the index of the transaction in the block
	TransactionIndex string `json:"transactionIndex" bson:"transactionIndex"`

	// Type is the transaction type
	Type string `json:"type" bson:"type"`

	// V is the V signature component
	V string `json:"v" bson:"v"`

	// Value is the transaction value
	Value string `json:"value" bson:"value"`
}

// GetTxsByHashResp represents the response for ankr_getTransactionsByHash
type GetTxsByHashResp struct {
	// Transactions is the list of transactions
	Transactions []Tx `json:"transactions" bson:"transactions"`
}

// GetTxsByAddressReq represents the request parameters for ankr_getTransactionsByAddress
type GetTxsByAddressReq struct {
	// Address is the address to search for transactions
	Address string `json:"address,omitempty" bson:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty" bson:"blockchain,omitempty"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty" bson:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty" bson:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty" bson:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty" bson:"toTimestamp,omitempty"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	IncludeLogs bool `json:"includeLogs,omitempty" bson:"includeLogs,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" bson:"descOrder,omitempty" default:"true"`

	// PageSize is the number of result pages you'd like to get (max: 10000; default: 100)
	PageSize int32 `json:"pageSize,omitempty" bson:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty" bson:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTxsByAddressReq) setPageToken(token string) {
	r.PageToken = token
}

// GetTxsByAddressResp represents the response for ankr_getTransactionsByAddress
type GetTxsByAddressResp struct {
	// Transactions is the list of transactions
	Transactions []Tx `json:"transactions" bson:"transactions"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken" bson:"nextPageToken"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTxsByAddressResp) getNextPageToken() string {
	return r.NextPageToken
}

// GetInteractionsReq represents the request parameters for ankr_getInteractions
type GetInteractionsReq struct {
	// Address is the address of the wallet or contract that created the logs
	Address string `json:"address,omitempty" bson:"address,omitempty"`
}

// GetInteractionsResp represents the response for ankr_getInteractions
type GetInteractionsResp struct {
	// Blockchains is the list of blockchains interacted with the address
	Blockchains []string `json:"blockchains" bson:"blockchains"`
}
