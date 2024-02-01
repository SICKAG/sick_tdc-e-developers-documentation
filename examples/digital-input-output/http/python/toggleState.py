import requests
from requests.structures import CaseInsensitiveDict
import json
import sys

# authentication for websockets
class BearerAuth(requests.auth.AuthBase):
    def __init__(self,token):
        self.token = token
    def __call__(self,r):
        # on call set token
        r.headers["authorization"] = "Bearer " + self.token
        return r
    
# gets authentication token for http requests  
def getToken():
    tokenUrl = "http://192.168.0.100:59801/user/Service/token"
    # set real password here
    password = "PASSWORD"
    myobj = {'password':password}
    headers = CaseInsensitiveDict()
    # setting header
    headers["Content-Type"] ="application/x-www-form-urlencoded"
    resp = requests.post(tokenUrl, data = myobj, headers=headers)
    token = json.loads(resp.text)
    return token["token"]

def setDios(token, js):
    #set LED state
    url = "http://192.168.0.100:59801/tdce/dio/SetStates"
    headers = CaseInsensitiveDict()
    headers["Content-Type"] = "application/json"
    # turn LED on; data set to js
    requests.post(url, data=js, auth=BearerAuth(token), headers=headers)

if __name__ == "__main__":
    payload = str(sys.argv[1])

    if(payload=="1"):
        js = '[{ "DioName": "DIO_A", "Value": 1, "Direction": "Output" }]'
    elif(payload=="0"):
        js = '[{ "DioName": "DIO_A", "Value": 0, "Direction": "Output" }]'
    else:
        print("Wrong value for DIO state.")
    token = getToken()
    setDios(token, js)
