package ankr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

// HTTPClient represents the HTTP client for Ankr Advanced API
type HTTPClient struct {
	uri         string
	httpClient  *http.Client
	rateLimiter *SimpleLimiter
}

type HTTPClientConfig struct {
	APIKey string
	// OnLimitExceeded RateLimitBehavior `default:"block"`
}

// NewHTTPClient creates a new HTTP client with the given configuration
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
	if config == nil {
		config = &HTTPClientConfig{}
	}

	// Create rate limiter
	rateLimiter := NewSimpleLimiter(time.Minute, 1000)

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 90 * time.Second,
		},
	}

	return &HTTPClient{
		uri:         "https://rpc.ankr.com/multichain/" + config.APIKey,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
	}
}

// post makes a JSON-RPC post request and returns the result with generic type
func post[Req any, Resp any](ctx context.Context, client *HTTPClient, method string, params Req) (result Resp, isRPCError bool, err error) {
	// Rate limiting
	client.rateLimiter.Wait(ctx)

	newParams, err := ApplyDefaults(params)
	if err != nil {
		return result, false, fmt.Errorf("failed to apply defaults: %w", err)
	}

	// Create JSON-RPC request
	request := RPCReqBody{
		ID:      1,
		JSONRPC: JSONRPC,
		Method:  method,
		Params:  newParams,
	}

	// Marshal request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return result, false, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", client.uri, bytes.NewReader(requestBody))
	if err != nil {
		return result, false, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return result, false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, false, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return result, false, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var apiResponse RPCRespBody[Resp]
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return result, false, fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResponse.Error != nil {
		return result, true, fmt.Errorf("ankr: rpc error: %+v", apiResponse.Error)
	}

	return apiResponse.Result, false, nil
}

func postWithRetries[Req any, Resp any](ctx context.Context, client *HTTPClient, method string, params Req, retries int) (result Resp, err error) {
	for range retries {
		result, _, err = post[Req, Resp](ctx, client, method, params)
		if err == nil {
			return
		}
		slog.Error("ankr: failed to post, retrying...", "error", err)
		time.Sleep(time.Second)
		continue

	}
	err = fmt.Errorf("ankr: failed to post after %d retries, last error: %w", retries, err)
	return
}

type nextPageFunc[Page any] func(ctx context.Context) (Page, bool, error)

type reqData interface {
	setPageToken(string)
}

type respData interface {
	getNextPageToken() string
}

func makeNextPageFunc[Req reqData, Resp respData](client *HTTPClient, method string, req Req) nextPageFunc[Resp] {
	return func(ctx context.Context) (resp Resp, hasNext bool, err error) {
		resp, err = postWithRetries[Req, Resp](ctx, client, method, req, 3)
		if err != nil {
			return resp, false, err
		}
		hasNext = resp.getNextPageToken() != ""
		if hasNext {
			req.setPageToken(resp.getNextPageToken())
		}
		return resp, hasNext, nil
	}
}

type Pages[Page any] struct {
	hasNext bool
	mu      sync.RWMutex
	next    nextPageFunc[Page]
}

func newPages[Page any](next nextPageFunc[Page]) *Pages[Page] {
	return &Pages[Page]{
		hasNext: true,
		next:    next,
	}
}

func (p *Pages[Page]) HasNext() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.hasNext
}

func (p *Pages[Page]) Next(ctx context.Context) (newPage Page, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.hasNext {
		// should never happen
		err = fmt.Errorf("ankr: Pages has no next func")
		return
	}
	newPage, ok, err := p.next(ctx)
	if err != nil {
		return
	}
	p.hasNext = ok
	return newPage, nil
}

// ============================================================================
// NFT API Methods
// ============================================================================

// GetNFTsByOwner retrieves account-associated NFTs with automatic pagination
//
// This method fetches NFTs owned by a specific wallet address across multiple chains.
// It supports pagination and filtering by contract addresses.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including wallet address, blockchain, page size, etc.
//
// Returns:
//   - *Pages[GetNFTsByOwnerResp]: Paginated response iterator
func (c *HTTPClient) GetNFTsByOwner(req GetNFTsByOwnerReq) *Pages[*GetNFTsByOwnerResp] {
	return newPages(makeNextPageFunc[*GetNFTsByOwnerReq, *GetNFTsByOwnerResp](c, "ankr_getNFTsByOwner", &req))
}

// GetNFTMetadata retrieves metadata of a particular NFT
//
// This method fetches detailed metadata for a specific NFT including attributes,
// traits, image URL, and other metadata information.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including blockchain, contract address, and token ID
//
// Returns:
//   - *GetNFTMetadataResp: Response containing NFT metadata
//   - error: Error if the request fails
func (c *HTTPClient) GetNFTMetadata(ctx context.Context, req GetNFTMetadataReq) (*GetNFTMetadataResp, error) {
	return postWithRetries[GetNFTMetadataReq, *GetNFTMetadataResp](ctx, c, "ankr_getNFTMetadata", req, 3)
}

// GetNFTHolders retrieves holders of a particular NFT with automatic pagination
//
// This method fetches all wallet addresses that hold a specific NFT collection.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including blockchain, contract address, and pagination
//
// Returns:
//   - *Pages[GetNFTHoldersResp]: Paginated response iterator
func (c *HTTPClient) GetNFTHolders(req GetNFTHoldersReq) *Pages[*GetNFTHoldersResp] {
	return newPages(makeNextPageFunc[*GetNFTHoldersReq, *GetNFTHoldersResp](c, "ankr_getNFTHolders", &req))
}

// GetNFTTransfers retrieves NFT transfers info with automatic pagination
//
// This method fetches NFT transfer transactions within a specified range,
// supporting filtering by addresses, blockchains, and time/block ranges.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including addresses, blockchain(s), and range filters
//
// Returns:
//   - *Pages[GetNFTTransfersResp]: Paginated response iterator
func (c *HTTPClient) GetNFTTransfers(req GetNFTTransfersReq) *Pages[*GetNFTTransfersResp] {
	return newPages(makeNextPageFunc[*GetNFTTransfersReq, *GetNFTTransfersResp](c, "ankr_getNftTransfers", &req))
}

// ============================================================================
// Query API Methods
// ============================================================================

// GetBlockchainStats retrieves blockchain statistics
//
// This method fetches statistics for one or more blockchains including
// transaction counts, event counts, latest block numbers, and native coin prices.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including blockchain(s) to query
//
// Returns:
//   - *GetBlockchainStatsResp: Response containing blockchain statistics
//   - error: Error if the request fails
func (c *HTTPClient) GetBlockchainStats(ctx context.Context, req GetBlockchainStatsReq) (*GetBlockchainStatsResp, error) {
	return postWithRetries[GetBlockchainStatsReq, *GetBlockchainStatsResp](ctx, c, "ankr_getBlockchainStats", req, 3)
}

// GetBlocks retrieves full info of blocks in a range
//
// This method fetches detailed information about blocks within a specified range.
// The maximum range is 100 blocks. Supports decoding of logs and transaction data.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including blockchain, block range, and decode options
//
// Returns:
//   - *GetBlocksResp: Response containing block information
//   - error: Error if the request fails
func (c *HTTPClient) GetBlocks(ctx context.Context, req GetBlocksReq) (*GetBlocksResp, error) {
	return postWithRetries[GetBlocksReq, *GetBlocksResp](ctx, c, "ankr_getBlocks", req, 3)
}

// GetLogs retrieves historical data for the specified range of blocks with automatic pagination
//
// This method fetches event logs within a specified block or timestamp range.
// Supports filtering by contract addresses and topics, with optional log decoding.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including blockchain, address filters, block/timestamp range
//
// Returns:
//   - *Pages[GetLogsResp]: Paginated response iterator
func (c *HTTPClient) GetLogs(req GetLogsReq) *Pages[*GetLogsResp] {
	return newPages(makeNextPageFunc[*GetLogsReq, *GetLogsResp](c, "ankr_getLogs", &req))
}

// GetTxsByHash retrieves the details of transactions by hash
//
// This method fetches detailed transaction information including decoded logs and
// transaction data for one or more transaction hashes across multiple chains.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including transaction hash and blockchain(s)
//
// Returns:
//   - *GetTxsByHashResp: Response containing transaction details
//   - error: Error if the request fails
func (c *HTTPClient) GetTxsByHash(ctx context.Context, req GetTxsByHashReq) (*GetTxsByHashResp, error) {
	return postWithRetries[GetTxsByHashReq, *GetTxsByHashResp](ctx, c, "ankr_getTransactionsByHash", req, 3)
}

// GetTxsByAddress retrieves transactions for a specific address with automatic pagination
//
// This method fetches all transactions involving a specific address within
// a specified block or timestamp range across multiple chains.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including address, blockchain(s), and range filters
//
// Returns:
//   - *Pages[GetTxsByAddressResp]: Paginated response iterator
func (c *HTTPClient) GetTxsByAddress(req GetTxsByAddressReq) *Pages[*GetTxsByAddressResp] {
	return newPages(makeNextPageFunc[*GetTxsByAddressReq, *GetTxsByAddressResp](c, "ankr_getTransactionsByAddress", &req))
}

// GetInteractions retrieves blockchains interacted with a particular wallet
//
// This method fetches a list of all blockchains that a specific wallet address
// has interacted with, providing a comprehensive view of the wallet's activity.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including the wallet address
//
// Returns:
//   - *GetInteractionsResp: Response containing list of blockchains
//   - error: Error if the request fails
func (c *HTTPClient) GetInteractions(ctx context.Context, req GetInteractionsReq) (*GetInteractionsResp, error) {
	return postWithRetries[GetInteractionsReq, *GetInteractionsResp](ctx, c, "ankr_getInteractions", req, 3)
}

// ============================================================================
// Token API Methods
// ============================================================================

// GetAccountBalances retrieves all account balances with automatic pagination
//
// This method fetches token balances for a specific wallet address across multiple chains.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including wallet address and blockchain(s)
//
// Returns:
//   - *Pages[GetAccountBalanceResp]: Paginated response iterator
func (c *HTTPClient) GetAccountBalances(req GetAccountBalanceReq) *Pages[*GetAccountBalanceResp] {
	return newPages(makeNextPageFunc[*GetAccountBalanceReq, *GetAccountBalanceResp](c, "ankr_getAccountBalance", &req))
}

// GetCurrencies retrieves info on currencies available for a particular blockchain
//
// This method fetches information about all supported currencies (tokens) on a specific blockchain,
// including their contract addresses, decimals, names, symbols, and thumbnails.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including the blockchain to query
//
// Returns:
//   - *GetCurrenciesResp: Response containing list of currencies
//   - error: Error if the request fails
func (c *HTTPClient) GetCurrencies(ctx context.Context, req GetCurrenciesReq) (*GetCurrenciesResp, error) {
	return postWithRetries[GetCurrenciesReq, *GetCurrenciesResp](ctx, c, "ankr_getCurrencies", req, 3)
}

// GetTokenPrice retrieves the price of a particular token
//
// This method fetches the current USD price of a specific token on a blockchain.
// If no contract address is provided, returns the native coin price of the blockchain.
//
// Args:
//   - ctx: Context for cancellation
//   - req: Request parameters including blockchain and optional contract address
//
// Returns:
//   - *GetTokenPriceResp: Response containing token price information
//   - error: Error if the request fails
func (c *HTTPClient) GetTokenPrice(ctx context.Context, req GetTokenPriceReq) (*GetTokenPriceResp, error) {
	return postWithRetries[GetTokenPriceReq, *GetTokenPriceResp](ctx, c, "ankr_getTokenPrice", req, 3)
}

// GetTokenHolders retrieves all token holders with automatic pagination
//
// This method fetches all wallet addresses that hold a specific token.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including contract address and blockchain
//
// Returns:
//   - *Pages[GetTokenHoldersResp]: Paginated response iterator
func (c *HTTPClient) GetTokenHolders(req GetTokenHoldersReq) *Pages[*GetTokenHoldersResp] {
	return newPages(makeNextPageFunc[*GetTokenHoldersReq, *GetTokenHoldersResp](c, "ankr_getTokenHolders", &req))
}

// GetTokenHolderCountHistories retrieves all token holder count data with automatic pagination
//
// This method fetches historical holder count data for a specific token over time.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including contract address and blockchain
//
// Returns:
//   - *Pages[GetTokenHoldersCountResp]: Paginated response iterator
func (c *HTTPClient) GetTokenHolderCountHistories(req GetTokenHoldersCountReq) *Pages[*GetTokenHoldersCountResp] {
	return newPages(makeNextPageFunc[*GetTokenHoldersCountReq, *GetTokenHoldersCountResp](c, "ankr_getTokenHoldersCount", &req))
}

// GetTokenTransfers retrieves all token transfers with automatic pagination
//
// This method fetches token transfer transactions within a specified range,
// supporting filtering by addresses, blockchains, and time/block ranges.
// Since the request contains PageToken, this method returns paginated results.
//
// Args:
//   - req: Request parameters including addresses, blockchain(s), and range filters
//
// Returns:
//   - *Pages[GetTokenTransfersResp]: Paginated response iterator
func (c *HTTPClient) GetTokenTransfers(req GetTokenTransfersReq) *Pages[*GetTokenTransfersResp] {
	return newPages(makeNextPageFunc[*GetTokenTransfersReq, *GetTokenTransfersResp](c, "ankr_getTokenTransfers", &req))
}
