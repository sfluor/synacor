## Code 5

Solution for the door enigma in coins package

## Code 6

Progression until the synacor headquarters is saved in the `src/moves.record` file.
To use it simply do `cat moves.record /dev/stdin | ./< vm binary > < path to the challenge.bin >` 


## Code 7

The strange book content says that we should try to modify the last register value to get access to a new area.

To do so, it's advised to extract the confirmation algorithm. Indeed the last register value needs to be confirmed by a computationally expensive algorithm.

Implementation of a logger to see the registers values (not so relevant)

Implementation of a new command "$" that can debug by printing the current register state. By forcing the R7 register to take a non zero value here we can trigger the confirmation process. We can notice that it calls the following memory addresses "6067 6056 6059 6061 6065 6027 6030 6034"

It would be better if we had a file describing every instructions