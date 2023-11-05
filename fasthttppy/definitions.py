import ctypes, json


class Request(ctypes.Structure):
    _fields_ = [("Data", ctypes.c_char_p),
                ("Body_Address", ctypes.POINTER(ctypes.c_ubyte)),
                ("Body_Len", ctypes.c_size_t),
                ]
    
    def GetURI(self):
        return json.loads(self.Data.decode()).get("URI")
    
    def GetData(self):
        if int(self.Body_Len) != 0:
            return ctypes.string_at(self.Body_Address,self.Body_Len)
        else:
            return None
    
    def GetHeaders(self):
        return json.loads(self.Data.decode()).get("Headers")




callback = ctypes.CFUNCTYPE(ctypes.c_void_p, ctypes.POINTER(Request))

