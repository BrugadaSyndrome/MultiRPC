# MultiRPC
___
I made this module to overcome the limitations of the builtin RPC module that I ran into while I was working on my
[DistributedMandelbrot](https://github.com/BrugadaSyndrome/DistributedMandelbrot) project. First off I noticed that you
can only have one connection open from a specific struct at any time. Trying to open more than one RPC connection at a
time for multiple objects of the same struct lead to all the RPC instances to no longer reference the correct object
any longer. Essentially the RPC instance now referenced a nilled object of the struct type involved.

I decided I wanted to implement the module to meet two requirements and work over the TCP and HTTP wire protocols. First 
I needed multiple objects of the same struct type to host their own RPC server without losing the reference to the
hosting object. Secondly I also needed to have the RPC connections stay alive throughout the entire execution of the
program, which could be hours or days.

I hope that you will find this module useful.