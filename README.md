# FastHttpPy
Python implementation of Golang FastHttp

![FastHttpPy](https://i.imgur.com/3SL6u9m.png)

63k requests/second
(for comparisons the fastest Python server (uvicorn) does 8k/s without multiprocess)

Benchmark made on MacBook Pro M2 Pro

```sh
pip install fasthttppy
```

Example:
```python
from fasthttppy import FastHttp, MakeResponse, Request

app = FastHttp()

app.Static("/static","static")

def my_callback(req:Request):
    print(req.contents.GetURI())
    print(req.contents.GetHeaders())
    print(req.contents.GetData())
    return MakeResponse(200,"Hello world")


app.Get("/py",my_callback)
app.Post("/py",my_callback)

app.Run()
```

Methods:
```python
class FastHttp:

    def Static(self,path:str,dir:str)
    #=-=-=Serve static folder=-=-=
    #path: url path
    #dir: folder
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def Get(self,path:str,function:Callable)
    #=-=-Bind path to function-=-=
    #path: url path
    #function: function
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def Post(self,path:str,function:Callable)
    #=-=-Bind path to function-=-=
    #path: url path
    #function: function
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def SetErrorPage(self,content:str)
    #=-=-=-=Set error page-=-=-=-=
    #content: html content
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def SetNotFoundPage(self,path:str,dir:str)
    #=-=-=-=-Set 404 page=-=-=-=-=
    #content: html content
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def Run(self,host:str,port:int)
    #=-=-=-=Start Go server=-=-=-=
    #host: server host
    #port: server port
    #=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

```

It is still possible to improve performance a lot, the golang/py communication is done through a *C.char Json that can be replaced by structs (removing the json encode decode step)