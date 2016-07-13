gogrok
======

**NOTE:** This readme uses the word "should" a lot, because as of right now, gogrok does
nothing that I intend it to do. I hope to flesh this out at GopherCon's hack day. Right now,
all this does is log http requests to a file.

----------------------------------------

Gogrok is basically a proxy server that records HTTP traffic as it is sent from your computer.

Sometimes, you don't have a web browser's inspector to show you your outgoing and incoming
HTTP requests. Gogrok should show them to you. Also, I noticed one day that large chunked
HTTP requests make the ngrok inspector's browser tab crash. It would be great to see the
chunked request as it comes in. Gogrok should do this.
