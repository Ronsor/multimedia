#!/bin/sh
ffmpeg -i $@ -vf fps=15  -c:v mjpeg -f image2pipe -huffman optimal -

