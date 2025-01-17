# game-manager

Golang module for managing game state, playing and storing game rounds and returning game results.


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


## Use of memory pools ##

Most structs have a corresponding ```NewXyz()``` function that instantiates it from a memory pool.

This pattern is intended to reduce the number of memory allocations, and makes working with slot machines blazing fast.

Always make sure to use the appropriate ```NewXyz()``` function and call ```xyz.ReturnToPool()``` when finished.

E.g.:
```go
r := validator.NewRound(opts...)
defer r.ReturnToPool()
```
