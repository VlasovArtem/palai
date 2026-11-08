# 0006. Dietary safety: deterministic allergen and strict-diet check on every generated recipe

- Status: Proposed
- Date: 2026-11-08

## Context

Users declare three kinds of dietary constraint:

- **Allergies** (strict): eating the allergen can cause harm.
- **Diets** (strict): vegan, vegetarian, pescatarian, halal, kosher and similar. Breaking them
  is a serious failure of trust (ethical, religious or medical), even when it is not dangerous.
- **Preferences** (soft): for example "not sweet" or "less spicy".

An LLM can ignore instructions, hallucinate, or miss that an ingredient is derived from a
restricted source. Butter, ghee, whey and casein are dairy. Soy sauce contains wheat. Pesto
often contains pine nuts and Parmesan. Gelatin and fish sauce are not vegan, and most Parmesan
is not vegetarian because of animal rennet. Serving a recipe that breaks a declared allergy or
diet is a safety failure, not a quality issue. The check must be explainable, testable and
reproducible.

## Decision

We will run a deterministic check in Go on every generated recipe before it is returned or
stored. The check covers both declared allergens and declared strict diets. The LLM is never
the authority on dietary safety.

### Taxonomy

- **Ingredient categories** are the shared building block: dairy, egg, peanut, tree nuts,
  wheat/gluten, soy, fish, shellfish, sesame, meat, pork, poultry, alcohol, honey, animal-derived
  additives (gelatin, animal rennet, carmine, ...), and others as needed. Each category maps to
  known ingredients and derivatives, for example `butter → dairy`, `ghee → dairy`,
  `mayonnaise → egg`, `tahini → sesame`, `fish sauce → fish`, `gelatin → animal-derived`.
- **Allergens** map to one or more categories. At minimum: the EU 14 and US major allergens.
- **Diets** are fixed codes, each mapped to the categories they forbid. For example:
  - vegan: meat, poultry, pork, fish, shellfish, dairy, egg, honey, animal-derived
  - vegetarian: meat, poultry, pork, fish, shellfish, animal-derived
  - pescatarian: meat, poultry, pork, animal-derived non-fish additives
  - halal: pork, alcohol, non-halal animal-derived additives
  - kosher: pork, shellfish, and mixing meat with dairy in one recipe
- Diets with rules beyond "forbidden category" (for example kosher meat/dairy separation) get
  explicit rule functions with their own tests. A diet we cannot check deterministically is not
  offered as strict.
- Allergen and diet codes in the dietary profile API come from this taxonomy. Unknown codes are
  rejected.

### Matching

- Ingredient names are normalized (case, plurals, punctuation, common qualifiers such as
  "unsalted" or "melted") before matching.
- Matching works on terms and phrases, not raw substrings. "Butternut squash" is not dairy,
  "peanut butter" is peanut, "coconut milk" is not dairy, "vegan mayonnaise" is not egg.
  Explicit "vegan", "plant-based" or "egg-free" qualifiers count only when listed in the
  taxonomy.
- The check runs on the structured ingredient list. Steps are also scanned for added
  ingredients ("finish with a knob of butter").
- `Check(recipe, allergies, diets) → []Violation`. Each violation names the ingredient, the
  allergen or diet, the category, and the matched term.

### Policy

- Fail closed. Any violation rejects the recipe. Ingredients that are ambiguous or unknown for
  a user with allergies or strict diets (for example "cheese" for a vegetarian, "stock" for a
  vegan) are flagged, not silently passed.
- On rejection we may regenerate a bounded, configurable number of times with the violations as
  feedback. After that we return a typed error (HTTP 422) that lists the unmet constraints.
- Preferences are soft. They influence the prompt and ranking but never cause a rejection.
- The taxonomy, matcher and diet rules have exhaustive table-driven tests: derivatives, false
  positives, false negatives, multi-constraint profiles. A regression test is added for every
  reported miss.
- Responses include which allergens and diets were checked. We do not claim certification
  (medical, halal, kosher). The API documentation states the check's limits.

## Consequences

- Dietary safety does not depend on model behaviour and is fully unit-testable.
- One taxonomy serves allergens and diets, so a new derivative such as `whey → dairy` protects
  dairy-allergic and vegan users at once.
- The taxonomy needs ongoing curation, and coverage gaps are the main residual risk. Hence fail
  closed on ambiguity, plus regression tests. Religious diets are the hardest to get right and
  may need review by people who follow them.
- Fail-closed rules will sometimes reject acceptable recipes (for example generic "cheese" for
  vegetarians). The generator prompt should ask for specific ingredients ("vegetarian
  Parmesan-style cheese") to reduce this.
- Rejections and regenerations add latency and LLM cost for users with constraints.
- Annotations and user-edited recipes should go through the same check when they change
  ingredients (follow-up spec).

## Status

Proposed.
