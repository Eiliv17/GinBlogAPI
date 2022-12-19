# Gin JWT Blog API

This is an example of a blog REST API for managing a blog, by default some routes are accessible by anyone like the blog posts, some routes however are protected by JWT Authentication, so you need to register and login to create a post or delete one of your posts.

Here are the available endpoints for now:

/user/signup
- POST - Register a new user, returns a result as JSON

/user/login
- POST - Logs in an already registered user, returns result as JSON and sets the JWT Auth token inside the Cookies.


/validate
- GET - Validation route to test the JWT Authentication

/posts
- GET - Get a list of all the blog posts, returned as JSON
- POST - Creates a blog post, requires authentication

/posts/:id
- GET - Get a post from its ID, returns the post as JSON
- DELETE - Deletes the post given its ID, requires authentication


## Setup
By default you'll need to provide a .env file or set up the environment variables for your operating system, these are the required variables that you need to setup:

Variable                    | Description
---                         | ---
PORT                        | The HTTP server port
MONGODB_URI                 | The URI for connecting to your MongoDB server
DB_NAME                     | The name of the database (inside MongoDB) where you want to store the users and posts
HMAC_SECRET                 | The HMAC secret string

Here's an example of the .env file:
```
PORT=30000
MONGODB_URI= mongodb://127.0.0.1:27017/
DB_NAME=Blog
HMAC_SECRET=Yp2s5v8y/B?E(H+MbQeThWmZq4t6w9z$C&F)J@NcRfUjXn2r5u8x/A%D*G-KaPdS
```

## User & Post Structure
Here's a sample of the basic structure of a user:
```javascript
{
    "_id": ObjectId("63a03df19d18e61bdd6b379e"),
    "username": "mark",
    "email": "mark@gmail.com",
    "password": "$2a$10$rBbnUBpE9IX2nvq3hTlAuC5o1AtYyqA1BzkvR41pdA5EavPeWFy2",
    "createdAt": ISODate("2022-12-19T10:33:21.755Z"),
    "updatedAt": ISODate("2022-12-19T10:33:21.755Z")
}
```

Here's a sample of the basic structure of a post:
```javascript
{
    _id: ObjectId("63a03f149d18e61bdd6b37a2"),
    postTitle: 'Blog Test',
    date: ISODate("2022-12-19T10:38:12.657Z"),
    authorID: ObjectId("63a03df19d18e61bdd6b379e"),
    tags: [ 'blog', "test", "api" ],
    shortDescription: 'This is the first test of a blog post',
    postData: 'This is the content of the test blog post'
}
```

## HTTP Requests Examples

### Register user

#### Request
POST `localhost:3000/user/signup`

```json
{
  "username":"james1",
  "email":"james1@gmail.com",
  "password":"james1"
}
```

#### Response

```json
{
  "result": "user created successfully"
}
```

----------------------------------
### Login User

#### Request

POST `localhost:3000/user/login`

```json
{
  "email":"james1@gmail.com",
  "password":"james1"
}
```

#### Response
```json
{
  "result": "logged in successfully"
}
```
This will set the Cookies for the Authentication of some routes

----------------------------------
### Get all posts

#### Request
GET `localhost:3000/posts`

#### Response
```json
[
  {
    "postID": "63a01f4fc6b7a77333f898d2",
    "postTitle": "My new post",
    "date": "2022-12-19T08:22:39.625Z",
    "authorID": "639f7ef5fcfaa2f78df20a54",
    "tags": [
      "api",
      "tech",
      "something"
    ],
    "shortDescription": "This is my first blog post",
    "postData": "This is an example of my first blog post"
  },
  {
    "postID": "63a03e769d18e61bdd6b379f",
    "postTitle": "test 1",
    "date": "2022-12-19T10:35:34.192Z",
    "authorID": "63a038798cf562a3996a9d11",
    "tags": [
      "test 1"
    ],
    "shortDescription": "test 1",
    "postData": "test 1"
  }
]
```

----------------------------------
### Get a specific post

#### Request

GET `localhost:3000/posts/63a01f4fc6b7a77333f898d2`

#### Response

```json
{
  "postID": "63a01f4fc6b7a77333f898d2",
  "postTitle": "My new post",
  "date": "2022-12-19T08:22:39.625Z",
  "authorID": "639f7ef5fcfaa2f78df20a54",
  "tags": [
    "api",
    "tech",
    "something"
  ],
  "shortDescription": "This is my first blog post",
  "postData": "This is an example of my first blog post"
}
```

----------------------------------
### Create a post

#### Request
POST `localhost:3000/posts/`

```json
{
  "postTitle":"Post test",
  "tags":["tag1", "tag2"],
  "shortDescription":"this is a short description",
  "postData":"this is the post data"
}
```

#### Response
```json
{
  "result": "post created"
}
```

----------------------------------
### Delete a post 

#### Request
DELETE `localhost:3000/posts/63a04d8236fdbf5c30ffa31d`

#### Response
```json
{
  "result": "post deleted"
}
```
