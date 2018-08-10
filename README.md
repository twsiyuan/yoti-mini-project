# Yoti

A mini project using Yoti for authentication

## How to Setup

- You must have a Yoti account, need to download Yoti on Android and register.
- Create a new application on [Yoti Dashboard](https://www.yoti.com/dashboard/login)
- Setup callback URL
- Get AppID, SDKID, and KeyFile (*.pem), and set up in ```main.go```
- Create a new database using ```sql/database.sql```
- Setup database connection string in ```main.go```

See [Yoti Go SDK](https://github.com/getyoti/yoti-go-sdk) for more detail.