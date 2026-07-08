# pitch.ideas

MVP for a school project.
Find the prod deployment here: https://pitch-ideas-wgg.vercel.app/

## Overview
- [pitch.ideas](#pitchideas)
  - [Overview](#overview)
  - [Idea](#idea)
  - [Usage](#usage)
    - [Go](#go)
    - [Locally](#locally)
    - [On Vercel](#on-vercel)
  - [Skipping Vercel deployment](#skipping-vercel-deployment)
  - [License](#license)
  - [Disclaimer](#disclaimer)


## Idea
Pitch your business idea, get feedback and find investors easily.
View ideas or create a pitch for your own idea.
Comment on other ideas and recieve feedback on your own.
Upvote ideas you like.


## Usage

### Go

```bash
go run ./cmd/server
```
or
```bash
go build -o ./server ./cmd/server
server.exe
```


### Locally
Requires Python installed.

```bash
git clone https://github.com/I-had-a-bad-idea/pitch.ideas.git
cd pitch.ideas
pip install -r local-requirements.txt
python server.py
```

Then open `http://localhost:4000`

### On Vercel
1. Import this project into Vercel.
2. Create a PostgreSQL DB (Neon is recommended) and set it up
3. Deploy like normal

## Skipping Vercel deployment
If you dont want to trigger a Vercel deploy add `[skip vercel]` to your commit message.

## License
This project is not licensed under any open-source or source-available license. All rights are reserved.
You may view and run this software for your own personal use. Except where required by applicable law, you may not copy, modify, distribute, sublicense, sell, or create derivative works from any part of this project without prior written permission from the copyright holder.

## Disclaimer
This is a student project. Content is fictional and for demonstration purposes only.
No real investment, funding or user data handling is implied.