#!/bin/sh

(
  cd /work/public && python3 -m http.server 8080 -b 127.0.0.1 &
)

exec chromium \
    --no-sandbox \
    --no-default-browser-check \
    --no-first-run \
    --disable-default-apps \
    --disable-popup-blocking \
    --disable-translate \
    --disable-background-timer-throttling \
    --headless \
    --disable-gpu \
    --print-to-pdf="pdfs/${1}.pdf" \
    "http://127.0.0.1:8080/${1}/_print/"
