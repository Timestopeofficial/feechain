#!/usr/bin/env bash
# This Script is for Testing the API functionality on both local and betanet.
# -l to run localnet, -b to run betanet(mutually exclusive)
# -v to see returns from each request
# Right now only tests whether a response is recieved
# You must have properly cloned into the dapp-examples repo and have installed nodejs
VERBOSE="FALSE"
TESTS_RAN=0
TESTS_PASSED=0

red=`tput setaf 1`
green=`tput setaf 2`
blue=`tput setaf 6`
white=`tput sgr0`
yellow=`tput setaf 11`
reset=`tput sgr0`

function response_test() {
	if [ "$1" != "" ]; then
		echo "${green}RESPONSE RECIEVED${reset}"
		return 1
	else
		echo "${red}NO RESPONSE${reset}"
		return 0
	fi
}

function isHashTest() {
	if [ "$TRANSACTION" != "null" ]; then
		if [[ "$TRANSACTION_HASH" =~ ^0x[0-9a-f]{64}$ ]]; then
			echo ${green}TRANSACTION HASH VALID${reset}
			echo
			return 1
		fi
	fi
	echo ${red}TRANSACTION HASH INVALID${reset}
	return 0
}

function isHexTest() {
	if [ "$1" != "null" ]; then
		if [[ "$1" =~ ^0x[0-9a-f]+$ ]]; then
			echo ${green}VALID HEX RECIEVED${reset}
			echo
			return 1
		fi
	fi
	echo ${red}INVALID HEX RECIEVED${reset}
	return 0
}

### SETUP COMMANDLINE FLAGS ###
while getopts "lbvp" OPTION; do
	case $OPTION in
	b)
		NETWORK="betanet"
		declare -A PORT=( [POST]="http://s0.t.timestope.net:9500/" [GET]="http://e0.t.timestope.net:5000/" )
		BLOCK_0_HASH=$(curl --location --request POST "http://l0.t.timestope.net:9500" \
			  --header "Content-Type: application/json" \
			  --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_getBlockByNumber\",\"params\":[\"0x1\", true],\"id\":1}" | jq -r '.result.hash')
		echo "BLOCK0HASH:"
		echo "$BLOCK_0_HASH"

		SIGNED_RAW_TRANSACTION=$(node ../dapp-examples/nodejs/apiTestSign.js)
		echo "RAWTX"
		echo "$SIGNED_RAW_TRANSACTION"
		TRANSACTION_HASH=$(curl  --location --request POST "http://l0.t.timestope.net:9500" \
			  --header "Content-Type: application/json" \
			    --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_sendRawTransaction\",\"params\":[\""$SIGNED_RAW_TRANSACTION"\"],\"id\":1}" | jq -r '.result')
		echo "TRANSACTION_HASH:"
		echo $TRANSACTION_HASH
		sleep 20s
		TRANSACTION=$(curl --location --request POST "http://l0.t.timestope.net:9500" \
			  --header "Content-Type: application/json" \
			  --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_getTransactionByHash\",\"params\":[\"$TRANSACTION_HASH\"],\"id\":1}")

		echo "TRANSACTION:"
		echo "$TRANSACTION"

		TRANSACTION_BLOCK_HASH=$(echo $TRANSACTION | jq -r '.result.blockHash')
		TRANSACTION_BLOCK_NUMBER=$(echo $TRANSACTION | jq -r '.result.blockNumber')
		TRANSACTION_INDEX=$(echo $TRANSACTION | jq -r '.result.transactionIndex')  #Needs to be get transaction Index


		TRANSACTION_BLOCK_ID=$(( $TRANSACTION_BLOCK_NUMBER ))
		echo TRANSACTION_BLOCK_ID
		echo $TRANSACTION_BLOCK_ID
		echo "TRANSACTION_BLOCK_HASH:"
		echo $TRANSACTION_BLOCK_HASH

		echo "TRANSACTION_BLOCK_NUMBER:"
		echo "$TRANSACTION_BLOCK_NUMBER"

		echo "TRANSACTION_INDEX:"
		echo $TRANSACTION_INDEX

		;;
	l)
		NETWORK="localnet"
		declare -A PORT=( [POST]="localhost:9500/" [GET]="localhost:5099/" )
		BLOCK_0_HASH=$(curl -s --location --request POST "localhost:9500" \
			  --header "Content-Type: application/json" \
			  --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_getBlockByNumber\",\"params\":[\"0x1\", true],\"id\":1}" | jq -r '.result.hash')

		echo "BLOCK0HASH:"
		echo "$BLOCK_0_HASH"

		SIGNED_RAW_TRANSACTION=$(node ../dapp-examples/nodejs/apiTestSign.js localnet)
		echo "RAWTX"
		echo "$SIGNED_RAW_TRANSACTION"
		TRANSACTION_HASH=$(curl  --location --request POST "localhost:9500" \
			  --header "Content-Type: application/json" \
			    --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_sendRawTransaction\",\"params\":[\""$SIGNED_RAW_TRANSACTION"\"],\"id\":1}" | jq -r '.result')
		echo "TRANSACTION_HASH:"
		echo $TRANSACTION_HASH
		sleep 20s
		TRANSACTION=$(curl --location --request POST "http://localhost:9500" \
			  --header "Content-Type: application/json" \
			  --data "{\"jsonrpc\":\"2.0\",\"method\":\"fch_getTransactionByHash\",\"params\":[\"$TRANSACTION_HASH\"],\"id\":1}")

		echo "TRANSACTION:"
		echo "$TRANSACTION"

		TRANSACTION_BLOCK_HASH=$(echo $TRANSACTION | jq -r '.result.blockHash')
		TRANSACTION_BLOCK_NUMBER=$(echo $TRANSACTION | jq -r '.result.blockNumber')
		TRANSACTION_INDEX=$(echo $TRANSACTION | jq -r '.result.transactionIndex')

		TRANSACTION_BLOCK_ID=$(( $TRANSACTION_BLOCK_NUMBER ))
		echo TRANSACTION_BLOCK_ID
		echo $TRANSACTION_BLOCK_ID
		echo "TRANSACTION_BLOCK_HASH:"
		echo $TRANSACTION_BLOCK_HASH

		echo "TRANSACTION_BLOCK_NUMBER:"
		echo "$TRANSACTION_BLOCK_NUMBER"

		echo "TRANSACTION_INDEX:"
		echo $TRANSACTION_INDEX
		;;
	v)
		VERBOSE="TRUE"
		;;
	p)
		PRETTY="TRUE"
		;;
	esac
		done

if [ $OPTIND -eq 1 ]; then echo "No options were passed, -l for localnet, -b for betanet, -v to view logs of either"; exit;  fi

declare -A GETDATA=( [GET_blocks]="blocks?from=$TRANSACTION_BLOCK_ID&to=$TRANSACTION_BLOCK_ID" [GET_tx]="tx?id=0" [GET_address]="address?id=0" [GET_node-count]="node-count" [GET_shard]="shard?id=0" [GET_committee]="committee?shard_id=0&epoch=0" )
declare -A POSTDATA

if [ "$NETWORK" == "localnet" ]; then
	POSTDATA[fch_getBlockByHash]="fch_getBlockByHash\",\"params\":[\"$TRANSACTION_BLOCK_HASH\", true]"
	POSTDATA[fch_getBlockByNumber]="fch_getBlockByNumber\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\", true]"
	POSTDATA[fch_getBlockTransactionCountByHash]="fch_getBlockTransactionCountByHash\",\"params\":[\"$TRANSACTION_BLOCK_HASH\"]"
	POSTDATA[fch_getBlockTransactionCountByNumber]="fch_getBlockTransactionCountByNumber\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\"]"
	POSTDATA[fch_getCode]="fch_getCode\",\"params\":[\"0x08AE1abFE01aEA60a47663bCe0794eCCD5763c19\", \"latest\"]"
	POSTDATA[fch_getTransactionByBlockHashAndIndex]="fch_getTransactionByBlockHashAndIndex\",\"params\":[\"$TRANSACTION_BLOCK_HASH\", \"$TRANSACTION_INDEX\"]"
	POSTDATA[fch_getTransactionByBlockNumberAndIndex]="fch_getTransactionByBlockNumberAndIndex\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\", \"$TRANSACTION_INDEX\"]"
	POSTDATA[fch_getTransactionByHash]="fch_getTransactionByHash\",\"params\":[\"$TRANSACTION_HASH\"]"
	POSTDATA[fch_getTransactionReceipt]="fch_getTransactionReceipt\",\"params\":[\"$TRANSACTION_HASH\"]"
	POSTDATA[fch_syncing]="fch_syncing\",\"params\":[]"
	POSTDATA[net_peerCount]="net_peerCount\",\"params\":[]"
	POSTDATA[fch_getBalance]="fch_getBalance\",\"params\":[\"one18t4yj4fuutj83uwqckkvxp9gfa0568uc48ggj7\", \"latest\"]"
	POSTDATA[fch_getStorageAt]="fch_getStorageAt\",\"params\":[\"0xD7Ff41CA29306122185A07d04293DdB35F24Cf2d\", \"0\", \"latest\"]"
	POSTDATA[fch_getAccountNonce]="fch_getAccountNonce\",\"params\":[\"0x806171f95C5a74371a19e8a312c9e5Cb4E1D24f6\", \"latest\"]"
	POSTDATA[fch_sendRawTransaction]="fch_sendRawTransaction\",\"params\":[\"$SIGNED_RAW_TRANSACTION\"]"
	POSTDATA[fch_getLogs]="fch_getLogs\", \"params\":[{\"BlockHash\": \"$TRANSACTION_BLOCK_HASH\"}]"
	POSTDATA[fch_getFilterChanges]="fch_getFilterChanges\", \"params\":[\"0x58010795a282878ed0d61da72a14b8b0\"]"
	POSTDATA[fch_newPendingTransactionFilter]="fch_newPendingTransactionFilter\", \"params\":[]"
	POSTDATA[fch_newBlockFilter]="fch_newBlockFilter\", \"params\":[]"
	POSTDATA[fch_newFilter]="fch_newFilter\", \"params\":[{\"BlockHash\": \"0x5725b5b2ab28206e7256a78cda4f9050c2629fd85110ffa54eacd2a13ba68072\"}]"
	POSTDATA[fch_call]="fch_call\", \"params\":[{\"to\": \"0x08AE1abFE01aEA60a47663bCe0794eCCD5763c19\"}, \"latest\"]"
	POSTDATA[fch_gasPrice]="fch_gasPrice\",\"params\":[]"
	POSTDATA[fch_blockNumber]="fch_blockNumber\",\"params\":[]"
	POSTDATA[net_version]="net_version\",\"params\":[]"
	POSTDATA[fch_protocolVersion]="fch_protocolVersion\",\"params\":[]"
fi

if [ "$NETWORK" == "betanet" ]; then
	POSTDATA[fch_getBlockByHash]="fch_getBlockByHash\",\"params\":[\"$TRANSACTION_BLOCK_HASH\", true]"
	POSTDATA[fch_getBlockByNumber]="fch_getBlockByNumber\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\", true]"
	POSTDATA[fch_getBlockTransactionCountByHash]="fch_getBlockTransactionCountByHash\",\"params\":[\"$TRANSACTION_BLOCK_HASH\"]"
	POSTDATA[fch_getBlockTransactionCountByNumber]="fch_getBlockTransactionCountByNumber\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\"]"
	POSTDATA[fch_getCode]="fch_getCode\",\"params\":[\"0x08AE1abFE01aEA60a47663bCe0794eCCD5763c19\", \"latest\"]"
	POSTDATA[fch_getTransactionByBlockHashAndIndex]="fch_getTransactionByBlockHashAndIndex\",\"params\":[\"$TRANSACTION_BLOCK_HASH\", \"$TRANSACTION_INDEX\"]"
	POSTDATA[fch_getTransactionByBlockNumberAndIndex]="fch_getTransactionByBlockNumberAndIndex\",\"params\":[\"$TRANSACTION_BLOCK_NUMBER\", \"$TRANSACTION_INDEX\"]"
	POSTDATA[fch_getTransactionByHash]="fch_getTransactionByHash\",\"params\":[\"$TRANSACTION_HASH\"]"
	POSTDATA[fch_getTransactionReceipt]="fch_getTransactionReceipt\",\"params\":[\"$TRANSACTION_HASH\"]"
	POSTDATA[fch_syncing]="fch_syncing\",\"params\":[]"
	POSTDATA[net_peerCount]="net_peerCount\",\"params\":[]"
	POSTDATA[fch_getBalance]="fch_getBalance\",\"params\":[\"one18t4yj4fuutj83uwqckkvxp9gfa0568uc48ggj7\", \"latest\"]"
	POSTDATA[fch_getStorageAt]="fch_getStorageAt\",\"params\":[\"0xD7Ff41CA29306122185A07d04293DdB35F24Cf2d\", \"0\", \"latest\"]"
	POSTDATA[fch_getAccountNonce]="fch_getAccountNonce\",\"params\":[\"0x806171f95C5a74371a19e8a312c9e5Cb4E1D24f6\", \"latest\"]"
	POSTDATA[fch_sendRawTransaction]="fch_sendRawTransaction\",\"params\":[\"$SIGNED_RAW_TRANSACTION\"]"
	POSTDATA[fch_getLogs]="fch_getLogs\", \"params\":[{\"BlockHash\": \"$TRANSACTION_BLOCK_HASH\"}]"
	POSTDATA[fch_getFilterChanges]="fch_getFilterChanges\", \"params\":[\"0x58010795a282878ed0d61da72a14b8b0\"]"
	POSTDATA[fch_newPendingTransactionFilter]="fch_newPendingTransactionFilter\", \"params\":[]"
	POSTDATA[fch_newBlockFilter]="fch_newBlockFilter\", \"params\":[]"
	POSTDATA[fch_newFilter]="fch_newFilter\", \"params\":[{\"BlockHash\": \"0x5725b5b2ab28206e7256a78cda4f9050c2629fd85110ffa54eacd2a13ba68072\"}]"
	POSTDATA[fch_call]="fch_call\", \"params\":[{\"to\": \"0x08AE1abFE01aEA60a47663bCe0794eCCD5763c19\"}, \"latest\"]"
	POSTDATA[fch_gasPrice]="fch_gasPrice\",\"params\":[]"
	POSTDATA[fch_blockNumber]="fch_blockNumber\",\"params\":[]"
	POSTDATA[net_version]="net_version\",\"params\":[]"
	POSTDATA[fch_protocolVersion]="fch_protocolVersion\",\"params\":[]"
fi

declare -A RESPONSES

RESPONSES[GET_blocks]=""
RESPONSES[GET_tx]=""
RESPONSES[GET_address]=""
RESPONSES[GET_node-count]=""
RESPONSES[GET_shard]=""
RESPONSES[GET_committee]=""
RESPONSES[fch_getBlockByHash]=""
RESPONSES[fch_getBlockByNumber]=""
RESPONSES[fch_getBlockTransactionCountByHash]=""
RESPONSES[fch_getBlockTransactionCountByNumber]=""
RESPONSES[fch_getCode]=""
RESPONSES[fch_getTransactionByBlockHashAndIndex]=""
RESPONSES[fch_getTransactionByBlockNumberAndIndex]=""
RESPONSES[fch_getTransactionByHash]=""
RESPONSES[fch_getTransactionReceipt]=""
RESPONSES[fch_syncing]=""
RESPONSES[net_peerCount]=""
RESPONSES[fch_getBalance]=""
RESPONSES[fch_getStorageAt]=""
RESPONSES[fch_getAccountNonce]=""
RESPONSES[fch_sendRawTransaction]=""
RESPONSES[fch_getLogs]=""
RESPONSES[fch_getFilterChanges]=""
RESPONSES[fch_newPendingTransactionFilter]=""
RESPONSES[fch_newBlockFilter]=""
RESPONSES[fch_newFilter]=""
RESPONSES[fch_call]=""
RESPONSES[fch_gasPrice]=""
RESPONSES[fch_blockNumber]=""
RESPONSES[net_version]=""
RESPONSES[fch_protocolVersion]=""

### Processes GET requests and stores reponses in RESPONSES ###
function GET_requests() {
	for K in "${!GETDATA[@]}";
	do
		RESPONSES[$K]=$(curl -s --location --request GET "${PORT[GET]}${GETDATA[$K]}" \
	  		--header "Content-Type: application/json" \
			--data "")
	done
}

### Processes POST requests and stores reponses in RESPONSES ###
function POST_requests() {
	for K in "${!POSTDATA[@]}";
	do
		RESPONSES[$K]="$(curl -s --location --request POST "${PORT[POST]}" \
	  		--header "Content-Type: application/json" \
			--data "{\"jsonrpc\":\"2.0\",\"method\":\"${POSTDATA[$K]},\"id\":1}")"
	done
}

function log_API_responses() {
	for K in "${!GETDATA[@]}";
	do
		echo "${yellow}$K"
		echo "${blue}REQUEST:"
		echo "${white}${GETDATA[$K]}"
		echo "${blue}RESPONSE:" ${white}
		echo ${RESPONSES[$K]} #| jq .
		echo
		echo
	done
	for K in "${!POSTDATA[@]}";
	do
		echo "${yellow}$K"
		echo "${blue}REQUEST:"
		echo "${white}${POSTDATA[$K]}"
		echo "${blue}RESPONSE: $white"
		echo ${RESPONSES[$K]} #| jq .
		echo
		echo
	done
}

GET_requests
POST_requests
### BASIC QUERY TESTS ###

function Explorer_getBlock_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET blocks(explorer) test:"
	response_test ${RESPONSES[GET_blocks]}
	if [ "$?"  == "1" ]; then
		BLOCKBYIDHASH=$(echo ${RESPONSES[GET_blocks]} | jq -r .[0].id)
		if [ "$BLOCKBYIDHASH" != "null" ]; then
			if [ "$BLOCKBYIDHASH" == "$TRANSACTION_BLOCK_HASH" ]; then
				TESTS_PASSED=$(( TESTS_PASSED + 1 ))
				echo ${green}BLOCK HASH MATCHES TX${reset}
				echo
				return
			fi
		fi
		echo ${red}BLOCK HASH DOES NOT MATCH TX OR IS NULL${reset}
	fi
	echo
}

#Needs updating - wtf does getTx do - no arguments?
function Explorer_getTx_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET tx(explorer) test:"
	response_test ${RESPONSES[GET_tx]}
	if [ "$?" == "1" ]; then
		TX_HASH=$(echo ${RESPONSES[GET_tx]} | jq -r .id) # fix agrs to jq
		if [ "$TX_HASH" != "null" ]; then
			if [ "$TX_HASH" == "$TX_HASH" ]; then
				TESTS_PASSED=$(( TESTS_PASSED + 1 ))
				echo ${green}BLOCK HASH MATCHES TX${reset}
				echo
				return
			fi
		fi
		echo ${red}BLOCK HASH DOES NOT MATCH TX OR IS NULL${reset}
	fi
	echo
}

function Explorer_getExplorerNodeAdress_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET address(explorer) test:"
	response_test ${RESPONSES[GET_address]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function Explorer_getExplorerNode_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET node-count(explorer) test:"
	response_test ${RESPONSES[GET_node-count]}
	if [ ${RESPONSES[GET_node-count]}="2" ]; then
		echo ${green}SANE VALUE, 2 explorer nodes reported $reset
		TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	else
		echo ${red}non 2 explorer nodes reported $reset
	fi
	echo
}

function Explorer_getShard_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET shard(explorer) test:"
	response_test ${RESPONSES[GET_shard]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function Explorer_getCommitte_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "GET committe(explorer) test:"
	response_test ${RESPONSES[GET_committee]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

### API POST REQUESTS ###

function API_getBlockByNumber_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getBlockByNumber test:"
	response_test ${RESPONSES[fch_getBlockByNumber]}
	BLOCKBYNUMBERHASH=$(echo ${RESPONSES[fch_getBlockByNumber]} | jq -r '.result.hash')

	if [ "$BLOCKBLOCKBYNUMBERHASH" != "null" ]; then
		if [ "$BLOCKBYNUMBERHASH" == "$TRANSACTION_BLOCK_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}BLOCK HASH MATCHES TX${reset}
			echo
			return
		fi
	fi
	echo ${red}BLOCK HASH DOES NOT MATCH TX OR IS NULL${reset}
	echo
}

function API_getBlockByHash_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getBlockByHash test:"
	response_test ${RESPONSES[fch_getBlockByHash]}
	BLOCKBYHASHHASH=$(echo ${RESPONSES[fch_getBlockByHash]} | jq -r '.result.hash')
	if [ "$BLOCKBYHASHBYHASH" != "null" ]; then
		if [ "$BLOCKBYHASHHASH" == "$TRANSACTION_BLOCK_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}BLOCK HASH MATCHES TX${reset}
			echo
			return
		fi
	fi
	echo ${red}BLOCK HASH DOES NOT MATCH TX OR IS NULL${reset}
	echo
}

function API_getBlockTransactionCountByHash_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getBlockTransactionCountByHash test:"
	response_test ${RESPONSES[fch_getBlockTransactionCountByHash]}
	TRANSACTIONCOUNTBYHASH=$(echo ${RESPONSES[fch_getBlockTransactionCountByHash]} | jq -r '.result')
	TRANSACTIONCOUNTBYHASH=$(( TRANSACTIONCOUNTBYHASH ))
	if [ "$TRANSACTIONCOUNTBYHASH" != "null" ]; then
		if [ $TRANSACTIONCOUNTBYHASH -gt 0 ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}NON ZERO TRANSACTION COUNT IN BLOCK${reset}
			echo
			return
		fi
	fi
	echo ${red}INVALID TRANSACTION COUNT IN BLOCK${reset}
	echo
}

function API_getBlockTransactionCountByNumber_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getBlockTransactionCountByNumber test:"
	response_test ${RESPONSES[fch_getBlockTransactionCountByNumber]}
	TRANSACTIONCOUNTBYNUMBER=$(echo ${RESPONSES[fch_getBlockTransactionCountByNumber]} | jq -r '.result')
	TRANSACTIONCOUNTBYNUMBER=$(( TRANSACTIONCOUNTBYNUMBER ))
	if [ "$BLOCKBYHASH" != "null" ]; then
		if [ $TRANSACTIONCOUNTBYNUMBER -gt 0 ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}NON ZERO TRANSACTION COUNT IN BLOCK${reset}
			echo
			return
		fi
	fi
	echo ${red}NON NATURAL TRANSACTION COUNT IN BLOCK${reset}
	echo
}

function API_getCode_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getCode test:"
	response_test ${RESPONSES[fch_getCode]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getTransactionByBlockHashAndIndex_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getTransactionByBlockHashAndIndex test:"
	response_test ${RESPONSES[fch_getTransactionByBlockHashAndIndex]}
	TRANSACTIONHASHBYHASHANDINDEX=$(echo ${RESPONSES[fch_getTransactionByBlockHashAndIndex]} | jq -r '.result.hash')
	if [ "$TRANSACTIONHASHBYHASHANDINDEX" != "null" ]; then
		if [ "$TRANSACTIONHASHBYHASHANDINDEX" == "$TRANSACTION_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}TRANSACTION FROM BLOCKHASH AND INDEX MATCH${reset}
			echo
			return
		fi
	fi
	echo ${red} TRANSACTION FROM BLOCKHASH AND INDEX MATCH${reset}
	echo
}

function API_getTransactionByBlockNumberAndIndex_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getTransactionByBlockNumberAndIndex test:"
	response_test ${RESPONSES[fch_getTransactionByBlockNumberAndIndex]}
	TRANSACTIONHASHBYNUMBERANDINDEX=$(echo ${RESPONSES[fch_getTransactionByBlockNumberAndIndex]} | jq -r '.result.hash')
	if [ "$TRANSACTIONHASHBYNUMBERANDINDEX" != "null" ]; then
		if [ "$TRANSACTIONHASHBYNUMBERANDINDEX" == "$TRANSACTION_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}TRANSACTION FROM BLOCKNUMBER AND INDEX MATCH${reset}
			echo
			return
		fi
	fi
	echo ${red} TRANSACTION FROM BLOCKNUMBER AND INDEX MISMATCH${reset}
	echo
}

function API_getTransactionByHash_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getTransactionByHash test:"
	TX_HASH=$(echo ${RESPONSES[fch_getTransactionByHash]} | jq -r '.result.hash')
	response_test ${RESPONSES[fch_getTransactionByHash]}
	if [ "$TX_HASH" != "null" ]; then
		if [ "$TX_HASH" == "$TRANSACTION_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}TRANSACTION HASH MATCH${reset}
			echo
			return
		fi
	fi
	echo ${red} TRANSACTION HASH MISMATCH${reset}
	echo
}

function API_getTransactionReceipt_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getTransactionReceipt test:"
	TX_HASH=$(echo ${RESPONSES[fch_getTransactionReceipt]} | jq -r '.result.transactionHash')
	response_test ${RESPONSES[fch_getTransactionReceipt]}
	if [ "$TX_HASH" != "null" ]; then
		if [ "$TX_HASH" == "$TRANSACTION_HASH" ]; then
			TESTS_PASSED=$(( TESTS_PASSED + 1 ))
			echo ${green}TRANSACTION HASH MATCH${reset}
			echo
			return
		fi
	fi
	echo ${red} TRANSACTION HASH MISMATCH${reset}
	echo
}

function API_syncing_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_syncing test:"
	response_test ${RESPONSES[fch_syncing]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_netPeerCount_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST net_peerCount test:"
	response_test ${RESPONSES[net_peerCount]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getBalance_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getBalance test:"
	response_test ${RESPONSES[fch_getBalance]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getStorageAt_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getStorageAt test:"
	response_test ${RESPONSES[fch_getStorageAt]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getAccountNonce_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getAccountNonce test:"
	response_test ${RESPONSES[fch_getAccountNonce]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_sendRawTransaction_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_sendRawTransaction test:"
	response_test ${RESPONSES[fch_sendRawTransaction]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getLogs_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getLogs test:"
	response_test ${RESPONSES[fch_getLogs]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_getFilterChanges_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_getFilterChanges test:"
	response_test ${RESPONSES[fch_getFilterChanges]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_newPendingTransactionFilter_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_sendRawTransaction test:"
	response_test ${RESPONSES[fch_newPendingTransactionFilter]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_newBlockFilter_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_newBlockFilter test:"
	response_test ${RESPONSES[fch_newBlockFilter]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_newFilter_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_newFilter test:"
	response_test ${RESPONSES[fch_newFilter]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_call_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_call test:"
	response_test ${RESPONSES[fch_call]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_gasPrice_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_gasPrice test:"
	response_test ${RESPONSES[fch_gasPrice]}
	if [ "$?" == "1" ]; then
		RESULT=$(echo ${RESPONSES[fch_gasPrice]} | jq -r '.result')
		isHexTest $RESULT
		[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	fi
}

function API_blockNumber_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_blockNumber test:"
	response_test ${RESPONSES[fch_blockNumber]}
	if [ "$?" == "1" ]; then
		RESULT=$(echo ${RESPONSES[fch_blockNumber]} | jq -r '.result')
		isHexTest $RESULT
		[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	fi
}

function API_net_version_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST net_version test:"
	response_test ${RESPONSES[net_version]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function API_protocolVersion_test() {
	TESTS_RAN=$(( TESTS_RAN + 1 ))
	echo "POST fch_protocolVersion test:"
	response_test ${RESPONSES[fch_protocolVersion]}
	[ "$?" == "1" ] && TESTS_PASSED=$(( TESTS_PASSED + 1 ))
	echo
}

function run_tests() {
	echo "### TESTING RPC CALLS ###"
	echo
	### Calls to the individual API method test ###
	Explorer_getBlock_test
	Explorer_getTx_test
	Explorer_getExplorerNodeAdress_test
	Explorer_getExplorerNode_test
	Explorer_getShard_test
	Explorer_getCommitte_test
	API_getBlockByNumber_test
	API_getBlockByHash_test
	API_getBlockTransactionCountByHash_test
	API_getBlockTransactionCountByNumber_test
	API_getCode_test
	API_getTransactionByBlockHashAndIndex_test
	API_getTransactionByBlockNumberAndIndex_test
	API_getTransactionByHash_test
	API_getTransactionReceipt_test
	API_syncing_test
	API_netPeerCount_test
	API_getBalance_test
	API_getStorageAt_test
	API_getAccountNonce_test
	API_sendRawTransaction_test
	API_getLogs_test
	API_getFilterChanges_test
	API_newPendingTransactionFilter_test
	API_sendRawTransaction_test
	API_newBlockFilter_test
	API_newFilter_test
	API_call_test
	API_gasPrice_test
	API_blockNumber_test
	API_net_version_test
	API_protocolVersion_test

	TESTS_FAILED=$(( $TESTS_RAN - $TESTS_PASSED ))
	echo -n ${red}
	[ $TESTS_FAILED -eq 0 ] && echo -n ${green}
	echo "PASSED $TESTS_PASSED/$TESTS_RAN: $TESTS_FAILED TESTS FAILED"${reset}
}



if [ "$VERBOSE" == "TRUE" ]; then
	log_API_responses
fi
### BETANET TESTS ###

run_tests
