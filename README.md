# pitch.ideas

MVP for a school project.
Find the prod deployment here: https://pitch-ideas-wgg.vercel.app/

## Overview
- [pitch.ideas](#pitchideas)
  - [Overview](#overview)
  - [Idea](#idea)
  - [Features](#features)
  - [Requirements](#requirements)
  - [Usage](#usage)
    - [Go](#go)
    - [On Vercel](#on-vercel)
      - [Env variables](#env-variables)
  - [Skipping Vercel deployment](#skipping-vercel-deployment)
  - [License](#license)
  - [Disclaimer](#disclaimer)


## Idea
Pitch your business idea, get feedback and find investors easily.
View ideas or create a pitch for your own idea.
Comment on other ideas and receive feedback on your own.
Upvote ideas you like.

## Features

- Create business pitches
- Browse pitches
- Search and sort pitches
- Upvote pitches
- Leave comments
- User authentication
- Edit/Delete comments & pitches


## Requirements
- Go
- Make (optional)

## Usage

### Go

```bash
make run
```
or
```bash
make build-and-run
```
or
```bash
make build
server.exe
```

Then open `http://localhost:4000`

### On Vercel
1. Import this project into Vercel.
2. Create a PostgreSQL DB (Neon is recommended) and set it up
3. Deploy like normal

#### Env variables

| Variable | Description |
| -------- | ----------- |
| DATABASE_URL | PostgreSQL connection string |

## Skipping Vercel deployment
If you don't want to trigger a Vercel deploy add `[skip vercel]` to your commit message.

## License
This project is not licensed under any open-source or source-available license. All rights are reserved.
You may view and run this software for your own personal use. Except where required by applicable law, you may not copy, modify, distribute, sublicense, sell, or create derivative works from any part of this project without prior written permission from the copyright holder.

## Disclaimer
This is a student project. Content is fictional and for demonstration purposes only.
No real investment, funding or user data handling is implied.