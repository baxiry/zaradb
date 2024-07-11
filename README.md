#### ZaraDB

is a lightweight fast open-licensed document database with no hidden restrictions.
Its goal is to be a lightweight alternative to common documentary databases.
It can also provide superior performance compared to Mongo in many common use cases.

**Installation**
download zara pre-compiled from [here](https://github.com/baxiry/zaradb/releases), extract it, and run it as any program.

Or compile it from soure:

```bash
git clone --depth 1 https://github.com/baxiry/zaradb.git && cd zaradb && go build . && ./zaradb
```

and then open : `localhost:1111`

#### Documentations:
this docs is for libs dev, but it good for everyone : 
[zara api](https://github.com/baxiry/zaradb/wiki/Zara-API)

> [!WARNING]
> It is possible that some things will change within the API, 
> The changes will be partial. Then it stabilizes completely in version 1.0.0
> This documentation is intended for library developers, but is useful to anyone interested in documentary databases 

#### Limitations:
> [!WARNING]
> Horizontal scale is not one of Zara's goals currently.

