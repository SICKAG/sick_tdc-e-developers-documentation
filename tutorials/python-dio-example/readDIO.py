import requests
from requests.structures import CaseInsensitiveDict
import json
import time

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
    password = "servicelevel"
    myobj = {'password':password}
    headers = CaseInsensitiveDict()
    # setting header
    headers["Content-Type"] ="application/x-www-form-urlencoded"
    resp = requests.post(tokenUrl, data = myobj, headers=headers)
    token = json.loads(resp.text)
    return token["token"]

# fetches DIO state
def fetchCurrState(token):
    resp = getDio(token, "http://192.168.0.100:59801/tdce/dio/GetState/DIO_B")
    return resp['Value']

# gets dio_b object from url 
def getDio(token, url):
    api_response = requests.get(url, auth=BearerAuth(token))
    data = api_response.text
    resp = json.loads(data)
    return resp


token = getToken()
while(1):
    print(fetchCurrState(token))
    time.sleep(5)