### Prerequisites

https://github.com/foundry-rs/foundry
need the tools for interacting with smart contracts

### Testnet 

#### bind w3bstream project
```bash
cast send 0xB564996622CE5610b9cF4ed35160f406185d7d0b "register(uint256)" 942 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```
```bash
cast send 0x7D3158166E9298fC47beA036fE162fEA17632E5D "updateConfig(uint256,string,bytes32)" 942 ipfs://ipfs.mainnet.iotex.io/QmUHfDnvWrr2wiC78dw85xfctzawNWAN1TEbzosxwHdzYC 0x8153291c230dd107f102f75e826a11d9d4a8ac3f0f4e1c3619e547f82a94410e --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```
```bash
cast send 0x7D3158166E9298fC47beA036fE162fEA17632E5D "resume(uint256)" 942 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```

#### bind dapp
```bash
cast send 0x19dD7163Ad80fE550C97Affef49E1995B24941B1 "bindDapp(uint256,address)" 942 0xB2Dda5D9E65E44749409E209d8b7b15fb4e82147 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```