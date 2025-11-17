# Flashback

Flashback is a command-line knowledge store.
It captures text, URLs, and commands, extracts structured metadata, and makes everything searchable.

It is designed for developers who want a fast, local, scriptable memory system.
![Demo GIF](demo.gif)

---

## Features

* CLI-first workflow (Cobra)
* Metadata extraction for URLs (OpenGraph, Twitter, JSON-LD)
* AI enrichment using Google Gemini (strict JSON schema)
* Local-first storage backed by Turso/libSQL
* TUI viewer built with Bubbletea
* Fast search with filters and tags
* Separate metadata table for incremental enrichment

---

## Usage

```bash
flashback help
```

TUI:

```bash
flashback
```

Add entries:

```bash
flashback add kubectl rollout restart deployment web
flashback add https://blog.bytebytego.com/p/understanding-load-balancers
```

Search:

```bash
flashback search load balancer
flashback search kubernetes
```

View entries:

```bash
flashback list
flashback show <id>
```

---

## How it works

Flashback processes inputs through a simple pipeline:

1. Detect type (text, URL, code)
2. Scrape metadata (OpenGraph, JSON-LD, fallbacks)
3. Run Gemini enrichment (tags, summary, normalization)
4. Store in SQLite/Turso with structured metadata

All AI output is enforced via JSON schema to ensure deterministic results.

---

## Data model

Each record is stored in the `flashbacks` table.
Metadata is stored separately in the `metadata` table as key/value pairs.
This allows multiple enrichment passes and avoids schema churn.

Example metadata fields:

```
tldr
description
tags (JSON array)
image (url)
```

---

## Configuration

Location: `~/.config/flashback/config.yaml`

---

## Install
### Prerequisites
- Go 1.24.4 or later
- Google AI API key (for Gemini embeddings)

### From source:

```bash
git clone https://github.com/yagnikpt/flashback
cd flashback
make install
```

### Direct Go Install
```bash
go install github.com/yagnikpt/flashback@latest
```

---

## Roadmap

* More site-specific extractors (GitHub, YouTube, Reddit, Medium)
* Configurable metadata schema
* Plugin system for custom enrichers
