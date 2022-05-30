# webclipClient

A command line client connect to webclip2 tool. webclip2 tool is web service the provide one time string retriving functionality. 
User can store a string to server with by design one time retriving mechanism. [[WebClip2](https://webclip2.mytools.express)].

# Installation 

```shell
go install github.com/xh-dev-go/webclipClient@latest
```

# Usage

Store String
```shell
webclipClient -post -from-clipboard # copy string clipboard from and store one the webclip server
# --------------- output
https://webclip2.mystools.express/#/get?id=123456
```

Retrieve String
```shell
webclipClient -get -code 123456 -to-clipboard # retrive code 123456 and copy the result to clipboard
```
