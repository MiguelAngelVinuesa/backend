#!/bin/bash
# #######################################################################################################
# certification.sh
# #######################################################################################################
#
# This script performs the steps required for a supervised build.
# The script will pause at critical points where checksums or other output needs to be checked.
# The script generates an evidences directory during the whole process, which includes critical files,
# generated checksums, etc. It also copies itself into the evidences directory.
# The last steps list the evidences directory, zip it, and print the checksum of the zip file.
#
# #######################################################################################################

AGENCY="iTech Labs"
DATE=`date +'%Y-%m-%d'`
EVIDENCE="TopGaming-RNG-${DATE}"
EVIDENCE_PATH="output/${EVIDENCE}"
CHACHA20="${GOROOT}/pkg/mod/github.com/aead/chacha20@v0.0.0-20180709150244-8b13a72661da"
GIT_HASH=`git rev-parse HEAD`
PLUGIN_FILE="libprng.so"
PLUGIN_PATH="${EVIDENCE_PATH}/bin/${PLUGIN_FILE}"

echo ""
echo "============================================================="
echo "supervised build script in collaboration with ${AGENCY}"
echo "============================================================="
echo "running in `pwd`"
echo "date: ${DATE}"
echo "GIT hash: ${GIT_HASH}"
echo ""

echo "============================================================="
echo "step 1: create directory to store required evidences"
echo "mkdir ${EVIDENCE_PATH}"
rm -Rf ${EVIDENCE_PATH}
mkdir -p ${EVIDENCE_PATH}
mkdir ${EVIDENCE_PATH}/cards
mkdir ${EVIDENCE_PATH}/rng
mkdir ${EVIDENCE_PATH}/chacha20
mkdir -p ${EVIDENCE_PATH}/cmd/sharedlib
mkdir ${EVIDENCE_PATH}/bin
echo ""

echo "============================================================="
echo "step 2: copy critical source files"
echo "cards/shuffle.go"
echo "rng/rng.go"
echo "chacha/chacha.go"
echo "cmd/sharedlib/*.go"
echo "go.mod & go.sum"
cp cards/shuffle.go ${EVIDENCE_PATH}/cards/
cp rng/rng.go ${EVIDENCE_PATH}/rng/
cp "${CHACHA20}/chacha/chacha.go" ${EVIDENCE_PATH}/chacha20/
cp cmd/sharedlib/main.go ${EVIDENCE_PATH}/cmd/sharedlib/
cp cmd/sharedlib/rng.go ${EVIDENCE_PATH}/cmd/sharedlib/
cp cmd/sharedlib/shuffle.go ${EVIDENCE_PATH}/cmd/sharedlib/
cp go.mod ${EVIDENCE_PATH}/
cp go.sum ${EVIDENCE_PATH}/
echo ""

echo "============================================================="
echo "step 3 - sha1 checksums of critical source files"
cd ${EVIDENCE_PATH}
find -type f -exec sha1sum '{}' \; | grep -v "sha1*" > sha1_sources.txt
cat sha1_sources.txt
cd - > /dev/null
echo ""

echo "============================================================="
echo "step 4 - ${AGENCY}: verify sha1 checksums..."
echo ""
read -n 1 -s -r -p "<Press any key to continue>"
echo ""

echo "============================================================="
echo "step 5 - build binaries"
echo "GIT hash (build version): ${GIT_HASH}"
echo ""
echo "go build -buildmode=c-shared -ldflags="-s -w -X main.hash=${GIT_HASH}" -o ${PLUGIN_PATH} cmd/sharedlib/*.go"
go version
go build -buildmode=c-shared -ldflags="-s -w -X main.hash=${GIT_HASH}" -o ${PLUGIN_PATH} cmd/sharedlib/*.go
echo ""

echo "============================================================="
echo "step 6 - sha1 checksums of binaries"
cd ${EVIDENCE_PATH}
find -type f -exec sha1sum '{}' \; | grep "libprng" | grep -v "sha1*" > sha1_binaries.txt
cat sha1_binaries.txt
cd - > /dev/null
echo ""
read -n 1 -s -r -p "<Press any key to continue>"
echo ""

echo "============================================================="
echo "step 7 - copy build process script"
echo "tools/certification.sh"
cp tools/certification.sh ${EVIDENCE_PATH}/
echo ""

echo "============================================================="
echo "step 8 - list evidences directory"
cd ${EVIDENCE_PATH}
tree
cd - > /dev/null
echo ""
read -n 1 -s -r -p "<Press any key to continue>"
echo ""

echo "============================================================="
echo "step 9 - zip evidences directory with password"
cd output
rm -Rf ${EVIDENCE}.zip
zip -re ${EVIDENCE}.zip ${EVIDENCE}
cd - > /dev/null
echo ""

echo "============================================================="
echo "step 10 - sha1 checksum of evidences zip"
sha1sum output/${EVIDENCE}.zip
echo ""

echo "============================================================="
echo "finished"
echo "============================================================="
echo ""
