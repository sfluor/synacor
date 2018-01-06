## Code 5

Solution for the door enigma in coins package

## Code 6

Progression until the synacor headquarters is saved in the `src/moves.record` file.
To use it simply do `cat moves.record /dev/stdin | ./< vm binary > < path to the challenge.bin >` 


## Code 7

The strange book content says that we should try to modify the last register value to get access to a new area.

To do so, it's advised to extract the confirmation algorithm. Indeed the last register value needs to be confirmed by a computationally expensive algorithm.

Implementation of a logger to see the registers values (not so relevant)

Implementation of a new command "$" that can debug by printing the current register state. By forcing the R7 register to take a non zero value here we can trigger the confirmation process. 

There are patterns in the memory address subsequent calls:

`P0`: 5472 -> 5474 -> 5476 -> 5478 -> 5479 -> 5480 -> 5481 -> 5482 -> 5483 -> 5486 -> 5489

`P1`: 6027 -> 6035 -> 6048 -> 6050 -> 6054

`P2`: 6027 -> 6035 -> 6038 -> 6042 -> 6045

`P3`: 6027 -> 6030 -> 6034 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034

`P4`: 6067 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065

`P5` 6067 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034

`P6`: 6067 -> 6067 -> 6056 -> 6059 -> 6061 -> 6065

`P0` seems to happen only one time at the very beginning then we have `P1`

After `P1` there is always `P1` or `P2`

After `P2` there is always `P1` or `P3`

After `P3` there is always `P5`

After `P4` there is always `P1`

After `P5` there is always `P4`, `P5` or `P6`

After `P6` there is always `P1`

It would be better if we had a file describing every instructions

With the extractor package we extract this "assembly" code (it have been ordered by occurence)

for `P0` `P0`: 5472 -> 5474 -> 5476 -> 5478 -> 5479 -> 5480 -> 5481 -> 5482 -> 5483 -> 5486 -> 5489

```
(  5472) |  pop: [R2]
(  5474) |  pop: [R1]
(  5476) |  pop: [R0]
(  5478) | noop: []
(  5479) | noop: []
(  5480) | noop: []
(  5481) | noop: []
(  5482) | noop: []
(  5483) |  set: [R0 4]
(  5486) |  set: [R1 1]
(  5489) | call: [6027]
```

Which does:

```
R2 = stack.pop()
stack.pop()
stack.pop()
R0 = 4
R1 = 1
goto 6027
```

for `P1` 6027 -> 6035 -> 6048 -> 6050 -> 6054
```
(  6027) |   jt: [R0 6035]
(  6035) |   jt: [R1 6048]
(  6048) | push: [R0]
(  6050) |  add: [R1 R1 32767]
(  6054) | call: [6027]
```

In a more readable format:
- If R0 is not zero go to 6035 (here R0 is not zero)
- If R1 is not zero go to 6048 (here R1 is not zero)
- Push R0 on the stack
- Add 32767 to R1 (modulo 32768)
- Write the next instruction to the stack (`(  6056) |  set: [R1 R0]`) and loop again

It can be resumed to:
```
if R0 != 0 {
    while R1 != 0 {
        R1 = (R1 + 32767)%32768
        stack <- R0
        stack <- 6056
    }
}
```

for `P2` 6027 -> 6035 -> 6038 -> 6042 -> 6045
```
(  6027) |   jt: [R0 6035]
(  6035) |   jt: [R1 6048]
(  6038) |  add: [R0 R0 32767]
(  6042) |  set: [R1 R7]
(  6045) | call: [6027]
```

- If R0 is not zero (here it's not) go to 6035
- If R1 is not zero go to 6048 (here R1 is zero) so we go to the next instruction
- Add 32767 to R0 (modulo 32768)
- Set R1 to the value of R7
- Write the next instruction to the stack (`(  6047) |  ret: []`) and loop again

it can be resumed to:
```
R0 = (R0 + 32767)%32768
R1 = R7
stack <- 6047
```


for `P3`  6027 -> 6030 -> 6034 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034
```
(  6027) |   jt: [R0 6035]
(  6030) |  add: [R0 R1 1]
(  6034) |  ret: []
(  6047) |  ret: []
(  6056) |  set: [R1 R0]
(  6059) |  pop: [R0]
(  6061) |  add: [R0 R0 32767]
(  6065) | call: [6027]
(  6027) |   jt: [R0 6035]
(  6030) |  add: [R0 R1 1]
(  6034) |  ret: []
```

- If R0 is not 0 go to 6035 (but R0 is zero)
- Set R0 to R1 + 1
- Remove the top element of the stack (6047) and jump to it
- set R1 to the value of R0
- Remove the top element of the stack and write its value to R0
- Add 32767 to R0 (modulo 32767)
- Write the next instruction to the stack `(  6067) |  ret: []` and jump to 6027
- If R0 is not 0 (still not 0) go to 6035 (not done)
- Set R0 to R1 + 1
- Remove 6067 from the stack and go to 6067

it can be resumed to:
```
R1 = R1 + 1
stack.pop()
stack.pop()
R0 = R1 + 1
goto 6067
```

for `P4` 6067 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065
```
(  6067) |  ret: []
(  6047) |  ret: []
(  6056) |  set: [R1 R0]
(  6059) |  pop: [R0]
(  6061) |  add: [R0 R0 32767]
(  6065) | call: [6027]
```

- Remove the top element of the stack (6047) and jump to it
- Remove the top element of the stack (6056) and jump to it
- Set R1 to the value of R0
- Remove the top element of the stack and set R0 to its value
- Add 32767 to R0 (modulo 32767)
- Add `(  6067) |  ret: []` to the stack and jump to 6027

it can be resumed to:
```
stack.pop()
stack.pop()
R1 = R0
R0 = stack.pop()
R0 = (R0 + 32767)%32768
stack <- 6067
goto 6027
```

for `P5` 6067 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034
```
(  6067) |  ret: []
(  6056) |  set: [R1 R0]
(  6059) |  pop: [R0]
(  6061) |  add: [R0 R0 32767]
(  6065) | call: [6027]
(  6027) |   jt: [R0 6035]
(  6030) |  add: [R0 R1 1]
(  6034) |  ret: []
```

So the P5 pattern do:
- Remove and jump to the top element of the stack (6056)
- Set register R1 to the value of register R0
- Remove the top element of the stack and write it to R0
- Add 32767 to R0 (modulo 32768)
- Write the address 6067 to the stack and jump to address 6027
- If R0 is not 0 jump to 6035 (but in our case we go to 6030 so R0 is 0)
- Set R0 to R1 + 1
- Remove 6067 from the stack and jump to 6067 (loop)

it can be resumed to:
```
while R0 != 0 {
    stack.pop()
    R1 = R0
    stack.pop()
    R0 = R1 + 1
}
```

for `P6` 6067 -> 6067 -> 6056 -> 6059 -> 6061 -> 6065
```
(  6067) |  ret: []
(  6067) |  ret: []
(  6056) |  set: [R1 R0]
(  6059) |  pop: [R0]
(  6061) |  add: [R0 R0 32767]
(  6065) | call: [6027]
```
- Remove and jump to the top element of the stack (6067)
- Remove and jump to the top element of the stack (6067)
- Set R1 to R0
- Remove the last element of the stack and write its value to R0
- Add 32767 to R0 (modulo 32768)
- Add 6067 to the stack and go to 6027

it can be resumed to:
```
stack.pop()
stack.pop()
R1 = R0
R0 = stack.pop()
R0 = (R0 + 32767)%32768
stack <- 6067
goto 6027
```

---------

So if we sum up:

We have `P1` until `R1` == 0 (it pushes `R0 != 0` and 6056 to te stack plenty of times)

Then we have `P2`, that adds 32767 to `R0`, sets `R1` to `R7` and push 6047 to the stack

Then If `R0` is zero we have `P3` else we have `P1` that sets `R1` to 0 again, this (`n*P1 + P2`) repeats until we have `R0 == 0` (we then have `R0 == 0` and `R1 == R7`)

Then in `P3` `R1++` and `R0 = R1+1` (popping the stack twice) then we have `P5`

In `P5` we pop the stack multiple times until `R0 == 0` and `R1 = R0 - 1 == 32767`

Then he have `P4` or `P6` but they have the same behaviour: two stack.pop(), set `R1` to 0, stack.pop() into `R0` and adds 32767 to `R0` then we go back to `P1`

---------
Let's try to have a more global look on theses patterns (they are all between 6027 and 6067):

```
(  6027) |   jt: [R0 6035]
(  6030) |  add: [R0 R1 1]
(  6034) |  ret: []
(  6035) |   jt: [R1 6048]
(  6038) |  add: [R0 R0 32767]
(  6042) |  set: [R1 R7]
(  6045) | call: [6027]
(  6047) |  ret: []
(  6048) | push: [R0]
(  6050) |  add: [R1 R1 32767]
(  6054) | call: [6027]
(  6056) |  set: [R1 R0]
(  6059) |  pop: [R0]
(  6061) |  add: [R0 R0 32767]
(  6065) | call: [6027]
(  6067) |  ret: []
```

Can be transformed to:

```
    // P1 block
    if R0 == 0 {
        R0 = R1 + 1
        goto stack.pop()
    }
    // P2 block
    if R1 == 0 {
        R0 = R0 - 1
        R1 = R7
        goto 6027
        goto stack.pop()
    }

    stack.push(R0)
    R1 = R1 - 1
    goto 6027
    R1 = R0
    R0 = stack.pop()
    R0 = R0 - 1
    goto 6027
    goto stack.pop()
```

The strange book speaks of a "confirmation algorithm", if we assume this algorithm can be described as a function and that we consider `R0`, `R1` and `R7` as global variables we can rewrite the previous set of instructions like the following (all operations are implicitly modulo 32768):

```go
func f() {
    if R0 == 0 {
        // If R0 is zero set it to R1 + 1
        R0 = R1 + 1
        return
    } else if R1 == 0 {
        // If R1 is zero call the function on (R0 - 1, R7, R7)
        R0 = R0 - 1
        R1 = R7
        f()
        return
    } else {
        // Else save R0
        // Set R1' to the resulting R0' of the function on (R0, R1 - 1, R7)
        // Retrieve the old R0 (before applying the function) thanks to the stack
        // Call the function on (R0 - 1, R1', R7)
        stack.push(0)
        R1 = R1 - 1
        f()
        R1 = R0
        R0 = stack.pop()
        R0 = R0 - 1
        f()
        return        
    }
}
```

And we can see that in fact the patterns described before are only parts of this function that is recursively called with the `call: [6027]` instructions, and that this algorithm is computationally expensive beceause of recursions

--------

If we know look at what happens between the `use teleporter` command and the first pattern (`P0`) for the "confirmation", we have:

```
# R7 check

5451 [2708 5445 3 10 101 0 0 7]   <-- R7 check
5454 [2708 5445 3 10 101 0 0 7]
5456 [2708 5445 3 10 101 0 0 7]
5458 [2708 5445 3 10 101 0 0 7]
5460 [2708 5445 3 10 101 0 0 7]
5463 [28844 5445 3 10 101 0 0 7]
5466 [28844 1531 3 10 101 0 0 7]
5470 [28844 1531 30326 10 101 0 0 7]
1458 [28844 1531 30326 10 101 0 0 7]
1460 [28844 1531 30326 10 101 0 0 7]
1462 [28844 1531 30326 10 101 0 0 7]
1464 [28844 1531 30326 10 101 0 0 7]
1466 [28844 1531 30326 10 101 0 0 7]
1468 [28844 1531 30326 10 101 0 0 7]
1471 [28844 1531 30326 10 101 0 28844 7]
1474 [28844 1531 30326 10 101 1531 28844 7]
1477 [28844 1531 30326 10 169 1531 28844 7]
1480 [28844 0 30326 10 169 1531 28844 7]
1484 [28844 0 30326 1 169 1531 28844 7]
1488 [0 0 30326 1 169 1531 28844 7]
1491 [0 0 30326 1 169 1531 28844 7]
1495 [0 0 30326 28845 169 1531 28844 7]
1498 [30263 0 30326 28845 169 1531 28844 7]
1531 [30263 0 30326 28845 169 1531 28844 7]
1533 [30263 0 30326 28845 169 1531 28844 7]
1536 [30263 30326 30326 28845 169 1531 28844 7]
2125 [30263 30326 30326 28845 169 1531 28844 7]
2127 [30263 30326 30326 28845 169 1531 28844 7]
2129 [30263 30326 30326 28845 169 1531 28844 7]
2133 [30263 30326 30262 28845 169 1531 28844 7]
2136 [30263 30326 2505 28845 169 1531 28844 7]
2140 [30327 30326 2505 28845 169 1531 28844 7]
2144 [65 30326 2505 28845 169 1531 28844 7]
2146 [65 30326 30326 28845 169 1531 28844 7]
2148 [65 30326 30326 28845 169 1531 28844 7]

Print  A strange.... years."

1540 [10 30326 30326 29013 169 1531 28844 7]
1542 [10 168 30326 29013 169 1531 28844 7]
1500 [10 168 30326 29013 169 1531 28844 7]
1504 [10 169 30326 29013 169 1531 28844 7]
1480 [10 169 30326 29013 169 1531 28844 7]
1484 [10 169 30326 170 169 1531 28844 7]
1488 [1 169 30326 170 169 1531 28844 7]
1507 [1 169 30326 170 169 1531 28844 7]
1509 [1 169 30326 170 169 1531 0 7]
1511 [1 169 30326 170 169 0 0 7]
1513 [1 169 30326 170 101 0 0 7]
1515 [1 169 30326 10 101 0 0 7]
1517 [28844 169 30326 10 101 0 0 7]
5472 [28844 169 30326 10 101 0 0 7]
5474 [28844 169 3 10 101 0 0 7]
5476 [28844 5445 3 10 101 0 0 7]
5478 [2708 5445 3 10 101 0 0 7]
5479 [2708 5445 3 10 101 0 0 7]
5480 [2708 5445 3 10 101 0 0 7]
5481 [2708 5445 3 10 101 0 0 7]
5482 [2708 5445 3 10 101 0 0 7]
5483 [2708 5445 3 10 101 0 0 7]
5486 [4 5445 3 10 101 0 0 7]
5489 [4 1 3 10 101 0 0 7]

Beginning of P0

6027 [4 1 3 10 101 0 0 7]
6035 [4 1 3 10 101 0 0 7]
6048 [4 1 3 10 101 0 0 7]
6050 [4 1 3 10 101 0 0 7]
```

We can notice that the entry values for the function are (4, 1 and R7), here R7 is set to 7 thanks to the debugging commands.