# Resume AI Backend

## Requirements
- Docker Desktop `https://www.docker.com/products/docker-desktop/`
- Homebrew `https://brew.sh/`
- poppler (ppdftotext) `brew install poppler`

## 1. Start Mongo db locally
- `make run-mongo-db`

## 2. Run local api server
- `make run-dev-api-server`



---


### Notes for Production
- `RESUME_AI_ENV=production` env variable is required
- `base_url` now requires a full url including https://
- `session_cookie_domain` - now required. this is the domain cookies are created for.