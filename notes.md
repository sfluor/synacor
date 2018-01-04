## Code 5

Solution for the door enigma in coins package

## Code 6

Progression until the synacor headquarters is saved in the `src/moves.record` file.
To use it simply do `cat moves.record /dev/stdin | ./< vm binary > < path to the challenge.bin >` 


## Code 7

The strange book content says that we should try to modify the last register value to get access to a new area.

To do so, it's advised to extract the confirmation algorithm. Indeed the last register value needs to be confirmed by a computationally expensive algorithm.

Implementation of a logger to see the registers values.
