<h1 style="color:green; font-weight:bold; text-align:center;" >ZaraDB</h1>

is an open-licensed document database with no hidden restrictions and good performance.
Its goal is to be a lightweight alternative to common documentary databases.
It can also provide superior performance compared to Mongo in many common use cases.

**Installation**

```bash
git clone --depth 1 https://github.com/baxiry/zaradb.git && cd zaradb && go build . && ./zaradb
```

and then open : `localhost:1111`


> [!NOTE]
> It is possible that some things will change within the API, 
> The changes will be partial. Then it stabilizes completely in version 1.0.0.
> This documentation is intended for library developers, but is useful to anyone interested in documentary databases 


## API Documentation

***Zara receives queries in JSON format. However, you can use JavaScript objects through the web interface provided by Zara via: `localhost:1111`***

### Insert
***insert one object:***

```javascript
{
  collection: "users",
  action: "insert",
  data: {
    name: "adam",
    age: 23,
  },
}
```

***Insert many entries at one time:***

```javascript
{
  collection: "users",
  action: "insertMany",
  data: [
    { name: "Jalal", age: 23 },
    { name: "Alex", age: 30 },
    { name: "Adam", age: 35 },
  ],
}
```
### Find & match
***find objects matching specific conditions:***

```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    name: "adam",
    age: { $gt: 18 }, // Greater than 18
  },
}
```

***Operators works with string & number:***

 `$eq` equal, e.g `match:{age:{$eq: 20}}`

 `$ne` not equal  `match:{age:{$ne: 20}}`

 `$lt` less than, e.g `match:{age:{$lt: 20}}`

 `$gt` greater then, e.g `match:{age:{$gt: 20}}`

 `$ge` greater than or equal to, e.g `match:{age:{$ge: 20}}`

 `$le` less than or equal to, e.g `match:{age:{$le: 20}}`

***Operators works with string only:***

 `$st` string starts with, e.g `match:{name:{$st: "x"}}`

 `$ns` string not starts with, e.g `match:{name:{$ns: "x"}}`

 `$en` string ends with, e.g `match:{name:{$en: "x"}}`

 `$nen` string not ends with, e.g `match:{name:{$nen: "x"}}`

 `$c` string contains, e.g `match:{name:{$c: "x"}}`

 `$nc` string not contain, e.g `match:{name:{$nc: "x"}}`

***Contain.*** String Operators that works with lists:

 `$can` string contains any, e.g `match:{name:{$can: ['a','b','c']}}`

 `$nca` string not contain any, e.g `match:{name:{$nca: ['a','b','c']}}`

 `$cal` string contain all, e.g `match:{name:{$cal: ['a','b','c']}}`

 `$nca` string not contain all, e.g `match:{name:{$nca: ['a','b','c']}}`

***Start.*** String Operators that works with lists:

 `$san` string starts with any, e.g `match:{name:{$san: ['a','b','c']}}`

 `$nsa` string not starts with any, e.g `match:{name:{$nsa: ['a','b','c']}}`

 `$ean` string ends with any, e.g `match:{name:{$ean: ['a','b','c']}}`

 `$nea` string not ends with any, e.g `match:{name:{$nea: ['a','b','c']}}`

#### find objects matching any value in list:

```javascript
// Number list:
{
  collection: "users",
  action: "findMany",
  match: {
    age: { $in: [12, 23, 34] },
  },
}
```
```javascript
// String list:
{
  collection: "users",
  action: "findMany",
  match: {
    name: { $in: ["John", "Zaid"] },
  },
}
```

#### Select objects that do not match any value in list:
find any object whose age does not match any value in the list.
```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    age: { $nin: [12, 23, 34] },
  },
}
```
#### find any object whose name does not match any value in the list.
```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    name: { $nin: ["akram", "zaid"] },
  },
}
```

#### find objects matching any conditions by `$or` operator:

```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    $or: [
      { name: { $eq: "akram" } },
      { age: { $gt: 13 } },
    ],
  },
}
```

#### find objects that matching all conditions by `$and` operator:

```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    $and: [
      { name: { $eq: "akram" } },
      { age: { $gt: 13 } },
    ],
  },
}
```

#### find one object:

```javascript
{
  collection: "users",
  action: "findOne",
  match: {
    name: "adam",
  },
}
```

#### find objects matching conditions:

```javascript
{
  collection: "users",
  action: "findMany",
  match: {
    name: "Adam",
  },
}
```

### Sort by & reverse result:

```javascript
{
  collection: "users",
  action: "findMany",
  sort:{name:1, age:1},
}
// sort by name. names are equal then sort by age
// param 1 = Ascending, anything else = Descending
// Preferably use 0 for Descending
```
### Update

***Update by ID:***

```javascript
{
  collection: "users",
  action: "updateById",
  _id: 3,
  data: {
    name: "Alex",
    age: 10,
  },
}
```

***Update one or more documents matching criteria:***

```javascript
{
  collection: "users",
  action: "updateOne",
  match: { _id: { $gt: 33 } }, // greater than 33
  data: {
    name: "hosam",
    age: 20,
  },
}
```

### Delete

***Delete the first document matching conditions:***

```javascript
{
  collection: "users",
  action: "deleteOne",
  match: {
    name: "adam",
    age: { $gt: 12 },
  }
}
```

### Aggregation

we will appely aggregation on this data:
 
```javascript
{ 
  collection: "products",
  action: "insertMany",
  data: [
	{ item: "Americanos", price: 5,  size: Short",  "quantity": 22 },
	{ item: "Cappuccino", price: 6,  size: Short",  "quantity": 12 },
	{ item: "Lattes",     price: 15, size: Grande", "quantity": 25 },
	{ item: "Mochas",     price: 25, size: Tall",   "quantity": 11 },
	{ item: "Americanos", price: 10, size: Grande", "quantity": 12 },
	{ item: "Cappuccino", price: 7,  size: Tall",   "quantity": 20 },
	{ item: "Lattes",     price: 25, size: Tall",   "quantity": 30 },
	{ item: "Americanos", price: 10, size: Grande", "quantity": 24 },
	{ item: "Cappuccino", price: 10, size: Grande", "quantity": 25 },
	{ item: "Americanos", price: 8,  size: Tall",   "quantity": 28 }
  ]
}
```
####  `$min` `$max` `$count` `$avg` `$sum` examples:

```js
{
  collection: "products",
  action: "aggregate",  
  group: {
      _id: 'item',
      countItems: {$count: ''}, // count param should be zero value
      minPrice: {$min: 'price'},
      maxPrice: {$max: 'price'},
      sumPrice: {$sum: 'price'},
      averagePrice: { $avg: 'price'},
      averageAmount: {$avg: { $multiply: ['quantity','price']}}
    },
}
```

#### match & sort with aggregation

```js
{
  collection: "products",
  action: "aggregate",  
  gmatch:{price:{$gte:20}}
  group: {
      _id: 'item',
      countItems: {$count: ''}, // count param should be zero value
      sumPrice: {$sum: 'price'},
      averagePrice: { $avg: 'price'},
      averageAmount: {$avg: { $multiply: ['quantity','price']}}
    },
    gsort:{averageAmount:1},
}
```


