# NALOGE ZAPISKI

### NALOGA 1:

REQUIRED FLAGS:

`-file="/path/to/json/file"`

Kolikor sem razumel tematiko, ko sem bral dokumentacijo je jwk JSON struktura, ki predstavlja JAVNE KLJUČE,
zato nisem bil siguren kaj pomeni privaten ključ v jwk formatu. Najbolj smiselno se mi je zdelo narediti tako, da sem PRIVATE KEY shranil kot env sprmenljivko v obliki (v datoteki environment.env):

`PRIVATE_KEY="-----BEGIN RSA PRIVATE KEY-----\nMIIEoQIB...Bu8iIQ==\n-----END RSA PRIVATE KEY-----"
`
Kar se tiče podajanja ENV SPREMENLJIVK kot flag-ov v ukazni vrstici, 

### NALOGA 2:

REQUIRED FLAGS:

`-key_name="..."`

Privatni ključi se shranjujejo v mapco **/keys** in sicer v formatu .pem. Datoteke so poimenovane: `keyName_rsa.pem`

### NALOGA 3:

- /sign

QUERY PARAMS:

keyName string

BODY:

JSON STRUCTURE

```JSON
{
"id": 1,
"content": "NJES"
}
```

- /public

QUERY PARAMS:

`keyName` string

- /validate

QUERY PARAMS:

`keyName` string

BODY:

```JSON
{
    "jws_object": "eyJhbGciOiJSUzI1NiJ9.eyJpZCI6MSwiY29udGVudCI6Ik5KRVMifQ.gjt24zCaFWZcAiUuRQmC-D1Im2AviSqm-P58l7GR0dhGOIZeI9IuHQkwmcnKdhBTJq9B2ekhySI3Whxln1gfSglRyLQpcJu8I3R8J5xHmSXbIFdAhmi13Laa_ObMjwR6zJiflhRX0Sfu0YSyEnhn5LECjjJgpvFZ7w_zMMHubJJt5_oIlEfHcwKWFHsK32IEx5jBvU6fAA1u8icxmk3mix8K3B5vzW66UlowRo0baXovrNodGyHm6wZ-wz1Q_jBhEtxweuCmfOiXaoILNJdfFVqgr7Y2Nw64zP7Pcx7GOX-gVifcuAIdrG3Sy31Sg6Eya6FmaDVBQaCAzQNNo5Nx8A"
}
```
