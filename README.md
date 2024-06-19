**Title:** ZaraDB: A Lightweight and Fast Document Database

**Description:**

ZaraDB is a lightweight, simple, and fast document database currently under intensive development. It aims to be a user-friendly alternative to MongoDB, offering a streamlined API for interacting with your data. While official documentation is forthcoming upon stabilization, this README provides a foundational overview.

**Features:**

* **Effortless Data Storage:** Store and manage flexible data structures similar to JSON objects (documents) within collections.
* **Efficient Data Manipulation:**
    * `insert`: Insert single or multiple documents into a collection.
    * `find`: Retrieve documents based on various criteria, including:
        * Equality (`$eq`) and inequality (`$nq`) comparisons.
        * Comparisons for numeric data (`$lt`, `$gt`, `$ge`, `$le`).
        * Matching or excluding specific values (`$in`, `$nin`).
        * Combining conditions with logical operators (`$or`, `$and`).
    * `update`: Modify documents by ID or matching criteria.
    * `delete`: Remove documents from a collection, either the first matching one or all that meet the criteria. Options include skipping and limiting deletions.
    * `fields`: Specify which document fields to include or exclude during retrieval.
    * `rename`: Rename fields during retrieval for clarity.
* **Intuitive API (Examples):**

  ```js
  // Insert a document
  {
      action: "insert",
      collection: "users",
      data: {
          name: "adam",
          age: 12
      }
  }

  // Insert multiple documents
  {
      action: "insertMany",
      collection: "test",
      data: [
          { name: "jalal", age: 23 },
          { name: "akram", age: 30 },
          { name: "hasna", age: 35 }
      ]
  }

  // Find one document
  {
      action: "findOne",
      collection: "users",
      match: {
          name: "adam"
      }
  }

  // Find documents with conditions
  {
      action: "findMany",
      collection: "users",
      match: {
          age: { $gt: 12 } // Find users older than 12
      }
  }

  // Additional examples provided in the codebase.
  ```

**Getting Started (Once Available):**

Detailed instructions on installation, usage, and contributing will be provided upon ZaraDB's official release.

**License:**

[Enter the license used by ZaraDB here (e.g., MIT, Apache 2.0)]. You can find the license information in the project's codebase or from the developers.

**Contributing (Once Available):**

[If applicable, outline the contribution process here. Consider mentioning a code of conduct or contribution guidelines.]

**Community:**

[Provide links to communication channels (e.g., forums, chat) if available to foster a growing community around ZaraDB.]



