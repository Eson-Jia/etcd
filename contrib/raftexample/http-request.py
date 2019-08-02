import http.client
import requests
import urllib

index = 0


def Put(index):
    resp = requests.put('http://localhost:12380/foo', data='bar{0}'.format(index))
    print(resp)



def AddPeer():
    # curl -L http://127.0.0.1:12380/4 -XPOST -d http://127.0.0.1:42379
    print(requests.post('http://127.0.0.1:12380/4',data='http://127.0.0.1:42379'))

# for i in range(0,1000):
#     Put(i)

for i in range(0,20):
    AddPeer()

