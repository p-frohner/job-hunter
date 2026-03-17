# Job Hunter

A job board scraper that aggregates listings from LinkedIn, NoFluffJobs, and Profession.hu using a headless Chromium browser.

## Docker (recommended)

The easiest way to get started. Runs both the Go server and React client:

```
make docker-up
```

To stop:

```
make docker-down
```

## Local Development

Install dependencies for both server and client:

```
make install
```

Then run each in a separate terminal:

```
make run-server
make run-client
```

For more details on environment variables, prerequisites, and codegen, see the individual READMEs:

- [Server README](server/README.md)
- [Client README](client/README.md)
