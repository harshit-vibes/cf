# DSA Prep Application - Technical Decision Document

> **Date**: 2024-12-28
> **Status**: Approved

---

## Decision: Use Codeforces API as Primary Data Source

### Context

We evaluated multiple competitive programming platforms for building a DSA preparation application:

| Platform | API Status | Reliability | Data Quality |
|----------|------------|-------------|--------------|
| **Codeforces** | ✅ Official public API | High | Excellent |
| LeetCode | ❌ No official API | Low (third-party) | Good |
| HackerRank | ❌ No public API | N/A | N/A |
| CodeChef | ❌ No public API | N/A | N/A |
| AtCoder | ❌ No official API | Low (scraping) | Good |

### Decision

**Selected: Codeforces API**

### Rationale

1. **Official & Stable**: Only major platform with documented public API
2. **Rich Problem Set**: 8000+ problems with difficulty ratings (800-3500)
3. **Tag System**: Problems tagged by topic (dp, graphs, greedy, etc.)
4. **User Tracking**: Can track user progress, submissions, ratings
5. **No Authentication Required**: Public data accessible anonymously
6. **Free**: No API key required for public endpoints

### Trade-offs

| Pros | Cons |
|------|------|
| Official API support | Rate limited (5 req/sec) |
| Problem difficulty ratings | No interview-style problems |
| Comprehensive problem tags | Different format than LeetCode |
| User progress tracking | No premium problem access |
| Contest history | Learning curve for users |

### Future Considerations

- Add LeetCode via [alfa-leetcode-api](https://github.com/alfaarghya/alfa-leetcode-api) later
- Build local problem cache to reduce API calls
- Consider [Kontests.net](https://kontests.net/api) for contest aggregation

---

## Codeforces API Reference

### Base URL
```
https://codeforces.com/api/{methodName}
```

### Rate Limits
- **5 requests per second** (anonymous)
- JSON response format
- HTTPS required

### Response Format
```json
{
  "status": "OK" | "FAILED",
  "comment": "Error message (if FAILED)",
  "result": { /* method-specific data */ }
}
```

---

## API Methods

### Problem Set

#### `problemset.problems`
Returns all problems from the problem set.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `tags` | string | No | Semicolon-separated list of tags |
| `problemsetName` | string | No | Custom problemset name |

**Example:**
```
GET /api/problemset.problems?tags=dp;greedy
```

**Returns:**
```json
{
  "problems": [Problem],
  "problemStatistics": [ProblemStatistics]
}
```

#### `problemset.recentStatus`
Returns recent submissions.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `count` | integer | Yes | Number of submissions (max 1000) |

---

### User Methods

#### `user.info`
Returns information about users.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `handles` | string | Yes | Semicolon-separated handles (max 10000) |

**Example:**
```
GET /api/user.info?handles=tourist;Petr
```

#### `user.status`
Returns user's submissions.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `handle` | string | Yes | User handle |
| `from` | integer | No | 1-based index of first submission |
| `count` | integer | No | Number of submissions to return |

#### `user.rating`
Returns user's rating history.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `handle` | string | Yes | User handle |

#### `user.ratedList`
Returns list of all rated users.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `activeOnly` | boolean | No | Only users active in last month |

---

### Contest Methods

#### `contest.list`
Returns list of all contests.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `gym` | boolean | No | Include gym contests |

#### `contest.standings`
Returns contest standings.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `contestId` | integer | Yes | Contest ID |
| `from` | integer | No | 1-based starting rank |
| `count` | integer | No | Number of rows |
| `handles` | string | No | Semicolon-separated handles |
| `showUnofficial` | boolean | No | Include virtual participants |

#### `contest.status`
Returns contest submissions.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `contestId` | integer | Yes | Contest ID |
| `handle` | string | No | Filter by user |
| `from` | integer | No | 1-based index |
| `count` | integer | No | Number of submissions |

#### `contest.ratingChanges`
Returns rating changes after contest.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `contestId` | integer | Yes | Contest ID |

#### `contest.hacks`
Returns hacks in a contest.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `contestId` | integer | Yes | Contest ID |

---

## Data Objects

### Problem
```typescript
interface Problem {
  contestId: number;           // Contest ID
  problemsetName?: string;     // Problemset name
  index: string;               // Problem index (A, B, C1, etc.)
  name: string;                // Problem name
  type: "PROGRAMMING" | "QUESTION";
  points?: number;             // Max points for ICPC
  rating?: number;             // Problem difficulty (800-3500)
  tags: string[];              // Topic tags
}
```

### ProblemStatistics
```typescript
interface ProblemStatistics {
  contestId: number;
  index: string;
  solvedCount: number;         // Number of users who solved
}
```

### User
```typescript
interface User {
  handle: string;              // Username
  email?: string;              // Only for authorized requests
  firstName?: string;
  lastName?: string;
  country?: string;
  city?: string;
  organization?: string;
  contribution: number;        // Contribution points
  rank: string;                // e.g., "expert", "master"
  rating: number;              // Current rating
  maxRank: string;             // Highest rank achieved
  maxRating: number;           // Highest rating achieved
  lastOnlineTimeSeconds: number;
  registrationTimeSeconds: number;
  friendOfCount: number;
  avatar: string;              // URL to avatar
  titlePhoto: string;          // URL to title photo
}
```

### Submission
```typescript
interface Submission {
  id: number;
  contestId: number;
  creationTimeSeconds: number;
  relativeTimeSeconds: number;
  problem: Problem;
  author: Party;
  programmingLanguage: string;
  verdict?: Verdict;
  testset: Testset;
  passedTestCount: number;
  timeConsumedMillis: number;
  memoryConsumedBytes: number;
}

type Verdict =
  | "FAILED"
  | "OK"
  | "PARTIAL"
  | "COMPILATION_ERROR"
  | "RUNTIME_ERROR"
  | "WRONG_ANSWER"
  | "PRESENTATION_ERROR"
  | "TIME_LIMIT_EXCEEDED"
  | "MEMORY_LIMIT_EXCEEDED"
  | "IDLENESS_LIMIT_EXCEEDED"
  | "SECURITY_VIOLATED"
  | "CRASHED"
  | "INPUT_PREPARATION_CRASHED"
  | "CHALLENGED"
  | "SKIPPED"
  | "TESTING"
  | "REJECTED";

type Testset =
  | "SAMPLES"
  | "PRETESTS"
  | "TESTS"
  | "CHALLENGES";
```

### Contest
```typescript
interface Contest {
  id: number;
  name: string;
  type: "CF" | "IOI" | "ICPC";
  phase: "BEFORE" | "CODING" | "PENDING_SYSTEM_TEST" | "SYSTEM_TEST" | "FINISHED";
  frozen: boolean;
  durationSeconds: number;
  startTimeSeconds?: number;
  relativeTimeSeconds?: number;
  preparedBy?: string;
  websiteUrl?: string;
  description?: string;
  difficulty?: number;         // 1-5 scale
  kind?: string;
  icpcRegion?: string;
  country?: string;
  city?: string;
  season?: string;
}
```

### RatingChange
```typescript
interface RatingChange {
  contestId: number;
  contestName: string;
  handle: string;
  rank: number;
  ratingUpdateTimeSeconds: number;
  oldRating: number;
  newRating: number;
}
```

### Party
```typescript
interface Party {
  contestId?: number;
  members: Member[];
  participantType: "CONTESTANT" | "PRACTICE" | "VIRTUAL" | "MANAGER" | "OUT_OF_COMPETITION";
  teamId?: number;
  teamName?: string;
  ghost: boolean;
  room?: number;
  startTimeSeconds?: number;
}

interface Member {
  handle: string;
}
```

### RanklistRow
```typescript
interface RanklistRow {
  party: Party;
  rank: number;
  points: number;
  penalty: number;
  successfulHackCount: number;
  unsuccessfulHackCount: number;
  problemResults: ProblemResult[];
}

interface ProblemResult {
  points: number;
  penalty?: number;
  rejectedAttemptCount: number;
  type: "PRELIMINARY" | "FINAL";
  bestSubmissionTimeSeconds?: number;
}
```

---

## References

- [Codeforces API Documentation](https://codeforces.com/apiHelp)
- [Codeforces API Methods](https://codeforces.com/apiHelp/methods)
- [Codeforces API Objects](https://codeforces.com/apiHelp/objects)
- [PublicAPI - Codeforces](https://publicapi.dev/codeforces-api)
- [OpenGenus - Exploring Codeforces API](https://iq.opengenus.org/exploring-codeforces-api/)
