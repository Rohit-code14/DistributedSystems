# DistributedSystems

Implementations of the things that i learnt from the book "[Distributed Services with Go](https://www.amazon.in/Distributed-Services-Go-Travis-Jeffrey/dp/1680507605)" 

# Resources
- [The Log](https://engineering.linkedin.com/distributed-systems/log-what-every-software-engineer-should-know-about-real-time-datas-unifying)
- [Building a distributed log storage](https://bravenewgeek.com/building-a-distributed-log-from-scratch-part-1-storage-mechanics/)
- [Memory mapped file operations](https://blog.labix.org/2010/11/28/removing-seatbelts-with-the-go-language-for-mmap-support)

# Learnings

Compile GRPC files

protoc api/v1/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.