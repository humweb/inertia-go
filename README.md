# Inertia.js Go Adapter

[![GitHub Release](https://img.shields.io/github/v/release/humweb/inertia-go)](https://github.com/humweb/inertia-go/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/humweb/inertia-go.svg)](https://pkg.go.dev/github.com/humweb/inertia-go)
[![go.mod](https://img.shields.io/github/go-mod/go-version/humweb/inertia-go)](go.mod)
[![LICENSE](https://img.shields.io/github/license/humweb/inertia-go)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/humweb/inertia-go/build.yml?branch=main)](https://github.com/humweb/inertia-go/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/humweb/inertia-go)](https://goreportcard.com/report/github.com/humweb/inertia-go)
[![Codecov](https://codecov.io/gh/humweb/inertia-go/branch/main/graph/badge.svg)](https://codecov.io/gh/humweb/inertia-go)

The Inertia.js server-side adapter for Go. Visit [inertiajs.com](https://inertiajs.com) to learn more.


Example App https://github.com/humweb/inertia-go-vue-example

## Installation

Install the package using the `go get` command:

```
go get github.com/humweb/inertia-go
```

## Usage

### 1. Create new instance

```go
url := "http://inertia-app.test" // Application URL for redirect
rootTemplate := "./app.gohtml"   // Root template, see the example below
version := ""                    // Asset version

inertiaManager := inertia.New(url, rootTemplate, version)
```

Or create with `embed.FS` for root template:

```go
import "embed"

//go:embed template
var templateFS embed.FS

// ...

inertiaManager := inertia.NewWithFS(url, rootTemplate, version, templateFS)
```

### 2. Register the middleware

```go
mux := http.NewServeMux()
mux.Handle("/", inertiaManager.Middleware(homeHandler))
```

### 3. Render in handlers

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
    // ...

    err := inertiaManager.Render(w, r, "home/Index", nil)
    if err != nil {
        // Handle server error...
    }
}
```

Or render with props:

```go
// ...

err := inertiaManager.Render(w, r, "home/Index", inertiaManager.Props{
    "total": 32,
})

//...
```

### 4. Server-side Rendering (Optional)

First, enable SSR with the url of the Node server:

```go
inertiaManager.EnableSsrWithDefault() // http://127.0.0.1:13714
```

Or with custom url:

```go
inertiaManager.EnableSsr("http://ssr-host:13714")
```

This is a simplified example using Vue 3 and Laravel Mix.

```js
// resources/js/ssr.js

import { createInertiaApp } from '@inertiajs/vue3';
import createServer from '@inertiajs/vue3/server';
import { renderToString } from '@vue/server-renderer';
import { createSSRApp, h } from 'vue';

createServer(page => createInertiaApp({
    page,
    render: renderToString,
    resolve: name => require(`./pages/${name}`),
    setup({ App, props, plugin }) {
        return createSSRApp({
            render: () => h(App, props)
        }).use(plugin);
    }
}));
```

The following config creates the `ssr.js` file in the root directory, which should not be embedded in the binary.

```js
// webpack.ssr.mix.js

const mix = require('laravel-mix');
const webpackNodeExternals = require('webpack-node-externals');

mix.options({ manifest: false })
    .js('resources/js/ssr.js', '/')
    .vue({
        version: 3,
        options: {
            optimizeSSR: true
        }
    })
    .webpackConfig({
        target: 'node',
        externals: [
            webpackNodeExternals({
                allowlist: [
                    /^@inertiajs/
                ]
            })
        ]
    });
```

You can find the example for the SSR based root template below. For more information, please read the official Server-side Rendering documentation on [inertiajs.com](https://inertiajs.com).

## Examples

The following examples show how to use the package.

### Share a prop globally

```go
inertiaManager.Share("title", "Inertia App Title")
```

### Share a function with root template

```go
inertiaManager.ShareFunc("asset", assetFunc)
```

```html
<script src="{{ asset "js/app.js" }}"></script>
```

### Share a prop from middleware

```go
func authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    
        ctx := inertiaManager.WithProp(r.Context(), "authUserID", user.ID)
        
        // or
        
        ctx := inertiaManager.WithProps(r.Context(), inertiaManager.Props{
            "authUserID": user.ID,
        })
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Share data with root template

```go
ctx := inertiaManager.WithViewData(r.Context(), "meta", meta)
r = r.WithContext(ctx)
```

```html
<meta name="description" content="{{ .meta }}">
```

### Root template

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link href="css/app.css" rel="stylesheet">
        <link rel="icon" type="image/x-icon" href="favicon.ico">
    </head>
    <body>
        <div id="app" data-page="{{ marshal .page }}"></div>
        <script src="js/app.js"></script>
    </body>
</html>
```

### Root template with Server-side Rendering (SSR)

```html
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link href="css/app.css" rel="stylesheet">
        <link rel="icon" type="image/x-icon" href="favicon.ico">
        {{ if .ssr }}
            {{ raw .ssr.Head }}
        {{ end }}
    </head>
    <body>
        {{ if not .ssr }}
            <div id="app" data-page="{{ marshal .page }}"></div>
        {{ else }}
            {{ raw .ssr.Body }}
        {{ end }}
    </body>
</html>
```

## Reporting Issues

If you are facing a problem with this package or found any bug, please open an issue on [GitHub](https://github.com/humweb/inertia-go/issues).

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
