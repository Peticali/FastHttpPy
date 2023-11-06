import ctypes, json
from .definitions import Request, callback


class FastHttp:
    
    def __init__(self):
        self.lib = ctypes.CDLL("dist/fasthttppy.lib")
        
        self.lib.StartServer.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        self.lib.MountStatic.argtypes = [ctypes.c_char_p, ctypes.c_char_p]
        self.lib.SetNotFoundPage.argtypes = [ctypes.c_char_p]
        self.lib.SetErrorPage.argtypes = [ctypes.c_char_p]
        self.lib.RegisterCallback.argtypes = [ctypes.c_char_p, ctypes.POINTER(callback),ctypes.c_char_p]
        
        self.func_list = []
    
    def Static(self,path:str, dir:str):
        self.lib.MountStatic(path.encode(),dir.encode())

    def Get(self,path:str,functionnn):
        c_callback = callback(functionnn)
        self.func_list.append(c_callback)
        
        self.lib.RegisterCallback(path.encode(), c_callback,b"GET")
        
    def Post(self,path:str,functionnn):
        c_callback = callback(functionnn)
        self.func_list.append(c_callback)
        
        self.lib.RegisterCallback(path.encode(), c_callback,b"POST")
    
    def Run(self,host="",port=80):
        self.lib.StartServer(host.encode(),str(port).encode())
    
    def SetErrorPage(self,content):
        self.lib.SetErrorPage(content.encode())
        
    def SetNotFoundPage(self,content):
        self.lib.SetNotFoundPage(content.encode())
        

class Cookie:
    def __init__(self,name:str,expiration:int,value:str):
        self.name = name
        self.expiration = expiration
        self.value = value


#Not used bc python GC keeps deleting the self str
class ___response:
    def __init__(self, status_code=200,content="",headers={}):
        # self.obj = {}
        
        # self.obj["status_code"] = status_code    
        # self.obj["headers"] = headers
        # self.obj["content"] = content
        # self.obj["cookies"] = []
        
        self.obj = json.dumps({"status_code":status_code,"content":content,"headers":headers,"cookies":[]})
        

    def set_cookie(self,cookie:Cookie):
        cookies = self.obj["cookies"]
        cookies.append({"name":cookie.name,"exp":cookie.expiration,"value":cookie.value})
        self.obj["cookies"] = cookies
    
    def make(self):
        self.final = ctypes.cast(ctypes.c_char_p(self.obj.encode()), ctypes.c_void_p).value
        
        return self.final

def MakeResponse(status_code=200,content="",headers={},cookies=[]):
    
    obj = json.dumps({"status_code":status_code,"content":content,"headers":headers,"cookies":cookies})
    
    return ctypes.cast(ctypes.c_char_p(obj.encode()), ctypes.c_void_p).value
    

