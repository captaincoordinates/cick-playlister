# CICK Playlister

Tool to support data entry on the CICK site's "Create Playlist" page. See [usage](./USAGE.md) for more information on using the tool.

## Overview

The tool has three components: a server, a client, and a bookmarklet. 

The server component manages communication with streaming service API(s) to retrieve track data for playlists, albums, and tracks. It exposes a number of endpoints that can be viewed with Swagger at [http://localhost:8123/docs/swagger/](http://localhost:8123/docs/swagger/). Credentials are required to interact with streaming service API(s) (see below). The server component is written in Go.

The client component presents a simple modal to the user that accepts URLs for playlists, albums, and tracks. It communicates with the server component to retrieve track data, and fills input fields on the "Create Playlist" page. The client component is written in TypeScript.

The bookmarklet launches the client component. It will only proceed if the current `window.location.href` is either the CICK website or a `file://` path (indicating local development). The bookmarklet is written in JavaScript.

## Testing, CI/CD

Automated testing is still outstanding. This project has limited resources and rapid output has so far been prioritised over rigorous testing. 

_Eventually_ automated testing and a proper CI/CD pipeline will be in place.

## Release

> [!NOTE]
> Requires Bash, Docker

> [!IMPORTANT]  
> A `credentials.json` file is required to provide credentials for the streaming service API(s). This must be present in `./cmd/cick-playlister` before creating a release. The format is as follows:

```json
{
    "spotify": {
        "client_id": "...",
        "client_secret": "..."
    }
}
```

Output is generated in `./dist/{today's date}` and compiled for Windows to suit the CICK station computer:

```sh
scripts/release.sh
```

A file called `bookmarklet.js` in `./dist/{today's date}` contains code required for the bookmarklet that triggers the input modal. The `credentials.json` file will also be copied to the output location so that the release directory contains all necessary files.

## Development

> [!NOTE]
> Requires Bash, Docker, Go 1.22.3, Node 20

### Client Component

Start a local watch-build to transpile TypeScript to JavaScript and Less to CSS during development of the client component:

```sh
scripts/watch-cilent.sh
```

`./research/songs.html` provides an example of the HTML in the "Create Playlist" page and can be opened in a browser to support development.

### Server Component

When the server component is built it bundles the content of `./internal/client/dist` (the client component) and `internal/docs` (OpenAPI & Swagger) in the binary. The server component must be rebuilt and rerun after a change to the client or OpenAPI. To build the client locally (not using Docker, like the release build does):

```sh
scripts/build.sh
```

#### Debugging

If you want to debug the server component, either for problem-solving or simply to avoid the two-step process of rebuilding _and_ rerunning the binary, the following `./vscode/launch.json` configuration can be used (debugging may require some additional tooling in VSCode). You will still have to restart the debugger after a change to the client or server component, but you will not have to explicitly rebuild it:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/cick-playlister/"
        }
    ]
}
```
