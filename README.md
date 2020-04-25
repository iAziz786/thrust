# 🔥Thrust

A load testing tool

The tester takes three arguments:
- second - Number of seconds to run the test
- open-sockets - Total number of active connection any time
- endpoint - The endpoint where you want to test the loads

Example

```sh
./main 5 1000 https://example.com
```

This example will poll https://example.com with 1000 active connection for 5 seconds.