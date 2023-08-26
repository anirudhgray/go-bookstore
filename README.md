[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/LECuYE4o)

# Deployed Links
- API baseurl: http://bookstore.anrdhmshr.tech/api/v1
- Health Check: http://bookstore.anrdhmshr.tech/api/health
- API Docs: https://documenter.getpostman.com/view/19697822/2s9Y5Wxifq
# Table of Contents
- [Deployed Links](#deployed-links)
- [Table of Contents](#table-of-contents)
- [How to run locally](#how-to-run-locally)
- [Features](#features)
    - [1. Auth/User:](#1-authuser)
    - [2. Access Control:](#2-access-control)
    - [3. Secure Password Reset Flow:](#3-secure-password-reset-flow)
    - [4. Forgot Password.](#4-forgot-password)
    - [5. Acc Deactivation:](#5-acc-deactivation)
    - [6. Acc Deletion:](#6-acc-deletion)
    - [7. Admin:](#7-admin)
    - [8. SuperAdmin:](#8-superadmin)
    - [9. Books Catalog:](#9-books-catalog)
    - [10. Shopping Cart:](#10-shopping-cart)
    - [11. User Library:](#11-user-library)
    - [12. Reviews:](#12-reviews)
    - [12. Logging with Retention (Rotating Log)](#12-logging-with-retention-rotating-log)
    - [13. Reverse Proxy on prod using nginx.](#13-reverse-proxy-on-prod-using-nginx)
- [Recommendation Engine](#recommendation-engine)
    - [Rough flow for giving users recommendations:](#rough-flow-for-giving-users-recommendations)
- [Where I ran into issues:](#where-i-ran-into-issues)
- [Project Structure](#project-structure)
  - [Explanation](#explanation)
  - [ERD](#erd)
# How to run locally
1. Clone the repo.
2. Make sure you have docker/docker desktop installed, and a docker daemon running.
3. Rename `sample.env` to `.env`, and set JWT signing secret (API_SECRET) with `openssl rand -hex 32`.
4. For MAILTRAP_API_TOKEN, obtain a free api token from https://mailtrap.io/, or run the app without it (confirmation mails, etc, will not be sent in that case.)
5. `make dev`
6. App will be served at http://0.0.0.0:8000/ by default. Check http://0.0.0.0:8000/health to see if everything is OK on serverside.
7. PG Admin will be served at http://0.0.0.0:5050/browser (only in development). Hostname: `postgres_db`, Username and Password as your env config.
8. API Docs at: https://documenter.getpostman.com/view/19697822/2s9Y5Wxifq
   
# Features
### 1. Auth/User: 
JWT Auth. Role Based (Base User and Admin) access to resources. Hashed Password. Password Strength Check. After implementing my registration controllers I realised that a more secure way to do it would have been by asking users to first specify an email, confirm that email, and then ask them to set a password. This would guard against user enumeration attacks. Right now, I'm guarding against it by essentially lying — I say that the confirmation email got sent if you try making an account with an existing email. This could be confusing, however, if the user has genuinely forgotten if they had an account associated with a particular email or not. Alternatively, I could send a warning email in such cases.
### 2. Access Control:
setup doesn't seem too complex (base user and admin), so will not use a library like casbin.
### 3. Secure Password Reset Flow:
Require current and new password. On change, mail email associated with user about the change ("If this was not you, we request you to change your password via the Email-based forgot password option"). On failure to reset password due to incorrect current password, mail associated email about attempt (in case it was a malicious use trying to transfer ownership).
### 4. Forgot Password.
### 5. Acc Deactivation:
(Not implemented) The user's profile is no longer accessible. User's name in their comments is replaced by "Generic User". All user data is retained (such as comment contents, shopping cart, wishlists, purchased books, etc) and restored on next login.
### 6. Acc Deletion:
The user's profile is no longer accessible. User can request an account deletion, and will then recieve a confirmation email, with a link to complete the process. The process must be completed within a set amount of time.
### 7. Admin:
Admin users can add books to the catalog, as well as edit their details. In addition, they can delete reviews and books. They can also ban users (essentially, deactivate plus ban boolean). Deleting books should ordinarily NOT remove it from libraries of users who have already purchased it.
### 8. SuperAdmin:
By default, if there are no admins in the database, the next user who signs up is set as the SuperAdmin (ie, intended to be used when deploying the applicatiod in prod for the first time). SuperAdmin can perform all normal admin functions, with the added permission of promoting base users to Admin, and vice verse.
### 9. Books Catalog:
Users can fuzzy search entire books catalog for title, author and/or category. They can filter by category, and sort by price.
### 10. Shopping Cart:
Users can add any book not already purchased by them to their chopping cart. They can remove books from their cart as well. On checking out their cart, they "buy" all books in the cart, and those books get added to their library of bought books. A transaction record is created. The cart gets cleared of all books on a successful transaction.
### 11. User Library:
Library of books bought by a user. The user can download any of them as many times as they want. Currently stored on server locally, will shift to GCP Buckets or AWS S3.
### 12. Reviews:
Users can only review a book that they have bought (ie, which is in their library). Reviews have a comment, and a rating (which is used to calc avg rating for the book).
### 12. Logging with Retention (Rotating Log)
Currently logging to a local logfile.
### 13. Reverse Proxy on prod using nginx.
# Recommendation Engine
I will be implementing a simple collaborative filtering based recommendations engine. Ref: https://www.toptal.com/algorithms/predicting-likes-inside-a-simple-recommendation-engine

The general idea is that we will not care about the specific attributes of books, and then using some sort of ML Algorithm to figure out what kind of books our user will like (using the user's existing library, reviews, etc etc). Instead, our system will look **similarities** between users.

We keep a record of each user's likes and dislikes (let's base this on review rating — maybe 4 or 5 is "like", everything else is "dislike"). These are two sets that exist for every user. We're going to use something called the [Jaccard Coefficient](https://en.wikipedia.org/wiki/Jaccard_index) to calculate how similar two such sets are. For example, two duplicate sets will be completely similar (coeff of 1) while two sets with nothing in common will have a coeff of 0 (no similarity or overlap between the sets).

```J(A, B) = |A ∩ B| / |A ∪ B|```

### Rough flow for giving users recommendations:
1. Find user's current likes and dislikes.
2. Get all users who have interacted with those items (both liked and unliked).
3. For each of those users, calculate their similarity to our current user using the Jaccard Coefficient. (modified, on a -1 to +1 scale)
4. Get a set of items the current user has not yet liked/disliked.
5. For each of those items, calculate the probability that it should be recommended to the current user. How do we do this:
   1. Result = Numerator / Denominator
   2. Get other users who liked or disliked that item.
   3. Numerator is the sum of the similarity indices of all these users, with the current user.
   4. Denominator is simply the total number of users of liked or disliked the item.
6. Now we can rank the items based on this calculated probability, and give X number of recommendations.
# Where I ran into issues:
- Did not initially realise that gorm's auto-migrations do not, in fact, drop unused columns. While it does make sense as a default so that we don't lose data... well, anyway, spent some time trying to debug why my many2many join table had an unrelated column in it. Ended up dropping the table and then running migrations, will make sure to use my own migration scripts or a more full fledged library like goose.
- Needed to enter associations mode to delete properly, otherwise only the reference would be yeeted.
- The whole flow of user deletion. The culprit? Gorm, yet again.
# Project Structure
```
.
├── Dockerfile
├── Dockerfile-dev
├── Makefile
├── README.md
├── app.log
├── config
│   ├── config.go
│   ├── db.go
│   └── server.go
├── controllers
│   ├── admin.controller.go
│   ├── auth.controller.go
│   ├── books.controller.go
│   ├── cart.controller.go
│   ├── checkout.controller.go
│   ├── review.controller.go
│   └── user.controller.go
├── docker-compose-dev.yml
├── docker-compose-prod.yml
├── go.mod
├── go.sum
├── helpers
│   ├── response.go
│   └── search.go
├── infra
│   ├── database
│   │   └── database.go
│   └── logger
│       └── logger.go
├── main.go
├── migrations
│   └── migration.go
├── models
│   ├── book.model.go
│   ├── cart.model.go
│   ├── deletion_confirmation.model.go
│   ├── forgot_password.model.go
│   ├── review.model.go
│   ├── transaction.model.go
│   ├── user.model.go
│   ├── user_library.model.go
│   └── verification.model.go
├── routers
│   ├── index.go
│   ├── middleware
│   │   ├── auth.middleware.go
│   │   └── cors.go
│   └── router.go
├── tests
│   └── auth.utils_test.go
├── tmp
│   ├── build-errors.log
│   └── main
├── uploads
│   └── resume_anirudh.pdf
└── utils
    ├── auth
    │   └── auth.go
    ├── email
    │   └── email.go
    └── token
        └── token.go
```
## Explanation
- `models/`: The model structs for each table in my DB, along with their relations.
- `controllers/`: Combined handlers for each route, along with controllers for the business logic, as well as data accessing repository functions.
- `routers/`: API routes and auth+cors middleware.
- `utils/`: Reusable utility functions for business logic in controllers. Eg, for mailing, password validation, etc.
- `config/`: Initial reading in config from .env, setting up server and DB configurations.
## ERD
**Note:** This does not show 1:1 relations properly (shown as 1:N). Will update.
![erd](erd.png)
