# Search Engine — Build Roadmap

Building a search engine from scratch in Go. Every phase ends with a **passing test**, not a
feeling. Tick a box only when `go test ./...` is green.

Rule for each phase: read the *concept*, then write the code yourself. If you're stuck for
more than an hour on the mechanics, ask for a hint about the specific stuck point — not the
whole solution.

---

## Phase 0 — Foundations ✅ DONE

- [x] Load documents from a folder → `internal/document/loader.go`
- [x] Tokenizer: lowercase, strip punctuation, split → `internal/tokenizer/tokenizer.go`
- [x] Stop word removal via a set → `internal/tokenizer/removeWords.go`
- [x] Inverted index `term → []docID` → `internal/index/index.go`
- [x] Forward index `docID → Document` (needed to render results)
- [x] Single shared `Analyze()` so index and query can never drift
- [x] Single-term search → `internal/search/search.go`
- [x] Test suite → `internal/search/search_test.go`

**What you learned:** the five-stage pipeline (load → tokenize → filter → index → query), and
why the analyzer must be identical on both sides.

---

## Phase 1 — Boolean Queries

*Goal: `"corona world"` returns only docs containing **both** words.*

- [ ] **1.1 — AND (intersection)**
  - Change `Search` to look up *every* query term, not just `words[0]`
  - Intersect the postings lists
  - Delete the `t.Skip` line in `TestSearchMultiTermAND` — it should now pass
- [ ] **1.2 — Optimize the intersection**
  - Sort the postings lists shortest-first before intersecting
- [ ] **1.3 — OR (union)**
  - Add `SearchOr`, or a mode flag on `Search`
- [ ] **1.4 — NOT**
  - Exclude docs containing a term (`corona -virus`)

**Key insight:** your postings lists are already **sorted ascending** (you build documents in
ID order), so intersection is a two-pointer walk in O(n+m) — no nested loop, no map. Advance
whichever pointer holds the smaller ID; when they're equal, that's a hit. This one property is
why real engines work so hard to keep postings sorted.

**Gotchas**
- An empty postings list means the whole AND result is empty — bail early.
- Decide now: does a missing term mean "no results" or "ignore that term"? Test it either way.
- Union has to stay sorted too, or Phase 2's ranking gets messy.

**Done when:** `TestSearchMultiTermAND` passes with the skip removed, plus your own OR/NOT cases.

---

## Phase 2 — Ranking (the big one)

*Goal: results come back **best first**, not doc-ID first.*

- [ ] **2.1 — Restructure postings to carry term frequency**
  - `Words map[string][]Posting` where `Posting{DocID int, TF int}`
  - Drop the `seen` map — you're counting occurrences now, not deduping
  - Store `DocLength map[int]int` and the corpus average on the index
  - Fix the compile errors in `search.go` this causes (expect a few)
- [ ] **2.2 — TF-IDF**
  - `score = TF × log(N / DF)` where N = total docs, DF = docs containing the term
  - Sum scores across query terms; sort descending
  - `Search` now returns scored results, not bare `[]int`
- [ ] **2.3 — BM25**
  - Replace TF-IDF. Constants `k1 ≈ 1.2`, `b ≈ 0.75`
- [ ] **2.4 — Top-K with a heap**
  - Use `container/heap` to keep only the best K instead of sorting everything

**Key insight:** TF-IDF has two flaws BM25 fixes. First, a word appearing 100 times isn't 50×
more relevant than one appearing twice — BM25 *saturates* term frequency along a curve that
flattens out. Second, long documents win unfairly just by containing more words — BM25
normalizes by document length against the corpus average. That's what `k1` and `b` control:
`k1` how fast saturation kicks in, `b` how hard length is penalized. Set `b = 0` and length
normalization switches off entirely — try it and watch what breaks.

**Gotchas**
- `log(N/DF)` is 0 when a term is in every document. That's correct — a universal term carries
  no signal — but don't be surprised by zero scores.
- Compute IDF from the index, never per-query in a loop. It's a property of the corpus.
- Return `[]Result{DocID, Score}`, not `[]int`. Your API changes here; that's expected.
- Ties should break deterministically (by doc ID) or tests get flaky.

**Done when:** searching a term that appears 5× in one doc and 1× in another puts the 5× doc
first, and a padded-out long document does *not* outrank a short focused one.

---

## Phase 3 — Better Text Analysis

*Goal: `"running"` finds `"run"`.*

- [ ] **3.1 — Porter stemmer**, written by hand, in `internal/tokenizer/stemmer.go`
- [ ] **3.2 — Wire it into `Analyze()`** (one line — this is the payoff for Phase 0)
- [ ] **3.3 — Unicode normalization** for accented text
- [ ] **3.4 — Expand the stop word list**, or derive it from document frequency

**Key insight:** stemming is lossy and *deliberately wrong* sometimes — `university` and
`universe` both stem to `univers`. You accept bad precision for better recall. Note where it
misfires; that tension is the whole subject of information retrieval.

**Gotchas**
- Stem at index time *and* query time. `Analyze()` already guarantees this — verify it does.
- Your existing tests will break. That's the system working: the index genuinely changed.
- Re-run Phase 2's ranking tests. Stemming merges terms and shifts every DF value.

**Done when:** `"running"`, `"ran"`, and `"runs"` all return the doc containing `"run"`.

---

## Phase 4 — Phrase Search

*Goal: `"lock down"` matches the phrase, not the two words scattered apart.*

- [ ] **4.1 — Store positions**: `Posting{DocID int, Positions []int}` (TF becomes `len(Positions)`)
- [ ] **4.2 — Phrase matching**: intersect docs, then check positions are consecutive
- [ ] **4.3 — Quoted query syntax** so `"lock down"` parses as a phrase
- [ ] **4.4 — Proximity search**: within N words of each other
- [ ] **4.5 — Measure it**: print index size before and after positions

**Key insight:** for a two-word phrase, you need a position `p` in the first term's list where
`p+1` exists in the second term's list. Both lists are sorted, so it's the same two-pointer
walk as Phase 1 — one level down. Recognizing that the same primitive solves both is the point.

**Gotchas**
- Positions must count tokens **after** stop word removal, or offsets shift and phrases break.
  Alternatively, keep stop word *slots* — real engines do this. Decide and document it.
- Index size will jump significantly. That's Phase 5's motivation — feel the pain first.

**Done when:** `"lock down"` matches doc 1, but a doc containing both words far apart does not.

---

## Phase 5 — Persistence & Compression

*Goal: stop rebuilding the index on every run.*

- [ ] **5.1 — Save/load with `encoding/gob`**, `Save(path)` / `Load(path)`
- [ ] **5.2 — Measure**: index build time vs. load time, and file size on disk
- [ ] **5.3 — Delta encoding**: `[3,7,9,20]` → `[3,4,2,11]`
- [ ] **5.4 — Varint encoding**: `encoding/binary`'s `PutUvarint`
- [ ] **5.5 — Measure again**, compare against 5.2
- [ ] **5.6 — Segments**: write immutable segment files, merge them, support incremental adds

**Key insight:** deltas are small numbers, and varint stores small numbers in fewer bytes — 1
byte under 128 instead of 8. The two techniques only work *together*, and only on sorted data.
That's the third time sorted postings have paid off; it's not a coincidence, it's the design.

**Gotchas**
- Delta encoding requires strictly ascending IDs. Assert it — a subtle bug here corrupts silently.
- Segment merging is genuinely hard. It's the LSM-tree pattern behind Lucene. Take your time.
- Write a round-trip test: index → save → load → search returns identical results.

**Done when:** loading a saved index is measurably faster than rebuilding, and search results
are byte-identical before and after a save/load cycle.

---

## Phase 6 — Query Quality

*Goal: it feels like a real search box.*

- [ ] **6.1 — Levenshtein distance** (classic DP — write it yourself)
- [ ] **6.2 — Make it fast**: BK-tree or trigram index, so you don't scan the whole vocabulary
- [ ] **6.3 — "Did you mean?"** suggestions for zero-result queries
- [ ] **6.4 — Autocomplete** via a trie over the vocabulary
- [ ] **6.5 — Snippets**: extract the matching region and highlight the terms

**Key insight:** naive fuzzy search compares the query against every term — O(vocabulary), which
dies at scale. A BK-tree exploits the triangle inequality of edit distance to prune most of the
tree without computing distances at all. Build the slow one first, measure it, then fix it.

**Gotchas**
- Snippets need the original text, not tokens — that's another job your forward index does.
- Cap edit distance at 1 or 2. Beyond that, suggestions become noise.
- Don't fuzzy-match short terms; `cat`/`car`/`can` are all distance 1.

**Done when:** a typo'd query returns a suggestion, and results display with highlighted snippets.

---

## Phase 7 — Serve It

*Goal: it's a service, not a `main()`.*

- [ ] **7.1 — HTTP API**: `GET /search?q=...` returning JSON
- [ ] **7.2 — Concurrency**: `sync.RWMutex` around the index
- [ ] **7.3 — Prove it's safe**: `go test -race`
- [ ] **7.4 — Benchmarks**: `go test -bench` on query latency
- [ ] **7.5 — Pagination**: `?page=` and `?limit=`
- [ ] **7.6 — Optional**: a minimal HTML search page

**Key insight:** searches are reads and can run fully in parallel; only reindexing writes. That
asymmetry is exactly what `RWMutex` is for — any number of concurrent readers, or one writer.

**Gotchas**
- Go maps are **not** safe for concurrent read+write. Without the mutex you get a runtime panic
  under load, not a wrong answer.
- `-race` catches what tests won't. Run it before you believe the locking is right.
- Benchmark with a realistically sized corpus. Three files prove nothing.

**Done when:** `go test -race ./...` is clean and concurrent requests return correct results.

---

## Ongoing habits

- [ ] Add a test with every phase — never tick a box on eyeballed output
- [ ] Commit per sub-phase, so you can bisect when ranking mysteriously changes
- [ ] Record measurements (index size, query latency) before and after each optimization
- [ ] Grow the corpus — 3 files hide bugs that 1000 files expose

## Stretch goals

- [ ] Field-based search (`title:corona`) with per-field weighting
- [ ] Multi-word synonyms
- [ ] Query result caching with an LRU
- [ ] Parallel indexing with goroutines + worker pool
- [ ] A relevance evaluation set: fixed queries with known-correct rankings

---

## Reading, when you want the theory

- *Introduction to Information Retrieval* (Manning, Raghavan, Schütze) — free online, chapters
  1–2 cover Phases 1–2 and chapter 6 covers ranking
- The original BM25 papers by Robertson & Spärck Jones
- Lucene's `index` package, once you reach Phase 5 — you'll recognize the structures
