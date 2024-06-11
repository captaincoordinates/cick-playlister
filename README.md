# CICK Playlister

Tool to support data entry on CICK website. See [usage](./USAGE.md) for more information.

## Testing

Automated testing is still outstanding. This project has limited resources and rapid output has so far been prioritised over rigorous testing.

## Environment Setup

### API
The API backend is written in Go 1.22.3. 
```sh
scripts/build.sh    # requires Docker
```

### Client
The client is written in TypeScript and assumes >= Node 20
```
scripts/build-client.sh
```
