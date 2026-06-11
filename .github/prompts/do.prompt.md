---
name: "do"
description: "General coding task with caveman mode and strict code quality. Use for any implementation, refactor, or fix task."
---

Follow the caveman communication skill exactly as defined in [SKILL.md](../../.agents/skills/caveman/SKILL.md). All responses must be terse, drop filler, use fragments. Active for every response in this session.

Before implementation, ALWAYS invoke [SKILL.md](../../.agents/skills/clean-code/SKILL.md). This is mandatory for every task.

Before you start, detect ALL the relevant skills for the task exactly as defined in:

- Production engineering skills orchestrator: [SKILL.md](../../.agents/skills/using-agent-skills/SKILL.md). Only use skills relevant to the task.
- Go specific skills orchestrator: [SKILL.md](../../.agents/skills/golang-how-to/SKILL.md). Only use skills relevant to the task.

## STRICT SCOPE CONSTRAINT

Implement ONLY what is explicitly requested in $input. Nothing more.

Rules:

- No scaffolding, wiring, or glue code outside the described scope
- No modifications to files not mentioned in the task
- If something is missing to make the app compile or run — that is intentional. Leave the gap.

Exceptions:

- You can provide stubs, placeholders, or hooks for future work with comments

If the task is partially implementable by design, implement the defined part and stop. Acknowledge what's missing.

## Model Selection

**Do NOT use GPT-5 mini, GPT-5.4 mini, or any other "mini" OpenAI variant** unless the user has explicitly requested one of those models by name or selected it via the model picker.

If Copilot is in Auto mode and classifies the task as simple, prefer **Claude Haiku** or **Gemini Flash** over any OpenAI mini model.

## TASK

$input
