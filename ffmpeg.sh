#!/usr/bin/env bash

ffmpeg -f v4l2 -framerate 25 -video_size 640x480 -i /dev/video0 \
    -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp rtp://127.0.0.1:5004>local.sdp
