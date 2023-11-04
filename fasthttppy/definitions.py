import ctypes, json


class Request(ctypes.Structure):
    _fields_ = [("Data", ctypes.c_char_p)]
    
    def GetURI(self):
        return json.loads(self.Data.decode()).get("URI")
    
    def GetHeaders(self):
        return json.loads(self.Data.decode()).get("Headers")

    def GetBody(self):
        return json.loads(self.Data.decode()).get("Body")


callback = ctypes.CFUNCTYPE(ctypes.c_void_p, (Request))

