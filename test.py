from fasthttppy import FastHttp, MakeResponse, Request
import ctypes,random

app = FastHttp()




def my_callback(req: Request):
    print(req.GetURI())
    print(req.GetHeaders())
    print(req.GetBody())
    
    return MakeResponse(200,"asd")


app.Get("/py",my_callback)
app.Post("/py",my_callback)


app.Run()