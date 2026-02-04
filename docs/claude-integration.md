# Claude Integration Guide

This guide explains how to integrate **Qurio** with **Claude** (Claude Code, Claude Desktop, or any Claude-based agent) to create a powerful AI assistant with persistent memory and context awareness.

---

## Table of Contents
1. [Overview](#overview)
2. [Architecture: Markdown-First + Qurio MCP](#architecture-markdown-first--qurio-mcp)
3. [Setup Instructions](#setup-instructions)
4. [The Memory System](#the-memory-system)
5. [Ready-to-Use System Prompt](#ready-to-use-system-prompt)

---

## Overview

Claude, by default, has no memory between sessions. Each conversation starts fresh. This integration solves that by combining:

1. **Local Markdown Files** - Fast, lightweight files for daily context and personality
2. **Qurio MCP** - A searchable knowledge base for detailed documentation and long-term context

This creates a **two-tier memory system**:
- **Tier 1 (Markdown):** Immediate context - read at the start of every session
- **Tier 2 (Qurio):** Deep context - queried on-demand when more detail is needed

---

## Architecture: Markdown-First + Qurio MCP

```
┌─────────────────────────────────────────────────────────────────┐
│                        Claude Agent                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│   1. START SESSION                                               │
│      ↓                                                           │
│   ┌─────────────────────────────────────────────────────┐       │
│   │  Read Local Markdown Files (Tier 1 - Always)        │       │
│   │  • SOUL.md     → Agent personality & approach       │       │
│   │  • USER.md     → User preferences & context         │       │
│   │  • MEMORY.md   → Durable facts & learnings          │       │
│   │  • DAILY.md    → Today's session log                │       │
│   └─────────────────────────────────────────────────────┘       │
│      ↓                                                           │
│   2. PROCESS USER REQUEST                                        │
│      ↓                                                           │
│   ┌─────────────────────────────────────────────────────┐       │
│   │  Need more context? Missing information?            │       │
│   │                    ↓                                │       │
│   │  Query Qurio MCP (Tier 2 - On-Demand)               │       │
│   │  • qurio_search  → Find relevant documentation      │       │
│   │  • qurio_read_page → Get full page content          │       │
│   │  • qurio_ingest  → Save important context           │       │
│   └─────────────────────────────────────────────────────┘       │
│      ↓                                                           │
│   3. END SESSION                                                 │
│      ↓                                                           │
│   ┌─────────────────────────────────────────────────────┐       │
│   │  Update Markdown Files                              │       │
│   │  • Append to DAILY.md (session log)                 │       │
│   │  • Update MEMORY.md (new learnings)                 │       │
│   │  • Optionally ingest to Qurio (long-term storage)   │       │
│   └─────────────────────────────────────────────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Why Markdown-First?

| Aspect | Markdown Files | Qurio MCP |
|--------|---------------|-----------|
| **Speed** | Instant (local file read) | Requires API call |
| **Scope** | Current project/user context | All indexed knowledge |
| **Use Case** | Daily logs, preferences, personality | Documentation, archives, detailed context |
| **Persistence** | Git-tracked, portable | Database-backed |

---

## Setup Instructions

### Step 1: Configure Qurio MCP

Add Qurio to your Claude MCP settings:

```json
{
  "mcpServers": {
    "qurio": {
      "httpUrl": "http://localhost:8081/mcp"
    }
  }
}
```

### Step 2: Create the Memory Directory

In your project root (or a dedicated location), create the following structure:

```
.claude/
├── SOUL.md      # Agent personality and approach
├── USER.md      # User preferences and context
├── MEMORY.md    # Durable facts and learnings
└── logs/
    └── DAILY.md # Today's session log (or date-based files)
```

### Step 3: Initialize the Memory Files

#### `.claude/SOUL.md` - Agent Personality

```markdown
# SOUL - Agent Identity

## Core Personality
- I am a thoughtful, precise coding assistant
- I prefer clarity over brevity
- I ask clarifying questions before making assumptions
- I celebrate small wins with the user

## Communication Style
- Use markdown formatting for all responses
- Break down complex tasks into steps
- Provide context for decisions
- Be honest about limitations

## Technical Preferences
- Prefer TypeScript over JavaScript
- Use modern ES6+ syntax
- Write tests for critical functionality
- Document non-obvious code

## Memory Protocol
1. At session start: Read USER.md, MEMORY.md, and today's DAILY.md
2. During session: Query Qurio when I need external documentation
3. At session end: Update DAILY.md with summary, update MEMORY.md with learnings
```

#### `.claude/USER.md` - User Context

```markdown
# USER - About the Human

## Identity
- Name: [Your Name]
- Role: [Your Role]
- Timezone: [Your Timezone]

## Preferences
- Editor: [VS Code / Cursor / etc.]
- Shell: [zsh / bash / fish]
- Preferred languages: [TypeScript, Go, Python, etc.]

## Current Projects
- [Project Name]: [Brief description]

## Communication Preferences
- [e.g., "Be concise", "Explain in detail", "Use analogies"]

## Known Context
- [Any persistent context Claude should always know]
```

#### `.claude/MEMORY.md` - Durable Memory

```markdown
# MEMORY - Durable Facts & Learnings

## Project Knowledge
<!-- Facts about the codebase that persist across sessions -->

## User Preferences Learned
<!-- Preferences discovered through interaction -->

## Important Decisions
<!-- Architectural decisions, tech choices, etc. -->

## Common Patterns
<!-- Patterns the user frequently uses -->

---
*Last updated: [Date]*
```

#### `.claude/logs/DAILY.md` - Session Log

```markdown
# Daily Log - [Date]

## Session 1 - [Time]
### Context
- [What the user was working on]

### Accomplished
- [Tasks completed]

### Decisions Made
- [Key decisions]

### Open Items
- [Things to continue later]

### Learnings
- [New information to remember]
```

---

## The Memory System

### Daily Logs (`DAILY.md`)

**Purpose:** Track what happened in each session for continuity.

**When to update:**
- At the end of each session
- When switching major contexts
- When making important decisions

**What to log:**
- Tasks accomplished
- Decisions made with rationale
- Open items / TODOs
- New learnings about the project or user

### Durable Memory (`MEMORY.md`)

**Purpose:** Store facts that should persist indefinitely.

**When to update:**
- When learning something new about the project architecture
- When user expresses a strong preference
- When a decision is made that affects future work

**What to store:**
- Project structure insights
- User preferences
- Architectural decisions
- Common patterns and conventions

### User Profile (`USER.md`)

**Purpose:** Understand who you're working with.

**When to update:**
- When user shares personal context
- When preferences change
- When starting new projects

### Agent Soul (`SOUL.md`)

**Purpose:** Define how you behave and communicate.

**When to update:**
- When refining your approach based on feedback
- When adding new capabilities
- When adjusting communication style

---

## Ready-to-Use System Prompt

Copy and paste this into your Claude CLAUDE.md or system prompt:

````markdown
# Claude Agent - Memory-Enhanced Assistant

## Identity
You are a memory-enhanced Claude assistant. You maintain context across sessions using a two-tier memory system.

## Memory System

### Tier 1: Local Markdown (Always Read First)
At the start of EVERY session, read these files from `.claude/`:
1. `SOUL.md` - Your personality and approach
2. `USER.md` - User preferences and context
3. `MEMORY.md` - Durable facts and learnings
4. `logs/DAILY.md` - Recent session history

### Tier 2: Qurio MCP (Query When Needed)
Use Qurio tools when you need:
- External documentation not in markdown files
- Detailed context about specific topics
- Historical information beyond recent logs

Available Qurio tools:
- `qurio_search` - Find relevant documentation
- `qurio_list_sources` - See available knowledge sources
- `qurio_read_page` - Get full page content
- `qurio_ingest` - Save important context for future sessions

## Session Protocol

### On Session Start
1. Read all Tier 1 markdown files
2. Greet user with awareness of recent context
3. Reference any open items from previous sessions

### During Session
1. Use markdown files as primary context
2. Query Qurio only when:
   - User asks about external documentation
   - You need specific technical details
   - Context is missing from markdown files
3. Prefer local knowledge over external queries

### On Session End (When User Says Goodbye)
1. Update `logs/DAILY.md` with session summary
2. Update `MEMORY.md` if new durable facts were learned
3. Optionally use `qurio_ingest` for important context that should be searchable

## File Management

### Creating Memory Files
If `.claude/` directory doesn't exist, offer to create it:
```bash
mkdir -p .claude/logs
touch .claude/SOUL.md .claude/USER.md .claude/MEMORY.md .claude/logs/DAILY.md
```

### Updating Files
- Use file editing tools to update markdown files
- Keep files concise but comprehensive
- Use consistent formatting
- Include timestamps for logs

## Response Guidelines
1. Always acknowledge context from memory files
2. Reference recent sessions when relevant
3. Ask before making assumptions
4. Celebrate progress with the user
5. Be honest about what you remember vs. what you're inferring

## Example Behaviors

### Good Start
"Good morning! I see from your daily log that yesterday you were working on the authentication module. You mentioned wanting to add OAuth support today. Should we continue with that?"

### Using Qurio Appropriately
"I remember you're using NextAuth from our memory, but I need to check the latest configuration options. Let me search Qurio for the NextAuth documentation..."

### Session Wrap-up
"Great session! I'll update the daily log with what we accomplished:
- Implemented OAuth with Google provider
- Fixed the session persistence bug
- TODO: Add GitHub OAuth tomorrow

I'll also note in MEMORY.md that you prefer the JWT strategy over database sessions."
````

---

## Tips for Best Results

1. **Initialize Once:** Set up the `.claude/` directory at the start of each project
2. **Review Periodically:** Clean up old daily logs and consolidate into MEMORY.md
3. **Use Qurio for Archives:** Ingest important daily logs into Qurio for long-term searchability
4. **Customize SOUL.md:** Tune the agent personality to match your working style
5. **Keep USER.md Current:** Update as your preferences evolve

---

## Troubleshooting

### Claude Doesn't Read Memory Files
- Ensure the files exist in `.claude/`
- Check that file paths are correct
- Verify Claude has file system access

### Qurio Queries Fail
- Ensure Qurio is running (`docker-compose up -d`)
- Check that MCP is configured correctly
- Verify the endpoint is accessible at `localhost:8081/mcp`

### Memory Gets Stale
- Set reminders to review MEMORY.md weekly
- Archive old daily logs monthly
- Use `qurio_ingest` for important context

---

<p align="center">
  <strong>Happy Coding with Persistent Memory! 🧠</strong>
</p>
