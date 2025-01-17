# game-engine

Repository for game related golang modules.


## Makefile ##
Getting help:
```shell
make help
```
Run all tests:
```shell
make test
```
View test coverage in your browser:
```shell
make coverhtml
```

The test code contains a considerable number of benchmarks.
Running the benchmarks:
```shell
make bench
```
The benchmarks cover all the critical hot paths of the code that can influence the speed of simulations.


## Use of memory pools ##

Most structs have a corresponding ```AcquireXyz()``` function that instantiates it from a memory pool.

This pattern is intended to reduce the number of memory allocations, and makes working with slot machines blazing fast.

Always make sure to use the appropriate ```AcquireXyz()``` function and call ```xyz.Release()``` when finished.

E.g.:
```go
spin := slots.AcquireSpin(slots, prng)
defer spin.Release()
```

Some structs take ownership of embedded structs, so don't call ```Release()``` on the embedded structs.
An example of this are the *Winline*'s in a *Spin* result.
These are automatically returned to the memory pool when you call ```spin.Release()```.
