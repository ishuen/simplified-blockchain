Simplified blockchain
====

This is a practice project based on the tutorial [Building Blockchain in Go](https://jeiwan.net/posts/building-blockchain-in-go-part-1/) series.

## CLI interface usage

- Add Block
    - Add a new block with specified data
    ```
    // Syntax
    go run . addblock -data "<data>"

    // Sample execution
    go run . addblock -data "send 1 BTC to Alice"
    
    // Sample output
    Mining the block containing "send 1 wld to eric"
    00000039eb5f8a5c6254013ae5a6d11ff5a0b9d0ccb7a064b1b945a32b98fec1
    Successfully added a new block.
    ```
- Print Chain
    - Print chain blocks from newest to oldest
    - When there is no blockchain, it creates a new one with genesis block
    ```
    go run . printchain

    // Sample output when there is no data
    Mining the block containing "Genesis Block"
    000000dd9c976918deb937b19bc6f87894efb92a5a26ba92abee7b4dbbd999a1
    Prev. hash: 
    Data: Genesis Block
    Hash: 000000dd9c976918deb937b19bc6f87894efb92a5a26ba92abee7b4dbbd999a1
    PoW: true


    // Sample output when the chain exists
    Prev. hash: 000000dd9c976918deb937b19bc6f87894efb92a5a26ba92abee7b4dbbd999a1
    Data: send 1 wld to eric
    Hash: 00000039eb5f8a5c6254013ae5a6d11ff5a0b9d0ccb7a064b1b945a32b98fec1
    PoW: true

    Prev. hash: 
    Data: Genesis Block
    Hash: 000000dd9c976918deb937b19bc6f87894efb92a5a26ba92abee7b4dbbd999a1
    PoW: true
    ```


## Data persistence

BoltDB is used for data persistence. When the program starts running for the first time, a file **blockchain.db** is created.