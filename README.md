# posts-example

Example backend for basic posting app

Has / will have endpoints:
- Create Post 
- Get Post
- List Posts (with pagination)
- Update Post
- Delete Post


## TODO
- proper pagination - probably timestamp of first request, offset, all the query parameters? If we put all this stuff in, token needs to be secure / encrypted

## Run jaeger
From tutorial https://www.scalyr.com/blog/jaeger-tracing-tutorial/
```
docker run -d --name jaeger -p 16686:16686 -p 6831:6831/udp jaegertracing/all-in-one:1.9
```

go to localhost:16686