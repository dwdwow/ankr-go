package ankr

func PtrFalse() *bool {
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
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type RPCRespBody[Result any] struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int64         `json:"id"`
	Result  Result        `json:"result"`
	Error   *RPCRespError `json:"error"`
}

type RPCRespError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// GetNFTsByOwnerRequest represents the request parameters for ankr_getNFTsByOwner
type GetNFTsByOwnerRequest struct {
	// WalletAddress is the account address to query for NFTs; supports ENS
	WalletAddress string `json:"walletAddress,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// PageSize is the number of page results (default=10, max=50)
	PageSize int32 `json:"pageSize,omitempty" default:"50"`

	// PageToken is provided at the end of response and used to fetch next page
	PageToken string `json:"pageToken,omitempty"`

	// Filter filters the request by smart contract address and optional NFT ID
	// Format: map[string][]string
	// Example: {"0xd8682bfa6918b0174f287b888e765b9a1b4dc9c3": []} - all NFTs from address
	// Example: {"0xd8682bfa6918b0174f287b888e765b9a1b4dc9c3": ["8937"]} - specific NFT
	Filter map[string][]string `json:"filter,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTsByOwnerRequest) setPageToken(token string) {
	r.PageToken = token
}

// NFTTrait represents a trait/attribute of an NFT
type NFTTrait struct {
	// TraitType is the trait's descriptive name
	TraitType string `json:"trait_type"`

	// Value is the value description
	Value string `json:"value"`
}

// NFT represents an NFT asset
type NFT struct {
	// Blockchain is one of the supported chains
	Blockchain string `json:"blockchain"`

	// CollectionName is the collection name the NFT belongs to
	CollectionName string `json:"collectionName,omitempty"`

	// ContractAddress is the NFT collection's EVM-compatible contract address
	ContractAddress string `json:"contractAddress"`

	// ContractType is the type of contract - either ERC721 or ERC1155
	ContractType string `json:"contractType"`

	// Name is the name of the NFT asset
	Name string `json:"name,omitempty"`

	// TokenID is the ID of the NFT asset
	TokenID string `json:"tokenId"`

	// ImageURL is a URL that points to the actual digital file, usually an IPFS link
	ImageURL string `json:"imageUrl,omitempty"`

	// Symbol is the symbol of the NFT asset
	Symbol string `json:"symbol,omitempty"`

	// Traits are the attributes of the NFT asset
	Traits []NFTTrait `json:"traits,omitempty"`
}

// GetNFTsByOwnerResponse represents the response for ankr_getNFTsByOwner
type GetNFTsByOwnerResponse struct {
	// Assets is the list of NFT assets
	Assets []NFT `json:"assets"`

	// NextPageToken is provided at the end of response for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTsByOwnerResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetNFTMetadataRequest represents the request parameters for ankr_getNFTMetadata
type GetNFTMetadataRequest struct {
	// Blockchain is the supported chain for the NFT
	Blockchain Chain `json:"blockchain,omitempty"`

	// ContractAddress is the address of the NFT contract the metadata belongs to; supports ENS
	ContractAddress string `json:"contractAddress,omitempty"`

	// ForceFetch determines the source of NFT metadata
	// true: fetch from the contract
	// false: fetch from the database
	ForceFetch bool `json:"forceFetch,omitempty" default:"false"`

	// SkipSyncCheck if set to true, the info will be returned regardless of indexer health
	SkipSyncCheck bool `json:"skipSyncCheck,omitempty" default:"false"`

	// TokenID is the token ID of the NFT the metadata belongs to
	// Created by the contract when minting the NFT
	TokenID string `json:"tokenId,omitempty"`
}

// NFTMetadataAttributes represents additional information on the NFT
type NFTMetadataAttributes struct {
	// ContractType is the contract type of the NFT (e.g., ERC721, ERC1155)
	ContractType string `json:"contractType"`

	// TokenURL is a URL that points to the place storing an NFT's metadata
	TokenURL string `json:"tokenUrl,omitempty"`

	// ImageURL is a URL that points to the actual digital file, usually an IPFS link
	ImageURL string `json:"imageUrl,omitempty"`

	// Name is the name of the token
	Name string `json:"name,omitempty"`

	// Description is the description of the NFT
	Description string `json:"description,omitempty"`

	// Traits is an array of pre-defined NFT traits
	Traits []NFTTrait `json:"traits,omitempty"`
}

// NFTMetadata represents the metadata of an NFT
type NFTMetadata struct {
	// Blockchain is one of the supported chains
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address of the NFT Collection; supports ENS
	ContractAddress string `json:"contractAddress"`

	// ContractType is the contract type of the NFT (e.g., ERC721, ERC1155)
	ContractType string `json:"contractType"`

	// TokenID is the token ID of the NFT the metadata belongs to
	TokenID string `json:"tokenId"`

	// Attributes contains additional information on the NFT
	Attributes NFTMetadataAttributes `json:"attributes"`
}

// GetNFTMetadataResponse represents the response for ankr_getNFTMetadata
type GetNFTMetadataResponse struct {
	// Metadata contains the NFT metadata and attributes
	Metadata NFTMetadata `json:"metadata"`
}

// GetNFTHoldersRequest represents the request parameters for ankr_getNFTHolders
type GetNFTHoldersRequest struct {
	// Blockchain is the supported blockchain for the NFT
	Blockchain Chain `json:"blockchain,omitempty"`

	// ContractAddress is the contract address of the NFT collection; supports ENS
	ContractAddress string `json:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000, default: 1000)
	PageSize int32 `json:"pageSize,omitempty" default:"1000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTHoldersRequest) setPageToken(token string) {
	r.PageToken = token
}

// GetNFTHoldersResponse represents the response for ankr_getNFTHolders
type GetNFTHoldersResponse struct {
	// Holders is a list of wallet addresses that hold the NFT
	Holders []string `json:"holders"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTHoldersResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetNFTTransfersRequest represents the request parameters for ankr_getNftTransfers
type GetNFTTransfersRequest struct {
	// Address is an address (or list of addresses) to search for transactions
	Address []string `json:"address,omitempty"`

	// Blockchain is the supported blockchains to search in
	Blockchain []Chain `json:"blockchain,omitempty"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" default:"true"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	FromBlock int64 `json:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	ToBlock int64 `json:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get (max: 10000, default: 100)
	PageSize int32 `json:"pageSize,omitempty" default:"100"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetNFTTransfersRequest) setPageToken(token string) {
	r.PageToken = token
}

// NFTTransfer represents an NFT transfer transaction
type NFTTransfer struct {
	// BlockHeight is the block number where the transfer occurred
	BlockHeight int64 `json:"blockHeight"`

	// Blockchain is the blockchain where the transfer occurred
	Blockchain string `json:"blockchain"`

	// CollectionName is the name of the NFT collection
	CollectionName string `json:"collectionName,omitempty"`

	// CollectionSymbol is the symbol of the NFT collection
	CollectionSymbol string `json:"collectionSymbol,omitempty"`

	// ContractAddress is the contract address of the NFT
	ContractAddress string `json:"contractAddress"`

	// FromAddress is the address that sent the NFT
	FromAddress string `json:"fromAddress"`

	// ImageURL is the URL of the NFT image
	ImageURL string `json:"imageUrl,omitempty"`

	// Name is the name of the NFT
	Name string `json:"name,omitempty"`

	// Timestamp is the timestamp when the transfer occurred
	Timestamp int64 `json:"timestamp"`

	// ToAddress is the address that received the NFT
	ToAddress string `json:"toAddress"`

	// TokenID is the token ID of the NFT
	TokenID string `json:"tokenId"`

	// TransactionHash is the hash of the transaction
	TransactionHash string `json:"transactionHash"`

	// Type is the type of the NFT contract (e.g., ERC721, ERC1155)
	Type string `json:"type"`

	// Value is the amount/value transferred (for ERC1155 tokens)
	Value string `json:"value"`
}

// GetNFTTransfersResponse represents the response for ankr_getNftTransfers
type GetNFTTransfersResponse struct {
	// Transfers is a list of NFT transfer transactions
	Transfers []NFTTransfer `json:"transfers"`

	// NextPageToken is provided at the end of response for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetNFTTransfersResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetAccountBalanceRequest represents the request parameters for ankr_getAccountBalance
type GetAccountBalanceRequest struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// NativeFirst sets sorting order - native network token first (true) or not (false)
	NativeFirst *bool `json:"nativeFirst,omitempty" default:"true"`

	// OnlyWhitelisted shows only tokens listed on CoinGecko (true) or all tokens (false)
	OnlyWhitelisted *bool `json:"onlyWhitelisted,omitempty" default:"true"`

	// PageSize is the number of results you'd like to get (max: all; default: all)
	PageSize int32 `json:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`

	// WalletAddress is the account address to query for balance; supports ENS
	WalletAddress string `json:"walletAddress,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetAccountBalanceRequest) setPageToken(token string) {
	r.PageToken = token
}

// TokenAsset represents a token asset in account balance
type TokenAsset struct {
	// Balance is the token balance
	Balance string `json:"balance"`

	// BalanceRawInteger is the raw integer balance
	BalanceRawInteger string `json:"balanceRawInteger"`

	// BalanceUsd is the USD value of the balance
	BalanceUsd string `json:"balanceUsd"`

	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain"`

	// ContractAddress is the token contract address
	ContractAddress string `json:"contractAddress"`

	// HolderAddress is the address holding the token
	HolderAddress string `json:"holderAddress"`

	// Thumbnail is the token thumbnail image URL
	Thumbnail string `json:"thumbnail,omitempty"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals"`

	// TokenName is the name of the token
	TokenName string `json:"tokenName"`

	// TokenPrice is the USD price of the token
	TokenPrice string `json:"tokenPrice"`

	// TokenSymbol is the symbol of the token
	TokenSymbol string `json:"tokenSymbol"`

	// TokenType is the type of the token
	TokenType string `json:"tokenType"`
}

// GetAccountBalanceResponse represents the response for ankr_getAccountBalance
type GetAccountBalanceResponse struct {
	// Assets is the list of token assets
	Assets []TokenAsset `json:"assets"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`

	// TotalBalanceUsd is the total USD value of all assets
	TotalBalanceUsd string `json:"totalBalanceUsd"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetAccountBalanceResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetCurrenciesRequest represents the request parameters for ankr_getCurrencies
type GetCurrenciesRequest struct {
	// Blockchain is the supported chain to get currencies for
	Blockchain Chain `json:"blockchain,omitempty"`
}

// Currency represents a currency on a blockchain
type Currency struct {
	// Address is the contract address of the currency
	Address string `json:"address"`

	// Blockchain is the blockchain where the currency exists
	Blockchain string `json:"blockchain"`

	// Decimals is the number of decimals for the currency
	Decimals int32 `json:"decimals"`

	// Name is the name of the currency
	Name string `json:"name"`

	// Symbol is the symbol of the currency
	Symbol string `json:"symbol"`

	// Thumbnail is the currency thumbnail image URL
	Thumbnail string `json:"thumbnail,omitempty"`
}

// GetCurrenciesResponse represents the response for ankr_getCurrencies
type GetCurrenciesResponse struct {
	// Currencies is the list of currencies on the blockchain
	Currencies []Currency `json:"currencies"`
}

// GetTokenPriceRequest represents the request parameters for ankr_getTokenPrice
type GetTokenPriceRequest struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	// If not provided, returns the native coin price of the blockchain specified
	ContractAddress string `json:"contractAddress,omitempty"`
}

// GetTokenPriceResponse represents the response for ankr_getTokenPrice
type GetTokenPriceResponse struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress"`

	// UsdPrice is the USD price of the token
	UsdPrice string `json:"usdPrice"`
}

// GetTokenHoldersRequest represents the request parameters for ankr_getTokenHolders
type GetTokenHoldersRequest struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	ContractAddress string `json:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenHoldersRequest) setPageToken(token string) {
	r.PageToken = token
}

// TokenHolder represents a token holder
type TokenHolder struct {
	// Balance is the token balance held by the address
	Balance string `json:"balance"`

	// BalanceRawInteger is the raw integer balance
	BalanceRawInteger string `json:"balanceRawInteger"`

	// HolderAddress is the address holding the tokens
	HolderAddress string `json:"holderAddress"`
}

// GetTokenHoldersResponse represents the response for ankr_getTokenHolders
type GetTokenHoldersResponse struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress"`

	// Holders is the list of token holders
	Holders []TokenHolder `json:"holders"`

	// HoldersCount is the total number of holders
	HoldersCount int64 `json:"holdersCount"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenHoldersResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetTokenHoldersCountRequest represents the request parameters for ankr_getTokenHoldersCount
type GetTokenHoldersCountRequest struct {
	// Blockchain is the supported chain for the token
	Blockchain Chain `json:"blockchain,omitempty"`

	// ContractAddress is the address of the token contract; supports ENS
	ContractAddress string `json:"contractAddress,omitempty"`

	// PageSize is the number of results you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenHoldersCountRequest) setPageToken(token string) {
	r.PageToken = token
}

// HolderCountHistory represents the holder count history
type HolderCountHistory struct {
	// HolderCount is the number of holders at this point in time
	HolderCount int64 `json:"holderCount"`

	// LastUpdatedAt is the timestamp when this data was last updated
	LastUpdatedAt string `json:"lastUpdatedAt"`

	// TotalAmount is the total amount of tokens
	TotalAmount string `json:"totalAmount"`

	// TotalAmountRawInteger is the raw integer total amount
	TotalAmountRawInteger string `json:"totalAmountRawInteger"`
}

// GetTokenHoldersCountResponse represents the response for ankr_getTokenHoldersCount
type GetTokenHoldersCountResponse struct {
	// Blockchain is the blockchain where the token exists
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress"`

	// HolderCountHistory is the history of holder counts over time
	HolderCountHistory []HolderCountHistory `json:"holderCountHistory"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenHoldersCountResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetTokenTransfersRequest represents the request parameters for ankr_getTokenTransfers
type GetTokenTransfersRequest struct {
	// Address is an address or list of addresses to search for token transfers
	Address []string `json:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" default:"true"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the UNIX timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the UNIX timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get (max: 10000; default: 10000)
	PageSize int32 `json:"pageSize,omitempty" default:"10000"`

	// PageToken is the current page token provided in the response
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTokenTransfersRequest) setPageToken(token string) {
	r.PageToken = token
}

// TokenTransfer represents a token transfer transaction
type TokenTransfer struct {
	// BlockHeight is the block number where the transfer occurred
	BlockHeight int64 `json:"blockHeight"`

	// Blockchain is the blockchain where the transfer occurred
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address of the token
	ContractAddress string `json:"contractAddress"`

	// FromAddress is the address that sent the tokens
	FromAddress string `json:"fromAddress"`

	// Thumbnail is the token thumbnail image URL
	Thumbnail string `json:"thumbnail,omitempty"`

	// Timestamp is the timestamp when the transfer occurred
	Timestamp int64 `json:"timestamp"`

	// ToAddress is the address that received the tokens
	ToAddress string `json:"toAddress"`

	// TokenDecimals is the number of decimals for the token
	TokenDecimals int32 `json:"tokenDecimals"`

	// TokenName is the name of the token
	TokenName string `json:"tokenName"`

	// TokenSymbol is the symbol of the token
	TokenSymbol string `json:"tokenSymbol"`

	// TransactionHash is the hash of the transaction
	TransactionHash string `json:"transactionHash"`

	// Value is the amount of tokens transferred
	Value string `json:"value"`

	// ValueRawInteger is the raw integer value
	ValueRawInteger string `json:"valueRawInteger"`
}

// GetTokenTransfersResponse represents the response for ankr_getTokenTransfers
type GetTokenTransfersResponse struct {
	// Transfers is a list of token transfer transactions
	Transfers []TokenTransfer `json:"transfers"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTokenTransfersResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetBlockchainStatsRequest represents the request parameters for ankr_getBlockchainStats
type GetBlockchainStatsRequest struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`
}

// BlockchainStat represents blockchain statistics
type BlockchainStat struct {
	// Blockchain is the blockchain identifier
	Blockchain string `json:"blockchain"`

	// TotalTransactionsCount is the total number of transactions
	TotalTransactionsCount int64 `json:"totalTransactionsCount"`

	// TotalEventsCount is the total number of events
	TotalEventsCount int64 `json:"totalEventsCount"`

	// LatestBlockNumber is the latest block number
	LatestBlockNumber int64 `json:"latestBlockNumber"`

	// BlockTimeMs is the average block time in milliseconds
	BlockTimeMs int64 `json:"blockTimeMs"`

	// NativeCoinUsdPrice is the USD price of the native coin
	NativeCoinUsdPrice string `json:"nativeCoinUsdPrice"`
}

// GetBlockchainStatsResponse represents the response for ankr_getBlockchainStats
type GetBlockchainStatsResponse struct {
	// Stats is the list of blockchain statistics
	Stats []BlockchainStat `json:"stats"`
}

// GetBlocksRequest represents the request parameters for ankr_getBlocks
type GetBlocksRequest struct {
	// Blockchain is the supported chain for the blocks
	Blockchain Chain `json:"blockchain,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" default:"false"`

	// DecodeTxData sets to true to decode transaction data, or to false if not interested in it
	DecodeTxData bool `json:"decodeTxData,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" default:"true"`

	// FromBlock is the first block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty"`

	// ToBlock is the last block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	// Note that logs are stored inside transactions, so make sure includeTxs is also set to true
	IncludeLogs bool `json:"includeLogs,omitempty" default:"false"`

	// IncludeTxs sets to true to include transactions, or to false to exclude them
	IncludeTxs *bool `json:"includeTxs,omitempty" default:"true"`
}

// EthBlockDetails represents Ethereum block details
type EthBlockDetails struct {
	// Difficulty is the block difficulty
	Difficulty string `json:"difficulty"`

	// ExtraData is the extra data field
	ExtraData string `json:"extraData"`

	// GasLimit is the gas limit for the block
	GasLimit int64 `json:"gasLimit"`

	// GasUsed is the gas used in the block
	GasUsed int64 `json:"gasUsed"`

	// Miner is the miner address
	Miner string `json:"miner"`

	// Nonce is the block nonce
	Nonce string `json:"nonce"`

	// Sha3Uncles is the SHA3 of uncles
	Sha3Uncles string `json:"sha3Uncles"`

	// Size is the block size
	Size string `json:"size"`

	// StateRoot is the state root hash
	StateRoot string `json:"stateRoot"`

	// TotalDifficulty is the total difficulty
	TotalDifficulty string `json:"totalDifficulty"`
}

// BlockDetails represents block details
type BlockDetails struct {
	// EthBlock contains Ethereum-specific block details
	EthBlock EthBlockDetails `json:"ethBlock,omitempty"`
}

// Block represents a blockchain block
type Block struct {
	// BlockHash is the hash of the block
	BlockHash string `json:"hash"`

	// BlockHeight is the height of the block
	BlockHeight string `json:"number"`

	// BlockchainLogo is the logo URL of the blockchain
	BlockchainLogo string `json:"blockchainLogo,omitempty"`

	// BlockchainName is the name of the blockchain
	BlockchainName string `json:"blockchain"`

	// Details contains blockchain-specific block details
	Details BlockDetails `json:"details"`

	// LogsBloom is the bloom filter for the logs of the block
	LogsBloom string `json:"logsBloom"`

	// MixHash is the mixed hash
	MixHash string `json:"mixHash"`

	// Nonce is the block nonce
	Nonce string `json:"nonce"`

	// ParentHash is the hash of the parent block
	ParentHash string `json:"parentHash"`

	// ReceiptsRoot is the root hash of the receipts trie
	ReceiptsRoot string `json:"receiptsRoot"`

	// Sha3Uncles is the SHA3 of uncles
	Sha3Uncles string `json:"sha3Uncles"`

	// StateRoot is the state root hash
	StateRoot string `json:"stateRoot"`

	// Miner is the miner address
	Miner string `json:"miner"`

	// Difficulty is the block difficulty
	Difficulty string `json:"difficulty"`

	// ExtraData is the extra data field
	ExtraData string `json:"extraData"`

	// Size is the block size
	Size string `json:"size"`

	// GasLimit is the gas limit for the block
	GasLimit string `json:"gasLimit"`

	// GasUsed is the gas used in the block
	GasUsed string `json:"gasUsed"`

	// Timestamp is the timestamp of the block
	Timestamp string `json:"timestamp"`

	// TransactionsRoot is the root hash of the transaction trie
	TransactionsRoot string `json:"transactionsRoot"`

	// TotalDifficulty is the total difficulty
	TotalDifficulty string `json:"totalDifficulty"`

	// TransactionsCount is the number of transactions in the block
	TransactionsCount int32 `json:"transactionsCount"`

	// Transactions contains the transactions in the block
	Transactions []Transaction `json:"transactions"`

	// Uncles contains the uncles block hashes
	Uncles []any `json:"uncles"`
}

// GetBlocksResponse represents the response for ankr_getBlocks
type GetBlocksResponse struct {
	// Blocks is the list of blocks
	Blocks []Block `json:"blocks"`
}

// GetLogsRequest represents the request parameters for ankr_getLogs
type GetLogsRequest struct {
	// Address is a contract address or list of addresses from which the logs originate
	// Supported value formats: hex or array of hexes
	Address []string `json:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" default:"true"`

	// FromBlock is the first block of the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty"`

	// FromTimestamp is the first timestamp of the range
	FromTimestamp int64 `json:"fromTimestamp,omitempty"`

	// PageSize is the number of result pages you'd like to get
	PageSize int32 `json:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`

	// ToBlock is the last block included in the range
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty"`

	// ToTimestamp is the last timestamp of the range
	ToTimestamp int64 `json:"toTimestamp,omitempty"`

	// Topics is the data the log contains
	Topics [][]string `json:"topics,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetLogsRequest) setPageToken(token string) {
	r.PageToken = token
}

// EventInput represents an event input parameter
type EventInput struct {
	// Indexed indicates if the input is indexed
	Indexed bool `json:"indexed"`

	// Name is the name of the input
	Name string `json:"name"`

	// Size is the size of the input
	Size int32 `json:"size"`

	// Type is the type of the input
	Type string `json:"type"`

	// ValueDecoded is the decoded value of the input
	ValueDecoded string `json:"valueDecoded"`
}

// Event represents a decoded event
type Event struct {
	// Anonymous indicates if the event is anonymous
	Anonymous bool `json:"anonymous"`

	// ID is the event ID
	ID string `json:"id"`

	// Inputs is the list of event inputs
	Inputs []EventInput `json:"inputs"`

	// Name is the name of the event
	Name string `json:"name"`

	// Signature is the event signature
	Signature string `json:"signature"`

	// String is the string representation
	String string `json:"string"`

	// Verified indicates if the event is verified
	Verified bool `json:"verified"`
}

// Log represents a blockchain log entry
type Log struct {
	// Address is the contract address that emitted the log
	Address string `json:"address"`

	// BlockHash is the hash of the block containing the log
	BlockHash string `json:"blockHash"`

	// BlockNumber is the number of the block containing the log
	BlockNumber string `json:"blockNumber"`

	// Data is the log data
	Data string `json:"data"`

	// Event contains decoded event information
	Event Event `json:"event,omitempty"`

	// LogIndex is the index of the log in the block
	LogIndex string `json:"logIndex"`

	// Removed indicates if the log was removed
	Removed bool `json:"removed"`

	// Topics is the list of log topics
	Topics []string `json:"topics"`

	// TransactionHash is the hash of the transaction containing the log
	TransactionHash string `json:"transactionHash"`

	// TransactionIndex is the index of the transaction in the block
	TransactionIndex string `json:"transactionIndex"`
}

// GetLogsResponse represents the response for ankr_getLogs
type GetLogsResponse struct {
	// Logs is the list of log entries
	Logs []Log `json:"logs"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetLogsResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetTransactionsByHashRequest represents the request parameters for ankr_getTransactionsByHash
type GetTransactionsByHashRequest struct {
	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// TransactionHash is the hash of the transaction you'd like to request details for
	TransactionHash string `json:"transactionHash,omitempty"`

	// DecodeLogs sets to true to decode logs, or to false if you don't need this kind of info
	DecodeLogs bool `json:"decodeLogs,omitempty" default:"false"`

	// DecodeTxData sets to true to decode transaction data, or to false if not interested in it
	DecodeTxData bool `json:"decodeTxData,omitempty" default:"false"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	IncludeLogs bool `json:"includeLogs,omitempty" default:"false"`
}

// MethodInput represents a method input parameter
type MethodInput struct {
	// Name is the name of the input
	Name string `json:"name"`

	// Size is the size of the input
	Size int32 `json:"size"`

	// Type is the type of the input
	Type string `json:"type"`

	// ValueDecoded is the decoded value of the input
	ValueDecoded string `json:"valueDecoded"`
}

// Method represents a decoded method
type Method struct {
	// ID is the method ID
	ID string `json:"id"`

	// Inputs is the list of method inputs
	Inputs []MethodInput `json:"inputs"`

	// Name is the name of the method
	Name string `json:"name"`

	// Signature is the method signature
	Signature string `json:"signature"`

	// String is the string representation
	String string `json:"string"`

	// Verified indicates if the method is verified
	Verified bool `json:"verified"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	// BlockHash is the hash of the block containing the transaction
	BlockHash string `json:"blockHash"`

	// BlockNumber is the number of the block containing the transaction
	BlockNumber string `json:"blockNumber"`

	// Blockchain is the blockchain where the transaction occurred
	Blockchain string `json:"blockchain"`

	// ContractAddress is the contract address (for contract creation transactions)
	ContractAddress string `json:"contractAddress,omitempty"`

	// CumulativeGasUsed is the cumulative gas used
	CumulativeGasUsed string `json:"cumulativeGasUsed"`

	// From is the sender address
	From string `json:"from"`

	// Gas is the gas limit
	Gas string `json:"gas"`

	// GasPrice is the gas price
	GasPrice string `json:"gasPrice"`

	// GasUsed is the gas used
	GasUsed string `json:"gasUsed"`

	// Hash is the transaction hash
	Hash string `json:"hash"`

	// Input is the transaction input data
	Input string `json:"input"`

	// Logs is the list of logs emitted by the transaction
	Logs []Log `json:"logs,omitempty"`

	// LogsBloom is the logs bloom filter
	LogsBloom string `json:"logsBloom"`

	// Method contains decoded method information
	Method Method `json:"method,omitempty"`

	// Nonce is the transaction nonce
	Nonce string `json:"nonce"`

	// R is the R signature component
	R string `json:"r"`

	// S is the S signature component
	S string `json:"s"`

	// Status is the transaction status
	Status string `json:"status"`

	// Timestamp is the timestamp of the transaction
	Timestamp string `json:"timestamp"`

	// To is the recipient address
	To string `json:"to"`

	// TransactionHash is the transaction hash
	TransactionHash string `json:"transactionHash"`

	// TransactionIndex is the index of the transaction in the block
	TransactionIndex string `json:"transactionIndex"`

	// Type is the transaction type
	Type string `json:"type"`

	// V is the V signature component
	V string `json:"v"`

	// Value is the transaction value
	Value string `json:"value"`
}

// GetTransactionsByHashResponse represents the response for ankr_getTransactionsByHash
type GetTransactionsByHashResponse struct {
	// Transactions is the list of transactions
	Transactions []Transaction `json:"transactions"`
}

// GetTransactionsByAddressRequest represents the request parameters for ankr_getTransactionsByAddress
type GetTransactionsByAddressRequest struct {
	// Address is the address to search for transactions
	Address string `json:"address,omitempty"`

	// Blockchain is a chain or combination of chains to query
	// Single chain: use Chain constants
	// Multiple chains: use []Chain
	// All chains: leave empty
	Blockchain Chain `json:"blockchain,omitempty"`

	// FromBlock narrows your search indicating the block number to start from (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	FromBlock any `json:"fromBlock,omitempty"`

	// ToBlock narrows your search indicating the block number to end with (inclusive; >= 0)
	// Supported value formats: hex, decimal, "earliest", "latest"
	ToBlock any `json:"toBlock,omitempty"`

	// FromTimestamp narrows your search indicating the timestamp to start from (inclusive; >= 0)
	FromTimestamp int64 `json:"fromTimestamp,omitempty"`

	// ToTimestamp narrows your search indicating the timestamp to end with (inclusive; >= 0)
	ToTimestamp int64 `json:"toTimestamp,omitempty"`

	// IncludeLogs sets to true to include logs, or to false to exclude them
	IncludeLogs bool `json:"includeLogs,omitempty" default:"false"`

	// DescOrder chooses data order, either descending (if true) or ascending (if false)
	DescOrder *bool `json:"descOrder,omitempty" default:"true"`

	// PageSize is the number of result pages you'd like to get
	PageSize int32 `json:"pageSize,omitempty"`

	// PageToken is the current page token provided at the end of the response body
	// Can be referenced in the request to fetch the next page
	PageToken string `json:"pageToken,omitempty"`
}

// setPageToken sets the page token for pagination
func (r *GetTransactionsByAddressRequest) setPageToken(token string) {
	r.PageToken = token
}

// GetTransactionsByAddressResponse represents the response for ankr_getTransactionsByAddress
type GetTransactionsByAddressResponse struct {
	// Transactions is the list of transactions
	Transactions []Transaction `json:"transactions"`

	// NextPageToken is provided at the end of the response body for pagination
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// getNextPageToken returns the next page token for pagination
func (r *GetTransactionsByAddressResponse) getNextPageToken() string {
	return r.NextPageToken
}

// GetInteractionsRequest represents the request parameters for ankr_getInteractions
type GetInteractionsRequest struct {
	// Address is the address of the wallet or contract that created the logs
	Address string `json:"address,omitempty"`
}

// GetInteractionsResponse represents the response for ankr_getInteractions
type GetInteractionsResponse struct {
	// Blockchains is the list of blockchains interacted with the address
	Blockchains []string `json:"blockchains"`
}
