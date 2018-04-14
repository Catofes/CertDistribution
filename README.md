## Cert Distribution Web

#### Build
```
glide i
make
```
or
```
docker build .
```


#### Add a cert
```
curl -X PUT https://example.com -F "Cert=@0003_chain.pem"
```
A uuid will return.

#### Show a cert
```
curl https://example.com/46ae2b8c-24ef-4c28-961e-e1abd9e31c55 | jq
```

#### Show a raw cert
```
curl https://example.com/46ae2b8c-24ef-4c28-961e-e1abd9e31c55/raw
```

#### Update a cert
```
curl -X POST https://example.com/46ae2b8c-24ef-4c28-961e-e1abd9e31c55 -F "Cert=@0003_chain.pem"
```