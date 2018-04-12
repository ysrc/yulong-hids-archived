#!/bin/sh

DIRNAME=`dirname $0`
if [ -z $IDS_ELASTICSEARCH ] || [ -z $IDS_MONGODB ]; then
    echo "You need to specify IDS_MONGODB and IDS_ELASTICSEARCH"
    exit 1
fi

for i in $(seq 1 10);
do
    echo "Wait for ElasticSearch Start...$i"
    curl -s http://$IDS_ELASTICSEARCH> /dev/null
    if [ $? = 0 ];then
        break;
    fi
    sleep 5
    if [ $i -ge 10 ]; then
        echo "Too much retry, exit."
        exit 2
    fi
done

$DIRNAME/server -db $IDS_MONGODB -es $IDS_ELASTICSEARCH