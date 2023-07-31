.PHONY: node-label
node-label:
	kubectl label nodes bcsvr02 ethbaas_node=node0 --overwrite
	kubectl label nodes bcsvr03 ethbaas_node=node1 --overwrite
	kubectl label nodes bcsvr04 ethbaas_node=node2 --overwrite

.PHONY: contract-complie
contract-complie:
	docker run --rm -v $(shell pwd)/contract:/contract ethereum/solc:0.4.24 \
		--abi --bin /contract/store/Store.sol -o /contract/store --overwrite

	docker run --rm -v $(shell pwd)/contract:/contract ethereum/client-go:alltools-v1.10.18-amd64 abigen \
		--bin /contract/store/Store.bin --abi /contract/store/Store.abi --pkg=store --out /contract/store/store.go
