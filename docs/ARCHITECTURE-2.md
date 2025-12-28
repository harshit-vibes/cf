# DSA Prep - Architecture 2 (SSH Server)

> SSH-based DSA practice platform powered by Wish + Bubble Tea
> Users connect via `ssh dsaprep.io` - zero installation required

---

## Table of Contents

1. [Overview](#overview)
2. [System Architecture](#system-architecture)
3. [Infrastructure](#infrastructure)
4. [Tech Stack](#tech-stack)
5. [Database Design](#database-design)
6. [Authentication](#authentication)
7. [Application Structure](#application-structure)
8. [TUI Components](#tui-components)
9. [External APIs](#external-apis)
10. [Data Flow](#data-flow)
11. [Deployment](#deployment)
12. [Security](#security)
13. [Monitoring](#monitoring)

---

## Overview

### What We're Building

A terminal-based DSA practice platform accessible via SSH. Users connect with `ssh dsaprep.io`, authenticate via SSH keys, and get an interactive TUI for practicing coding problems.

### Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Access Method | SSH | Zero installation, works everywhere |
| Auth | SSH Public Keys | No passwords, cryptographically secure |
| TUI Framework | Bubble Tea | Best-in-class terminal UI for Go |
| Server | Wish | SSH server with Bubble Tea integration |
| Database | PostgreSQL | Relational data, self-hosted or RDS |
| Problem Source | Codeforces API | 8000+ problems, free, reliable |
| Hosting | AWS EC2 t3.micro | 12-month free tier, reliable |

---

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              USERS                                      │
│                                                                         │
│    macOS Terminal    Linux Terminal    Windows Terminal    Mobile SSH   │
│         │                  │                  │                │        │
└─────────┴──────────────────┴──────────────────┴────────────────┴────────┘
                                    │
                                    │ SSH (Port 22)
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                              AWS                                        │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────────┐  │
│  │                  EC2 t3.micro (Free Tier)                         │  │
│  │                                                                   │  │
│  │  ┌─────────────────────────────────────────────────────────────┐ │  │
│  │  │                    WISH SSH SERVER                          │ │  │
│  │  │                                                             │ │  │
│  │  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │ │  │
│  │  │  │  Logging    │  │  ActiveTerm │  │  BubbleTea  │         │ │  │
│  │  │  │  Middleware │─▶│  Middleware │─▶│  Middleware │         │ │  │
│  │  │  └─────────────┘  └─────────────┘  └─────────────┘         │ │  │
│  │  │                                           │                 │ │  │
│  │  │                                           ▼                 │ │  │
│  │  │  ┌─────────────────────────────────────────────────────┐   │ │  │
│  │  │  │                 TUI APPLICATION                      │   │ │  │
│  │  │  │                                                      │   │ │  │
│  │  │  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌─────────┐ │   │ │  │
│  │  │  │  │Dashboard │ │ Problems │ │ Practice │ │Settings │ │   │ │  │
│  │  │  │  └──────────┘ └──────────┘ └──────────┘ └─────────┘ │   │ │  │
│  │  │  │                                                      │   │ │  │
│  │  │  └─────────────────────────────────────────────────────┘   │ │  │
│  │  └─────────────────────────────────────────────────────────────┘ │  │
│  │                               │                                   │  │
│  │                               ▼                                   │  │
│  │  ┌─────────────────────────────────────────────────────────────┐ │  │
│  │  │                    PostgreSQL                                │ │  │
│  │  │           (Self-hosted on EC2 or RDS Free Tier)             │ │  │
│  │  └─────────────────────────────────────────────────────────────┘ │  │
│  │                                                                   │  │
│  │                    EBS Storage (30GB Free)                        │  │
│  └───────────────────────────────────────────────────────────────────┘  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    │ HTTPS
                                    ▼
                    ┌───────────────────────────────┐
                    │       CODEFORCES API          │
                    │   codeforces.com/api/         │
                    └───────────────────────────────┘
```

### Component Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           APPLICATION                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  cmd/server/main.go                                                     │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ internal/server/                                                 │   │
│  │ ├── server.go      # Wish server setup                          │   │
│  │ ├── auth.go        # SSH key authentication                     │   │
│  │ └── handler.go     # Session handler                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ internal/tui/                                                    │   │
│  │ ├── app.go         # Main application model                     │   │
│  │ ├── views/                                                       │   │
│  │ │   ├── dashboard.go                                            │   │
│  │ │   ├── problems.go                                             │   │
│  │ │   ├── practice.go                                             │   │
│  │ │   ├── leaderboard.go                                          │   │
│  │ │   ├── settings.go                                             │   │
│  │ │   └── onboarding.go                                           │   │
│  │ ├── components/                                                  │   │
│  │ │   ├── header.go                                               │   │
│  │ │   ├── footer.go                                               │   │
│  │ │   ├── problemcard.go                                          │   │
│  │ │   └── statsbar.go                                             │   │
│  │ └── styles/                                                      │   │
│  │     └── styles.go  # Lip Gloss styles                           │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ internal/service/                                                │   │
│  │ ├── user.go        # User business logic                        │   │
│  │ ├── problem.go     # Problem selection logic                    │   │
│  │ ├── progress.go    # Progress tracking                          │   │
│  │ └── sync.go        # Codeforces sync                            │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ internal/repository/                                             │   │
│  │ ├── user.go        # User CRUD                                  │   │
│  │ ├── problem.go     # Problem queries                            │   │
│  │ ├── progress.go    # Progress storage                           │   │
│  │ └── session.go     # Session management                         │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│       │                                                                 │
│       ▼                                                                 │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ internal/api/                                                    │   │
│  │ └── codeforces.go  # Codeforces API client                      │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Infrastructure

### AWS Free Tier (12 Months)

| Resource | Free Allowance | Usage |
|----------|----------------|-------|
| **EC2 t3.micro** | 750 hours/month | App server (always-on) |
| **EBS Storage** | 30GB gp2/gp3 | System + data |
| **Data Transfer** | 100GB/month | SSH traffic |
| **Elastic IP** | 1 (if attached to running instance) | Static IP |
| **RDS db.t3.micro** | 750 hours/month (optional) | Managed PostgreSQL |

### Why AWS EC2?

| Benefit | Description |
|---------|-------------|
| **Reliable** | Industry-standard cloud infrastructure |
| **Free Tier** | 12 months free, perfect for MVPs |
| **SSH Native** | Direct port 22 access, no proxies |
| **Scalable** | Upgrade instance type when needed |
| **Ecosystem** | Route53, SES, CloudWatch available |

### Architecture on AWS

```
┌─────────────────────────────────────────────────────────────────┐
│                           AWS VPC                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │                  EC2 t3.micro                            │   │
│  │              (2 vCPU, 1GB RAM, Free Tier)                │   │
│  ├─────────────────────────────────────────────────────────┤   │
│  │                                                         │   │
│  │  ┌─────────────────────────────────────────────────┐   │   │
│  │  │              DSA Prep Server                     │   │   │
│  │  │              (Go Binary + systemd)               │   │   │
│  │  │                                                  │   │   │
│  │  │   • Wish SSH Server (Port 22)                   │   │   │
│  │  │   • Bubble Tea TUI                              │   │   │
│  │  │   • Codeforces API Client                       │   │   │
│  │  └─────────────────────────────────────────────────┘   │   │
│  │                         │                               │   │
│  │                         │ localhost:5432                │   │
│  │                         ▼                               │   │
│  │  ┌─────────────────────────────────────────────────┐   │   │
│  │  │              PostgreSQL 16                       │   │   │
│  │  │          (Self-hosted on EC2)                    │   │   │
│  │  │                                                  │   │   │
│  │  │   • Users, Progress, Stats                      │   │   │
│  │  │   • Problems Cache                              │   │   │
│  │  │   • EBS-backed storage                          │   │   │
│  │  └─────────────────────────────────────────────────┘   │   │
│  │                                                         │   │
│  │  EBS Volume: 30GB gp3                                   │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  Security Group:                                                │
│  ├── Inbound: Port 22 (SSH) from 0.0.0.0/0                     │
│  └── Outbound: All traffic                                      │
│                                                                 │
│  Elastic IP: xxx.xxx.xxx.xxx                                    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### AWS Services Used

```
┌─────────────────────────────────────────────────────────────────┐
│                       AWS SERVICES                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  EC2 Instance:                                                  │
│  ├── Type: t3.micro (2 vCPU, 1GB RAM)                          │
│  ├── AMI: Amazon Linux 2023 or Ubuntu 24.04                    │
│  ├── Storage: 30GB gp3 EBS                                     │
│  └── Region: us-east-1 (or closest)                            │
│                                                                 │
│  Networking:                                                    │
│  ├── VPC: Default VPC                                          │
│  ├── Subnet: Public subnet                                     │
│  ├── Security Group: dsa-prep-sg                               │
│  └── Elastic IP: For static address                            │
│                                                                 │
│  Database (Option A - Self-hosted):                            │
│  └── PostgreSQL 16 on EC2 (included in instance)               │
│                                                                 │
│  Database (Option B - RDS Free Tier):                          │
│  ├── Engine: PostgreSQL 16                                     │
│  ├── Instance: db.t3.micro                                     │
│  └── Storage: 20GB gp2                                         │
│                                                                 │
│  Optional:                                                      │
│  ├── Route53: Custom domain                                    │
│  ├── SES: Email notifications                                  │
│  └── CloudWatch: Monitoring & logs                             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### DNS Setup

```
dsaprep.io (or your domain)
    │
    ├── A Record ──► EC2 Elastic IP
    └── (Optional) AAAA Record ──► IPv6 if enabled

# Or use EC2 public DNS:
ec2-xx-xx-xx-xx.compute-1.amazonaws.com
```

### Security Group Rules

```
┌─────────────────────────────────────────────────────────────────┐
│               Security Group: dsa-prep-sg                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Inbound Rules:                                                 │
│  ┌─────────┬──────────┬─────────────┬─────────────────────┐    │
│  │ Type    │ Protocol │ Port        │ Source              │    │
│  ├─────────┼──────────┼─────────────┼─────────────────────┤    │
│  │ SSH     │ TCP      │ 22          │ 0.0.0.0/0           │    │
│  │ SSH     │ TCP      │ 22          │ ::/0 (IPv6)         │    │
│  └─────────┴──────────┴─────────────┴─────────────────────┘    │
│                                                                 │
│  Outbound Rules:                                                │
│  ┌─────────┬──────────┬─────────────┬─────────────────────┐    │
│  │ Type    │ Protocol │ Port        │ Destination         │    │
│  ├─────────┼──────────┼─────────────┼─────────────────────┤    │
│  │ All     │ All      │ All         │ 0.0.0.0/0           │    │
│  └─────────┴──────────┴─────────────┴─────────────────────┘    │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### AWS CLI Commands

```bash
# Configure AWS CLI
aws configure

# Create Security Group
aws ec2 create-security-group \
    --group-name dsa-prep-sg \
    --description "DSA Prep SSH access"

# Add SSH inbound rule
aws ec2 authorize-security-group-ingress \
    --group-name dsa-prep-sg \
    --protocol tcp \
    --port 22 \
    --cidr 0.0.0.0/0

# Create Key Pair
aws ec2 create-key-pair \
    --key-name dsa-prep-key \
    --query 'KeyMaterial' \
    --output text > ~/.ssh/dsa-prep-key.pem

chmod 400 ~/.ssh/dsa-prep-key.pem

# Launch EC2 Instance (Amazon Linux 2023)
aws ec2 run-instances \
    --image-id ami-0c7217cdde317cfec \
    --instance-type t3.micro \
    --key-name dsa-prep-key \
    --security-groups dsa-prep-sg \
    --block-device-mappings '[{"DeviceName":"/dev/xvda","Ebs":{"VolumeSize":30,"VolumeType":"gp3"}}]' \
    --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=dsa-prep}]'

# Allocate Elastic IP
aws ec2 allocate-address --domain vpc

# Associate Elastic IP (replace with actual values)
aws ec2 associate-address \
    --instance-id i-xxxxxxxxxx \
    --allocation-id eipalloc-xxxxxxxxxx

# SSH into instance
ssh -i ~/.ssh/dsa-prep-key.pem ec2-user@<elastic-ip>
```

---

## Tech Stack

### Core Technologies

| Layer | Technology | Version | Purpose |
|-------|------------|---------|---------|
| **Language** | Go | 1.22+ | Backend |
| **SSH Server** | Wish | latest | SSH handling |
| **TUI Framework** | Bubble Tea | latest | Terminal UI |
| **Components** | Bubbles | latest | UI widgets |
| **Styling** | Lip Gloss | latest | Terminal styling |
| **Forms** | Huh | latest | User input |
| **Markdown** | Glamour | latest | Render problems |
| **Logging** | Charm Log | latest | Structured logs |
| **Database** | PostgreSQL | 16 | Data storage |
| **Cache** | Redis | 7 (optional) | Session cache |

### Go Dependencies

```go
// go.mod
module github.com/yourusername/dsa-prep

go 1.22

require (
    // Charm stack
    github.com/charmbracelet/wish      v1.4.0
    github.com/charmbracelet/bubbletea v0.26.0
    github.com/charmbracelet/bubbles   v0.18.0
    github.com/charmbracelet/lipgloss  v0.11.0
    github.com/charmbracelet/huh       v0.4.0
    github.com/charmbracelet/glamour   v0.7.0
    github.com/charmbracelet/log       v0.4.0
    github.com/charmbracelet/ssh       v0.0.0

    // Database
    github.com/jackc/pgx/v5            v5.6.0
    github.com/redis/go-redis/v9       v9.5.0

    // Utilities
    github.com/joho/godotenv           v1.5.1
    github.com/google/uuid             v1.6.0
)
```

---

## Database Design

### Entity Relationship Diagram

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│     users       │       │   user_keys     │       │  user_progress  │
├─────────────────┤       ├─────────────────┤       ├─────────────────┤
│ id (PK)         │──┐    │ id (PK)         │       │ id (PK)         │
│ cf_handle       │  │    │ user_id (FK)    │───┐   │ user_id (FK)    │───┐
│ display_name    │  │    │ fingerprint     │   │   │ problem_id      │   │
│ created_at      │  │    │ public_key      │   │   │ status          │   │
│ updated_at      │  │    │ created_at      │   │   │ attempts        │   │
│ last_seen_at    │  │    │ last_used_at    │   │   │ time_spent      │   │
│ settings (JSON) │  │    └─────────────────┘   │   │ solved_at       │   │
└─────────────────┘  │                          │   │ created_at      │   │
         │           │                          │   └─────────────────┘   │
         │           └──────────────────────────┴─────────────────────────┘
         │
         │           ┌─────────────────┐       ┌─────────────────┐
         │           │    problems     │       │     tags        │
         │           ├─────────────────┤       ├─────────────────┤
         │           │ id (PK)         │───┐   │ id (PK)         │
         │           │ cf_contest_id   │   │   │ name            │
         │           │ cf_index        │   │   │ slug            │
         │           │ name            │   │   └─────────────────┘
         │           │ rating          │   │            │
         │           │ solved_count    │   │            │
         │           │ url             │   │   ┌─────────────────┐
         │           │ created_at      │   │   │  problem_tags   │
         │           │ updated_at      │   │   ├─────────────────┤
         │           └─────────────────┘   └───│ problem_id (FK) │
         │                                     │ tag_id (FK)     │
         │                                     └─────────────────┘
         │
         │           ┌─────────────────┐       ┌─────────────────┐
         │           │    sessions     │       │  daily_stats    │
         │           ├─────────────────┤       ├─────────────────┤
         └───────────│ id (PK)         │       │ id (PK)         │
                     │ user_id (FK)    │       │ user_id (FK)    │
                     │ started_at      │       │ date            │
                     │ ended_at        │       │ problems_solved │
                     │ problems_viewed │       │ time_spent      │
                     │ problems_solved │       │ streak_day      │
                     └─────────────────┘       └─────────────────┘
```

### Table Schemas

```sql
-- migrations/001_init.sql

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cf_handle VARCHAR(50) UNIQUE,
    display_name VARCHAR(100),
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_seen_at TIMESTAMPTZ DEFAULT NOW()
);

-- SSH public keys (users can have multiple)
CREATE TABLE user_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    fingerprint VARCHAR(100) UNIQUE NOT NULL,
    public_key TEXT NOT NULL,
    name VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_used_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_keys_fingerprint ON user_keys(fingerprint);

-- Problems (cached from Codeforces)
CREATE TABLE problems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cf_contest_id INTEGER NOT NULL,
    cf_index VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    rating INTEGER,
    solved_count INTEGER DEFAULT 0,
    url VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(cf_contest_id, cf_index)
);

CREATE INDEX idx_problems_rating ON problems(rating);

-- Tags
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) UNIQUE NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL
);

-- Problem-Tag junction
CREATE TABLE problem_tags (
    problem_id UUID REFERENCES problems(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (problem_id, tag_id)
);

-- User progress per problem
CREATE TABLE user_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    problem_id UUID REFERENCES problems(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'unseen', -- unseen, attempted, solved, skipped
    attempts INTEGER DEFAULT 0,
    time_spent INTEGER DEFAULT 0, -- seconds
    solved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, problem_id)
);

CREATE INDEX idx_user_progress_user ON user_progress(user_id);
CREATE INDEX idx_user_progress_status ON user_progress(status);

-- Sessions (for analytics)
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    problems_viewed INTEGER DEFAULT 0,
    problems_solved INTEGER DEFAULT 0
);

-- Daily stats (for streaks)
CREATE TABLE daily_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    problems_solved INTEGER DEFAULT 0,
    time_spent INTEGER DEFAULT 0, -- seconds
    streak_day INTEGER DEFAULT 0,
    UNIQUE(user_id, date)
);

CREATE INDEX idx_daily_stats_user_date ON daily_stats(user_id, date);

-- User settings JSON structure
COMMENT ON COLUMN users.settings IS '{
    "difficulty": {
        "min": 800,
        "max": 1400
    },
    "focusTags": ["dp", "graphs"],
    "dailyGoal": 5,
    "notifications": {
        "email": true,
        "streakReminder": true
    }
}';
```

### Data Types

```go
// internal/domain/user.go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID          uuid.UUID         `json:"id"`
    CFHandle    *string           `json:"cf_handle"`
    DisplayName *string           `json:"display_name"`
    Settings    UserSettings      `json:"settings"`
    CreatedAt   time.Time         `json:"created_at"`
    UpdatedAt   time.Time         `json:"updated_at"`
    LastSeenAt  time.Time         `json:"last_seen_at"`
}

type UserSettings struct {
    Difficulty    DifficultyRange `json:"difficulty"`
    FocusTags     []string        `json:"focus_tags"`
    DailyGoal     int             `json:"daily_goal"`
    Notifications NotificationSettings `json:"notifications"`
}

type DifficultyRange struct {
    Min int `json:"min"`
    Max int `json:"max"`
}

type NotificationSettings struct {
    Email          bool `json:"email"`
    StreakReminder bool `json:"streak_reminder"`
}

// internal/domain/problem.go
type Problem struct {
    ID          uuid.UUID  `json:"id"`
    CFContestID int        `json:"cf_contest_id"`
    CFIndex     string     `json:"cf_index"`
    Name        string     `json:"name"`
    Rating      *int       `json:"rating"`
    SolvedCount int        `json:"solved_count"`
    URL         string     `json:"url"`
    Tags        []Tag      `json:"tags"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type Tag struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
    Slug string    `json:"slug"`
}

// internal/domain/progress.go
type ProgressStatus string

const (
    StatusUnseen    ProgressStatus = "unseen"
    StatusAttempted ProgressStatus = "attempted"
    StatusSolved    ProgressStatus = "solved"
    StatusSkipped   ProgressStatus = "skipped"
)

type UserProgress struct {
    ID        uuid.UUID      `json:"id"`
    UserID    uuid.UUID      `json:"user_id"`
    ProblemID uuid.UUID      `json:"problem_id"`
    Status    ProgressStatus `json:"status"`
    Attempts  int            `json:"attempts"`
    TimeSpent int            `json:"time_spent"` // seconds
    SolvedAt  *time.Time     `json:"solved_at"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
}
```

---

## Authentication

### SSH Key Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      AUTHENTICATION FLOW                                │
└─────────────────────────────────────────────────────────────────────────┘

1. User connects
   ┌──────────┐                         ┌──────────────┐
   │  User    │ ── ssh dsaprep.io ───▶  │  Wish Server │
   │ Terminal │                         │              │
   └──────────┘                         └──────┬───────┘
                                               │
2. SSH handshake                               │
   ┌──────────────────────────────────────────┐│
   │ Client offers public key                 ││
   │ Server challenges with random data       ││◄───┐
   │ Client signs with private key            ││    │
   │ Server verifies signature                ││────┘
   └──────────────────────────────────────────┘│
                                               │
3. Key lookup                                  ▼
   ┌───────────────────────────────────────────────────────────────────┐
   │                                                                   │
   │   fingerprint := ssh.FingerprintSHA256(publicKey)                │
   │                          │                                        │
   │                          ▼                                        │
   │   ┌─────────────────────────────────────────────────────────┐    │
   │   │              SELECT * FROM user_keys                     │    │
   │   │              WHERE fingerprint = ?                       │    │
   │   └─────────────────────────────────────────────────────────┘    │
   │                          │                                        │
   │              ┌───────────┴───────────┐                           │
   │              ▼                       ▼                           │
   │        ┌─────────┐             ┌─────────┐                       │
   │        │  Found  │             │Not Found│                       │
   │        └────┬────┘             └────┬────┘                       │
   │             │                       │                            │
   │             ▼                       ▼                            │
   │      Load user from DB       Create new user                     │
   │      Update last_used_at     Store key + user                    │
   │                                                                   │
   └───────────────────────────────────────────────────────────────────┘
                                               │
4. Session created                             ▼
   ┌───────────────────────────────────────────────────────────────────┐
   │  ctx.SetValue("user", user)                                       │
   │  ctx.SetValue("session", session)                                 │
   │  Launch Bubble Tea TUI with user context                          │
   └───────────────────────────────────────────────────────────────────┘
```

### Auth Implementation

```go
// internal/server/auth.go
package server

import (
    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
)

type AuthHandler struct {
    userRepo repository.UserRepository
}

func (h *AuthHandler) PublicKeyHandler(ctx ssh.Context, key ssh.PublicKey) bool {
    fingerprint := ssh.FingerprintSHA256(key)

    // Look up existing key
    userKey, err := h.userRepo.GetKeyByFingerprint(ctx, fingerprint)
    if err != nil && !errors.Is(err, repository.ErrNotFound) {
        log.Error("Failed to lookup key", "error", err)
        return false
    }

    var user *domain.User

    if userKey != nil {
        // Existing user
        user, err = h.userRepo.GetByID(ctx, userKey.UserID)
        if err != nil {
            log.Error("Failed to get user", "error", err)
            return false
        }

        // Update last used
        h.userRepo.UpdateKeyLastUsed(ctx, userKey.ID)
    } else {
        // New user - create account
        user, err = h.userRepo.Create(ctx, &domain.User{
            Settings: domain.DefaultSettings(),
        })
        if err != nil {
            log.Error("Failed to create user", "error", err)
            return false
        }

        // Store the key
        _, err = h.userRepo.AddKey(ctx, user.ID, fingerprint, string(key.Marshal()))
        if err != nil {
            log.Error("Failed to store key", "error", err)
            return false
        }

        log.Info("New user created", "user_id", user.ID)
    }

    // Store user in context
    ctx.SetValue("user", user)

    log.Info("User authenticated",
        "user_id", user.ID,
        "fingerprint", fingerprint[:16]+"...",
    )

    return true
}
```

---

## Application Structure

### Project Layout

```
dsa-prep/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
│
├── internal/
│   ├── server/
│   │   ├── server.go               # Wish server setup
│   │   ├── auth.go                 # SSH authentication
│   │   ├── handler.go              # Session handler
│   │   └── middleware.go           # Custom middleware
│   │
│   ├── tui/
│   │   ├── app.go                  # Main Bubble Tea model
│   │   ├── keymap.go               # Key bindings
│   │   ├── messages.go             # Custom messages
│   │   │
│   │   ├── views/
│   │   │   ├── dashboard.go        # Home dashboard
│   │   │   ├── problems.go         # Problem browser
│   │   │   ├── practice.go         # Practice session
│   │   │   ├── leaderboard.go      # Rankings
│   │   │   ├── stats.go            # User statistics
│   │   │   ├── settings.go         # User settings
│   │   │   └── onboarding.go       # First-time setup
│   │   │
│   │   ├── components/
│   │   │   ├── header.go           # App header
│   │   │   ├── footer.go           # Help/nav footer
│   │   │   ├── problemcard.go      # Problem display
│   │   │   ├── progressbar.go      # Progress indicator
│   │   │   └── timer.go            # Practice timer
│   │   │
│   │   └── styles/
│   │       └── styles.go           # Lip Gloss styles
│   │
│   ├── service/
│   │   ├── user.go                 # User business logic
│   │   ├── problem.go              # Problem selection
│   │   ├── progress.go             # Progress tracking
│   │   ├── stats.go                # Statistics
│   │   └── sync.go                 # Codeforces sync
│   │
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user.go             # User queries
│   │   │   ├── problem.go          # Problem queries
│   │   │   ├── progress.go         # Progress queries
│   │   │   └── stats.go            # Stats queries
│   │   └── interfaces.go           # Repository interfaces
│   │
│   ├── api/
│   │   ├── codeforces/
│   │   │   ├── client.go           # HTTP client
│   │   │   ├── types.go            # Response types
│   │   │   └── problems.go         # Problem fetching
│   │   └── cache.go                # Response caching
│   │
│   └── domain/
│       ├── user.go                 # User entity
│       ├── problem.go              # Problem entity
│       ├── progress.go             # Progress entity
│       └── errors.go               # Domain errors
│
├── migrations/
│   ├── 001_init.sql
│   └── 002_add_indexes.sql
│
├── scripts/
│   ├── setup.sh                    # Server setup script
│   ├── deploy.sh                   # Deployment script
│   └── sync_problems.sh            # Problem sync cron
│
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
└── go.sum
```

---

## TUI Components

### View Hierarchy

```
App (main model)
│
├── Header (persistent)
│   ├── Logo/Title
│   ├── User Info (handle, rating)
│   └── Streak indicator
│
├── Content (switches based on active view)
│   │
│   ├── Dashboard View
│   │   ├── Today's Challenge card
│   │   ├── Weekly Progress bar
│   │   ├── Weak/Strong Topics
│   │   └── Recent Activity list
│   │
│   ├── Problems View
│   │   ├── Search input
│   │   ├── Filter controls
│   │   ├── Problem table
│   │   └── Pagination
│   │
│   ├── Practice View
│   │   ├── Timer
│   │   ├── Problem description (viewport)
│   │   └── Action buttons
│   │
│   ├── Leaderboard View
│   │   ├── Tab bar (Daily/Weekly/All-time)
│   │   └── Rankings table
│   │
│   ├── Stats View
│   │   ├── Solved by difficulty chart
│   │   ├── Solved by tag chart
│   │   └── Rating history
│   │
│   └── Settings View
│       ├── Codeforces handle input
│       ├── Difficulty range
│       ├── Focus tags multi-select
│       └── Daily goal
│
└── Footer (persistent)
    ├── Navigation keys
    └── Context help
```

### Screen Mockups

```
┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Practice                                    tourist │ 1847 │ 🔥 12 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─ Dashboard ─────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │  Today's Challenge                                              │   │
│  │  ┌───────────────────────────────────────────────────────────┐ │   │
│  │  │  1850A - Two Sum                                          │ │   │
│  │  │  Rating: 800  │  Tags: math, implementation               │ │   │
│  │  │                                              [Enter] Start │ │   │
│  │  └───────────────────────────────────────────────────────────┘ │   │
│  │                                                                 │   │
│  │  Weekly Progress                                                │   │
│  │  ████████████████░░░░░░░░ 16/25 problems (64%)                 │   │
│  │                                                                 │   │
│  │  ┌─ Improve ─────────────┐  ┌─ Strong ──────────────┐          │   │
│  │  │  • DP (23% solved)    │  │  • Greedy (78%)       │          │   │
│  │  │  • Graphs (31%)       │  │  • Implementation (82%)│          │   │
│  │  │  • Trees (29%)        │  │  • Math (75%)         │          │   │
│  │  └───────────────────────┘  └───────────────────────┘          │   │
│  │                                                                 │   │
│  │  Recent Activity                                                │   │
│  │  ✓ Binary Search Basics      1200   2 hours ago                │   │
│  │  ✓ Greedy Algorithm          1100   5 hours ago                │   │
│  │  ✗ DP on Trees               1600   yesterday                  │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [d]ashboard  [p]roblems  [l]eaderboard  [s]tats  [?]help  [q]uit      │
└─────────────────────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  DSA Practice                                    tourist │ 1847 │ 🔥 12 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌─ Problems ──────────────────────────────────────────────────────┐   │
│  │                                                                 │   │
│  │  Search: ________________    [f]ilter: 800-1400 │ dp, greedy   │   │
│  │                                                                 │   │
│  │  ┌───────┬────────────────────────────┬────────┬─────────┬───┐ │   │
│  │  │  ID   │  Name                      │ Rating │ Tags    │ ✓ │ │   │
│  │  ├───────┼────────────────────────────┼────────┼─────────┼───┤ │   │
│  │  │►1850A │ Two Sum                    │   800  │ math    │   │ │   │
│  │  │ 1851B │ Colored Segments           │  1000  │ greedy  │ ✓ │ │   │
│  │  │ 1852C │ Binary Search              │  1200  │ binary  │ ✓ │ │   │
│  │  │ 1853A │ Maximum Subarray           │  1100  │ dp      │   │ │   │
│  │  │ 1854D │ Graph Traversal            │  1400  │ graphs  │   │ │   │
│  │  │ 1855B │ Tree Diameter              │  1300  │ trees   │ ✓ │ │   │
│  │  │ 1856A │ Palindrome Check           │   800  │ strings │ ✓ │ │   │
│  │  │ 1857C │ Knapsack Problem           │  1500  │ dp      │   │ │   │
│  │  │ ...   │ ...                        │  ...   │ ...     │   │ │   │
│  │  └───────┴────────────────────────────┴────────┴─────────┴───┘ │   │
│  │                                                                 │   │
│  │  Page 1 of 156                                   4521 problems │   │
│  │                                                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [/]search  [f]ilter  [r]andom  [Enter]solve  [←→]navigate  [q]back    │
└─────────────────────────────────────────────────────────────────────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  Practice Session                                           ⏱️  08:34   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1850A - Two Sum                                                        │
│  Rating: 800  │  Tags: math, implementation                             │
│  ───────────────────────────────────────────────────────────────────── │
│                                                                         │
│  You are given an array of n integers a₁, a₂, ..., aₙ and an           │
│  integer target. Find two indices i and j (1 ≤ i < j ≤ n) such         │
│  that aᵢ + aⱼ = target.                                                 │
│                                                                         │
│  Input                                                                  │
│  The first line contains two integers n and target                      │
│  (2 ≤ n ≤ 10⁵, 1 ≤ target ≤ 2·10⁹).                                    │
│  The second line contains n integers a₁, a₂, ..., aₙ.                  │
│                                                                         │
│  Output                                                                 │
│  Print two indices i and j. If there are multiple answers,             │
│  print any of them.                                                     │
│                                                                         │
│  Example                                                                │
│  ┌─────────────────────┬─────────────────────┐                         │
│  │ Input               │ Output              │                         │
│  ├─────────────────────┼─────────────────────┤                         │
│  │ 4 9                 │ 1 2                 │                         │
│  │ 2 7 11 15           │                     │                         │
│  └─────────────────────┴─────────────────────┘                         │
│                                                                         │
│  ▼ Scroll for more                                                     │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  [o]pen in browser  [✓]solved  [s]kip  [h]int  [q]uit                  │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## External APIs

### Codeforces API Integration

```
Base URL: https://codeforces.com/api/
Rate Limit: 5 requests/second
```

#### Endpoints Used

| Endpoint | Purpose | Frequency |
|----------|---------|-----------|
| `problemset.problems` | Get all problems | Daily sync |
| `user.info` | Validate CF handle | On user setup |
| `user.status` | Get user submissions | On demand |
| `user.rating` | Get rating history | On demand |

#### Sync Strategy

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      PROBLEM SYNC WORKFLOW                              │
└─────────────────────────────────────────────────────────────────────────┘

Daily Cron (2:00 AM UTC)
         │
         ▼
┌─────────────────────┐
│ Fetch problemset    │
│ /api/problemset.    │
│      problems       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐     ┌─────────────────────┐
│ Parse response      │────▶│ problems: []        │
│ Extract problems    │     │ problemStatistics:[]│
└──────────┬──────────┘     └─────────────────────┘
           │
           ▼
┌─────────────────────┐
│ For each problem:   │
│ - Upsert to DB      │
│ - Update tags       │
│ - Update stats      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Log sync results    │
│ - New problems: X   │
│ - Updated: Y        │
│ - Total: Z          │
└─────────────────────┘
```

#### API Client

```go
// internal/api/codeforces/client.go
package codeforces

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "golang.org/x/time/rate"
)

type Client struct {
    baseURL    string
    httpClient *http.Client
    limiter    *rate.Limiter
}

func NewClient() *Client {
    return &Client{
        baseURL: "https://codeforces.com/api",
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        limiter: rate.NewLimiter(rate.Every(200*time.Millisecond), 1), // 5 req/sec
    }
}

func (c *Client) GetProblems(ctx context.Context) (*ProblemsResponse, error) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    resp, err := c.httpClient.Get(c.baseURL + "/problemset.problems")
    if err != nil {
        return nil, fmt.Errorf("fetch problems: %w", err)
    }
    defer resp.Body.Close()

    var result ProblemsResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("decode response: %w", err)
    }

    if result.Status != "OK" {
        return nil, fmt.Errorf("API error: %s", result.Comment)
    }

    return &result, nil
}
```

---

## Data Flow

### User Session Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         USER SESSION FLOW                               │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────┐         ┌──────────────┐         ┌──────────────┐
│   User   │──SSH───▶│  Wish Server │────────▶│   Auth       │
└──────────┘         └──────────────┘         │   Handler    │
                                              └──────┬───────┘
                                                     │
                     ┌───────────────────────────────┘
                     ▼
              ┌─────────────┐
              │ User exists?│
              └──────┬──────┘
                     │
         ┌───────────┴───────────┐
         ▼                       ▼
   ┌───────────┐           ┌───────────┐
   │  Yes      │           │  No       │
   └─────┬─────┘           └─────┬─────┘
         │                       │
         ▼                       ▼
   ┌───────────┐           ┌───────────┐
   │Load User  │           │Create User│
   │Load Prefs │           │Show       │
   └─────┬─────┘           │Onboarding │
         │                 └─────┬─────┘
         │                       │
         └───────────┬───────────┘
                     │
                     ▼
              ┌─────────────┐
              │ Create      │
              │ Session     │
              │ Record      │
              └──────┬──────┘
                     │
                     ▼
              ┌─────────────┐
              │ Launch      │
              │ Bubble Tea  │
              │ TUI         │
              └──────┬──────┘
                     │
    ┌────────────────┼────────────────┐
    ▼                ▼                ▼
┌────────┐    ┌───────────┐    ┌───────────┐
│Dashboard│   │ Problems  │    │ Practice  │
└────────┘    └───────────┘    └───────────┘
    │                │                │
    └────────────────┼────────────────┘
                     │
                     ▼
              ┌─────────────┐
              │ User quits  │
              └──────┬──────┘
                     │
                     ▼
              ┌─────────────┐
              │ End session │
              │ Save stats  │
              │ Close conn  │
              └─────────────┘
```

### Problem Solving Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      PROBLEM SOLVING FLOW                               │
└─────────────────────────────────────────────────────────────────────────┘

┌───────────────┐
│ Select Problem│
│ from list     │
└───────┬───────┘
        │
        ▼
┌───────────────┐         ┌───────────────┐
│ Create/Update │────────▶│ user_progress │
│ Progress      │         │ status=       │
│ Record        │         │ "attempted"   │
└───────┬───────┘         └───────────────┘
        │
        ▼
┌───────────────┐
│ Show Practice │
│ View          │
│ Start Timer   │
└───────┬───────┘
        │
        │◄──────────────────────────────────┐
        │                                   │
        ▼                                   │
┌───────────────┐                          │
│ User action:  │                          │
│ - Continue    │──────────────────────────┘
│ - Mark Solved │──────┐
│ - Skip        │──────┼───┐
│ - Quit        │──────┼───┼───┐
└───────────────┘      │   │   │
                       │   │   │
        ┌──────────────┘   │   │
        ▼                  │   │
┌───────────────┐         │   │
│ Update status │         │   │
│ = "solved"    │         │   │
│ Record time   │         │   │
│ Update streak │         │   │
└───────┬───────┘         │   │
        │                 │   │
        │   ┌─────────────┘   │
        │   ▼                 │
        │ ┌───────────────┐   │
        │ │ Update status │   │
        │ │ = "skipped"   │   │
        │ └───────┬───────┘   │
        │         │           │
        └────┬────┘           │
             │                │
             ▼                │
      ┌─────────────┐         │
      │ Show next   │         │
      │ problem?    │         │
      └──────┬──────┘         │
             │                │
   ┌─────────┴─────────┐      │
   ▼                   ▼      │
┌──────┐          ┌────────┐  │
│ Yes  │─────────▶│  No    │◄─┘
└──────┘          └────┬───┘
                       │
                       ▼
                ┌─────────────┐
                │ Return to   │
                │ Problems or │
                │ Dashboard   │
                └─────────────┘
```

---

## Deployment

### EC2 Server Setup Script

```bash
#!/bin/bash
# scripts/setup.sh - Run on EC2 instance after launch

set -e

echo "=== DSA Prep Server Setup ==="

# Update system
sudo dnf update -y  # Amazon Linux 2023
# sudo apt update && sudo apt upgrade -y  # Ubuntu

# Install dependencies
sudo dnf install -y git golang postgresql16-server  # Amazon Linux
# sudo apt install -y git golang postgresql  # Ubuntu

# Initialize PostgreSQL
sudo postgresql-setup --initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Create database and user
sudo -u postgres psql <<EOF
CREATE USER dsaprep WITH PASSWORD 'your-secure-password';
CREATE DATABASE dsaprep OWNER dsaprep;
GRANT ALL PRIVILEGES ON DATABASE dsaprep TO dsaprep;
EOF

# Configure PostgreSQL for local connections
sudo sed -i "s/#listen_addresses = 'localhost'/listen_addresses = 'localhost'/" /var/lib/pgsql/data/postgresql.conf

# Create app directory
sudo mkdir -p /opt/dsaprep
sudo mkdir -p /opt/dsaprep/.ssh
sudo chown -R ec2-user:ec2-user /opt/dsaprep

# Generate SSH host key for Wish server
ssh-keygen -t ed25519 -f /opt/dsaprep/.ssh/host_key -N ""

# Clone and build application
cd /opt/dsaprep
git clone https://github.com/yourusername/dsa-prep.git .
go build -o server ./cmd/server

# Create environment file
cat > /opt/dsaprep/.env <<EOF
DATABASE_URL=postgres://dsaprep:your-secure-password@localhost:5432/dsaprep?sslmode=disable
SSH_PORT=22
SSH_HOST_KEY_PATH=/opt/dsaprep/.ssh/host_key
EOF

# Run migrations
./server migrate

echo "=== Setup Complete ==="
```

### Systemd Service

```ini
# /etc/systemd/system/dsaprep.service
[Unit]
Description=DSA Prep SSH Server
After=network.target postgresql.service
Wants=postgresql.service

[Service]
Type=simple
User=root
Group=root
WorkingDirectory=/opt/dsaprep
EnvironmentFile=/opt/dsaprep/.env
ExecStart=/opt/dsaprep/server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/dsaprep

[Install]
WantedBy=multi-user.target
```

```bash
# Install and enable service
sudo cp dsaprep.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable dsaprep
sudo systemctl start dsaprep

# Check status
sudo systemctl status dsaprep

# View logs
sudo journalctl -u dsaprep -f
```

### Dockerfile (for local development)

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata openssh-keygen

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations

# Create data directory
RUN mkdir -p /data/.ssh

EXPOSE 22

CMD ["./server"]
```

### Deploy Script

```bash
#!/bin/bash
# scripts/deploy.sh - Deploy to EC2

set -e

EC2_HOST="${EC2_HOST:-your-elastic-ip}"
EC2_USER="${EC2_USER:-ec2-user}"
EC2_KEY="${EC2_KEY:-~/.ssh/dsa-prep-key.pem}"

echo "=== Building ==="
GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

echo "=== Uploading ==="
scp -i $EC2_KEY bin/server $EC2_USER@$EC2_HOST:/opt/dsaprep/server.new
scp -i $EC2_KEY -r migrations $EC2_USER@$EC2_HOST:/opt/dsaprep/

echo "=== Deploying ==="
ssh -i $EC2_KEY $EC2_USER@$EC2_HOST << 'EOF'
    sudo systemctl stop dsaprep
    mv /opt/dsaprep/server.new /opt/dsaprep/server
    chmod +x /opt/dsaprep/server
    /opt/dsaprep/server migrate
    sudo systemctl start dsaprep
    sudo systemctl status dsaprep
EOF

echo "=== Done ==="
```

### Makefile

```makefile
# Makefile
.PHONY: dev build deploy logs status migrate ssh

# Variables
EC2_HOST ?= your-elastic-ip
EC2_USER ?= ec2-user
EC2_KEY ?= ~/.ssh/dsa-prep-key.pem

# Local development
dev:
	go run ./cmd/server

# Build for Linux
build:
	GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

# Build for local
build-local:
	go build -o bin/server ./cmd/server

# Deploy to EC2
deploy: build
	./scripts/deploy.sh

# View logs
logs:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "sudo journalctl -u dsaprep -f"

# Check status
status:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "sudo systemctl status dsaprep"

# Run migrations
migrate:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "/opt/dsaprep/server migrate"

# SSH into server
ssh:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST)

# Restart service
restart:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "sudo systemctl restart dsaprep"

# Sync problems from Codeforces
sync:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "/opt/dsaprep/server sync"

# Database backup
backup:
	ssh -i $(EC2_KEY) $(EC2_USER)@$(EC2_HOST) "pg_dump -U dsaprep dsaprep" > backup_$(shell date +%Y%m%d).sql
```

### Local Development with Docker

```yaml
# docker-compose.yml (for local dev)
version: '3.8'

services:
  app:
    build: .
    ports:
      - "2222:22"
    environment:
      - DATABASE_URL=postgres://dsaprep:password@db:5432/dsaprep?sslmode=disable
      - SSH_PORT=22
      - SSH_HOST_KEY_PATH=/data/.ssh/host_key
    volumes:
      - ./data:/data
    depends_on:
      - db

  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=dsaprep
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=dsaprep
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data

volumes:
  pgdata:
```

```bash
# Local development
docker-compose up -d db          # Start Postgres
go run ./cmd/server              # Run server locally

# Or full Docker setup
docker-compose up --build
ssh -p 2222 localhost            # Connect locally
```

### CI/CD with GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy to EC2

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'

      - name: Build
        run: |
          GOOS=linux GOARCH=amd64 go build -o bin/server ./cmd/server

      - name: Deploy to EC2
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.EC2_HOST }}
          username: ec2-user
          key: ${{ secrets.EC2_SSH_KEY }}
          script: |
            sudo systemctl stop dsaprep || true

      - name: Copy binary
        uses: appleboy/scp-action@v0.1.7
        with:
          host: ${{ secrets.EC2_HOST }}
          username: ec2-user
          key: ${{ secrets.EC2_SSH_KEY }}
          source: "bin/server,migrations/"
          target: "/opt/dsaprep/"
          strip_components: 1

      - name: Start service
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.EC2_HOST }}
          username: ec2-user
          key: ${{ secrets.EC2_SSH_KEY }}
          script: |
            chmod +x /opt/dsaprep/server
            /opt/dsaprep/server migrate
            sudo systemctl start dsaprep
```

### GitHub Secrets Required

```
EC2_HOST: Your EC2 Elastic IP address
EC2_SSH_KEY: Contents of your dsa-prep-key.pem file
```

### Cron Jobs (on EC2)

```bash
# /etc/cron.d/dsaprep
# Sync problems from Codeforces daily at 2 AM UTC
0 2 * * * root /opt/dsaprep/server sync >> /var/log/dsaprep-sync.log 2>&1

# Database backup weekly
0 3 * * 0 root pg_dump -U dsaprep dsaprep | gzip > /opt/dsaprep/backups/backup_$(date +\%Y\%m\%d).sql.gz
```

---

## Security

### Security Measures

| Layer | Measure | Implementation |
|-------|---------|----------------|
| **Network** | Firewall | Only port 22 open |
| **Auth** | SSH Keys | No password auth |
| **Data** | Encryption at rest | PostgreSQL encryption |
| **Secrets** | Environment vars | No secrets in code |
| **Rate Limit** | Per-user limits | Middleware |
| **Input** | Sanitization | Validate all inputs |

### SSH Hardening

```go
// Only allow public key auth, no passwords
wish.WithPasswordAuth(nil),  // Disabled

// Use strong host key
wish.WithHostKeyPath("/opt/dsaprep/.ssh/host_key"),

// Idle timeout
wish.WithIdleTimeout(30 * time.Minute),

// Max concurrent connections per IP
// (implement in middleware)
```

### AWS Security Group Rules

```
Inbound Rules (dsa-prep-sg):
┌─────────────┬──────────┬─────────────┬─────────────────────┐
│ Type        │ Protocol │ Port        │ Source              │
├─────────────┼──────────┼─────────────┼─────────────────────┤
│ SSH         │ TCP      │ 22          │ 0.0.0.0/0           │
│ SSH         │ TCP      │ 22          │ ::/0 (IPv6)         │
└─────────────┴──────────┴─────────────┴─────────────────────┘

Outbound Rules:
┌─────────────┬──────────┬─────────────┬─────────────────────┐
│ Type        │ Protocol │ Port        │ Destination         │
├─────────────┼──────────┼─────────────┼─────────────────────┤
│ All traffic │ All      │ All         │ 0.0.0.0/0           │
└─────────────┴──────────┴─────────────┴─────────────────────┘
```

---

## Monitoring

### Logging

```go
// Structured logging with charm/log
log.Info("User connected",
    "user_id", user.ID,
    "handle", user.CFHandle,
    "ip", session.RemoteAddr(),
)

log.Info("Problem solved",
    "user_id", user.ID,
    "problem_id", problem.ID,
    "time_taken", duration,
)

log.Error("Database error",
    "operation", "get_user",
    "error", err,
)
```

### Metrics to Track

| Metric | Type | Description |
|--------|------|-------------|
| `active_sessions` | Gauge | Current connected users |
| `total_connections` | Counter | All-time connections |
| `problems_solved` | Counter | Total problems solved |
| `session_duration` | Histogram | Session length distribution |
| `auth_failures` | Counter | Failed auth attempts |
| `api_latency` | Histogram | Codeforces API response time |

### Health Checks

```go
// Health check endpoint (internal)
func healthCheck() bool {
    // Check DB connection
    if err := db.Ping(); err != nil {
        return false
    }

    // Check Redis (if used)
    if err := redis.Ping(); err != nil {
        return false
    }

    return true
}
```

### Alerting (Future)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         ALERTING PIPELINE                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  App Logs ──▶ journald ──▶ Loki (optional) ──▶ Grafana ──▶ Alerts      │
│                                                                         │
│  Metrics  ──▶ Prometheus (optional) ──────────▶ Grafana ──▶ Alerts     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Next Steps

1. **Set up AWS EC2** - Launch t3.micro instance
2. **Initialize Go project** - Create project structure
3. **Implement core SSH server** - Wish setup + auth
4. **Build basic TUI** - Dashboard + problem list
5. **Database setup** - PostgreSQL + migrations
6. **Codeforces sync** - Fetch and cache problems
7. **User progress tracking** - Solve/skip/stats
8. **Deploy** - Ship to AWS EC2

---

## References

- [Wish Documentation](https://github.com/charmbracelet/wish)
- [Bubble Tea Guide](https://github.com/charmbracelet/bubbletea)
- [AWS Free Tier](https://aws.amazon.com/free/)
- [EC2 User Guide](https://docs.aws.amazon.com/ec2/)
- [Codeforces API](https://codeforces.com/apiHelp)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
