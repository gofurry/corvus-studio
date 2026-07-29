# Corvus Studio ExecPlan Standard

An ExecPlan is a self-contained, continuously maintained implementation plan. It must let a developer who has no prior repository context complete the described phase without guessing scope, paths, commands, validation, or recovery steps.

## Source of truth and phase control

- Name every source document used by the plan and state what it decides.
- Record the document precedence that applies to the phase.
- Exclude later-phase behavior even when a future document describes it in detail.
- Do not silently resolve conflicting requirements. Record the conflict, recommendation, rationale, and status as `Accepted` or `Open Question`.
- Treat an accepted phase boundary as an invariant. New discoveries may refine execution but may not expand the phase without an explicit decision-log entry.

## Required qualities

Every ExecPlan must be:

- **Self-contained:** define relevant terms, repository roots, module boundaries, tool versions, and assumptions in the plan itself.
- **Executable:** give repository-relative paths and exact commands from a stated working directory.
- **Observable:** describe expected output or files for every validation step.
- **Evidence-based:** distinguish intended results from results actually observed.
- **Idempotent:** explain how reruns detect existing work and avoid overwrites, duplicate initialization, and lockfile drift.
- **Recoverable:** provide a safe recovery procedure for every milestone without destructive Git operations.
- **Living:** update progress, discoveries, decisions, and outcomes as implementation proceeds.

## Required sections

Use these headings in this order unless the phase has an explicitly approved reason to add a more specific subsection:

1. Purpose and user-visible outcome
2. Current repository assessment
3. Source-of-truth documents
4. Scope
5. Non-goals
6. Document conflicts and decisions
7. Proposed repository tree
8. Milestones
9. Go workspace strategy
10. pnpm workspace strategy
11. Build and embedding strategy
12. CI strategy
13. Validation matrix
14. Idempotence and recovery
15. Risks and blockers
16. Completion criteria
17. Progress, discoveries and decision log

## Milestone format

Each milestone must state:

- the goal and user/developer-visible result;
- repository-relative paths created or modified;
- ordered operations, including preconditions and boundary checks;
- exact commands and the directory from which to run them;
- expected output or generated files;
- acceptance evidence;
- safe recovery steps if the milestone stops halfway.

Milestones should be few, independently verifiable, and ordered so that a failed milestone does not invalidate evidence from an earlier one.

## Command and evidence rules

- Use commands that work from the repository root unless another working directory is stated.
- Provide PowerShell commands for required Windows-local verification. CI may use the native shell of its runner.
- Pin tools that affect generated or locked output. Never regenerate a lockfile with an unapproved package-manager version.
- Record command, timestamp, platform, exit code, and a short result in Progress when validation is run.
- A planned command is not evidence. Mark unexecuted checks as pending.
- If a check cannot run, record the exact blocker and any read-only evidence gathered; do not report it as passed.

## Safe maintenance rules

- Before changing a pre-existing file, inspect it and preserve a recoverable copy outside the repository when Git cannot restore the user's current content.
- Before relocating files, resolve absolute source and target paths, check collisions, and hash important content.
- Never use `git reset --hard`, `git checkout --`, or another destructive cleanup command as recovery.
- Update the ExecPlan immediately when an assumption is disproved, a new dependency appears, a command differs from documentation, or a milestone changes shape.
- Keep unresolved decisions in `Open Question` state and stop before the affected mutation if choosing implicitly could change product or architecture behavior.

## Living sections

Every ExecPlan ends with these subsections:

### Progress

Use timestamped checkboxes. Split partially completed work so completed and remaining work are explicit.

### Surprises & Discoveries

Record unexpected repository state, tool behavior, incompatibilities, and evidence. Do not erase discoveries after resolving them.

### Decision Log

For each decision record the date, status, decision, rationale, and affected milestone or interface.

### Outcomes & Retrospective

Before implementation, state `Not started — awaiting user review`. At completion, summarize delivered outcomes, validation evidence, remaining risks, and deviations from the original plan.

## Plan revision rule

When the plan changes, update all affected sections rather than appending a contradictory note. Add a Decision Log entry describing why the plan changed. A replacement plan must remain self-contained.
