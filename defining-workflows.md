# Defining a Workflow

Create a `workflow.yml` file in your project root. Conductor reads this file and orchestrates each task as a Temporal workflow.

```bash
conductor run workflow.yml
```

Start with [`workflow-example.yaml`](../workflow-example.yaml) as a template.

---

## Overview

Every workflow has three required sections:

| Section | Purpose |
|---------|---------|
| `name` | Display label for this pipeline (e.g., "PR Validation", "Deploy") |
| `tasks` | Catalog of all available tasks and their scripts |
| `execution` | How and in what order tasks should run |

**Key concept:** `tasks` defines *what exists*, `execution` defines *what runs*.

---

## Tasks

Define all available tasks in the `tasks` section. Each task has a unique `name` and a `script`.

### Script Types

**Inline command:**
```yaml
tasks:
  - name: test
    script: go test ./...
```

**Shell script file:**
```yaml
tasks:
  - name: test
    script: ./scripts/test.sh
```

**Requirements:**
- Scripts must be executable (`chmod +x ./scripts/test.sh`)
- Paths are relative to project root
- Work identically on your machine and in Conductor

---

## Execution

Define how tasks run in the `execution` section. Stages execute **top to bottom**. Within each stage, the `mode` determines task behavior.

### Execution Modes

| Mode | Behavior |
|------|----------|
| `parallel` | All tasks start simultaneously. Stage completes when all tasks finish. |
| `sequential` | Tasks run one at a time, in order. Next task waits for previous to succeed. |

### Stage Structure

```yaml
execution:
  - stage: validation
    tasks: [lint, test]           # Task names from tasks catalog
    mode: parallel                # or 'sequential'
    requires_approval: false      # Optional: pause before this stage
```

### Approval Gates

Add `requires_approval: true` to pause a stage. Workflow waits for manual approval before continuing.

```yaml
execution:
  - stage: deploy
    tasks: [deploy]
    mode: sequential
    requires_approval: true       # Pause here, wait for approval
```

**Use case:** Production deploys, security reviews, or any manual gate.

---

## Complete Example

```yaml
name: PR Validation

tasks:
  - name: lint
    script: ./scripts/lint.sh

  - name: test
    script: ./scripts/test.sh

  - name: build
    script: ./scripts/build.sh

  - name: security_scan
    script: ./scripts/security.sh

  - name: deploy
    script: ./scripts/deploy.sh

execution:
  - stage: validation
    description: Lint and test in parallel
    tasks: [lint, test]
    mode: parallel

  - stage: build
    description: Build the application
    tasks: [build]
    mode: sequential

  - stage: security
    description: Run security checks
    tasks: [security_scan]
    mode: sequential

  - stage: deploy
    description: Deploy to production (requires approval)
    tasks: [deploy]
    mode: sequential
    requires_approval: true
```

### Runtime Behavior

1. **validation** — `lint` and `test` run in parallel
2. **build** — `build` runs only if validation passed
3. **security** — `security_scan` runs only if build passed
4. **deploy** — Workflow pauses here. After approval, `deploy` runs

---

