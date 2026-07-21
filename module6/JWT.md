Here's a complete JWT (JSON Web Token) Zero to Hero roadmap in Markdown that you can save as JWT-Zero-to-Hero.md.

# JWT (JSON Web Token) - Zero to Hero

> A complete guide to understanding, implementing, and securing JWT authentication.

---

# Table of Contents

1. What is JWT?
2. Why JWT?
3. Traditional Session Authentication
4. JWT Authentication Flow
5. JWT Structure
6. Header
7. Payload
8. Signature
9. JWT Example
10. Encoding vs Encryption
11. Signing Algorithms
12. Creating JWT
13. Verifying JWT
14. Access Tokens
15. Refresh Tokens
16. Authentication Flow
17. Authorization
18. JWT Best Practices
19. Security Risks
20. Token Storage
21. Logout Strategy
22. JWT in REST APIs
23. JWT in Microservices
24. JWT with OAuth
25. Common Mistakes
26. Interview Questions
27. Practical Examples

---

# 1. What is JWT?

JWT (JSON Web Token) is a compact, URL-safe token format used to securely transmit information between parties.

It is mainly used for:

- Authentication
- Authorization
- Information exchange

Official Format:

```
xxxxx.yyyyy.zzzzz
```

Three parts separated by dots.

---

# 2. Why JWT?

Without JWT:

```
User
 ↓
Login
 ↓
Server
 ↓
Stores Session
 ↓
Database
```

With JWT:

```
User
 ↓
Login
 ↓
Server
 ↓
Creates Token
 ↓
Client Stores Token
 ↓
Every Request Includes Token
```

Server becomes mostly stateless.

Advantages

- Stateless
- Scalable
- Fast
- Works across domains
- Ideal for APIs

---

# 3. Traditional Session Authentication

```
User
  |
Login
  |
Server
  |
Creates Session
  |
Stores in Memory/DB
  |
Returns Session ID Cookie
```

Every request:

```
Cookie
 ↓
Server
 ↓
Database Lookup
```

Problem:

- Memory usage
- Scaling issues
- Sticky sessions

---

# 4. JWT Authentication Flow

```
User
 |
Login
 |
Server
 |
Verify Password
 |
Generate JWT
 |
Return JWT
 |
Client Stores JWT
 |
Future Requests

Authorization:
Bearer <token>
```

---

# 5. JWT Structure

Example

```
xxxxx.yyyyy.zzzzz
```

Meaning

```
HEADER.PAYLOAD.SIGNATURE
```

---

# 6. Header

Example

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

Meaning

alg

Algorithm used.

Examples

- HS256
- HS384
- HS512
- RS256

typ

Always JWT.

---

# 7. Payload

Contains claims.

Example

```json
{
  "id": 101,
  "name": "John",
  "role": "admin"
}
```

Claims Types

Registered

```
iss
sub
aud
exp
nbf
iat
jti
```

Public

Custom agreed names.

Private

Application-specific.

Example

```json
{
  "userId": 22,
  "isPremium": true
}
```

---

# 8. Signature

Generated using

```
Base64(Header)
+
Base64(Payload)
+
Secret Key
```

Formula

```
HMACSHA256(
base64UrlEncode(header) + "." +
base64UrlEncode(payload),
secret
)
```

Signature prevents tampering.

---

# 9. JWT Example

Header

```json
{
 "alg":"HS256",
 "typ":"JWT"
}
```

Payload

```json
{
 "id":10,
 "name":"Alice"
}
```

Secret

```
mySecretKey
```

Generated Token

```
eyJhbGc...
.
eyJpZCI...
.
LxP2Fa...
```

---

# 10. Encoding vs Encryption

JWT is NOT encrypted.

It is Base64URL encoded.

Anyone can decode:

```
Header
Payload
```

Only Signature cannot be forged.

Never store:

- Password
- OTP
- Credit card
- Secret API keys

Inside JWT.

---

# 11. Signing Algorithms

HS256

```
Shared Secret
```

Server signs

Server verifies

---

RS256

```
Private Key
```

Signs

```
Public Key
```

Verifies

Used by Google, Auth0, Firebase.

---

# 12. Creating JWT

Pseudo Code

```javascript
jwt.sign(payload, secret, options)
```

Example

```javascript
jwt.sign(
{
 id:1,
 role:"admin"
},
SECRET,
{
 expiresIn:"1h"
}
)
```

---

# 13. Verifying JWT

Example

```javascript
jwt.verify(token, SECRET)
```

If signature matches

Token is valid.

Else

Unauthorized.

---

# 14. Access Token

Short-lived.

Example

```
15 minutes
```

Used for API requests.

---

# 15. Refresh Token

Long-lived.

Example

```
30 days
```

Purpose

Generate new access tokens.

Never send refresh token on every request.

---

# 16. Authentication Flow

```
Login
   |
Access Token
Refresh Token
   |
Access expires
   |
Refresh Endpoint
   |
New Access Token
```

---

# 17. Authorization

Authentication

```
Who are you?
```

Authorization

```
What can you access?
```

Payload

```json
{
 "role":"admin"
}
```

Middleware

```
if(role=="admin")
allow
else
deny
```

---

# 18. JWT Best Practices

✅ HTTPS only

✅ Short expiration

✅ Refresh tokens

✅ Rotate refresh tokens

✅ Strong secret

✅ Verify issuer

✅ Verify audience

✅ Validate expiration

✅ Validate signature

---

# 19. Security Risks

## Token Theft

Use HTTPS.

---

## XSS

Avoid Local Storage.

Use HttpOnly cookies if possible.

---

## CSRF

Use SameSite cookies.

Use CSRF tokens.

---

## Weak Secret

Bad

```
123456
```

Good

```
VeryRandomLongSecretGeneratedKey
```

---

# 20. Token Storage

Option 1

Local Storage

Pros

Easy

Cons

XSS vulnerable

---

Option 2

Session Storage

Pros

Temporary

Cons

Still XSS vulnerable

---

Option 3

HttpOnly Cookie

Pros

Cannot access from JavaScript

Best choice

---

# 21. Logout Strategy

JWT is stateless.

Server cannot instantly destroy token.

Solutions

- Short expiry
- Blacklist
- Refresh token rotation

---

# 22. JWT in REST APIs

Example

```
GET /users
```

Headers

```
Authorization:
Bearer eyJhb...
```

Middleware

```
Verify Token
↓
Extract User
↓
Continue
```

---

# 23. JWT in Microservices

```
Gateway

↓

JWT

↓

Service A

↓

Service B

↓

Service C
```

Every service verifies token independently.

---

# 24. JWT with OAuth

Example

```
Login with Google
```

Google authenticates user.

Returns ID Token.

Backend verifies token.

Creates its own JWT.

---

# 25. Common Mistakes

❌ Storing passwords in JWT

❌ No expiration

❌ Weak secret

❌ Blindly trusting payload

❌ Ignoring signature verification

❌ Long-lived access tokens

---

# 26. Interview Questions

## What is JWT?

A compact token format for authentication and authorization.

---

## Is JWT encrypted?

No.

Encoded only.

---

## Difference between Access and Refresh Token?

Access

Short-lived.

Refresh

Generates new access token.

---

## Can JWT be revoked?

Not directly.

Need blacklist or short expiry.

---

## Why use RS256?

Public/private key separation.

Better for distributed systems.

---

## What is Bearer Token?

```
Authorization:
Bearer <JWT>
```

---

## What happens if payload changes?

Signature becomes invalid.

---

# 27. Practical Example

User Login

```
POST /login
```

↓

Server verifies credentials.

↓

Creates JWT.

↓

Returns token.

↓

Client stores token.

↓

Future request

```
GET /profile

Authorization:
Bearer JWT
```

↓

Middleware verifies token.

↓

Controller executes.

---

# JWT Lifecycle

```
Login
   |
Generate JWT
   |
Store JWT
   |
Send JWT
   |
Receive JWT
   |
Verify JWT
   |
Allow Request
   |
Token Expires
   |
Refresh Token
   |
Generate New JWT
```

---

# Summary

JWT consists of:

```
Header
+
Payload
+
Signature
```

Authentication

```
Login → JWT
```

Authorization

```
JWT → Access Protected Resources
```

Security Checklist

- HTTPS
- Strong Secret
- Short Expiry
- Refresh Tokens
- Verify Signature
- Validate Claims
- Avoid Sensitive Data in Payload

---

# Key Takeaways

- JWT is stateless and ideal for APIs.
- The payload is encoded, not encrypted.
- Always verify the signature before trusting a token.
- Use short-lived access tokens and longer-lived refresh tokens.
- Prefer HttpOnly cookies for browser-based applications when appropriate.
- Never store secrets or passwords inside a JWT payload.


This serves as a comprehensive beginner-to-advanced reference covering JWT concepts, authentication flows, security practices, and common interview topics.