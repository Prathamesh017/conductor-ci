# Conductor CI

YAML-driven CI orchestrator with a terminal UI, powered by [Temporal](https://temporal.io/).

Define a `workflow.yaml` in **your** project, then run Conductor from that directory. Stages can run sequentially or in parallel, with approval gates, auto-retries, and manual retry.

## Dependencies

| Dependency | Why |
|------------|-----|
| [Go](https://go.dev/dl/) 1.22+ | Build / run Conductor |
| [Docker](https://docs.docker.com/get-docker/) | Runs the Temporal server (and UI) |
| Temporal | Orchestrates workflows — started via `docker compose` in this repo |

Optional: [Task](https://taskfile.dev/) for shortcuts (`task run`, `task start:temporal-server`).

## Local setup

### 1. Clone and start Temporal

```bash
git clone <this-repo-url>
cd conductor-ci
docker compose up -d
```

- Temporal server: `localhost:7233`
- Temporal UI: http://localhost:8080

Leave this running while you use Conductor.

### 2. Build

```bash
cd conductor-ci
go build -o conductor ./cmd/main.go
```

### 3. Use it on your project

Conductor looks for `workflow.yml` or `workflow.yaml` in the **current directory**.

```bash
cd /path/to/your-project
cp /path/to/conductor-ci/workflow-example.yaml ./workflow.yaml
# edit workflow.yaml for your scripts

/path/to/conductor-ci/conductor
```

In the TUI:

1. **validate-workflow** — check your YAML  
2. **start-temporal-server** — run the pipeline  
3. While running: **`a`** approve · **`r`** retry failed task · **`q`** quit (terminates the run)

Do **not** use `go run /path/to/conductor-ci/cmd/main.go` from another folder — Go won’t resolve this module correctly. Build the binary and run that instead.

### Develop Conductor itself

From this repo (with Temporal already up):

```bash
task run
# or
go run ./cmd/main.go
```

That uses `workflow.yaml` in the Conductor repo.

## Global setup

Install the binary once, then run `conductor` from any project:

```bash
cd conductor-ci
go build -o conductor ./cmd/main.go
mkdir -p ~/bin
cp conductor ~/bin/conductor
# ensure ~/bin is on your PATH, e.g. in ~/.zshrc:
# export PATH="$HOME/bin:$PATH"
```

Then:

```bash
cd /path/to/your-project   # must contain workflow.yaml
conductor
```

Temporal still needs to be running (`docker compose up -d` from the Conductor repo, or any Temporal on `localhost:7233`).

After you change Conductor code, rebuild and copy again:

```bash
cd conductor-ci
go build -o conductor ./cmd/main.go
cp conductor ~/bin/conductor
```

## Workflow file

See [`workflow-example.yaml`](./workflow-example.yaml) and [`defining-workflows.md`](./defining-workflows.md) for the full YAML shape.

Minimal idea:

- **`tasks`** — named scripts (and optional `retries`)
- **`execution`** — stages with `mode: sequential | parallel` and optional `requires_approval: true`

Scripts run with your project directory as the working directory.
