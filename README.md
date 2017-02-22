# nano-ddns
Plain and dead-simple dynamic DNS server

# Building

```bash
make clean
make
```

# Usage

```bash
Usage of ddns.windows_386.exe:
  -database string
        Database file (default "db.json")
  -dns-address string
        DNS Address to listen to (TCP and UDP) (default ":5354")
  -http-address string
        HTTP Address to listen to (default ":8080")
  -update-url string
        URL for updating ddns info (default "/updateip.php")
```

# Example

```bash
ddns.windows_386.exe -dns-address=":5353" -http-address=":8080"
```

```bash
curl -XGET 'localhost:8080/updateip.php?modem=100&ip=192.168.1.100'
```

```bash
dig @localhost -p 5353 modem100. A

; <<>> DiG 9.8.3-P1 <<>> @localhost -p 5353 modem100. A
; (2 servers found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 5183
;; flags: qr rd; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;modem100.			IN	A

;; ANSWER SECTION:
modem100.		0	IN	A	192.168.1.100

;; Query time: 0 msec
;; SERVER: ::1#5353(::1)
;; WHEN: Wed Feb 22 14:10:55 2017
;; MSG SIZE  rcvd: 50
```
