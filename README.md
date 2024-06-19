**Title:** ZaraDB: A Lightweight and Fast Document Database

**Description:**

ZaraDB is a lightweight, simple, and fast document database currently under intensive development. It aims to be a user-friendly alternative to MongoDB, offering a streamlined API for interacting with your data. While official documentation is forthcoming upon stabilization, this README provides a foundational overview.

**Features:**


**examples:**

             <p>insert one data object </p>
             {action:"insert", collection:"users", data:{name:"adam", age:12}}

             <p>insertMany inserts many data objects at one time 'bulk'</p>
             {action:"insertMany", collection:"test", data:[{name:"jalal", age:23},{name:"akram", age:30},{name:"hasna", age:35}]}

             <p>  select one object</p>
             {action:"findOne", collection:"users", match:{name:"adam"}}

             <p>  select objects match conditions</p>
             {action:"findMany", collection:"users", match:{name:"adam"}}

             <p>select objects that match the conditions</p>
             {action:"findMany", collection:"users", match:{name:"adam", age:{$gt:12}}}
             <p>match numeric data by $eq $nq $lt $gt $ge $le</p>
             <p>match text data by $eq $nq $lt $gt $ge $le </p>

             <p>select objects that match any value </p>
             {action:"findMany", collection:"users", match:{ age:{$in:[12, 23, 34]}}}
             {action:"findMany", collection:"users", match:{ name:{$in:["akram", "zaid"]}}}

             <p>select objects that do not match any value</p>
             {action:"findMany", collection:"users", match:{ age:{$in:[12, 23, 34]}}}
             {action:"findMany", collection:"users", match:{ name:{$nin:["akram", "zaid"]}}}

             <p>select objects that match any conditions</p>
             {action:"findMany", collection:"users", match:{ $or:[name:{$eq:"akram", age:$gt:13}]}}

             <p>select objects that match all conditions</p>
             {action:"findMany", collection:"users", match:{ $and:[name:{$eq:"akram", age:$gt:13}]}}


             <p>updateById </p>
             {action:"updateById", collection:"test", _id:3, data:{name:"hosam", age:10}}

             <p>updateOne </p>
             {action:"updateById", collection:"test", match:{_id{$gt:33}}, data:{name:"hosam", age:10}}

             <p>updateMany </p>
             {action:"updateById", collection:"test",  match:{_id{$gt:33}}, data:{name:"hosam", age:10}}



             <p>delete first objects that match the conditions</p>
             {action:"deleteOne", collection:"users", match:{name:"adam", age:{$gt:12}}}

             <p>delete all objects that match the conditions </p>
             {action:"deleteMany", collection:"users", match:{name:"adam", age:{$gt:12}}}

             <p>skip or ignor some matching objects</p>
             {action:"deleteMany", collection:"users", match:{name:"adam", age:{$gt:12}}, skip: 3}

             <p>Limited to a number of matching objects</p>
             {action:"deleteMany", collection:"users", match:{name:"adam", age:{$gt:12}}, skip: 3, limit:3}

             <p>deleteMany</p>
             {action:"deleteMany", collection:"users", match:{name:"adam", age:{$gt:12}}, skip: 3, limit:3}



             <p>exclode fields</p>
            {action:"findMany", collection:"test", fields:{_id:0, name:0}}


             <p>rename fields</p>
              {action:"findMany", collection:"test", fields:{_id:0, name:"full_name"}}


