# WLX_Bridge_Bot
A Go-based crypto bridge bot enabling seamless cross-chain swaps integrated with Wallex, simplifying secure token exchanges across multiple blockchains.

# BridgeBot - Cross-chain Bridge Automation in Go

BridgeBot is an automated bot implemented in Go that seamlessly interacts with The Bridgers cross-chain bridge API. It enables efficient token swaps across different blockchain networks, ensuring smooth interoperability and user-friendly transactions.

---

## 🚀 Overview

The BridgeBot leverages the robust Bridgers API to automate token exchanges across blockchains, including functionalities like token quotes, swaps, approvals, and fetching token details.

**Core functionalities include:**

- **Cross-chain Token Swaps**
- **Real-time Token Quotes**
- **Smart Contract Approvals**
- **Token Information Retrieval**
- **Structured & Colorized Logging**
- **Automatic Retry and Error Handling**

---

## 📂 Project Structure

```plaintext
go_bridgebot
├── cmd
│   └── bridge_bot      # Entry point for the BridgeBot application
├── internal
│   └── utils
│       ├── httpClient  # HTTP request utilities with retries
│       └── logger      # Structured and colorful logging
└── pkg
    └── prettylog       # Enhanced and structured logging implementation
```


---

## 🔧 Getting Started

### Clone the repository

```bash
git clone https://github.com/yourusername/go_bridgebot.git
cd go_bridgebot
```

### Install Dependencies

```bash
go mod tidy
```

---

## ▶️ Running BridgeBot

Execute the bot using the following command:

```bash
go run cmd/bridge_bot/main.go
```

---

## 📦 Bridgers API Endpoints

The bot utilizes the following API endpoints:

- **Token Swap Endpoint**
  ```plaintext
  POST https://api.bridgers.xyz/api/sswap/swap
```

- **Quote API Endpoint**
```plaintext
https://api.bridgers.xyz/api/sswap/quote
```

- **Token List API**
```plaintext
https://api.bridgers.xyz/api/exchangeRecord/getToken
```

---

## ⚙️ Key API Interactions

### Example Swap Request

```go
swapRequest := bridgers.SwapRequest{
  FromTokenAddress: "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
  ToTokenAddress:   "0x5aA96F60C1aFf555c43552931a177728f32fcA27",
  FromAddress:      "TMGcWzEDiECVCwAxoprCedtXSeuJthq4AA",
  ToAddress:        "0x5aA96F60C1aFf555c43552931a177728f32fcA27",
  FromTokenChain:   "TRX",
  ToTokenChain:     "BSC",
  FromTokenAmount:  "100000000",
  AmountOutMin:     "99000000",
  FromCoinCode:     "USDT(TRX)",
  ToCoinCode:       "USDT(BSC)",
  EquipmentNo:      "0000000000000000000000TMGcWzEDiECV",
  SourceFlag:       "bridgebot",
  Slippage:         "1",
}
```

### Quote Request

```go
quoteRequest := bridgers.QuoteRequest{
  FromTokenAddress: "TXYZopYRdj2D9XRtbG411XZZ3kM5VkAeBf",
  ToTokenAddress:   "0x5aA96F60C1aFf555c43552931a177728f32fcA27",
  FromTokenAmount:  "100000000",
  FromTokenChain:   "TRX",
  ToTokenChain:     "BSC",
  UserAddr:         "TMGcWzEDiECVCwAxoprCedtXSeuJthq4AA",
  EquipmentNo:      "0000000000000000000000TMGcWzEDiECV",
  SourceFlag:       "bridgebot",
}
```

### Fetch Tokens List

```go
tokens, err := bridgers.FetchTokens(ctx, bridgers.RequestBody{Chain: "TRX"})
```

---

## 📖 Logging and Debugging

This project uses a custom structured logging implementation with clear and colorful outputs:

- **Structured logging with Go's `slog` library**
- **Colorized outputs for better readability**

Example log output:
```
[15:04:05.000] INFO: Sending GET request to URL: https://api.bridgers.xyz/api/...
```

---

## 🔐 Smart Contract Approvals (TRON Example)

The BridgeBot checks and manages contract approvals for tokens on supported networks (e.g., TRON):

```go
isNeeded, allowance, err := tron.IsApprovalNeeded(ctx, client, walletAddr)
if isApprovalNeeded {
  txHash, err := tron.ApproveContract(ctx, client, privateKey)
  fmt.Println("Approval TX Hash:", txHash)
}
```

---
