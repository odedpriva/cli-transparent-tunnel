#!/usr/bin/env bash

set -x

ROOT_DIR=$(git rev-parse --show-toplevel)

# general variables.
ESC="\x1b["
RESET=$ESC"39;49;00m"
RED=$ESC"31;01m"
GREEN=$ESC"32;01m"
YELLOW=$ESC"33;01m"
MAGENTA=$ESC"35;01m"
LOG_ERROR_PREFIX="ERROR: "
LOG_INFO_PREFIX="INFO: "
LOG_DEBUG_PREFIX="DEBUG: "
LOG_FATAL_PREFIX="FATAL: "

export TOOL_NAME=ctt
export DOCKER='docker'

export COMMANDS_DIR=${ROOT_DIR}/command-tunneler/commands
export TEST_DIR=${ROOT_DIR}/for-tests
export TEST_CONFIG_DIR=${TEST_DIR}/generated-config/
mkdir -p "${TEST_CONFIG_DIR}"
export BIN_DIR=${ROOT_DIR}/bin
mkdir -p "${BIN_DIR}"

export SSH_SERVER_CONTAINER_NAME=openssh-server
export SSH_SERVER_PORT=2222
export SSH_SERVER_USER=linuxserver.io
export SSH_PRIVATE_KEY_PATH=${ROOT_DIR}/for-tests/scripts/id_rsa
export SSH_PUBLIC_KEY=${ROOT_DIR}/for-tests/scripts/id_rsa.pub

export GO_TEMPLATE=${BIN_DIR}/gomplate

declare -a commands_list
commands_list=(kubectl)

OS=""
case $(uname -s) in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
esac

ARCH=""
case $(uname -m) in
    x86_64) ARCH="amd64" ;;
esac

if [[ "$OS" == 'linux' ]]; then
   LOCAL_IP=$(ifconfig eth0  | grep inet | cut -d: -f2 | awk '{print $2}')
   export LOCAL_IP
elif [[ "$OS" == 'darwin' ]]; then
   LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | cut -d\  -f2 | head -1)
   export LOCAL_IP
fi


function LOG() {
	local LOG_LEVEL_PREFIX=""

	case $1 in
	ERROR) LOG_LEVEL_PREFIX=$YELLOW$LOG_ERROR_PREFIX$RESET ;;
	INFO) LOG_LEVEL_PREFIX=$GREEN$LOG_INFO_PREFIX$RESET ;;
	DEBUG) LOG_LEVEL_PREFIX=$MAGENTA$LOG_DEBUG_PREFIX$RESET ;;
	FATAL) LOG_LEVEL_PREFIX=$RED$LOG_FATAL_PREFIX$RESET ;;
	*) LOG_LEVEL_PREFIX=${MAGENTA}UNKNOWN:$RESET ;;
	esac

	# Return if attempting to print a debug message without debug mode enabled
	test "$1" == "DEBUG" && [ "$DEBUG_MODE" != "TRUE" ] && return 0

	echo "${@:2}" | while read LINE; do
		printf "$(date "+%F %H:%M:%S") ${LOG_LEVEL_PREFIX} %s\n" "${LINE}" | tee -a ${LOGFILE}
	done

	if [ "$1" == "FATAL" ]; then
		exit 1
	fi
}

function GO_GET() {
  src=$1
  file_name=${2:-non-exist}
#  set +e
  [[ -f "${BIN_DIR}/${file_name}" ]] && return
#  set -e
  (TMP_DIR=$(mktemp -d) ;\
  cd "$TMP_DIR" || exit 1 ;\
  LOG INFO "Downloading ${src} to ${BIN_DIR}" ;\
  GOBIN=${BIN_DIR} go install "${src}" ;)
}

function CURL_GET() {
    src=${1}
    file_name=${2}
    [[ -f "${BIN_DIR}/${file_name}" ]] && return
    curl -L "${src}" -o "${BIN_DIR}/${file_name}"
    chmod +x "${BIN_DIR}/${file_name}"
}