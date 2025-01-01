# MailNuke

`mail_nuke` is a Go program that helps users clean up their Gmail inboxes by removing unwanted emails from specific senders. It does this by interacting with the Gmail API, allowing users to authorize the application and delete emails from spammy or unwanted sources.

## Features

- **Authenticate with Google**: Uses OAuth2 to authenticate the user with their Google account.
- **Extract Email Senders**: Parses a JSON file (`mails.json`) to get a list of email addresses to delete emails from.
- **Batch Delete Emails**: Deletes emails from specified senders in batches, using the Gmail API.
- **Simple and Efficient**: Easily removes unwanted emails with a minimal setup.

## Prerequisites

- **Go 1.18+**: Ensure that Go is installed and up-to-date.
- **Google Developer Console Project**: You need to create a project on the Google Developer Console, enable the Gmail API, and download the credentials JSON file.

## Installation

1. Clone or download the repository:

```bash
$ git clone https://github.com/mariadriana-deemaze/mail-nuke
$ cd mail_nuke
```

2. Make sure you have Go installed. If not, install it from here.

3. Install dependencies.

```bash
$ go mod get
```

4. Set up your Google Developer Console project:

- Go to Google Developer Console.
- Create a new project.
- Enable the Gmail API.
- Create OAuth 2.0 credentials, with the mail scopes, and download the credentials.json file.
- Save the credentials.json file in the root directory of your project.

## Usage

1. Create a JSON file (mails.json) with an array of email addresses. Example:

```json
[
    "spam@example.com", 
    "junk@example.com", 
    "promo@example.com"
]
```

2. Run the program:

```bash
$ go run main.go
```

> ℹ️ During the first run, the program will ask you to visit a URL and authenticate with your Google account. Enter the authorization code in the callback URL to proceed.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
