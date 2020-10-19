#!/bin/bash

# VER="v1.0.8_ff422"
# TZ='Asia/Seoul'
# export TZ
# FFMPEG=/usr/bin/ffmpeg
# UPDIR=/dev/shm
# TMPDIR=/dev/shm
# FTPURL_A="ftp://org1"
# FTPURL_B="ftp://org2"
# FTPURL=${FTPURL_A}
# FTPUSER="11st:Asd098123!"
# JSONDIR="/.json"
# ORGDIR="/.org"
# ERR_FILE=/var/log/11st.log
# LogDir=/home/11st/.log
# Err_Log=/var/log/11st_err.log
# DATE=$(date +%Y-%m-%d" "%H:%M:%S)
# UploadErrStatus=0
# RawDuration=0
# chk_upload_mode=$(curl -q 'http://ovp.myskcdn.net/EncodeMode/index.json' 2>/dev/null | jq '.Upload' 2>/dev/null)

# echo "++++++++++++++++++++++++++" >>/tmp/time.log
# date >>/tmp/time.log

function duration() {
    RawDuration=$(ffprobe -v quiet -print_format json -show_entries stream=duration $1 | jq '.streams[0].duration' | sed -e 's/"//g')
    #echo $RawDuration

    if [ $RawDuration == null ]; then
        duration=$(ffprobe $1 2>&1 | grep Duration)
        tmp=${duration#*: }
        duration=${tmp%%.*}

        #RawDuration=`echo ${duration} | nawk -F: '{seconds=(\$1*60)*60; seconds=seconds+(\$2*60); seconds=seconds+\$3; print seconds}'`
        DuurationToSecs $duration
        RawDuration=$?

        echo $RawDuration

    else
        RawDuration=${RawDuration%%.*}
        tmp=$(ffprobe -v quiet -print_format json -show_entries stream=duration -sexagesimal $1 | jq '.streams[0].duration' | sed -e 's/"//g')
        duration=${tmp%%.*}
    fi

    thb_sec=1
    [ $RawDuration -lt $thb_sec ] && duration="0:00:01"

    echo $RawDuration
    duration="$h:$m:$s"
    duration=$(ffprobe $1 2>&1 | grep Duration)
}

function make_json() {
    # metadata=$(ffprobe -hide_banner -i $file 2>&1)
    metadata=$(ffprobe -v quiet -print_format json -show_format -i $file 2>&1)
    cat <<EOM >$json
{
  "response": {
    "convertCount": "3",
    "runTime": "$duration2",
    "FileName": "$filename",
    "FileSize": "$Ol",
    "orgMetadata": $metadata,
    "md5": "$Om",
    "convertDataList": [
      {
        "originSeq": $seq,
        "runTime": "$duration2",
        "fileSize": "$QAl",
        "filePath": "http://video.11st.co.kr/11stvod/_definst_/movie/item/www/mp4:$QAn",
        "thumbnailPath": "http://snsvideo.11st.co.kr/thb/movie/item/www/$n.jpg",
        "cntinfo": "resolution=$s_mq^fps=29.97^codec_id=libx264^bitrate=2500000^audid=aac",
        "md5": "$QAm",
        "presetName": "QA"
      },
      {
        "runTime": "$duration2",
        "fileSize": "$HQl",
        "filePath": "http://video.11st.co.kr/11stvod/_definst_/movie/item/www/mp4:$HQn",
        "thumbnailPath": "http://snsvideo.11st.co.kr/thb/movie/item/www/$n.jpg",
        "cntinfo": "resolution=$s_hq^fps=29.97^codec_id=libx264^bitrate=2500000^audid=aac",
        "md5": "$HQm",
        "presetName": "HQ"
      },
      {
        "runTime": "$duration2",
        "fileSize": "$MQl",
        "filePath": "http://video.11st.co.kr/11stvod/_definst_/movie/item/www/mp4:$MQn",
        "thumbnailPath": "http://snsvideo.11st.co.kr/thb/movie/item/www/$n.jpg",
        "cntinfo": "resolution=$s_mq^fps=29.97^codec_id=libx264^bitrate=2500000^audid=aac",
        "md5": "$MQm",
        "presetName": "MQ"
      },
    ],
    "rtMsg": "OK",
    "rt": "100"
  }
}
EOM
}

file="/tmp/sample.mp4"

duration $file

if [ ${#duration} -lt 2 ]; then
    sleep 3
    exit 2
fi

#duration2=${duration#*: }
duration2=${duration%%0000}
duration2=${duration2%%000}

echo duration2
