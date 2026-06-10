# Claude Code custom commands

This directory contains custom slash commands for the /specify workflow.

These commands implement a structured approach to feature development:

- `/constitution` - Create/update project architectural principles
- `/specify` - Generate feature specification from natural language
- `/clarify` - Ask targeted questions about underspecified areas
- `/plan` - Create implementation plan with research and design
- `/tasks` - Generate actionable task breakdown
- `/analyze` - Perform cross-artifact consistency analysis
- `/implement` - Execute tasks from tasks.md

## Usage

These commands are automatically available in Claude Code when working in this project.

Example workflow:

1. `/constitution` - Establish project principles once
1. `/specify "Add validation for email and callsign fields"` - Create spec
1. `/clarify` - If needed, resolve ambiguities
1. `/plan` - Generate implementation plan
1. `/tasks` - Break down into actionable tasks
1. `/implement` - Execute the tasks
1. `/analyze` - Verify consistency across artifacts

## Files

- `commands/*.md` - Command definitions
- `settings.local.json` - Local configuration for commands

See `.specify/README.md` for documentation on how these commands were used in this project.
