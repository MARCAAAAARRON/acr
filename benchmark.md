# ACR Model Benchmarks

Informal comparison notes from local testing on RTX 4050 6GB VRAM, via LM Studio.
Same test query used across all models for fair comparison: `--query "what does main.go do"`

## Test Setup
- Repo: ACR itself (self-referential test)
- Retrieval: keyword + filename-match scoring
- Scheduler budget: 500 tokens
- Query: "what does main.go do"

## Results

| Model | Size / Quant | Accuracy | Style | Notes |
|---|---|---|---|---|
| bonsai-27b | 27B, 1-bit | Accurate | Detailed | Slow — 778 reasoning tokens on a simple question. Reasons by default (toggleable via thinking_budget_tokens). 262K context window inherited from Qwen3.6 27B base. |
| deepseek-r1-0528-qwen3-8b | 8B | **Hallucinated** | Fast, structured | Claimed retrieval "uses embeddings" — false, retrieval is keyword/filename only. Don't fully trust for precise code claims without verification. |
| deepreinforce-ai_ornith-1.0-9b | 9B | Accurate | Best formatting — self-organized into a table + example usage | No hallucinations detected. Strong structured output. |
| qwen2.5-coder-7b-instruct | 7B, Q4_K_M | Accurate | Most concise — no fluff | Purpose-built for code. Best fit for precise, quick code-description tasks. |
| glm-4-9b-0414 | 9B, Q3_K_L | Accurate | Verbose but careful | Honestly hedged an uncertain claim ("likely stands for...") rather than stating it as fact. Good epistemic behavior. |

## Early Observations
- Reasoning models (Bonsai, DeepSeek) are not automatically better — DeepSeek was fastest but least trustworthy on factual code claims.
- Model verbosity varies a lot independent of accuracy — Qwen Coder 7B and Ornith 9B were both accurate but styled very differently (terse vs. richly formatted).
- All five models correctly described the pipeline structure; the retrieval/scheduler quality (what content they actually received) mattered more than model choice, once accuracy issues are excluded.
- Filename-match scoring bug (fixed) previously caused main.go itself to be excluded from context — a reminder that upstream retrieval quality can matter more than which model answers.

## Open Questions / TODO
- [ ] Test speed with actual timestamps, not just impressions
- [ ] Test on a real multi-file query (not just "what does X do")
- [ ] Try glm-z1-9b-0414 (reasoning variant of glm-4-9b) for comparison
- [ ] Try Qwen 3.5 4B once downloaded — different context/precision tradeoff
- [ ] Revisit Bonsai with thinking_budget_tokens lowered to test speed vs quality tradeoff

## Token Budget Impact (real repo test)

Tested against nestjs-ratelimiter-gateway (22 files), same query across three budgets:
`"why does the rate limiter use Lua scripts instead of separate Redis calls"`

| Budget | Chunks scheduled | Real prompt tokens | Answer quality |
|---|---|---|---|
| 500 | 2 of 18 | ~600 | Generic — correct concept, no filenames, no bug details |
| 3000 | 3 of 18 | 2,923 | Correct — named `incr-and-expire.lua`, explained the crash-window bug |
| 8000 | 18 of 18 (all matched) | 6,418 | Comprehensive — race conditions, token bucket complexity, performance, distributed consistency |

**Takeaway:** on a real 22-file repo, the original 500-token default was severely starving the model — missing the exact file (`incr-and-expire.lua`) that mattered most for this question. At 8000, every matched chunk fit; there was nothing left to gain from raising the budget further on this repo/query. Default changed to 4000 as a middle ground — well past the "generic answer" zone without assuming every repo maxes out around 6-8K tokens.

This is a concrete demonstration of ACR's core thesis: context management quality drives output quality more directly than model choice does, once accuracy floors are met.

## Retrieval Limitation: Short-File Scoring Bias

Discovered via `--explain` on nestjs-ratelimiter-gateway, budget 500:

`src\ratelimiter\scripts\incr-and-expire.lua` — the file the README names as the
core atomicity fix for this project — scored only 54 (near the bottom of 18
matched chunks) and was cut before the model ever saw it. Meanwhile `README.md`
scored highest (155) but was too large (2,318 tokens) to fit in a 500 budget at all.
Neither the best summary nor the actual implementation made it into context.

Result: the model's answer was accurate but generic — correct on atomicity in the
abstract, but never named `incr-and-expire.lua`, never referenced the specific
crash-window bug, never mentioned token bucket's refill math. Raising the budget
to 3000+ (see Token Budget Impact above) fixed this for this specific query, but
that's a workaround, not a fix to the underlying issue.

**Root cause:** retrieval scores by raw keyword *count*, which structurally
penalizes short files. A 49-token Lua script mentioning "atomic"/"lua"-relevant
terms twice will always lose to a 600-token file mentioning them four times,
regardless of which file is actually more relevant. Same underlying issue as the
earlier `scanner.go` self-naming problem — this is a systemic scoring flaw, not a
one-off bug specific to `.lua` files.

**Rejected fix:** adding a `.lua`-specific score boost. This would only patch this
one repo/file-type combination and wouldn't generalize to the next repo's
equivalent short-but-critical file (a small config file, a short but crucial
utility function, etc).

**Real fix, deferred:** switch from raw keyword *count* to keyword *density*
(matches relative to file length), so short, highly relevant files stop being
structurally penalized against longer, loosely relevant ones. Not yet implemented.