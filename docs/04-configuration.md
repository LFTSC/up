---
title: Configuration
---

Configuration for your app lives in the `up.json` within your project's directory. This section details each of the options available.

## Name

The name of the application, which is used to name resources such as the Lambda function or API Gateway.

```json
{
  "name": "api"
}
```

## Profile

The `profile` property is equivalent to setting `AWS_PROFILE` for referencing AWS credentials in the `~/.aws` directory. Use of this property is preferred as it prevents accidents with environment variables.

```json
{
  "profile": "someapp"
}
```

## Regions

You may specify one or more target regions for deployment using the `regions` array. Glob style patterns may be used to match region ids. By default "us-west-2" is used unless the `AWS_REGION` environment variable is defined.

A single region:

```json
{
  "regions": ["us-west-2"]
}
```

Several regions:

```json
{
  "regions": ["us-west-2", "us-east-1", "ca-central-1"]
}
```

USA and Canada only:

```json
{
  "regions": ["us-*", "ca-*"]
}
```

Western USA only:

```json
{
  "regions": ["us-west-*"]
}
```

All regions like a boss:

```json
{
  "regions": ["*"]
}
```

**WARNING**: multi-region support won't be complete until https://github.com/apex/up/issues/134 is closed.

## Lambda Settings

The following Lambda-specific settings are available:

- `role` – IAM role ARN, defaulting to the one Up creates for you
- `memory` – Function memory in mb (Default `512`, Min `128`, Max `1536`)
- `accelerate` – Enable [S3 Transfer Acceleration](https://docs.aws.amazon.com/AmazonS3/latest/dev/transfer-acceleration.html) (Default `false`)

For example:

```json
{
  "name": "api",
  "lambda": {
    "memory": 512,
    "accelerate": true
  }
}
```

Lambda timeout is implied from the [Reverse Proxy](#configuration.reverse_proxy) `timeout` setting.

Lambda `memory` also scales the CPU, if your app is slow, or for cases such as larger Node applications with many `require()`s you may need to increase this value. View the [Lambda Pricing](https://aws.amazon.com/lambda/pricing/) page for more information regarding the `memory` setting.

Lambda uploads may be optionally accelerated for Up Pro users, using [S3 Transfer Acceleration](https://docs.aws.amazon.com/AmazonS3/latest/dev/transfer-acceleration.html), this is especially helpful for improving overseas deploy times. Try the [speedtest tool](http://s3-accelerate-speedtest.s3-accelerate.amazonaws.com/en/accelerate-speed-comparsion.html) to get an idea of the improvements, and view [S3 Pricing](https://aws.amazon.com/s3/pricing/).

Changes to Lambda configuration do not require a `up stack apply`, just deploy and these changes are picked up!

### Active Warming

An AWS Lambda "cold start" happens when there are no freely available containers to serve the request, Lambda must fetch your code and create a new container, after this it is "warm" and remains in the Lambda cluster to serve subsequent requests for roughly an hour.

If a container does not receive any traffic within the hour, it is removed from the AWS Lambda cluster, and thus a new request may incur a cold start. Up Pro's "active warming" feature mitigates this by periodically requesting against your app, at the specified concurrency.

For example when a user visits your web application, a cold start may occur for each resource, say you have one JavaScript file, CSS file, and the HTML itself, then this will be 3 concurrent containers. By default Up will warm `15` containers, however you may want to adjust `warm_count` this for your use-case.

Note that if your application receives steady traffic this may not be an issue at all in practice, as containers will already be warm and re-used.

- `warm` – Enable active warming (Default: false)
- `warm_count` – Number of concurrent containers to warm (Default: `15`)
- `warm_rate` – Rate at which to perform the warming (Default: `"15m"`)

Here's a example specifying `50` idle containers:

```json
{
  "name": "app",
  "lambda": {
    "warm": true,
    "warm_count": 50,
    "warm_rate": "30m"
  }
}
```

Run `up stack plan` and `up stack apply` to make the changes to your stack!

Another way to mitigate cold starts is to use an uptime monitoring tool, such as [Apex Ping](https://apex.sh/ping/) which also monitors global performance, so it's a win-win! Use the "up" coupon for 15% off your first year.

[![Uptime Monitoring Tool](https://apex-software.imgix.net/ping/marketing/overview.png?compress=auto)](https://apex.sh/ping)

## Hook Scripts

Up provides "hooks" which are commands invoked at certain points within the deployment workflow for automating builds, linting and so on. The following hooks are available:

- `prebuild` – Run before building
- `build` – Run before building. Overrides inferred build command(s)
- `postbuild` – Run after building
- `predeploy` – Run before deploying
- `postdeploy` – Run after deploying
- `clean` – Run after a deploy to clean up artifacts. Overrides inferred clean command(s)

Here's an example using Browserify to bundle a Node application. Use the `-v` verbose log flag to see how long each hook takes.

```json
{
  "name": "app",
  "hooks": {
    "build": "browserify --node app.js > server.js",
    "clean": "rm server.js"
  }
}
```

Up performs runtime inference to discover what kind of application you're using, and does its best to provide helpful defaults – see the [Runtimes](#runtimes) section.

Multiple commands be provided by using arrays, and are run in separate shells:

```json
{
  "name": "app",
  "hooks": {
    "build": [
      "mkdir -p build",
      "cp -fr static build",
      "browserify --node index.js > build/client.js"
    ],
    "clean": "rm -fr build"
  }
}
```

To get a better idea of when hooks run, and how long the command(s) take, you may want to deploy with `-v` for verbose debug logs.

## Static File Serving

Up ships with a robust static file server, to enable it specify the app `type` as `"static"`.

```json
{
  "type": "static"
}
```

By default the current directory (`.`) is served, however you can change this using the `dir` setting. The following configuration restricts only serving of files in `./public/*`, any attempts to read files from outside of this root directory will fail.

```json
{
  "name": "app",
  "type": "static",
  "static": {
    "dir": "public"
  }
}
```

Note that `static.dir` only tells Up which directory to serve – it does not exclude other files from the directory – see [Ignoring Files](#configuration.ignoring_files). For example you may want an `.upignore` containing:

```
*
!public/**
```

## Environment Variables

The `environment` object may be used for plain-text environment variables. Note that these are not encrypted, and are stored in up.json which is typically committed to GIT, so do not store secrets here.

```json
{
  "name": "api",
  "environment": {
    "API_FEATURE_FOO": "1",
    "API_FEATURE_BAR": "0"
  }
}
```

These become available to you via `process.env.API_FEATURES_FOO`, `os.Getenv("API_FEATURES_FOO")` or similar in your language of choice.

The following environment variables are provided by Up:

- `PORT` – port number such as "3000"
- `UP_STAGE` – stage name such as "staging" or "production"

Up Pro offers encrypted environment variables via the [up env](#commands.env) sub-command.

## Header Injection

The `headers` object allows you to map HTTP header fields to paths. The most specific pattern takes precedence.

Here's an example of two header fields specified for `/*` and `/*.css`:

```json
{
  "name": "app",
  "type": "static",
  "headers": {
    "/*": {
      "X-Something": "I am applied to everything"
    },
    "/*.css": {
      "X-Something-Else": "I am applied to styles"
    }
  }
}
```

Requesting `GET /` will match the first pattern, injecting `X-Something`:

```
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 200
Content-Type: text/html; charset=utf-8
Last-Modified: Fri, 21 Jul 2017 20:42:51 GMT
X-Powered-By: up
X-Something: I am applied to everything
Date: Mon, 31 Jul 2017 20:49:33 GMT
```

Requesting `GET /style.css` will match the second, more specific pattern, injecting `X-Something-Else`:

```json
HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 50
Content-Type: text/css; charset=utf-8
Last-Modified: Fri, 21 Jul 2017 20:42:51 GMT
X-Powered-By: up
X-Something-Else: I am applied to styles
Date: Mon, 31 Jul 2017 20:49:35 GMT
```

## Error Pages

By default Up will serve a minimalistic error page for requests accepting `text/html`. The following settings are available:

- `disable` — remove the error page feature and default pages
- `dir` — the directory where the error pages are located
- `variables` — vars available to the pages

The default template's `color` and optionally provide a `support_email` to allow customers to contact your support team, for example:

```json
{
  "name": "site",
  "type": "static",
  "error_pages": {
    "variables": {
      "support_email": "support@apex.sh",
      "color": "#228ae6"
    }
  }
}
```

If you'd like to provide custom templates you may create one or more of the following files. The most specific file takes precedence.

- `error.html` – Matches any 4xx or 5xx
- `5xx.html` – Matches any 5xx error
- `4xx.html` – Matches any 4xx error
- `CODE.html` – Matches a specific code such as 404.html

Variables specified via `variables`, as well as `.StatusText` and `.StatusCode` may be used in the template.

```html
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>{{.StatusText}} - {{.StatusCode}}</title>
    <link rel="stylesheet" href="/css/style.css">
  </head>
  <body>
    <h1>{{.StatusText}}</h1>
    {{with .Variables.support_email}}
      <span class="message">Please try your request again or <a href="mailto:{{.}}">contact support</a>.</span>
    {{else}}
      <span class="message">Please try your request again or contact support.</span>
    {{end}}
  </body>
</html>
```

## Script Injection

Scripts, styles, and other tags may be injected to HTML pages before the closing `</head>` tag or closing `</body>` tag.

In the following example the `<link rel="/style.css">` is injected to the head, as well as the inlining the `scripts/config.js` file. A `<script src="/app.js"></script>` is then injected into the body.


```json
{
  "name": "site",
  "type": "static",
  "inject": {
    "head": [
      {
        "type": "style",
        "value": "/style.css"
      },
      {
        "type": "inline script",
        "file": "scripts/config.js"
      }
    ],
    "body": [
      {
        "type": "script",
        "value": "/app.js"
      }
    ]
  }
}
```

Currently you may specify the following types:

- `literal` – A literal string
- `comment` – An html comment
- `style` – A style `href`
- `script` – A script `src`
- `inline style` – An inline style
- `inline script` – An inline script
- `google analytics` – Google Analytics snippet with API key
- `segment` – Segment snippet with API key

All of these require a `value`, which sets the `src`, `href`, or inline content. Optionally you can populate `value` via a `file` path to a local file on disk, this is typically more convenient for inline scripts or styles. For example:

- `{ "type": "literal", "value": "<meta name=...>" }`
- `{ "type": "comment", "value": "Just a boring comment" }`
- `{ "type": "script", "value": "/feedback.js" }`
- `{ "type": "style", "value": "/feedback.css" }`
- `{ "type": "inline script", "file": "/feedback.js" }`
- `{ "type": "inline style", "file": "/feedback.css" }`
- `{ "type": "script", "value": "var config = {};" }`
- `{ "type": "google analytics", "value": "API_KEY" }`
- `{ "type": "segment", "value": "API_KEY" }`

## Redirects and Rewrites

Up supports redirects and URL rewriting via the `redirects` object, which maps path patterns to a new location. If `status` is omitted (or 200) then it is a rewrite, otherwise it is a redirect.

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/blog": {
      "location": "https://blog.apex.sh/",
      "status": 301
    },
    "/docs/:section/guides/:guide": {
      "location": "/help/:section/:guide",
      "status": 302
    },
    "/store/*": {
      "location": "/shop/:splat"
    }
  }
}
```

In the previous example `/blog` will redirect to a different site, while `/docs/ping/guides/alerting` will redirect to `/help/ping/alerting`. Finally `/store/ferrets` and nested paths such as `/store/ferrets/tobi` will redirect to `/shop/ferrets/tobi` and so on.

A common use-case for rewrites is for SPAs or Single Page Apps, where you want to serve the `index.html` file regardless of the path. The other common requirement for SPAs is that you of course can serve scripts and styles, so by default if a file is found, it will not be rewritten to `location`.

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/*": {
      "location": "/",
      "status": 200
    }
  }
}
```

If you wish to force the rewrite regardless of a file existing, set `force` to `true` as shown here:

```json
{
  "name": "app",
  "type": "static",
  "redirects": {
    "/*": {
      "location": "/",
      "status": 200,
      "force": true
    }
  }
}
```

Note that more specific target paths take precedence over those which are less specific, for example `/blog` will win over and `/*`.

## Cross-Origin Resource Sharing

CORS is a mechanism which allows requests originating from a different host to make requests to your API. Several options are available to restrict this access, if the defaults are appropriate simply enable it as shown below.

```json
{
  "cors": {
    "enable": true
  }
}
```

Suppose you have `https://api.myapp.com`, you may want to customize `cors` to permit access only from `https://myapp.com` so that other sites cannot call your API directly.

```json
{
  "cors": {
    "allowed_origins": ["https://myapp.com"],
    "allowed_methods": ["HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"],
    "allowed_headers": ["Content-Type", "Authorization"],
    "allow_credentials": true
  }
}
```

- `allowed_origins` – A list of origins a cross-domain request can be executed from. Use `*` to allow any origin, or a wildcard such as `http://*.domain.com` (Default: `["*"]`)
- `allowed_methods` – A list of methods the client is allowed to use with cross-domain requests. (Default: `["HEAD", "GET", "POST"]`)
- `allowed_headers` – A list of headers the client is allowed to use with cross-domain requests. If the special `*` value is present in the list, all headers will be allowed. (Default: `[]`)
- `exposed_headers` – A list of headers which are safe to expose to the API of a CORS response.
- `max_age` – A number indicating how long (in seconds) the results of a preflight request can be cached.
- `allow_credentials` – A boolean indicating whether the request can include user credentials such as cookies, HTTP authentication or client side SSL certificates. (Default: `true`)

## Reverse Proxy

Up acts as a reverse proxy in front of your server, this is how CORS, redirection, script injection and other middleware style features are provided.

The following settings are available:

- `command` – Command run through the shell to start your server (Default `./server`)
  - When `package.json` is detected `npm start` is used
  - When `app.js` is detected `node app.js` is used
  - When `app.py` is detected `python app.py` is used
- `backoff` – Backoff configuration object described in "Crash Recovery"
- `retry` – Retry idempotent requests upon 5xx or server crashes. (Default `true`)
- `timeout` – Timeout in seconds per request (Default `15`, Max `25`)
- `listen_timeout` – Timeout in seconds Up will wait for your app to boot and listen on `PORT` (Default `15`, Max `25`)
- `shutdown_timeout` – Timeout in seconds Up will wait after sending a SIGINT to your server, before sending a SIGKILL (Default `15`)

```json
{
  "proxy": {
    "command": "node app.js",
    "timeout": 10,
    "listen_timeout": 5,
    "shutdown_timeout": 5
  }
}
```

Lambda's function timeout is implied from the `.proxy.timeout` setting.

### Crash Recovery

Another benefit of using Up as a reverse proxy is performing crash recovery. Up will retry idempotent requests upon failure, and upon crash it will restart your server and re-attempt before responding to the client.

By default the back-off is configured as:

- `min` – Minimum time before retrying (Default `100ms`)
- `max` – Maximum time before retrying (Default `500ms`)
- `factor` – Factor applied to each attempt (Default `2`)
- `attempts` – Attempts made before failing (Default `3`)
- `jitter` – Apply jitter (Default `false`)

A total of 3 consecutive attempts will be made before responding with an error, in the default case this will be a total of 700ms for the three attempts.

Here's an example tweaking the default behaviour:

```json
{
  "proxy": {
    "command": "node app.js",
    "backoff": {
      "min": 500,
      "max": 1500,
      "factor": 1.5,
      "attempts": 5,
      "jitter": true
    }
  }
}
```

Since Up's purpose is to proxy your http traffic, Up will treat network errors as a crash.  When Up detects this, it will allow the server to cleanly close by sending a SIGINT, it the server does not close within `proxy.shutdown_timeout` seconds, it will forcibly close it with a SIGKILL.

## DNS Zones & Records

Up allows you to configure DNS zones and records. One or more zones may be provided as keys in the `dns` object ("myapp.com" here), with a number of records defined within it.

```json
{
  "name": "gh-polls",
  "dns": {
    "gh-polls.com": [
      {
        "name": "app.gh-polls.com",
        "type": "CNAME",
        "value": ["gh-polls.netlify.com"]
      }
    ]
  }
}
```

The record `type` must be one of:

- A
- AAAA
- CNAME
- MX
- NAPTR
- NS
- PTR
- SOA
- SPF
- SRV
- TXT

## Stages

Up supports the concept of "stages" for configuration, such as mapping of custom domains, or tuneing the size of Lambda function to use.

By default the following stages are defined:

- `development` — local development environment
- `staging` — remote environment for staging new features or releases
- `production` — remote environment for production

To create a new stage, first add it to your configuration, in this case we'll call it "beta":

```json
{
  "name": "app",
  "lambda": {
    "memory": 128
  },
  "stages": {
    "beta": {

    }
  }
}
```

Now you'll need to plan your stack changes, which will set up a new API Gateway and permissions:


```
$ up stack plan

Add api deployment
  id: ApiDeploymentBeta

Add lambda permission
  id: ApiLambdaPermissionBeta
```

Apply those changes:

```
$ up stack apply
```

Now you can deploy to your new stage by passing the name `beta` and open the end-point in the browser:

```
$ up beta
$ up url -o beta
```

To delete a stage, simply remove it from the `up.json` configuration and run `up stack plan` again, and `up stack apply` after reviewing the changes.

You may of course assign a custom domain to these stages as well, let's take a look at that next!

## Stages & Custom Domains

By defining a stage and its `domain`, Up knows it will need to create a free SSL certificate for `gh-polls.com`, setup the DNS records, and map the domain to API Gateway.

```json
{
  "stages": {
    "production": {
      "domain": "gh-polls.com"
    }
  }
}
```

 Here's another example mapping each stage to a domain, note that the domains do not need to be related, you could use `stage-gh-polls.com` for example.


```json
{
  "stages": {
    "production": {
      "domain": "gh-polls.com"
    },
    "staging": {
      "domain": "stage.gh-polls.com"
    }
  }
}
```

You may also provide an optional base path, for example to prefix your API with `/v1`. Note that currently your application will still receive "/v1" in its request path, for example Node's `req.url` will be  "/v1/users" instead of "/users".

```json
{
  "stages": {
    "production": {
      "domain": "api.gh-polls.com",
      "path": "/v1"
    }
  }
}
```


Plan the changes via `up stack plan` and `up stack apply` to perform the changes. Note that CloudFront can take up to ~40 minutes to distribute this configuration globally, so grab a coffee while these changes are applied.

You may [purchase domains](#guides.development_to_production_workflow.purchasing_a_domain) from the command-line, or map custom domains from other registrars. Up uses Route53 to purchase domains using your AWS account credit card. See `up help domains`.


## Stage Overrides

Up allows some configuration properties to be overridden at the stage level. The following example illustrates how you can tune lambda memory and hooks per-stage.

```json
{
  "name": "app",
  "hooks": {
    "build": "parcel index.html --no-minify -o build",
    "clean": "rm -fr build"
  },
  "stages": {
    "production": {
      "hooks": {
        "build": "parcel index.html -o build"
      },
      "lambda": {
        "memory": 1024
      }
    }
  }
}
```

Currently the following properties may be specified at the stage level:

- `hooks`
- `lambda`
- `proxy.command`

For example you may want to override `proxy.command` for development, which is the env `up start` uses. In the following example [gin](https://github.com/codegangsta/gin) is used for hot reloading of Go programs:

```json
{
  "name": "app",
  "stages": {
    "development": {
      "proxy": {
        "command": "gin --port $PORT"
      }
    }
  }
}
```

## Logs

By default Up treats stdout as `info` level logs, and stderr as `error` level. If your logger uses stderr, such as Node's `debug()` module and you'd like to change this behaviour you may override these levels:

```json
{
  "name": "app",
  "environment": {
    "DEBUG": "myapp"
  },
  "logs": {
    "stdout": "info",
    "stderr": "info"
  }
}
```

## Ignoring Files

Up supports gitignore style pattern matching for omitting files from deployment via the `.upignore` file.

An example `.upignore` to omit markdown and `.go` source files might look like this:

```
*.md
*.go
```

### Negation

By default dotfiles are ignored, if you wish to include them, you may use `!` to negate a pattern in `.upignore`:

```
!.myfile
```

Another use-case for negation is to ignore everything and explicitly include a number of files instead, to be more specific:

```
*
!app.js
!package.json
!node_modules/**
!src/**
```

### Inspecting

To get a better idea of which files are being filtered or added, use `up -v` when deploying, and you may also find it useful to `grep` in some cases:

```
$ up -v 2>&1 | grep filtered
DEBU filtered .babelrc – 25
DEBU filtered .git – 408
DEBU filtered .gitignore – 13
DEBU filtered node_modules/ansi-regex/readme.md – 1749
DEBU filtered node_modules/ansi-styles/readme.md – 1448
DEBU filtered node_modules/binary-extensions/readme.md – 751
DEBU filtered node_modules/chalk/readme.md – 6136
```

You may also wish to use `up build --size` to view the largest files within the zip.

### Pattern matching

Note that patterns are matched much like `.gitignore`, so if you have the following `.upignore` contents even `node_modules/debug/src/index.js` will be ignored since it contains `src`.

```
src
```

You can be more specific with a leading `./`:

```
./src
```

Files can be matched recursively using `**`, for example ignoring everything except the files in `dist`:

```
*
!dist/**
```

## Alerting

The Pro version of Up supports defining alerts which can notify your team when your service is failing, responding slowly, or receiving traffic spikes.

```json
{
  "name": "app",
  "actions": [
    {
      "name": "email.backend",
      "type": "email",
      "emails": ["tj@apex.sh"]
    },
    {
      "name": "text.backend",
      "type": "sms",
      "numbers": ["+12508183100"]
    }
  ],
  "alerts": [
    {
      "metric": "http.count",
      "statistic": "sum",
      "threshold": 100,
      "action": "email.backend"
    },
    {
      "metric": "http.5xx",
      "statistic": "sum",
      "threshold": 1,
      "period": "1m",
      "action": "email.backend"
    },
    {
      "metric": "http.4xx",
      "statistic": "sum",
      "threshold": 50,
      "period": "5m",
      "action": "email.backend"
    },
    {
      "metric": "http.latency",
      "statistic": "avg",
      "threshold": 1000,
      "period": "5m",
      "action": "email.backend",
      "description": "Large traffic spike"
    }
  ]
}
```

### Defining Actions

An action must be defined in order to notify your team of triggered and resolved  alerts. The action requires a `name`, which can be any string such as `backend`, `frontend_team`, `email.backend`, `email backend` – whichever you prefer.

#### Email

The email action notifies your team via email.

```json
{
  "name": "email.backend",
  "type": "email",
  "emails": ["tj@apex.sh"]
}
```

#### SMS

The sms action notifies your team via sms text message.

```json
{
  "name": "text.backend",
  "type": "sms",
  "numbers": ["+12508183100"]
}
```

#### Slack

The slack action notifies your team via Slack message.

```json
{
  "name": "slack.backend",
  "type": "slack",
  "url": "https://hooks.slack.com/services/T0YS6H6S5/..."
}
```

Optionally `gifs` can be enabled and a `channel` can be specified instead of using the default for the `url`.

```json
{
  "name": "slack.backend",
  "type": "slack",
  "url": "https://hooks.slack.com/services/T0YS6H6S5/...",
  "channel": "alerts",
  "gifs": true
}
```

### Defining Alerts

An alert requires a `metric` such as request count or latency, statistic such as sum or svg, threshold, and an action to perform when the alert is triggered. Here's a simple example emailing the backend team when we encounter a spike of over 1000 requests.

```json
{
  "metric": "http.count",
  "statistic": "sum",
  "threshold": 1000,
  "action": "email.backend"
}
```

#### Required settings

- `metric` – Metric to alert against
  - `http.count` – Request count
  - `http.latency` – Request latency in milliseconds
  - `http.4xx` – HTTP 4xx client errors
  - `http.5xx` – HTTP 5xx server errors
- `statistic` – Statistic name ("sum", "min", "max", "avg", "count")
- `threshold` – Threshold which is compared to `operator`
- `action` – Name of the action to perform

#### Optional settings

- `period` – Period is the alert query time-span (default: `1m`)
- `evaluation_periods` – Number of periods to evaluate over (default: `1`)
- `operator` – Operator is the comparison operator (default `>`)
- `namespace` – Metric namespace (example: "AWS/ApiGateway")
- `missing`– How to treat missing data (default: `notBreaching`)
- `disable` – Disable or mute the alert
- `description` – Description of the alert
