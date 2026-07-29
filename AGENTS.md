# Corvus Studio Repository Rules

- Treat the design documents under `docs/` as implementation sources of truth. Resolve conflicts using the priority recorded in the active ExecPlan.
- Respect phase boundaries strictly. Do not introduce a later-phase domain, dependency, schema, API, workflow, or deployment artifact for completeness.
- Use Trunk Based Development. Conventional Commits are optional and must not be enforced.
- After each logical repository modification, inspect status and diffs, run relevant validation, and create a local commit with explicit paths. Do not push unless explicitly requested.
- Do not add Issue templates, Pull Request templates, or `CONTRIBUTING.md` unless the repository policy is explicitly changed.
- Windows, macOS, and Linux are supported platforms. GitHub Actions is the only CI system.
- Preserve existing user work and inspect the worktree before changing files. Never overwrite unrelated or uncommitted content.
- Every change must have an executable validation path. Record the command and observed result.
- Never claim that a build, test, lint, migration, or workflow passed without current evidence.
- Keep ExecPlans in `.agent/` and maintain their Progress, Surprises & Discoveries, Decision Log, and Outcomes & Retrospective sections while work proceeds.
