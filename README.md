[![Review Assignment Due Date](https://classroom.github.com/assets/deadline-readme-button-24ddc0f5d75046c5622901739e7c5dd533143b0c8e959d652212380cedb1ea36.svg)](https://classroom.github.com/a/LECuYE4o)

### Notes:
- Access Control setup doesn't seem too complex (base user and admin), so will not use a library like casbin.

# Features
1. Auth/User: JWT Auth. Role Based (Base User and Admin) access to resources. Hashed Password. Password Strength Check. After implementing my registration controllers I realised that a more secure way to do it would have been by asking users to first specify an email, confirm that email, and then ask them to set a password. This would guard against user enumeration attacks. Right now, I'm guarding against it by essentially lying — I say that the confirmation email got sent if you try making an account with an existing email. This could be confusing, however, if the user has genuinely forgotten if they had an account associated with a particular email or not. Alternatively, I could send a warning email in such cases.
2. Deactivation: The user's profile is no longer accessible. User's name in their comments is replaced by "Generic User". All user data is retained (such as comment contents, shopping cart, wishlists, purchased books, etc) and restored on next login.
3. Deletion: The user's profile is no longer accessible. User's name in their comments is replaced by "Generic User". Purchase records are kept for auditing purposes. All other data removed after grace period of 1 day (in which they can restore account by simply logging in).

TODO add ERD
