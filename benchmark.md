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