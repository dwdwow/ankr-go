package ankr

import (
	"context"
	"os"
	"testing"
	"time"
)

// Helper function to create HTTP client for testing
func createTestClient(t *testing.T) *HTTPClient {
	apiKey := os.Getenv("ANKR_API_KEY")
	if apiKey == "" {
		t.Skip("ANKR_API_KEY environment variable not set, skipping integration tests")
	}

	return NewHTTPClient(&HTTPClientConfig{
		APIKey:          apiKey,
		OnLimitExceeded: RateLimitBlock,
	})
}

// TestGetBlockchainStats tests the GetBlockchainStats API method
func TestGetBlockchainStats(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetBlockchainStatsRequest{
		Blockchain: ChainEthereum,
	}

	resp, err := client.GetBlockchainStats(ctx, req)
	if err != nil {
		t.Fatalf("GetBlockchainStats failed: %v", err)
	}

	if len(resp.Stats) == 0 {
		t.Error("Expected at least one blockchain stat")
	}

	stat := resp.Stats[0]
	if stat.Blockchain == "" {
		t.Error("Expected non-empty blockchain name")
	}
	if stat.LatestBlockNumber == 0 {
		t.Error("Expected non-zero latest block number")
	}

	t.Logf("Blockchain stats: %+v", stat)
}

// TestGetAccountBalances tests the GetAccountBalances API method
func TestGetAccountBalances(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetAccountBalanceRequest{
		WalletAddress: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // Vitalik's address
		Blockchain:    ChainEthereum,
		PageSize:      10,
	}

	pages := client.GetAccountBalances(req)

	pageCount := 0
	totalAssets := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next page: %v", err)
		}

		totalAssets += len(page.Assets)
		t.Logf("Page %d: Found %d assets, Total USD: %s", pageCount+1, len(page.Assets), page.TotalBalanceUsd)

		if len(page.Assets) > 0 {
			asset := page.Assets[0]
			if asset.TokenName == "" {
				t.Error("Expected non-empty token name")
			}
			if asset.Balance == "" {
				t.Error("Expected non-empty balance")
			}
		}

		pageCount++
	}

	t.Logf("Total assets found: %d", totalAssets)
}

// TestGetTokenPrice tests the GetTokenPrice API method
func TestGetTokenPrice(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test native token price (ETH)
	req := GetTokenPriceRequest{
		Blockchain: ChainEthereum,
	}

	resp, err := client.GetTokenPrice(ctx, req)
	if err != nil {
		t.Fatalf("GetTokenPrice failed: %v", err)
	}

	if resp.UsdPrice == "" {
		t.Error("Expected non-empty USD price")
	}
	if resp.Blockchain == "" {
		t.Error("Expected non-empty blockchain")
	}

	t.Logf("ETH price: $%s", resp.UsdPrice)

	// Test ERC-20 token price (WETH)
	req.ContractAddress = "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2" // WETH on Ethereum
	resp2, err := client.GetTokenPrice(ctx, req)
	if err != nil {
		t.Logf("GetTokenPrice for WETH failed (this might be expected): %v", err)
	} else {
		t.Logf("WETH price: $%s", resp2.UsdPrice)
	}
}

// TestGetCurrencies tests the GetCurrencies API method
func TestGetCurrencies(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetCurrenciesRequest{
		Blockchain: ChainEthereum,
	}

	resp, err := client.GetCurrencies(ctx, req)
	if err != nil {
		t.Fatalf("GetCurrencies failed: %v", err)
	}

	if len(resp.Currencies) == 0 {
		t.Error("Expected at least one currency")
	}

	currency := resp.Currencies[0]
	if currency.Symbol == "" {
		t.Error("Expected non-empty currency symbol")
	}
	if currency.Name == "" {
		t.Error("Expected non-empty currency name")
	}
	if currency.Decimals == 0 {
		t.Error("Expected non-zero decimals")
	}

	t.Logf("Found %d currencies on Ethereum", len(resp.Currencies))
	t.Logf("Sample currency: %s (%s) - %d decimals", currency.Name, currency.Symbol, currency.Decimals)
}

// TestGetNFTsByOwner tests the GetNFTsByOwner API method
func TestGetNFTsByOwner(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetNFTsByOwnerRequest{
		WalletAddress: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // Vitalik's address
		Blockchain:    ChainEthereum,
		PageSize:      5,
	}

	pages := client.GetNFTsByOwner(req)

	pageCount := 0
	totalNFTs := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next NFT page: %v", err)
		}

		totalNFTs += len(page.Assets)
		t.Logf("NFT Page %d: Found %d NFTs", pageCount+1, len(page.Assets))

		if len(page.Assets) > 0 {
			nft := page.Assets[0]
			if nft.ContractAddress == "" {
				t.Error("Expected non-empty contract address")
			}
			if nft.TokenID == "" {
				t.Error("Expected non-empty token ID")
			}
			t.Logf("Sample NFT: %s (ID: %s, Contract: %s)", nft.Name, nft.TokenID, nft.ContractAddress)
		}

		pageCount++
	}

	t.Logf("Total NFTs found: %d", totalNFTs)
}

// TestGetBlocks tests the GetBlocks API method
func TestGetBlocks(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetBlocksRequest{
		Blockchain: ChainEthereum,
		FromBlock:  23583308,
		ToBlock:    23583308,
		// IncludeTxs: PtrFalse(),
	}

	resp, err := client.GetBlocks(ctx, req)
	if err != nil {
		t.Fatalf("GetBlocks failed: %v", err)
	}

	if len(resp.Blocks) == 0 {
		t.Error("Expected at least one block")
	}

	block := resp.Blocks[0]
	if block.BlockHash == "" {
		t.Error("Expected non-empty block hash")
	}
	if block.BlockHeight == "" {
		t.Error("Expected non-empty block height")
	}
	if block.BlockchainName == "" {
		t.Error("Expected non-empty blockchain name")
	}

	t.Logf("Found %d blocks", len(resp.Blocks))
	t.Logf("Sample block: Height %s, Hash %s", block.BlockHeight, block.BlockHash)
	t.Logf("Transactions: %d", len(block.Transactions))
}

// TestGetLogs tests the GetLogs API method
func TestGetLogs(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetLogsRequest{
		Blockchain: ChainEthereum,
		FromBlock:  "latest",
		ToBlock:    "latest",
		PageSize:   5,
	}

	pages := client.GetLogs(req)

	pageCount := 0
	totalLogs := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next logs page: %v", err)
		}

		totalLogs += len(page.Logs)
		t.Logf("Logs Page %d: Found %d logs", pageCount+1, len(page.Logs))

		if len(page.Logs) > 0 {
			log := page.Logs[0]
			if log.Address == "" {
				t.Error("Expected non-empty log address")
			}
			if log.BlockNumber == "" {
				t.Error("Expected non-empty block number")
			}
			t.Logf("Sample log: Address %s, Block %s", log.Address, log.BlockNumber)
		}

		pageCount++
	}

	t.Logf("Total logs found: %d", totalLogs)
}

// TestGetTransactionsByHash tests the GetTransactionsByHash API method
func TestGetTransactionsByHash(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a placeholder hash - this will likely fail but tests the API call
	req := GetTransactionsByHashRequest{
		TransactionHash: "0x728db5670765cefed19446cb4440e0a162c218cf0f345d7cf3a6ab2a5ad46a29",
		Blockchain:      ChainEthereum,
		DecodeLogs:      false,
		DecodeTxData:    false,
		IncludeLogs:     false,
	}

	txs, err := client.GetTransactionsByHash(ctx, req)
	if err != nil {
		t.Logf("GetTransactionsByHash failed as expected (transaction not found): %v", err)
	} else {
		t.Log("GetTransactionsByHash succeeded")
	}
	if len(txs.Transactions) == 0 {
		t.Error("Expected at least one transaction")
	}
	t.Logf("Transactions: %d", len(txs.Transactions))
	t.Logf("Sample transaction: %s", txs.Transactions[0].Hash)
}

// TestGetTransactionsByAddress tests the GetTransactionsByAddress API method
func TestGetTransactionsByAddress(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetTransactionsByAddressRequest{
		Address:     "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // Vitalik's address
		Blockchain:  ChainEthereum,
		IncludeLogs: false,
		PageSize:    5,
	}

	pages := client.GetTransactionsByAddress(req)

	pageCount := 0
	totalTxs := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next transactions page: %v", err)
		}

		totalTxs += len(page.Transactions)
		t.Logf("Transactions Page %d: Found %d transactions", pageCount+1, len(page.Transactions))

		if len(page.Transactions) > 0 {
			tx := page.Transactions[0]
			if tx.Hash == "" {
				t.Error("Expected non-empty transaction hash")
			}
			if tx.From == "" {
				t.Error("Expected non-empty from address")
			}
			t.Logf("Sample transaction: %s", tx.Hash)
		}

		pageCount++
	}

	t.Logf("Total transactions found: %d", totalTxs)
}

// TestGetInteractions tests the GetInteractions API method
func TestGetInteractions(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetInteractionsRequest{
		Address: "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045", // Vitalik's address
	}

	resp, err := client.GetInteractions(ctx, req)
	if err != nil {
		t.Fatalf("GetInteractions failed: %v", err)
	}

	if len(resp.Blockchains) == 0 {
		t.Error("Expected at least one blockchain interaction")
	}

	t.Logf("Interactions found on %d blockchains: %v", len(resp.Blockchains), resp.Blockchains)
}

// TestGetTokenHolders tests the GetTokenHolders API method
func TestGetTokenHolders(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use USDC contract address
	req := GetTokenHoldersRequest{
		Blockchain:      ChainEthereum,
		ContractAddress: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
		PageSize:        5,
	}

	pages := client.GetTokenHolders(req)

	pageCount := 0
	totalHolders := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next holders page: %v", err)
		}

		totalHolders += len(page.Holders)
		t.Logf("Holders Page %d: Found %d holders", pageCount+1, len(page.Holders))

		if len(page.Holders) > 0 {
			holder := page.Holders[0]
			if holder.HolderAddress == "" {
				t.Error("Expected non-empty holder address")
			}
			if holder.Balance == "" {
				t.Error("Expected non-empty balance")
			}
			t.Logf("Sample holder: %s, Balance: %s", holder.HolderAddress, holder.Balance)
		}

		pageCount++
	}

	t.Logf("Total holders found: %d", totalHolders)
}

// TestGetTokenHolderCountHistories tests the GetTokenHolderCountHistories API method
func TestGetTokenHolderCountHistories(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetTokenHoldersCountRequest{
		Blockchain:      ChainEthereum,
		ContractAddress: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
		PageSize:        5,
	}

	pages := client.GetTokenHolderCountHistories(req)

	pageCount := 0
	totalHistories := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next holder count history page: %v", err)
		}

		totalHistories += len(page.HolderCountHistory)
		t.Logf("History Page %d: Found %d history entries", pageCount+1, len(page.HolderCountHistory))

		if len(page.HolderCountHistory) > 0 {
			history := page.HolderCountHistory[0]
			if history.HolderCount == 0 {
				t.Error("Expected non-zero holder count")
			}
			if history.LastUpdatedAt == "" {
				t.Error("Expected non-empty last updated timestamp")
			}
			t.Logf("Sample history: %d holders at %s", history.HolderCount, history.LastUpdatedAt)
		}

		pageCount++
	}

	t.Logf("Total history entries found: %d", totalHistories)
}

// TestGetTokenTransfers tests the GetTokenTransfers API method
func TestGetTokenTransfers(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetTokenTransfersRequest{
		Address:       []string{"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"}, // Vitalik's address
		Blockchain:    ChainEthereum,
		FromTimestamp: 1640995200, // 2022-01-01
		ToTimestamp:   1672531200, // 2023-01-01
		PageSize:      5,
	}

	pages := client.GetTokenTransfers(req)

	pageCount := 0
	totalTransfers := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next transfers page: %v", err)
		}

		totalTransfers += len(page.Transfers)
		t.Logf("Transfers Page %d: Found %d transfers", pageCount+1, len(page.Transfers))

		if len(page.Transfers) > 0 {
			transfer := page.Transfers[0]
			if transfer.TransactionHash == "" {
				t.Error("Expected non-empty transaction hash")
			}
			if transfer.FromAddress == "" {
				t.Error("Expected non-empty from address")
			}
			if transfer.ToAddress == "" {
				t.Error("Expected non-empty to address")
			}
			t.Logf("Sample transfer: %s -> %s, Amount: %s %s",
				transfer.FromAddress, transfer.ToAddress, transfer.Value, transfer.TokenSymbol)
		}

		pageCount++
	}

	t.Logf("Total transfers found: %d", totalTransfers)
}

// TestGetNFTMetadata tests the GetNFTMetadata API method
func TestGetNFTMetadata(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a known NFT contract and token ID
	req := GetNFTMetadataRequest{
		Blockchain:      ChainEthereum,
		ContractAddress: "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", // BAYC contract
		TokenID:         "1",
		ForceFetch:      false,
		SkipSyncCheck:   false,
	}

	resp, err := client.GetNFTMetadata(ctx, req)
	if err != nil {
		t.Logf("GetNFTMetadata failed (this might be expected if NFT not found): %v", err)
	} else {
		if resp.Metadata.ContractAddress == "" {
			t.Error("Expected non-empty contract address")
		}
		if resp.Metadata.TokenID == "" {
			t.Error("Expected non-empty token ID")
		}
		// if resp.Metadata.Attributes.Name == "" {
		// 	t.Error("Expected non-empty NFT name")
		// }
		// t.Logf("NFT Metadata: %s (ID: %s)", resp.Metadata.Attributes.Name, resp.Metadata.TokenID)
	}
}

// TestGetNFTHolders tests the GetNFTHolders API method
func TestGetNFTHolders(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a known NFT contract
	req := GetNFTHoldersRequest{
		Blockchain:      ChainEthereum,
		ContractAddress: "0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D", // BAYC contract
		PageSize:        5,
	}

	pages := client.GetNFTHolders(req)

	pageCount := 0
	totalHolders := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next NFT holders page: %v", err)
		}

		totalHolders += len(page.Holders)
		t.Logf("NFT Holders Page %d: Found %d holders", pageCount+1, len(page.Holders))

		if len(page.Holders) > 0 {
			holder := page.Holders[0]
			if holder == "" {
				t.Error("Expected non-empty holder address")
			}
			t.Logf("Sample holder: %s", holder)
		}

		pageCount++
	}

	t.Logf("Total NFT holders found: %d", totalHolders)
}

// TestGetNFTTransfers tests the GetNFTTransfers API method
func TestGetNFTTransfers(t *testing.T) {
	client := createTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := GetNFTTransfersRequest{
		Address:       []string{"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"}, // Vitalik's address
		Blockchain:    []Chain{ChainEthereum},
		FromTimestamp: 1640995200, // 2022-01-01
		ToTimestamp:   1672531200, // 2023-01-01
		PageSize:      5,
	}

	pages := client.GetNFTTransfers(req)

	pageCount := 0
	totalTransfers := 0
	for pages.HasNext() && pageCount < 2 { // Limit to 2 pages for testing
		page, err := pages.Next(ctx)
		if err != nil {
			t.Fatalf("Failed to get next NFT transfers page: %v", err)
		}

		totalTransfers += len(page.Transfers)
		t.Logf("NFT Transfers Page %d: Found %d transfers", pageCount+1, len(page.Transfers))

		if len(page.Transfers) > 0 {
			transfer := page.Transfers[0]
			if transfer.TransactionHash == "" {
				t.Error("Expected non-empty transaction hash")
			}
			if transfer.FromAddress == "" {
				t.Error("Expected non-empty from address")
			}
			if transfer.ToAddress == "" {
				t.Error("Expected non-empty to address")
			}
			if transfer.ContractAddress == "" {
				t.Error("Expected non-empty contract address")
			}
			if transfer.TokenID == "" {
				t.Error("Expected non-empty token ID")
			}
			t.Logf("Sample NFT transfer: %s -> %s, TokenID: %s, Contract: %s",
				transfer.FromAddress, transfer.ToAddress, transfer.TokenID, transfer.ContractAddress)
		}

		pageCount++
	}

	t.Logf("Total NFT transfers found: %d", totalTransfers)
}
