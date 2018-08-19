# Yoti

A mini project using Yoti for authentication

## How to Setup

- Clone the project
  ```
  go get github.com/twsiyuan/yoti-mini-project
  ```
- You must have a Yoti account, need to download Yoti on Android and register
- Create a new application on [Yoti Dashboard](https://www.yoti.com/dashboard/login)
- Setup the application domain: https://127.0.0.1:8080
- Create a scenario
  - Callback URL: https://127.0.0.1:8080/callback
  - User information required
    - ID Photo
    - Full Name
    - Mobile Number
    - Email Address
- Get Application ID, SDK ID, and KeyFile (*.pem), and set up in ```main.go```
- Create a new database and set up tables using ```sql/database.sql``` on MariaDB (MySQL)
- Setup database connection string in ```main.go```
  - ```user:password@tcp(example.com:3306)/dbname```
- Build application and run

See [Yoti Go SDK](https://github.com/getyoti/yoti-go-sdk) for Yoti more detail.

## Issues

### Anonymous design

Even if post as an anonymous comment, the user ID is stored in the database, but display unknown users in the front-end site.

Administrators can view all user information in anonymous comments in the database. 

Maybe need to remove user ID storage when posts an anonymous comment, for "truly anonymous". 

### Rendering design

Use ```html/template``` to render front-end pages currently, should consider client rendering solutions, such as Vue.js, React.js, or Angular.js. Build API for those front-end and mobile-client calls.

### Others

- Testing files
- HTTPS supports