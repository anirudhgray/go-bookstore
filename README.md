[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/LECuYE4o)

### Notes:
- Access Control setup doesn't seem too complex (base user and admin), so will not use a library like casbin.

# Features
1. Auth/User: JWT Auth. Role Based (Base User and Admin) access to resources. Hashed Password. Password Strength Check. After implementing my registration controllers I realised that a more secure way to do it would have been by asking users to first specify an email, confirm that email, and then ask them to set a password. This would guard against user enumeration attacks. Right now, I'm guarding against it by essentially lying — I say that the confirmation email got sent if you try making an account with an existing email. This could be confusing, however, if the user has genuinely forgotten if they had an account associated with a particular email or not. Alternatively, I could send a warning email in such cases.
2. Deactivation: The user's profile is no longer accessible. User's name in their comments is replaced by "Generic User". All user data is retained (such as comment contents, shopping cart, wishlists, purchased books, etc) and restored on next login.
3. Deletion: The user's profile is no longer accessible. User's name in their comments is replaced by "Generic User". Purchase records are kept for auditing purposes. All other data removed after grace period of 1 day (in which they can restore account by simply logging in).
4. Admin: Admin users can add books to the catalog, as well as edit their details. In addition, they can delete reviews and books. They can also ban users (essentially, deactivate plus ban boolean). Deleting books should ordinarily NOT remove it from libraries of users who have already purchased it.
1. Catalog: Users can fuzzy search entire books catalog for title, author and/or category. They can filter by category, and sort by price.
2. Shopping Cart: Users can add any book not already purchased by them to their chopping cart. They can remove books from their cart as well. On checking out their cart, they "buy" all books in the cart, and those books get added to their library of bought books. A transaction record is created. The cart gets cleared of all books on a successful transaction.
3. User Library: library of books bought by a user. The user can download any of them as many times as they want. Currently stored on server locally, will shift to GCP Buckets or AWS S3.
4. Reviews: Users can only review a book that they have bought (ie, which is in their library). Reviews have a comment, and a rating (which is used to calc avg rating for the book).
5.  Logging with Retention (Rotating Log): Currently logging to a local logfile.

## Where I ran into issues:
- Did not initially realise that gorm's auto-migrations do not, in fact, drop unused columns. While it does make sense as a default so that we don't lose data... well, anyway, spent some time trying to debug why my many2many join table had an unrelated column in it. Ended up dropping the table and then running migrations, will make sure to use my own migration scripts or a more full fledged library like goose.
- Needed to enter associations mode to delete properly, otherwise only the reference would be yeeted.
