## Tool Description

This tool will agregate metrics across a fleet of servers by hitting /status endpoint on each of them , this will be done in parallel using a pool of go routines. The server list will be read froma input file and fed into go channel, which will server as input for the pool of go routines that will write the results to a result channel. Once they are done hitting all servers results are agregated.

 
### Build Instructions on MAC with brew installed
git clone https://github.com/lrfurtado/metric_poller metrics_poller
cd metrics_poller
make deps
make

### Test instructions 

1. Open 1st terminal window and run **make servers**
1. Open 2nd terminal window and run **make test**
