# Konache

Konache is a simple, lightweight, and easy-to-use CLI tool for managing your passwords and auto-filling them in your browser.

## This is v2.0.0

- This is an extension of the original Konache, which was a simple CLI tool for managing your passwords. This version separates the logic onto a separate server, creates a web extension that adds the ability to auto-fill your passwords in your browser, and adds a few more features
- The server is written in Go, and the web extension is written in JavaScript
- Redis is used as the database for storing the passwords, it's hosted on a docker container and is managed by docker-compose
- The server is hosted on a docker container and is managed by docker-compose
