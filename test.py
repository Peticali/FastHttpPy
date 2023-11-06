from fasthttppy import FastHttp, MakeResponse, Request
import ctypes,random

app = FastHttp()

app.Static("/static","dist")

def my_callback(req:Request):
    print(req.contents.GetURI())
    print(req.contents.GetHeaders())
    print(req.contents.GetData())
    
    return MakeResponse(200,"asd")


app.Get("/py",my_callback)
app.Post("/py",my_callback)


app.Run()