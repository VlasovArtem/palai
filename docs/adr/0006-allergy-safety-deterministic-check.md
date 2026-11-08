# 0006. Allergy safety: deterministic allergen check on every generated recipe

- Status: Proposed
- Date: 2026-11-08

## Context

Users declare allergies (strict) and preferences (soft). An LLM can ignore instructions,
hallucinate, or miss that an ingredient is derived from an allergen (butter, ghee, whey and
casein are dairy; soy sauce contains wheat; pesto often contains pine nuts). Serving a recipe
that contains a declared allergen is a safety failure, not a quality issue. The check must be
explainable, testable and reproducible.

## Decision

We will run a deterministic allergen check in Go on every generated recipe before it is
returned or stored. The LLM is never the authority on allergen safety.

- An allergen taxonomy in Go maps each allergen category (at minimum the EU 14 / US major
  allergens: milk/dairy, egg, peanut, tree nuts, wheat/gluten, soy, fish, shellfish, sesame,
  plus others as needed) to known ingredients and derivatives, for example
  `butter → dairy`, `ghee → dairy`, `mayonnaise → egg`, `tahini → sesame`.
- Ingredient names are normalized (case, plurals, punctuation, common qualifiers such as
  "unsalted" or "melted") before matching. Matching works on terms and phrases, not raw
  substrings ("butternut squash" is not dairy, "peanut butter" is peanut).
- Fail closed. If any ingredient matches a declared allergen, the recipe is rejected. We may
  regenerate a bounded number of times with the violation as feedback, then return an error.
  Unknown or ambiguous ingredients for an allergic user are flagged, not silently passed.
- The check runs on the structured ingredient list. Steps are also scanned for added
  ingredients ("finish with a knob of butter").
- Preferences are soft. They influence the prompt and ranking but never cause a safety
  rejection.
- The taxonomy and matcher have exhaustive table-driven tests, including derivative,
  false-positive and false-negative cases. A regression test is added for every reported miss.
- Responses include which allergens were checked. We do not claim certification; the API
  documentation states the check's limits.

## Consequences

- Safety does not depend on model behaviour and is fully unit-testable.
- The taxonomy needs ongoing curation, and coverage gaps are the main residual risk. Hence
  fail closed on ambiguity, plus regression tests.
- Rejections and regenerations add latency and LLM cost for allergic users.
- Annotations and user-edited recipes should go through the same check when they change
  ingredients (follow-up spec).

## Status

Proposed.
