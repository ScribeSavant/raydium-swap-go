# Raydium Swap Example For Go

## THIS REPO STILL UNDER DEVELOPMENT USE AT YOUR OWN RISK!!!

## Disclaimer
This repository is still in the development stage and is not ready for production use. By downloading and using this repository, you assume full responsibility for your usage. i am not liable for any consequences that may arise due to any errors, omissions, or issues in the repository since it is in the development stage.

Additionally, i do not provide any guarantee regarding the security or accuracy of any content or code provided in this repository. It is important that you carefully review and appropriately test the code before usage.

In short, please use this repository with caution and understand that you are using it at your own risk. Feel free to reach out to me for any issues, feedback, or suggestions.

## Description
This repository is not an SDK. It merely serves as an example how to swap using Raydium in the Go. I am relatively new to Go, so I may not be fully proficient in its mechanics. Feel free to suggest any improvements or changes you think would enhance the code.

## How to use?

```shell
git clone https://github.com/ScribeSavant/raydium-swap-go.git
cd raydium-swap-go
export WALLET_PRIVATE_KEY = "someimportantprivatekey"
export RPC_URL = "https://..."
go run .
```
## How to change tokens and amount ?

```javascript
// All changed should be done in run.go

executeTransaction := false // change this to true if you want to execute real transaction. Yes Real.

inputToken := utils.NewToken("SOL", "So11111111111111111111111111111111111111112", 9) // You can change this to any token which you want to sell
outputToken := utils.NewToken("RAY", "4k3Dyjzvzp8eMZWUXbBCjEvwSkkk59S5iCNLY3QrkX6R", 6) // You can change this to any token which you want to buy

slippage := utils.NewPercent(1, 100) // // this means 1% slippage, if you want to set this to 0.5 change this to utils.NewPercent(5, 1000)


amount := utils.NewTokenAmount(inputToken, 0.1) // Swap amount (0.1 sol)
```

## Common erorrs
- Pool not found
  - Solution
    - change this 
      ```javascript
      poolKeys, err := raydium.Pool.GetPoolKeys(inputToken.Mint, outputToken.Mint)
      ```
      to this
      ```javascript
      poolKeys, err := raydium.Pool.GetPoolKeys(outputToken.Mint, inputToken.Mint)
      ```
- pool info unavailable
  - This means pool not ready for swap or just try again


## References
- https://github.com/precious-void/raydium-swap
- https://github.com/gagliardetto/solana-go
- https://github.com/raydium-io/raydium-sdk
- https://github.com/solana-labs/solana-web3.js


