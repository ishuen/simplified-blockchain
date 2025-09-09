Simplified blockchain
====

This is a practice project based on the tutorial [Building Blockchain in Go](https://jeiwan.net/posts/building-blockchain-in-go-part-1/) series.

## CLI interface usage

- Print Chain
    - Print chain blocks from newest to oldest
    - When there is no blockchain, it creates a new one with genesis block
    ```
    go run . printchain
    ```
- Add Block
    - Add a new block with specified data
    ```
    go run . addblock -data "send 1 BTC to Alice"
    ```

## Data persistence

BoltDB is used for data persistence. When the program starts running for the first time, a file **blockchain.db** is created.