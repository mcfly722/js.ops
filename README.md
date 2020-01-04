### JS.Ops

DevOps platform to delegate business logic to developers without access to sensitive data like passwords and certs.


Project not ready yet.

![alt tag](https://raw.githubusercontent.com/mcfly722/js.ops/master/doc/JS.Ops.png?raw=true)

Creating server authentication certificate
```
makecert.exe -pe -r -n "CN=test-server" -eku "1.3.6.1.5.5.7.3.1,1.3.6.1.5.5.7.3.2" -sky excanage -sv testServer.pvk testServer.cer
pvk2pfx.exe -pvk testServer.pvk -spc testServer.cer -pfx testServer.pfx
```
Creating client authentication certificate
```
makecert.exe -pe -r -n "CN=test-client" -eku "1.3.6.1.5.5.7.3.1,1.3.6.1.5.5.7.3.2" -sky excanage -sv testClient.pvk testClient.cer
pvk2pfx.exe -pvk testClient.pvk -spc testClient.cer -pfx testClient.pfx
```
After certificates creating you have to import it to certificates store.
