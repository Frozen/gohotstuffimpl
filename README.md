This repository contains the code for the paper [HotStuff: BFT Consensus in the Lens of Blockchain](https://arxiv.org/pdf/1803.05069).
Key differency between Harmony and HotStuff implementation are:
- Hotstuff provides new leader for each new block.
- Hotstuff validator signs PBFT message and generates new blocks at the same time, which saves at least 1 communication. 