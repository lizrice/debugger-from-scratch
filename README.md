# debugger-from-scratch

This is a very simple debugger that shows how ptrace can be used to set breakpoints and generate stack traces in a target process. You can choose where to breakpoint the target process, and single step or continue to the next breakpoint (or target exit). 

Big caveat: it assumes there is only one file in the target! 

I fully expect this only to work on Linux. 

This is a more detailed version of the code I wrote in my [Debuggers From Scratch talk at dotGo Paris](https://youtu.be/TBrv17QyUE0).
