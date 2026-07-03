---
name: test-driven-development
description: Drives development with tests. Use when implementing logic, fixing bugs, or changing behavior. Define behavior first, prove with failing test, implement minimally, then refactor safely.
---

# Test-Driven Development

## Overview

Test-Driven Development (TDD) means behavior-first delivery:

1. Write test for desired behavior.
2. Confirm test fails for right reason.
3. Write minimal implementation to pass.
4. Refactor while keeping tests green.

TDD goal: reduce regression risk, improve design feedback, keep changes verifiable.

## Scope of This Skill

This skill describes TDD process only.

- No language-specific code examples here.
- No framework-specific snippets here.

For Go testing examples, patterns, and tooling, use:

- `.agents/skills/golang-testing/SKILL.md`

## When to Use

- New behavior
- Bug fix
- Refactor touching behavior
- Edge-case hardening
- Any change where proof needed

When not needed:

- Pure docs/content updates
- Non-behavioral config edits

## Core Loop

RED -> GREEN -> REFACTOR

- RED: write focused failing test for single behavior.
- GREEN: make smallest change that passes test.
- REFACTOR: improve structure/naming without behavior change.
- Repeat per behavior slice.

## Red Phase Rules

- Test must fail before implementation.
- Failure reason must match target behavior.
- Test name must describe expected outcome.
- Keep scope narrow: one behavior per test.

## Green Phase Rules

- Implement minimum viable change.
- Avoid speculative architecture.
- Avoid extra features not required by failing test.

## Refactor Phase Rules

- Refactor only with green tests.
- Change one thing at time.
- Re-run relevant tests after each refactor step.

## Bug Fix Pattern (Prove-It)

1. Capture bug as reproducible test.
2. Confirm test fails pre-fix.
3. Implement fix.
4. Confirm bug test passes.
5. Run broader suite to guard regressions.

## Test Quality Principles

- Test outcomes, not implementation internals.
- Prefer deterministic tests.
- Keep tests isolated and order-independent.
- Avoid over-mocking; mock boundaries only.
- Keep tests readable (DAMP over DRY for test intent).

## Anti-Patterns

- Writing code before failing test
- Passing test without meaningful assertion
- Large multi-behavior tests
- Tests coupled to internal call sequence
- Flaky timing/order dependent tests
- Skipping or disabling failing tests instead of fixing

## Verification Checklist

After behavior change:

- Every new behavior has corresponding test.
- Bug fixes include regression test.
- Relevant test suite passes.
- No newly skipped/disabled tests.
- Test names remain behavior-descriptive.

## Go-Specific Follow-Up

After applying this TDD workflow in Go projects, consult:

- `.agents/skills/golang-testing/SKILL.md`

Use that skill for:

- table-driven test structure
- integration test isolation
- fuzzing/benchmark guidance
- race/leak detection and CI commands
