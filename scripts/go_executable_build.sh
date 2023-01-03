#!/usr/bin/env bash

export GO111MODULE=on

declare -A SRC
SRC[feechain]=./cmd/feechain
SRC[bootnode]=./cmd/bootnode

BINDIR=bin
REL=
GOOS=linux
GOARCH=$(uname -m)
GOARCH=${GOARCH/x86_64/amd64}
GOARCH=${GOARCH/aarch64/arm64}
RACE=
TRACEPTR=
VERBOSE=
if [ "$(uname -s)" == "Darwin" ]; then
   GO_GCFLAGS=""
else
   GO_GCFLAGS="all=-c 2"
fi
DEBUG=false
STATIC=true

unset -v progdir
case "${0}" in
*/*) progdir="${0%/*}";;
*) progdir=.;;
esac

. "${progdir}/setup_bls_build_flags.sh"

declare -A LIB

if [ "$(uname -s)" == "Darwin" ]; then
   MD5='md5 -r'
   GOOS=darwin
   LIB[libbls384_256.dylib]=${BLS_DIR}/lib/libbls384_256.dylib
   LIB[libmcl.dylib]=${MCL_DIR}/lib/libmcl.dylib
   LIB[libgmp.10.dylib]=/usr/local/opt/gmp/lib/libgmp.10.dylib
   LIB[libgmpxx.4.dylib]=/usr/local/opt/gmp/lib/libgmpxx.4.dylib
   LIB[libcrypto.1.1.dylib]=/usr/local/opt/openssl/lib/libcrypto.1.1.dylib
else
   MD5=md5sum
   LIB[libbls384_256.so]=${BLS_DIR}/lib/libbls384_256.so
   LIB[libmcl.so]=${MCL_DIR}/lib/libmcl.so
fi

function usage
{
   ME=$(basename $0)
   cat<<EOF

Usage: $ME [OPTIONS] ACTION

OPTIONS:
   -h             print this help message
   -p profile     aws profile name
   -a arch        set build arch (default: $GOARCH)
   -o os          set build OS (default: $GOOS, windows is supported)
   -r             enable -race build option (default: $RACE)
   -t             full analysis on {pointer} build option (default: $TRACEPTR)
   -v             verbose build process (default: $VERBOSE)
   -s             build static linux executable (default: $STATIC)
   -S             build non-static linux executable


ACTION:
   build       build binaries only (default action)

   feechain|bootnode|
               only build the specified binary

EXAMPLES:

# build linux binaries only by default
   $ME

# build windows binaries
   $ME -o windows

EOF
   exit 1
}

function build_only
{
   if [[ "$STATIC" == "true" && "$GOOS" == "darwin" ]]; then
      echo "static build only supported on Linux platform"
      exit 2
   fi

   VERSION=$(git rev-list --count HEAD)
   COMMIT=$(git describe --always --long --dirty)
   BUILTAT=$(date +%FT%T%z)
   BUILTBY=${USER}@
   local build=$1

   set_gcflags
   set -e

   for bin in "${!SRC[@]}"; do
      if [[ -z "$build" || "$bin" == "$build" ]]; then
         rm -f $BINDIR/$bin
         echo "building ${SRC[$bin]}"
         if [ "$DEBUG" == "true" ]; then
            env GOOS=$GOOS GOARCH=$GOARCH go build $VERBOSE -gcflags="${GO_GCFLAGS}" -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}" -o $BINDIR/$bin $RACE $TRACEPTR ${SRC[$bin]}
         else
            if [ "$STATIC" == "true" ]; then
               env GOOS=$GOOS GOARCH=$GOARCH go build $VERBOSE -gcflags="${GO_GCFLAGS}" -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}  -w -extldflags \"-static -lm\"" -o $BINDIR/$bin $RACE $TRACEPTR ${SRC[$bin]}
            else
               env GOOS=$GOOS GOARCH=$GOARCH go build $VERBOSE -gcflags="${GO_GCFLAGS}" -ldflags="-X main.version=v${VERSION} -X main.commit=${COMMIT} -X main.builtAt=${BUILTAT} -X main.builtBy=${BUILTBY}" -o $BINDIR/$bin $RACE $TRACEPTR ${SRC[$bin]}
            fi
         fi
         if [ "$(uname -s)" == "Linux" ]; then
            if [ $bin = feechain ]; then
               $BINDIR/$bin version || $BINDIR/$bin version
            else
               $BINDIR/$bin --version || $BINDIR/$bin version
            fi
         fi
         if [ "$(uname -s)" == "Darwin" -a "$GOOS" == "darwin" -a -e $BINDIR/$bin ]; then
            if [ $bin = feechain ]; then
               $BINDIR/$bin version || $BINDIR/$bin version
            else
               $BINDIR/$bin --version || $BINDIR/$bin version
            fi
         fi
      fi
   done

   pushd $BINDIR
   if [ "$STATIC" == "true" ]; then
      $MD5 "${!SRC[@]}" > md5sum.txt
      cp -pf ../scripts/node.sh .
      $MD5 node.sh >> md5sum.txt
   else
      for lib in "${!LIB[@]}"; do
         if [ -e ${LIB[$lib]} ]; then
            cp -pf ${LIB[$lib]} .
         fi
      done

      $MD5 "${!SRC[@]}" "${!LIB[@]}" > md5sum.txt
      # hardcode the prebuilt libcrypto to md5sum.txt
      if [ "$(uname -s)" == "Linux" ]; then
         echo '771150db04267126823190c873a96e48  libcrypto.so.10' >> md5sum.txt
      fi
   fi
   popd
}

function set_gcflags
{
   if [[ ! -z "$RACE" ]]; then
      if [ "$DEBUG" == "true" ]; then
         GO_GCFLAGS="all=-N -l"
      else
         GO_GCFLAGS=""
      fi
   else
      if [ "$DEBUG" == "true" ]; then
         GO_GCFLAGS="all=-N -l -c 2"
      fi
   fi
}

################################ MAIN FUNCTION ##############################
while getopts "hp:a:o:b:f:rtvsdS" option; do
   case $option in
      h) usage ;;
      p) PROFILE=$OPTARG ;;
      a) GOARCH=$OPTARG ;;
      o) GOOS=$OPTARG ;;
      r) RACE=-race ;;
      t) TRACEPTR='-gcflags=all=-d=checkptr' ;;
      v) VERBOSE='-v -x' ;;
      d) DEBUG=true ;;
      s) STATIC=true ;;
      S) STATIC=false ;;
   esac
done

mkdir -p $BINDIR

shift $(($OPTIND-1))

ACTION=${1:-build}

case "$ACTION" in
   "build") build_only ;;
   "feechain"|"bootnode") build_only $ACTION ;;
   *) usage ;;
esac
