# Specify workflow documentation

This directory contains design documentation generated using the `/specify` workflow with Claude Code.

## Learning project note

This is a learning project to explore spec-kit methodology. The documentation in this directory demonstrates
a structured approach to feature development:

1. **Feature specification** (specs/*/spec.md) - What the feature does and why
1. **Implementation plan** (specs/*/plan.md) - Technical approach and design decisions
1. **Research notes** (specs/*/research.md) - Investigation of patterns and best practices
1. **Data models** (specs/*/data-model.md) - Entity structures and validation rules
1. **Task breakdown** (specs/*/tasks.md) - Step-by-step implementation tasks
1. **Testing scenarios** (specs/*/quickstart.md) - Manual test cases

Review these documents to understand how spec-kit structures software development from initial concept
through implementation.

## Directory structure

```text
.specify/
├── memory/
│   └── constitution.md       # Project principles and constraints
├── templates/                # Templates for spec generation
│   ├── spec-template.md
│   ├── plan-template.md
│   ├── tasks-template.md
│   └── agent-file-template.md
└── scripts/
    └── bash/
        └── check-prerequisites.sh

specs/
└── 001-add-input-validation/  # Example feature
    ├── spec.md
    ├── plan.md
    ├── research.md
    ├── data-model.md
    ├── quickstart.md
    └── tasks.md
```

## Key concepts

- **Constitution**: Non-negotiable architectural principles (pure Go, offline-first, simplicity)
- **TDD approach**: Tests written first, must fail before implementation
- **Parallel execution**: Tasks marked [P] can run independently
- **Validation gates**: Each phase has checkpoints to ensure quality
- **Constitutional compliance**: Every feature verified against project principles

## How this feature was developed

The input validation feature (001-add-input-validation) followed this workflow:

1. `/constitution` - Established project principles
1. `/specify` - Created feature specification from natural language description
1. `/plan` - Generated implementation plan with research and design artifacts
1. `/tasks` - Broke down implementation into 13 sequential tasks
1. Implementation - Followed TDD approach with tests first
1. Validation - All tests pass, constitutional principles maintained

Total time: ~2 hours from concept to working, tested code.

## Further reading

- [Specify workflow documentation](https://docs.anthropic.com/en/docs/claude-code/specify-workflow)
- [Project constitution](memory/constitution.md)
- [Example feature: Input validation](../specs/001-add-input-validation/)
