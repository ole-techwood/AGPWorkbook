---
name: "spec"
description: "Create a technical specification from a feature idea or requirement. Always applies project React rules first, interviews the user if vague, refines if partial, then produces a structured spec via spec-driven-development. Saves output to docs/specs/."
argument-hint: "Describe the feature or requirement to spec out"
agent: "Tech Writer"
---

Follow the caveman communication skill exactly as defined in [SKILL.md](../../.agents/skills/caveman-ultra/). All responses must be terse, drop all filler, use fragments. This mode is active for **every response in this session**. The final spec document itself must NOT use caveman mode — use clear, standard technical writing for the document only.

## Skill Selection Logic

Assess the input before doing anything else.

**VAGUE** — input is a single sentence, buzzword soup, or lacks measurable outcomes (e.g., "The UI shall load and display the current active configuration"):
→ Execute [interview-me SKILL.md](../../.agents/skills/interview-me/SKILL.md) first — one question at a time until ~95% confidence.
→ Then execute [idea-refine SKILL.md](../../.agents/skills/idea-refine/SKILL.md) if gaps remain.
→ Then execute [spec-driven-development SKILL.md](../../.agents/skills/spec-driven-development/SKILL.md).

**PARTIAL** — input has some detail but ambiguous acceptance criteria, missing edge cases, or unclear scope:
→ Execute [idea-refine SKILL.md](../../.agents/skills/idea-refine/SKILL.md) first.
→ Then execute [spec-driven-development SKILL.md](../../.agents/skills/spec-driven-development/SKILL.md).

**DETAILED** — input has clear objective, measurable criteria, constraints, and scope:
→ Execute [spec-driven-development SKILL.md](../../.agents/skills/spec-driven-development/SKILL.md) immediately.

**Never skip [spec-driven-development SKILL.md](../../.agents/skills/spec-driven-development/SKILL.md). It is always the final step.**

---

## Clarifying Questions

Begin by understanding what you want to build. Ask clarifying questions about:

1. **Objective and target users** — What problem does this solve? Who uses it?
2. **Core features and acceptance criteria** — What must it do? How do you know it's complete?
3. **Tech stack preferences and constraints** — Language, framework, database, deployment targets?
4. **Known boundaries** — What to always do, what requires user approval, what never to do?
5. **Existing context** — Related code, prior work, dependencies?
6. **Timeline and scope** — MVP vs. full feature? Phases?

## Spec Structure

The specification must cover these six core areas:

1. **Objective** — What, why, and success criteria
2. **Scope** — What's in, what's out, phases if applicable
3. **Architecture & Design** — Tech stack, key components, data flow
4. **Commands & Operations** — How to build, test, run, deploy
5. **Code Style & Conventions** — Naming, structure, testing expectations
6. **Boundaries** — What the team commits to, asks about, and never does

Surface assumptions immediately before writing spec content. List them clearly for user review.

## Spec Delivery

Once reviewed and approved by the user, save the final specification to the `docs/specs/` folder with a descriptive filename matching the project or feature name.

## TASK

$input
