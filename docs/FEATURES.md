# DSA Prep - Feature List

> Features for the SSH-based DSA practice platform

---

## Access Method

```bash
ssh dsaprep.io
```

Users connect via SSH, authenticate with their SSH key, and get an interactive TUI.

---

## Core Features

### 1. Problem Bank

| Feature | API Method | Priority |
|---------|------------|----------|
| Browse all problems | `problemset.problems` | P0 |
| Filter by difficulty (800-3500) | `problemset.problems` | P0 |
| Filter by tags (dp, graphs, etc.) | `problemset.problems?tags=` | P0 |
| View solve count per problem | `problemset.problems` (statistics) | P1 |
| Search problems by name | Server-side filter | P1 |
| Fuzzy search | Server-side | P2 |

**Available Tags:**
```
implementation, math, greedy, dp, data structures,
brute force, constructive algorithms, graphs, sortings,
binary search, dfs and similar, trees, strings,
number theory, combinatorics, geometry, bitmasks,
two pointers, dsu, shortest paths, probabilities,
divide and conquer, hashing, games, matrices, flows,
string suffix structures, expression parsing, fft,
2-sat, meet-in-the-middle, ternary search,
interactive, schedules, chinese remainder theorem
```

---

### 2. User Profile & Progress Tracking

| Feature | Implementation | Priority |
|---------|----------------|----------|
| SSH key authentication | Automatic on connect | P0 |
| Link Codeforces handle | Onboarding form | P0 |
| Display current rating & rank | `user.info` API | P0 |
| Track solved problems | PostgreSQL | P0 |
| Track attempted problems | PostgreSQL | P0 |
| Time spent per problem | Server-side timer | P1 |
| Submission sync from CF | `user.status` API | P1 |

---

### 3. Practice Sessions

| Feature | Implementation | Priority |
|---------|----------------|----------|
| Start practice session | TUI view | P0 |
| Problem timer | Bubble Tea stopwatch | P0 |
| Mark as solved | Database update | P0 |
| Skip problem | Database update | P0 |
| Open in browser link | Display URL | P0 |
| View problem description | Glamour markdown | P0 |
| Practice history | PostgreSQL queries | P1 |

---

### 4. Daily Challenge & Streaks

| Feature | Implementation | Priority |
|---------|----------------|----------|
| Daily challenge | Random problem selection | P1 |
| Streak tracking | `daily_stats` table | P1 |
| Streak reminders | Future: email via SES | P3 |
| Weekly goals | User settings | P2 |

---

### 5. Study Plans / Problem Lists

| Feature | Implementation | Priority |
|---------|----------------|----------|
| Curated lists by topic | Admin-defined lists | P1 |
| Difficulty ladders | Auto-generated | P1 |
| Custom problem lists | User-created | P2 |
| Recommended next problem | Algorithm based on weak areas | P2 |

**Difficulty Ladders:**

| Level | Rating Range | Target Audience |
|-------|--------------|-----------------|
| Beginner | 800-1000 | New to CP |
| Easy | 1000-1200 | Learning basics |
| Medium | 1200-1400 | Intermediate |
| Hard | 1400-1700 | Advanced |
| Expert | 1700-2000 | Competitive |
| Master | 2000+ | Elite |

---

### 6. Analytics Dashboard

| Feature | Data Source | Priority |
|---------|-------------|----------|
| Problems solved by difficulty | `user_progress` table | P1 |
| Problems solved by tag | `user_progress` + `problem_tags` | P1 |
| Weak topics identification | Comparison algorithm | P1 |
| Strong topics | Comparison algorithm | P1 |
| Solving speed trends | `user_progress.time_spent` | P2 |
| Weekly/monthly stats | `daily_stats` aggregation | P1 |
| Rating history graph | `user.rating` API | P2 |

---

### 7. Leaderboard & Social

| Feature | Implementation | Priority |
|---------|----------------|----------|
| Global leaderboard | All users ranked | P2 |
| Daily leaderboard | Problems solved today | P2 |
| Weekly leaderboard | Problems solved this week | P2 |
| Compare with friends | Side-by-side stats | P3 |
| Achievements/badges | Computed milestones | P3 |

**Leaderboard Metrics:**

| Metric | Description |
|--------|-------------|
| Problems Solved | Total count |
| Current Streak | Days in a row |
| Weekly Average | Problems/week |
| Difficulty Score | Weighted by rating |

---

### 8. Settings & Preferences

| Setting | Type | Priority |
|---------|------|----------|
| Codeforces handle | Text input | P0 |
| Difficulty range (min/max) | Range selector | P0 |
| Focus tags | Multi-select | P0 |
| Daily goal | Number input | P1 |
| Display name | Text input | P1 |

---

## TUI Views

### View Map

```
┌─────────────────────────────────────────────────────────────────┐
│                         APP VIEWS                               │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐                                               │
│  │ Onboarding  │ ◄── First-time users only                     │
│  └──────┬──────┘                                               │
│         │                                                       │
│         ▼                                                       │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────┐       │
│  │  Dashboard  │◄───▶│  Problems   │◄───▶│  Practice   │       │
│  └──────┬──────┘     └─────────────┘     └─────────────┘       │
│         │                                                       │
│         │            ┌─────────────┐     ┌─────────────┐       │
│         └───────────▶│ Leaderboard │     │   Stats     │       │
│                      └─────────────┘     └─────────────┘       │
│                                                                 │
│                      ┌─────────────┐                           │
│                      │  Settings   │                           │
│                      └─────────────┘                           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Keyboard Navigation

| Key | Action |
|-----|--------|
| `d` | Go to Dashboard |
| `p` | Go to Problems |
| `l` | Go to Leaderboard |
| `s` | Go to Stats |
| `,` | Go to Settings |
| `?` | Show help |
| `q` | Quit / Go back |
| `Enter` | Select / Confirm |
| `↑↓` | Navigate lists |
| `/` | Search |
| `f` | Filter |

---

## MVP Scope

### Phase 1: Core (Week 1-2)

| Feature | Component |
|---------|-----------|
| SSH server | Wish setup |
| User auth | SSH key → DB |
| Onboarding | Huh form |
| Dashboard | Basic stats |
| Problem list | Table with filters |
| Practice view | Timer + description |

### Phase 2: Features (Week 3-4)

| Feature | Component |
|---------|-----------|
| Progress tracking | Solve/skip/stats |
| Codeforces sync | API integration |
| Daily challenge | Random selection |
| Streaks | Daily tracking |
| Leaderboard | Rankings |

### Phase 3: Polish (Week 5+)

| Feature | Component |
|---------|-----------|
| Recommendations | Smart suggestions |
| Achievements | Badges system |
| Study plans | Curated lists |
| Advanced stats | Graphs, trends |

---

## Data Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        DATA SOURCES                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  Codeforces API                    PostgreSQL                   │
│  ┌─────────────┐                  ┌─────────────┐              │
│  │ Problems    │──── Sync ───────▶│ problems    │              │
│  │ User Info   │     (daily)      │ tags        │              │
│  │ Submissions │                  │ users       │              │
│  └─────────────┘                  │ user_keys   │              │
│                                   │ user_progress│              │
│                                   │ daily_stats │              │
│                                   │ sessions    │              │
│                                   └─────────────┘              │
│                                          │                      │
│                                          ▼                      │
│                                   ┌─────────────┐              │
│                                   │  TUI Views  │              │
│                                   └─────────────┘              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## API Usage

### Codeforces Endpoints

| Endpoint | Usage | Frequency |
|----------|-------|-----------|
| `problemset.problems` | Sync all problems | Daily cron |
| `user.info` | Validate CF handle | On setup |
| `user.status` | Sync submissions | On demand |
| `user.rating` | Rating history | On demand |
| `contest.list` | Upcoming contests | Hourly |

### Rate Limiting

```
Codeforces: 5 requests/second max

Strategy:
- Daily problem sync (1 request)
- User actions rate-limited per session
- Cache responses in PostgreSQL
- 200ms minimum between API calls
```

---

## Success Metrics

| Metric | Target |
|--------|--------|
| Daily active users | Track growth |
| Problems solved/day | Engagement |
| Average session time | 15+ minutes |
| Streak retention | 7+ days |
| User satisfaction | Feedback |

---

## Future Ideas

- [ ] Multiplayer problem races
- [ ] Team challenges
- [ ] Problem discussions
- [ ] AI-powered hints
- [ ] Code submission (via Codeforces)
- [ ] Virtual contests
- [ ] Mobile SSH client recommendations
