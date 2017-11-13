#!/bin/sh
ffmpeg -i $@ -f u8 -ac 1 -ar 8000 -acodec pcm_u8 -
