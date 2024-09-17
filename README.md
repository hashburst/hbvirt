# hbvirt
System to test for free Hashburst distribute node for Clusters for Cloud Computing for Crypto-Mining.

## Test:

                $ sudo go run main.go

# Compiling GO main routine:

## Windows:
                GOOS=windows GOARCH=amd64 go build -o hashburst.exe main.go
## MacOS:
                GOOS=darwin GOARCH=amd64 go build -o hashburst-darwin main.go
## Linux:
                GOOS=linux GOARCH=amd64 go build -o hashburst-linux main.go

# Run Hashburst Free Node:

## Windows:
Open like Administrator your "Prompt - Command Line":

                cd \Users\username\Desktop
                C:\\Users\username\Desktop\hashburst.exe

## MacOS:
Open Terminal (Shell):

                $ sudo ./hashburst-darwin

## Linux:
Open Terminal (Shell):

                $ sudo ./hashburst-linux

## Output:
                Network speed test output:
                PING 8.8.8.8 (8.8.8.8): 56 data bytes
                64 bytes from 8.8.8.8: icmp_seq=0 ttl=115 time=14.629 ms
                64 bytes from 8.8.8.8: icmp_seq=1 ttl=115 time=13.797 ms
                64 bytes from 8.8.8.8: icmp_seq=2 ttl=115 time=40.820 ms
                64 bytes from 8.8.8.8: icmp_seq=3 ttl=115 time=14.635 ms
                64 bytes from 8.8.8.8: icmp_seq=4 ttl=115 time=14.249 ms
                64 bytes from 8.8.8.8: icmp_seq=5 ttl=115 time=14.662 ms
                64 bytes from 8.8.8.8: icmp_seq=6 ttl=115 time=14.025 ms
                64 bytes from 8.8.8.8: icmp_seq=7 ttl=115 time=46.045 ms
                64 bytes from 8.8.8.8: icmp_seq=8 ttl=115 time=13.788 ms
                64 bytes from 8.8.8.8: icmp_seq=9 ttl=115 time=15.217 ms
                
                --- 8.8.8.8 ping statistics ---
                10 packets transmitted, 10 packets received, 0.0% packet loss
                round-trip min/avg/max/stddev = 13.788/20.187/46.045/11.689 ms
                
                2024/09/14 02:53:30 Download starting: HashburstMiner
                2024/09/14 02:53:30 Download complete: HashburstMiner
                2024/09/14 02:53:30 Starting server on :8080...

## Start Web GUI:

Open your favourite browser and go to this default url: [http://localhost:8080](http://localhost:8080)

![Web-GUI](https://github.com/user-attachments/assets/2325e744-545c-42bd-a83e-3bb2e50d5ca7)
