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

There is patterns in the memory address subsequent calls:

`P1`: 6027 -> 6035 -> 6048 -> 6050 -> 6054

`P2`: 6027 -> 6035 -> 6038 -> 6042 -> 6045

`P3`: 6027 -> 6030 -> 6034 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034

`P4`: 6067 -> 6047 -> 6056 -> 6059 -> 6061 -> 6065

`P5` 6067 -> 6056 -> 6059 -> 6061 -> 6065 -> 6027 -> 6030 -> 6034

`P6`: 6067 -> 6067 -> 6056 -> 6059 -> 6061 -> 6065

After `P1` there is always `P1` or `P2`
After `P2` there is always `P1` or `P3`
After `P3` there is always `P5`
After `P4` there is always `P1`
After `P5` there is always `P4`, `P5` or `P6`
After `P6` there is always `P1`

It would be better if we had a file describing every instructions

With the extractor package we extract this "assembly" code (it have been ordered by occurence)

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
while R0 != 0 {
    R0 = (R0 + 32767)%32768
    R1 = R7
    stack <- 6047
}
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
if R0 == 1 {
    stack.pop()
    R1 = 1
    R0 = 2
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