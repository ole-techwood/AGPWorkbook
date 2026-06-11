---
name: "plan"
description: "Plan creation from an attached spec using structured task breakdown and review checkpoints."
argument-hint: "Attach spec and describe planning scope"
agent: Tech Writer
---

## Caveman Rules

Follow the caveman communication skill exactly as defined in [SKILL.md](../../.agents/skills/caveman/SKILL.md). All responses must be terse, drop filler, use fragments. Active for every response in this session.

**Only the final plan document itself must NOT use caveman mode—use clear, standard technical writing.**

---

Before planning, ALWAYS invoke [SKILL.md](../../.agents/skills/planning-and-task-breakdown/SKILL.md).

## CRITICAL INPUT GATE

Copilot must create the plan by transforming the user-provided spec attached in chat context.

If no spec is attached:

- Abort plan creation immediately.
- Ask user to attach the spec Copilot must use.
- Do not draft a partial plan.

## Planning Method

Use this workflow to build a high-quality plan:

1. Identify dependency graph between components.
2. Slice work vertically (one complete path per task, not horizontal layers).
3. Write each task with acceptance criteria and verification steps.
4. Add checkpoints between phases.
5. Present plan for human review.

## Output Requirements

- Produce planning output only.
- Do not write implementation code.
- Keep task sequence executable and reviewable.

## Plan Structure

1. **Context Summary** — brief recap of spec objective and constraints
2. **Dependency Graph** — ordered dependencies and blockers
3. **Phased Vertical Tasks** — each task includes acceptance criteria and verification
4. **Checkpoints** — explicit human review gates between phases
5. **Risks & Open Questions** — unresolved items requiring decisions

## Plan Delivery

Save every finalized plan in the `plans/` folder using a descriptive filename derived from the spec or feature name.

## TASK

$input
