---
name: "plan"
description: "Break a spec file into small, atomic, vertically-sliced tasks for iterative development. Reads a spec from context, builds a dependency graph, writes tasks with acceptance criteria and verification steps, adds phase checkpoints, and saves the plan to docs/plans/."
argument-hint: "Attach or describe the spec file to plan"
agent: "Tech Writer"
---

Follow the caveman communication skill in exactly as defined in [SKILL.md](../../.agents/skills/caveman-ultra/SKILL.md). All responses must be terse, drop all filler, use fragments. This mode is active for **every response in this session**. The final plan document itself must NOT use caveman mode — use clear, standard technical writing for the document only.

## Planning Process

Execute [SKILL.md](../../.agents/skills/planning-and-task-breakdown/SKILL.md) in full. Follow each step below in order.

### Step 1 — Read

- Read the spec file provided in context
- Read relevant codebase sections referenced by the spec
- Note risks, unknowns, and hard constraints

**Do NOT write any code. Output is a plan document only.**

### Step 2 — Slice vertically

One complete feature path per task. Never horizontal layers (no "build all hooks, then all components").

Each task must deliver working, testable functionality end-to-end.

### Step 3 — Write tasks

Each task uses this structure:

```markdown
## Task [N]: [Short descriptive title]

**Description:** What this task accomplishes.

**Acceptance criteria:**

- [ ] [Specific, testable condition]
- [ ] [Specific, testable condition]

**Verification:**

- [ ] Tests pass: `pnpm run test`
- [ ] Build succeeds: `pnpm run build`
- [ ] Manual check: [what to verify]

**Dependencies:** [Task numbers, or "None"]

**Files likely touched:**

- `apps/decs-ui/src/...`

**Estimated scope:** [Small: 1-2 files | Medium: 3-5 files | Large: 5+ files]
```

### Step 4 — Add checkpoints between phases

After every 2–3 tasks, insert a checkpoint:

```markdown
## Checkpoint: After Tasks [N]–[M]

- [ ] All tests pass
- [ ] Application builds without errors
- [ ] [Feature-specific manual verification]
```

### Step 5 — Save

Save the plan to `docs/plans/<spec-name>-plan.md`. Drop `spec` suffix for the plan filename.

---

## TASK

$input
