---
name: "do"
description: "General coding task with caveman mode and strict code quality. Use for any implementation, refactor, or fix task."
agent: agent
---

Follow the caveman communication skill exactly as defined in [SKILL.md](../../.agents/skills/caveman/SKILL.md). All responses must be terse, drop filler, use fragments. Active for every response in this session.

## Skills Preload Flow (Deterministic)

Before implementation, load skills in this exact order:

1. [SKILL.md](../../.agents/skills/clean-code/SKILL.md)
2. [SKILL.md](../../.agents/skills/using-agent-skills/SKILL.md)
3. [SKILL.md](../../.agents/skills/golang-how-to/SKILL.md)

Rules:

- Step 1 is always mandatory.
- Step 2 and step 3 are mandatory skill orchestrators for relevant-skill detection.
- After step 2 and step 3, load only task-relevant skills they select.
- Do not start implementation before preload and detection complete.

## Loaded Skills Announcement (Required)

After loading relevant skills using:

- [SKILL.md](../../.agents/skills/using-agent-skills/SKILL.md)
- [SKILL.md](../../.agents/skills/golang-how-to/SKILL.md)

print this exact lead-in before implementation:

```markdown
Required skills loaded. Relevant set is:

- `[skill-name-1]`
- `[skill-name-2]`
- `[skill-name-3]`
```

Each loaded skill must appear as one bullet item wrapped in backticks.

## Model Selection

**Do NOT use GPT-5 mini, GPT-5.4 mini, or any other "mini" OpenAI variant** unless the user has explicitly requested one of those models by name or selected it via the model picker.

If Copilot is in Auto mode and classifies the task as simple, prefer **Claude Haiku** or **Gemini Flash** over any OpenAI mini model.

## TASK

$input
