#!/usr/bin/env bash
#   Use this script to test if a given TCP host/port are available

WAITFORIT_cmdname=${0##*/}
WAITFORIT_usage()
{
    cat << USAGE >&2
Usage:
    $WAITFORIT_cmdname host:port [-s] [-t timeout] [-- command args]
    -h HOST | --host=HOST       Host or IP under test
    -p PORT | --port=PORT       TCP port under test
                                Alternatively, you specify the host and port as host:port
    -s | --strict               Only execute subcommand if the test succeeds
    -q | --quiet                Don't output any status messages
    -t TIMEOUT | --timeout=TIMEOUT
                                Timeout in seconds, zero for no timeout
    -- COMMAND ARGS             Execute command with args after the test finishes
USAGE
    exit 1
}

WAITFORIT_wait_for()
{
    if [ "$WAITFORIT_TIMEOUT" -gt 0 ]; then
        echo "$WAITFORIT_cmdname: waiting $WAITFORIT_TIMEOUT seconds for $WAITFORIT_HOST:$WAITFORIT_PORT"
    else
        echo "$WAITFORIT_cmdname: waiting for $WAITFORIT_HOST:$WAITFORIT_PORT without a timeout"
    fi
    WAITFORIT_start_ts=$(date +%s)
    while :
    do
        if [ "$WAITFORIT_ISBUSY" = "true" ]; then
            nc -z $WAITFORIT_HOST $WAITFORIT_PORT
            WAITFORIT_result=$?
        else
            (echo > /dev/tcp/$WAITFORIT_HOST/$WAITFORIT_PORT) >/dev/null 2>&1
            WAITFORIT_result=$?
        fi
        if [ $WAITFORIT_result -eq 0 ]; then
            WAITFORIT_end_ts=$(date +%s)
            echo "$WAITFORIT_cmdname: $WAITFORIT_HOST:$WAITFORIT_PORT is available after $((WAITFORIT_end_ts-WAITFORIT_start_ts)) seconds"
            break
        fi
        sleep 1
    done
    return $WAITFORIT_result
}

WAITFORIT_hostport=$(echo $1 | sed -e 's|^tcp://||g')
WAITFORIT_hostport_array=(${WAITFORIT_hostport//:/ })
WAITFORIT_HOST=${WAITFORIT_hostport_array[0]}
WAITFORIT_PORT=${WAITFORIT_hostport_array[1]}
WAITFORIT_TIMEOUT=15
WAITFORIT_STRICT=false
WAITFORIT_QUIET=false
WAITFORIT_ISBUSY=false

# process arguments
shift
while [ $# -gt 0 ]
do
    case "$1" in
        *:* )
        WAITFORIT_hostport=$1
        WAITFORIT_hostport_array=(${WAITFORIT_hostport//:/ })
        WAITFORIT_HOST=${WAITFORIT_hostport_array[0]}
        WAITFORIT_PORT=${WAITFORIT_hostport_array[1]}
        shift 1
        ;;
        -q | --quiet)
        WAITFORIT_QUIET=true
        shift 1
        ;;
        -s | --strict)
        WAITFORIT_STRICT=true
        shift 1
        ;;
        -h)
        WAITFORIT_HOST=$2
        if [[ $WAITFORIT_HOST == "" ]]; then break; fi
        shift 2
        ;;
        --host=*)
        WAITFORIT_HOST=$(echo $1 | sed -e 's/^[^=]*=//g')
        shift 1
        ;;
        -p)
        WAITFORIT_PORT=$2
        if [[ $WAITFORIT_PORT == "" ]]; then break; fi
        shift 2
        ;;
        --port=*)
        WAITFORIT_PORT=$(echo $1 | sed -e 's/^[^=]*=//g')
        shift 1
        ;;
        -t)
        WAITFORIT_TIMEOUT=$2
        if [[ $WAITFORIT_TIMEOUT == "" ]]; then break; fi
        shift 2
        ;;
        --timeout=*)
        WAITFORIT_TIMEOUT=$(echo $1 | sed -e 's/^[^=]*=//g')
        shift 1
        ;;
        --)
        shift
        WAITFORIT_CLI=("$@")
        break
        ;;
        *)
        WAITFORIT_usage
        ;;
    esac
done

if [[ "$WAITFORIT_HOST" == "" || "$WAITFORIT_PORT" == "" ]]; then
    echo "Error: you need to provide a host and port to test."
    WAITFORIT_usage
fi

WAITFORIT_TIMEOUT=$(echo $WAITFORIT_TIMEOUT | grep -E -o '^[0-9]+$')
if [[ "$WAITFORIT_TIMEOUT" == "" ]]; then
    echo "Error: invalid timeout."
    WAITFORIT_usage
fi

if [ "$WAITFORIT_QUIET" = "true" ]; then
    exec 3>&1
    exec 1>/dev/null
    exec 2>/dev/null
else
    exec 3>&1
fi

if ! command -v nc >/dev/null 2>&1; then
    echo "wait-for-it: command not found: nc"
    WAITFORIT_ISBUSY=true
fi

WAITFORIT_wait_for
WAITFORIT_result=$?

exec 1>&3

if [[ $WAITFORIT_CLI != "" ]]; then
    if [ $WAITFORIT_result -ne 0 ] && [ "$WAITFORIT_STRICT" = "true" ]; then
        echo "$WAITFORIT_cmdname: strict mode, refusing to execute subprocess"
        exit $WAITFORIT_result
    fi
    exec "${WAITFORIT_CLI[@]}"
else
    exit $WAITFORIT_result
fi
