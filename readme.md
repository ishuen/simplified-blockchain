Simplified blockchain
====

This is a practice project based on the tutorial [Building Blockchain in Go](https://jeiwan.net/posts/building-blockchain-in-go-part-1/) series.

## CLI interface usage

- Get Balance
    - Get balance of the specified address
    - syntax: `getbalance -address <ADDRESS>`
      - Shows the balance of the specified address.
      ```
      go run . getbalance -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa

      // Sample output
      Balance of '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa': 8
      ```

- Create Blockchain
    - Create a blockchain and send genesis block reward to the specified address
    - syntax: `createblockchain -address <ADDRESS>`
    ```
    go run . createblockchain -address 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa
    
    // Sample output
    000000e5c2d2333a60f10e136b6c16aa38618080ff5b1ff4c9f968a2463a08dc
    Done!
    ```

- Send
    - Send a certain amount of coins from one address to the other address
    - syntax: `send -from <FROM_ADDRESS> -to <TO_ADDRESS> -amount <AMOUNT>`
    ```
    go run . send -from 1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa -to 9B1zP2wP5QGefi2DMetFmL49Lmv7viDfbB -amount 2
    // Sample output
    0000002e266352367cb319d073c8515ed10b2cdf01482b8e77b934e95910bd67
    Success!
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