## ProcessIsolator
A tool used for running specified process in isolated environment. Enjoy it, just for fun.

### Usage
To start, we provide two ways:

- Run it in front, you can use `./ProcessIsolator start`
- Run it in background, you can use `./ProcessIsolator -daemon -out logfilepath`

### Alert

1. If you need to change the RPC's port, see `constants/default.go`
2. For the first run, the path of `/run/ProcZygote/logs` should be available
3. Only root user can run it, because of the linux's restriction on namespace