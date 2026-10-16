# Working in Guardian

## Start here

Read `README.MD`, `go.mod`, and the relevant code before making changes.
Read the nearest component `AGENTS.md` if one exists. This project is being
initialized: distinguish implemented behavior from proposals, and do not invent
architecture, service contracts, or product requirements.

## Development workflow

- Inspect the working tree first. Preserve existing user edits and untracked
  files; do not reset or restore unrelated work.
- Understand the expected behavior, callers, and relevant edge cases before
  editing. Keep changes small, clear, and focused on the current task.
- For features and fixes, write meaningful behavior tests first, confirm they
  fail for the expected reason, then implement and refactor. For refactors, use
  existing tests or add characterization tests. Documentation-only edits do not
  need application tests.
- Reuse existing conventions. Add dependencies and abstractions only when needed.
- As setup becomes runnable, document actual install, run, and check commands in
  `README.MD`. Run relevant available checks; report missing checks honestly.
- Review the final diff against the requirements before reporting completion.

## Go conventions

- Keep executable startup in `cmd/guardian/` using `package main`. Keep server
  implementation packages under `internal/`. Add packages when responsibilities
  justify them; do not invent layers or interfaces in advance.
- Format Go changes with `gofmt` (or `make fmt`). Keep tests beside the code in
  `*_test.go` files. Run `make check` and `make build` for relevant Go changes.
- Use Go modules for dependencies. Run `go mod tidy` after dependency changes and
  review changes to `go.mod` and `go.sum`; do not change the module path casually.

## Branches and commits

- Start each new feature on its own branch: `<type>/<short-kebab-case-name>`.
- Use commit messages such as `feature: Add health endpoint` or
  `fix: Handle missing configuration`.
- Keep commits coherent and reviewable. Stage only intended changes and preserve
  unrelated work.

## Agentic workflow

- The main agent owns requirements, design decisions, coordination, integration
  review, and the final report. It may edit documentation directly; delegate
  application code, tests, scripts, and tooling changes to an implementation worker.
- Use one worker by default. The worker owns tests, implementation, focused
  verification, and self-review; do not split ordinary testing and coding across
  separate agents. Workers must not spawn further agents.
- Use `gpt-6-luna` with `high` reasoning for implementation, and `medium` for
  delegated documentation. Use `gpt-6.1-sol` with `low` reasoning for independent
  review, or `medium` for high-risk work. Report unavailable models instead of
  silently substituting. Use `fork_turns="none"` for model overrides.
- Give workers concise requirements, owned files, applicable instructions, and
  required checks. Agree on shared interfaces before parallel work and assign
  one owner per edited file.
- Only the main agent changes branches, stages files, and creates commits in a
  shared checkout. Serialize overlapping edits; use worktrees when useful.
- The main agent reviews worker changes. Request independent review for security,
  shared contracts, or broad lifecycle/concurrency changes. Return revisions to
  the existing worker while its context remains useful.
- Report what changed, checks actually run and their results, and remaining
  blockers. Rerun affected checks after relevant changes; avoid repeating checks
  on unchanged code without a reason.
