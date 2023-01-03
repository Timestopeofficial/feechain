#!/usr/bin/env bash
set -e

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
DATA="$DIR/data"
LOGS="$DATA/logs"
DATA_NAME="${DATA_NAME:=feechain_sharddb_0}"

case "$NETWORK" in
asadal)
  CONFIG_PATH="-c /root/feechain-asadal.conf"
  ;;
testnet)
  CONFIG_PATH="-c /root/feechain-pstn.conf"
  ;;
*)
  echo "unknown network"
  exit 1
  ;;
esac

if [ "$MODE" = "offline" ]; then
  BASE_ARGS=(--datadir "$DATA" --log.dir "$LOGS" --run.offline)
else
  BASE_ARGS=(--datadir "$DATA" --log.dir "$LOGS")
fi

mkdir -p "$LOGS"
echo -e NODE ARGS: \" $CONFIG_PATH "$@" "${BASE_ARGS[@]}" \"
echo "NODE VERSION: $($DIR/feechain --version)"

"$DIR/feechain" $CONFIG_PATH "$@" "${BASE_ARGS[@]}"
