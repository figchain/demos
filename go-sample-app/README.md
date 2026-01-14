# Go FigChain Client — Sample App (PoC)

A small proof-of-concept Go application that demonstrates how to use the FigChain Go client to listen for updates to a Fig and process them in real time.

---

## 🚀 Quick start

### Prerequisites

- Go 1.20+ installed and available on your PATH
- A `client-config.json` file exported from the FigChain UI (see notes about `environmentId` below)
- A schema and associated Fig created in [FigChain](https://app.figchain.io)

### Important

We have included a default schema model in `figchain/test.go`. You can use this model, but you will need to create
the schema in your account exactly as defined in `figchain/test.go`. This means creating a schema named `test` with a field `test` of type `string`.

### Fig Creation

Once you have defined your schema in the FigChain UI, create a Fig. Go to the "Figs" section and create a new Fig with key `test` and choose the `test` schema.
You have just defined the key for your data, which remains constant over the lifetime of your use of the Fig. Now you will create a version of the Fig with initial data.
Click the `Add Version` button and the value editor will appear. Double-click the value of the `test` field, and enter any string that you like. Click the check box to confirm,
and click the `Save` button. Now, you have a fig with a value attached to your current schema version. You are ready to test your client.

### Run it locally

1. Put your `client-config.json` in the project root (or the path expected by the app).
2. Run the sample:

```bash
go run .
```

Or build and run a binary:

```bash
go build -o figchain-sample-app .
./figchain-sample-app
```

The app will connect to the FigChain API, initialize the client from `client-config.json`, and register a type-safe listener for the `test` Fig.

---

## 🔧 How it works

- `main.go` boots the FigChain client using `NewClientFromConfig` (reads `client-config.json`).
- A listener is registered for the Fig key named `test` using the generated model in `models/test.go`.
- When the Fig updates, the listener callback receives a typed value and prints it to stdout.
- The app handles graceful shutdown on SIGINT/SIGTERM.

> Tip: You can change the Fig key being watched by editing `main.go` and re-running the app.

---

## 📁 Files of interest

- `main.go` — example client initialization and listener registration
- `client-config.json` — FigChain client configuration (download from FigChain UI)
- `figchain/test.go` — generated Avro model used for type-safe listeners
- `README.md` — this guide

---

## ⚠️ Important notes & troubleshooting

- environmentId is currently required in `client-config.json`. There is an option in the FigChain UI to include it when exporting the configuration.

- Regenerating models: If you change the Avro schema in FigChain, regenerate the Go models and recompile the app so listeners remain type-safe.

---

## 🔁 Development notes

- There are some other options available in the client builder and via environment variables, but `client-config.json` is the default way to initialize the client.
- This app showcases live configuration updates, which means changes to the Fig in FigChain will be reflected in the client in real time. This uses long polling by default, but can be switched to polling if your environment does not support long-lived outbound HTTP connections.

---

## ✅ Expected behavior

When running, you should see logs indicating the client initialized and that the listener is active. On Fig updates, you should see the typed payload printed to stdout by the listener callback.

---

## Contributing

PRs and issues welcome — keep changes small and focused.

---

## License

These sample apps are in the public domain. No license is required.