### Prerequisites

https://github.com/foundry-rs/foundry
need the tools for interacting with smart contracts

### Testnet 

#### bind dapp
```bash
cast send 0x74309Bc83fF7Ba8aBaB901936927976792a6d9B6 "bindDapp(uint256,address)" 942 0xB2Dda5D9E65E44749409E209d8b7b15fb4e82147 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```

#### bind w3bstream project
```bash
cast send 0x74309Bc83fF7Ba8aBaB901936927976792a6d9B6 "register(uint256)" 942 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```
```bash
cast send 0x0abec44FC786e8da12267Db5fdeB4311AD1A0A8A "updateConfig(uint256,string,bytes32)" 942 ipfs://ipfs.mainnet.iotex.io/QmUHfDnvWrr2wiC78dw85xfctzawNWAN1TEbzosxwHdzYC 0x8153291c230dd107f102f75e826a11d9d4a8ac3f0f4e1c3619e547f82a94410e --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```
```bash
cast send 0x0abec44FC786e8da12267Db5fdeB4311AD1A0A8A "resume(uint256)" 942 --private-key "your private key" --rpc-url "https://babel-api.testnet.iotex.io" --legacy
```