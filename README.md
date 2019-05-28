# webtracereplay

This project provides simple tools to replay http request traces to evaluate the performance of caching systems and webservers.

There are three components:

 - the client: which reads in the trace file and generates valid http requests
 - the server which you are trying to benchmark, e.g., nginx, Varnish, Apache TS, etc.
 - the origin: which emulates a database or storage server

These tools were built for the [AdaptSize project](https://github.com/dasebe/AdaptSize), see [References](#references) for more information.

## Install webtracereplay

Please install

 - development versions of these three libraries: curl, fcgi, and fcgi++
 - spawn-fcgi
 - gcc with a version > 5.0, ideally 6.0+

There are Makefiles for the client and origin:

    cd client
    make
    cd ..

    cd origin
    make


Also compile and install a webserver/caching system, such as nginx. The example config in the server subfolder assumes we are using nginx. You can either point nginx to use the configuration file server/nginx.conf or copy the directory to /usr/share/nginx .


## Request traces

We will need request traces for the client (to request objects) and the origin (to emulate objects).

### Client request trace

Request traces must be given in a space-separated format with three colums
- time is an unisgned int; not used currently, but can be used to schedule the replay
- url/key should be a long long int, used to uniquely identify web objects

| time |  id |
| ---- | --- |
|   1  |  1  |
|   2  |  2  |
|   3  |  1  |
|   4  |  3  |
|   4  |  1 |


### Origin request trace

The origin request trace is simply a mapping of ids/urls to object sizes.

Request traces must be given in a space-separated format with two colums
- url/key should be a long long int, used to uniquely identify web objects
- size should be a long long int, this is the object's size in bytes

|  id | size |
| --- | ---- |
|  1  |  120 |
|  2  |   64 |
|  1  |  120 |
|  3  |  14  |
|  1 |  120 |


## Run an experiment with webtracereplay

This example assumes you have nginx installed and everything set up.
The client request trace is called "client.tr" and the origin request trace is called "origin.tr".

We will need three VMs, or different terminals/screens on the same box.
Start these one after another.

#### Start nginx

    sudo nginx -c server/nginx.conf

#### Start the origin

    spawn-fcgi -a 127.0.0.1 -p 9000 -n origin/origin origin.tr

#### Start trace replaying and write throughput and delay measurements to through.log and histogram.log respectively

     client/client client.tr 20 127.0.0.1:7000 throughput.log histogram.log



## Contributors are welcome

Want to contribute? Great! We follow the [Github contribution work flow](https://help.github.com/articles/github-flow/).
This means that submissions should fork and use a Github pull requests to get merged into this code base.

If you come across a bug in webcachesim, please file a bug report by [creating a new issue](https://github.com/dasebe/webcachesim/issues/new).

## References

We ask academic works, which built on this code, to reference the AdaptSize paper:

    AdaptSize: Orchestrating the Hot Object Memory Cache in a CDN
    Daniel S. Berger, Ramesh K. Sitaraman, Mor Harchol-Balter
    To appear in USENIX NSDI in March 2017.
    
You can find more information on [USENIX NSDI 2017 here.](https://www.usenix.org/conference/nsdi17/technical-sessions)
