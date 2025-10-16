# Ankr Go SDK

A comprehensive Go SDK for interacting with the Ankr Advanced API, providing easy access to blockchain data across multiple networks.

## Requirements

- Go 1.25.0 or higher

## Features

- 🚀 **Complete API Coverage** - Support for NFT, Token, and Query APIs
- 🔄 **Automatic Pagination** - Built-in pagination handling for large datasets
- 🛡️ **Rate Limiting** - Configurable rate limiting with token bucket algorithm
- 🔁 **Retry Logic** - Automatic retry with exponential backoff
- 📊 **Type Safety** - Fully typed request/response structures
- 🎯 **Default Values** - Smart default value application
- 🔧 **Flexible Configuration** - Customizable HTTP client settings

## Installation

Make sure you have Go 1.25.0 or higher installed:

```bash
go version
```

Then install the SDK:

```bash
go get github.com/dwdwow/ankr
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/dwdwow/ankr"
)

func main() {
    // Create HTTP client with your API key
    client := ankr.NewHTTPClient(&ankr.HTTPClientConfig{
        APIKey: "your-api-key-here",
    })
    
    // Get blockchain statistics
    ctx := context.Background()
    req := ankr.GetBlockchainStatsReq{
        Blockchain: ankr.ChainEthereum,
    }
    
    resp, err := client.GetBlockchainStats(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Latest block: %d\n", resp.Stats[0].LatestBlockNumber)
}
```

### Paginated Results

```go
// Get NFTs by owner with automatic pagination
req := ankr.GetNFTsByOwnerReq{
    WalletAddress: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    Blockchain:    ankr.ChainEthereum,
    PageSize:      10,
}

pages := client.GetNFTsByOwner(req)

for pages.HasNext() {
    page, err := pages.Next(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d NFTs\n", len(page.Assets))
}
```

### Rate Limiting

```go
client := ankr.NewHTTPClient(&ankr.HTTPClientConfig{
    APIKey:          "your-api-key-here",
    OnLimitExceeded: ankr.RateLimitBlock, // or RateLimitError
})
```

### Pointer Fields

Some boolean fields use pointer types to support default values. Use the following patterns:

```go
// For boolean fields with default value true
field := &[]bool{true}[0]

// For boolean fields with default value false  
field := &[]bool{false}[0]

// Or use the helper function
field := ankr.TruePtr()
field := ankr.FalsePtr()
```

## API Modules

### NFT API

Retrieve NFT-related data across multiple blockchains:

- `GetNFTsByOwner()` - Get NFTs owned by an address
- `GetNFTMetadata()` - Get metadata for a specific NFT
- `GetNFTHolders()` - Get holders of a specific NFT collection
- `GetNFTTransfers()` - Get NFT transfer history

```go
// Get NFT metadata
req := ankr.GetNFTMetadataReq{
    Blockchain:      ankr.ChainEthereum,
    ContractAddress: "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", // BAYC
    TokenID:         "1",
    ForceFetch:      false, // Regular bool field (default: false)
    SkipSyncCheck:   false, // Regular bool field (default: false)
}

resp, err := client.GetNFTMetadata(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("NFT Name: %s\n", resp.Metadata.Attributes.Name)
```

### Token API

Access token and balance information:

- `GetAccountBalances()` - Get token balances for an address
- `GetTokenPrice()` - Get current token prices
- `GetCurrencies()` - List available currencies
- `GetTokenHolders()` - Get token holders
- `GetTokenTransfers()` - Get token transfer history

```go
// Get account balances
req := ankr.GetAccountBalanceReq{
    WalletAddress:   "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
    Blockchain:      ankr.ChainEthereum,
    OnlyWhitelisted: &[]bool{true}[0], // Pointer type for default value true
    NativeFirst:     &[]bool{true}[0], // Pointer type for default value true
}

pages := client.GetAccountBalances(req)
```

### Query API

Query blockchain data and transaction information:

- `GetBlockchainStats()` - Get blockchain statistics
- `GetBlocks()` - Get block information
- `GetLogs()` - Get event logs
- `GetTxsByHash()` - Get transaction by hash
- `GetTxsByAddress()` - Get transactions by address
- `GetInteractions()` - Get wallet interactions

```go
// Get recent blocks
req := ankr.GetBlocksReq{
    Blockchain: ankr.ChainEthereum,
    FromBlock:  "latest",
    ToBlock:    "latest",
    IncludeTxs: &[]bool{true}[0], // Pointer type for default value true
}

resp, err := client.GetBlocks(ctx, req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d blocks\n", len(resp.Blocks))
```

## Supported Blockchains

### Mainnet

- Ethereum (`eth`)
- BSC (`bsc`)
- Polygon (`polygon`)
- Arbitrum (`arbitrum`)
- Optimism (`optimism`)
- Avalanche (`avalanche`)
- Fantom (`fantom`)
- Base (`base`)
- And many more...

### Testnet

- Ethereum Sepolia (`eth_sepolia`)
- Ethereum Holesky (`eth_holesky`)
- BSC Testnet (`bsc_testnet`)
- Polygon Amoy (`polygon_amoy`)
- Base Sepolia (`base_sepolia`)
- And more...

## Configuration

### HTTP Client Configuration

```go
config := &ankr.HTTPClientConfig{
    APIKey:          "your-api-key",           // Required
    BaseURL:         "https://rpc.ankr.com",   // Optional, defaults to Ankr RPC
    Timeout:         30 * time.Second,         // Optional
    OnLimitExceeded: ankr.RateLimitBlock,      // Optional
}

client := ankr.NewHTTPClient(config)
```

### Rate Limiting Options

- `RateLimitBlock` - Block requests when rate limit exceeded
- `RateLimitError` - Return error when rate limit exceeded

## Error Handling

The SDK provides comprehensive error handling:

```go
resp, err := client.GetTokenPrice(ctx, req)
if err != nil {
    // Handle different error types
    switch e := err.(type) {
    case *ankr.RateLimitError:
        log.Printf("Rate limit exceeded: %v", e)
    case *ankr.APIError:
        log.Printf("API error: %v", e)
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Pagination

Many API endpoints support pagination. The SDK provides a convenient `Pages` interface:

```go
pages := client.GetNFTsByOwner(req)

for pages.HasNext() {
    page, err := pages.Next(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    // Process page data
    for _, nft := range page.Assets {
        fmt.Printf("NFT: %s\n", nft.Name)
    }
}
```

## Testing

Run the test suite:

```bash
# Set your API key
export ANKR_API_KEY="your-api-key-here"

# Run all tests
go test -v ./...

# Run specific test
go test -v -run TestGetNFTsByOwner
```

## Examples

Check out the test files for comprehensive usage examples:

- `http_client_test.go` - Complete API usage examples
- `defaults_test.go` - Default value handling examples

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- 📖 [Ankr Documentation](https://docs.ankr.com/)

## Changelog

### v1.0.0

- Initial release
- Complete NFT, Token, and Query API support
- Automatic pagination handling
- Rate limiting and retry logic
- Comprehensive type definitions
- Full test coverage
