# Categorical Flow Diagrams

Enhanced visualization of the Elm Architecture data flow using categorical concepts.

---

## Problem Statement

The existing BUBBLETEA-ARCHITECTURE.md shows the Elm flow as:

```
User Input -> Msg -> Update -> Model -> View -> Renderer -> Terminal
```

**Issues**:
- Doesn't show the categorical structure (functors, monads, products)
- Misses the quantitative aspect (morphism counts, coupling metrics)
- No clear representation of the fixed-point/recursive nature
- Coproduct structure of messages not visualized

**Goal**: Create visualizations that reveal categorical structure while remaining intuitive.

---

## Categorical Structure of Elm Architecture

### Type Signatures

```haskell
-- Core types
Model  : Type                           -- State object
Msg    : Type = KeyMsg | MouseMsg | ... -- Coproduct (sum type)
Cmd    : Type -> Type                   -- Free monad (IO description)

-- Core morphisms
init   : Unit -> (Model, Cmd Msg)       -- Initial state + effects
update : Msg -> Model -> (Model, Cmd Msg)  -- State transition
view   : Model -> VirtualDOM            -- Pure rendering (functor)
```

### Categorical Interpretation

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ELM ARCHITECTURE AS CATEGORY                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  OBJECTS (Types):                                                       │
│  ──────────────                                                        │
│  • Model     : State space object                                       │
│  • Msg       : Event space object (coproduct ∑)                        │
│  • Cmd       : Effect space (free monad F)                              │
│  • VirtualDOM: UI representation                                        │
│  • Unit      : Terminal object (1)                                      │
│                                                                         │
│  MORPHISMS (Functions):                                                 │
│  ─────────────────────                                                  │
│  • init   : 1 → Model × Cmd    (product introduction)                  │
│  • update : Msg × Model → Model × Cmd  (state machine)                 │
│  • view   : Model → VirtualDOM         (functor mapping)               │
│  • render : VirtualDOM → IO ()         (natural transformation)        │
│                                                                         │
│  CATEGORICAL STRUCTURES:                                                │
│  ──────────────────────                                                │
│  • Model × Cmd : Product (⊗) - state paired with effects              │
│  • Msg = ∑ MsgType_i : Coproduct - tagged union of events             │
│  • Cmd : Free Monad - composable effect descriptions                   │
│  • view : Functor - structure-preserving map                           │
│  • update : Coalgebra - state machine as unfold                        │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Enhanced Flow Diagram 1: Product Structure

Show how the product (Model × Cmd) flows through the architecture:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    ELM ARCHITECTURE: PRODUCT FLOW                              ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║     ┌─────────────────────────────────────────────────────────────────┐      ║
║     │                         INIT PHASE                              │      ║
║     │                                                                 │      ║
║     │      1 ──init──→ (Model₀ ⊗ Cmd₀)                               │      ║
║     │         unit     ↙ π₁       ↘ π₂                               │      ║
║     │                Model₀       Cmd₀                                │      ║
║     │                  │            │                                 │      ║
║     │                  ▼            ▼                                 │      ║
║     │              [STATE]     [EFFECTS]                              │      ║
║     │                                                                 │      ║
║     └─────────────────────────────────────────────────────────────────┘      ║
║                                    │                                         ║
║                                    ▼                                         ║
║     ┌─────────────────────────────────────────────────────────────────┐      ║
║     │                        UPDATE CYCLE                              │      ║
║     │                                                                  │      ║
║     │                 ┌─────────────┐                                  │      ║
║     │   Msg ─────────→│             │                                  │      ║
║     │                 │   update    │──→ (Model' ⊗ Cmd')               │      ║
║     │   Model ───────→│             │      ↙        ↘                  │      ║
║     │                 └─────────────┘   Model'      Cmd'               │      ║
║     │                                     │           │                │      ║
║     │                                     │           ▼                │      ║
║     │                                     │     ┌──────────┐           │      ║
║     │                                     │     │  RUNTIME │           │      ║
║     │                                     │     │ Execute  │           │      ║
║     │                                     │     │  Cmds    │           │      ║
║     │                                     │     └────┬─────┘           │      ║
║     │                                     │          │                 │      ║
║     │                                     │          ▼                 │      ║
║     │                                     │     ┌──────────┐           │      ║
║     │                                     │     │   Msg'   │           │      ║
║     │                                     │     │(feedback)│           │      ║
║     │                                     │     └────┬─────┘           │      ║
║     │                                     │          │                 │      ║
║     │                                     ▼          │                 │      ║
║     │                              ┌───────────┐     │                 │      ║
║     │                              │   view    │     │                 │      ║
║     │                              │ (Functor) │     │                 │      ║
║     │                              └─────┬─────┘     │                 │      ║
║     │                                    │          │                  │      ║
║     │                                    ▼          │                  │      ║
║     │                            ┌──────────────┐   │                  │      ║
║     │                            │ VirtualDOM   │   │                  │      ║
║     │                            └──────┬───────┘   │                  │      ║
║     │                                   │           │                  │      ║
║     │                                   ▼           │                  │      ║
║     │                            ┌──────────────┐   │                  │      ║
║     │                            │   render     │   │                  │      ║
║     │                            │   (η: nat.   │   │                  │      ║
║     │                            │    transf.)  │   │                  │      ║
║     │                            └──────┬───────┘   │                  │      ║
║     │                                   │           │                  │      ║
║     │                                   ▼           │                  │      ║
║     │                            ┌──────────────┐   │                  │      ║
║     │                            │   Terminal   │◄──┘                  │      ║
║     │                            │     I/O      │                      │      ║
║     │                            └──────┬───────┘                      │      ║
║     │                                   │                              │      ║
║     │                                   │ (User Input)                 │      ║
║     │                                   ▼                              │      ║
║     │                            ┌──────────────┐                      │      ║
║     │                            │    Msg       │◄─────────────────────│      ║
║     │                            │  (Input)     │         LOOP         │      ║
║     │                            └──────────────┘                      │      ║
║     │                                                                  │      ║
║     └──────────────────────────────────────────────────────────────────┘      ║
║                                                                               ║
║  CATEGORICAL EQUATION:                                                        ║
║  ═══════════════════                                                         ║
║                                                                               ║
║  Program = fix(λp. init >>= loop)                                            ║
║  where loop (m, c) = do                                                      ║
║          msg <- runCmd c <|> readInput                                       ║
║          let (m', c') = update msg m                                         ║
║          render (view m')                                                    ║
║          loop (m', c')                                                       ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Enhanced Flow Diagram 2: Message Coproduct

Visualize the sum type structure of messages:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    MESSAGE COPRODUCT DECOMPOSITION                            ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║                              ┌───────────────────┐                            ║
║                              │       Msg         │                            ║
║                              │   (Coproduct ∑)   │                            ║
║                              └─────────┬─────────┘                            ║
║                                        │                                      ║
║           ┌────────────┬───────────────┼───────────────┬────────────┐        ║
║           │            │               │               │            │         ║
║           ▼            ▼               ▼               ▼            ▼         ║
║     ┌──────────┐ ┌──────────┐   ┌──────────┐   ┌──────────┐ ┌──────────┐    ║
║     │ KeyMsg   │ │ MouseMsg │   │ SizeMsg  │   │ TickMsg  │ │CustomMsg │    ║
║     │ (keys)   │ │ (clicks) │   │ (resize) │   │ (timer)  │ │(app-def) │    ║
║     └────┬─────┘ └────┬─────┘   └────┬─────┘   └────┬─────┘ └────┬─────┘    ║
║          │            │              │              │            │           ║
║          ▼            ▼              ▼              ▼            ▼           ║
║     ┌─────────────────────────────────────────────────────────────────┐      ║
║     │                                                                 │      ║
║     │                       update : Msg × Model → Model × Cmd        │      ║
║     │                                                                 │      ║
║     │  Pattern Match on Msg Type:                                     │      ║
║     │                                                                 │      ║
║     │    case KeyMsg k:                                              │      ║
║     │      if k == 'q' then return (model, tea.Quit())               │      ║
║     │      else return (model.handleKey(k), tea.None())              │      ║
║     │                                                                 │      ║
║     │    case MouseMsg m:                                            │      ║
║     │      return (model.handleMouse(m), tea.None())                 │      ║
║     │                                                                 │      ║
║     │    case SizeMsg (w, h):                                        │      ║
║     │      return (model.resize(w, h), tea.None())                   │      ║
║     │                                                                 │      ║
║     │    case TickMsg t:                                             │      ║
║     │      return (model.tick(t), tea.None())                        │      ║
║     │                                                                 │      ║
║     │    case CustomMsg c:                                           │      ║
║     │      return app.handleCustom(c, model)                         │      ║
║     │                                                                 │      ║
║     └─────────────────────────────────────────────────────────────────┘      ║
║                                                                               ║
║  CATEGORICAL INTERPRETATION:                                                  ║
║  ═══════════════════════════                                                 ║
║                                                                               ║
║  Each Msg variant has an injection morphism:                                  ║
║                                                                               ║
║    ι₁ : KeyMsg → Msg      (keyboard injection)                               ║
║    ι₂ : MouseMsg → Msg    (mouse injection)                                  ║
║    ι₃ : SizeMsg → Msg     (resize injection)                                 ║
║    ι₄ : TickMsg → Msg     (timer injection)                                  ║
║    ι₅ : CustomMsg → Msg   (application injection)                            ║
║                                                                               ║
║  Universal Property:                                                          ║
║    ∀ handlers f₁...f₅ : MsgType_i → Result,                                  ║
║    ∃! h : Msg → Result such that h ∘ ι_i = f_i                               ║
║                                                                               ║
║  This is exactly what pattern matching provides!                              ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Enhanced Flow Diagram 3: Functor Chain

Show the view as a functor and render as natural transformation:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    VIEW & RENDER: FUNCTOR CHAIN                               ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  CATEGORY: Model                        CATEGORY: VirtualDOM                  ║
║  ═══════════════                        ════════════════════                 ║
║                                                                               ║
║  Objects:                               Objects:                              ║
║    Model                                  VNode, VText, VElement              ║
║    Model.counter                          VNode.text, VNode.div               ║
║    Model.items                            VNode.children                      ║
║                                                                               ║
║  Morphisms:                             Morphisms:                            ║
║    increment : Model → Model              render : VNode → VNode              ║
║    addItem : Model → Model                style : VNode → VNode               ║
║                                                                               ║
║                         view : F                                              ║
║          ┌─────────────────────────────────────────┐                         ║
║          │             FUNCTOR                     │                         ║
║          │                                         │                         ║
║          │  F(Model) = VirtualDOM                  │                         ║
║          │  F(f: A → B) = view f : VA → VB        │                         ║
║          │                                         │                         ║
║          │  Laws:                                  │                         ║
║          │    F(id) = id        (identity)        │                         ║
║          │    F(g ∘ f) = F(g) ∘ F(f)  (compose)  │                         ║
║          │                                         │                         ║
║          └─────────────────────────────────────────┘                         ║
║                          │                                                    ║
║                          ▼                                                    ║
║                                                                               ║
║  NATURAL TRANSFORMATION: η : VirtualDOM ⇒ IO                                 ║
║  ═══════════════════════════════════════════                                 ║
║                                                                               ║
║          ┌─────────────────────────────────────────┐                         ║
║          │          render : η                     │                         ║
║          │                                         │                         ║
║          │  For each VNode v:                      │                         ║
║          │    η_v : VirtualDOM(v) → IO(terminal)   │                         ║
║          │                                         │                         ║
║          │  Naturality Square:                     │                         ║
║          │                                         │                         ║
║          │    VirtualDOM(A) ──η_A──→ IO(A)        │                         ║
║          │         │                    │          │                         ║
║          │       F(f)               IO(f)         │                         ║
║          │         │                    │          │                         ║
║          │         ▼                    ▼          │                         ║
║          │    VirtualDOM(B) ──η_B──→ IO(B)        │                         ║
║          │                                         │                         ║
║          │  "Rendering commutes with structure"   │                         ║
║          │                                         │                         ║
║          └─────────────────────────────────────────┘                         ║
║                                                                               ║
║  CONCRETE EXAMPLE:                                                            ║
║  ════════════════                                                            ║
║                                                                               ║
║  Model { counter: 5, label: "Count" }                                        ║
║      │                                                                        ║
║      │ view                                                                   ║
║      ▼                                                                        ║
║  VNode.div([                                                                  ║
║    VNode.text("Count: "),                                                     ║
║    VNode.span({ style: bold }, "5")                                          ║
║  ])                                                                           ║
║      │                                                                        ║
║      │ render (η)                                                             ║
║      ▼                                                                        ║
║  "Count: \033[1m5\033[0m" → stdout                                           ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Enhanced Flow Diagram 4: Fixed Point Loop

The recursive nature of the program as a fixed point:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    PROGRAM AS FIXED POINT                                     ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║  Program = fix(F) where F is the iteration functor                           ║
║                                                                               ║
║  ┌─────────────────────────────────────────────────────────────────────────┐ ║
║  │                                                                         │ ║
║  │    ┌───────────────────────────────────────────────────────────────┐   │ ║
║  │    │                    ITERATION FUNCTOR F                        │   │ ║
║  │    │                                                               │   │ ║
║  │    │         ┌─────┐     ┌────────┐     ┌──────┐     ┌──────┐    │   │ ║
║  │    │  ──────→│Read │────→│ Update │────→│ View │────→│Render│────│   │ ║
║  │    │  │      │Input│     │        │     │      │     │      │    │   │ ║
║  │    │  │      └─────┘     └────────┘     └──────┘     └──────┘    │   │ ║
║  │    │  │                                                          │   │ ║
║  │    │  │                                                          │   │ ║
║  │    └──│──────────────────────────────────────────────────────────┘   │ ║
║  │       │                                                              │ ║
║  │       │  UNFOLD (coalgebra)                                          │ ║
║  │       │                                                              │ ║
║  │       │  State: (Model, Cmd)                                         │ ║
║  │       │                                                              │ ║
║  │       ▼                                                              │ ║
║  │    ┌──────────────────────────────────────────────────────────────┐  │ ║
║  │    │                    STATE STREAM                               │  │ ║
║  │    │                                                               │  │ ║
║  │    │  s₀ ──F──→ s₁ ──F──→ s₂ ──F──→ s₃ ──F──→ ...                 │  │ ║
║  │    │  │         │         │         │                              │  │ ║
║  │    │  ▼         ▼         ▼         ▼                              │  │ ║
║  │    │  UI₀      UI₁       UI₂       UI₃       ...                  │  │ ║
║  │    │                                                               │  │ ║
║  │    └──────────────────────────────────────────────────────────────┘  │ ║
║  │                                                                      │ ║
║  │    TERMINATION CONDITION:                                            │ ║
║  │    ════════════════════                                             │ ║
║  │                                                                      │ ║
║  │    When update returns (_, tea.Quit()):                             │ ║
║  │      - Break the fixed point loop                                   │ ║
║  │      - Cleanup terminal state                                       │ ║
║  │      - Return final model                                           │ ║
║  │                                                                      │ ║
║  └──────────────────────────────────────────────────────────────────────┘ ║
║                                                                               ║
║  HASKELL REPRESENTATION:                                                     ║
║  ═══════════════════════                                                    ║
║                                                                               ║
║  -- The fixed point combinator                                               ║
║  fix :: (a -> a) -> a                                                        ║
║  fix f = let x = f x in x                                                    ║
║                                                                               ║
║  -- Elm program as fixed point                                               ║
║  program :: Model -> IO ()                                                   ║
║  program = fix $ \loop model -> do                                           ║
║      msg <- readInput                                                        ║
║      let (model', cmd) = update msg model                                    ║
║      effects <- runCmd cmd                                                   ║
║      render (view model')                                                    ║
║      case cmd of                                                             ║
║          Quit -> return ()                                                   ║
║          _    -> loop model'                                                 ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Quantitative Flow Annotation

Add metrics from analysis.json to the flow:

```
╔═══════════════════════════════════════════════════════════════════════════════╗
║                    QUANTIFIED ELM ARCHITECTURE                                ║
╠═══════════════════════════════════════════════════════════════════════════════╣
║                                                                               ║
║                           ┌───────────────────────┐                          ║
║                           │     User Input        │                          ║
║                           │   (External Event)    │                          ║
║                           └───────────┬───────────┘                          ║
║                                       │                                       ║
║                                       │ 12 morphisms                          ║
║                                       ▼                                       ║
║                           ┌───────────────────────┐                          ║
║                           │       key.go          │                          ║
║                           │       E: 18           │                          ║
║                           │  (Input Decoding)     │                          ║
║                           └───────────┬───────────┘                          ║
║                                       │                                       ║
║                                       │ Msg                                   ║
║                                       ▼                                       ║
║            ┌──────────────────────────────────────────────────────┐          ║
║            │                    tea.go                             │          ║
║            │                    E: 53                              │          ║
║            │              (Program Orchestrator)                   │          ║
║            │                                                       │          ║
║            │  ┌─────────────┐    ┌──────────────┐                 │          ║
║            │  │   Update    │───→│ Model × Cmd  │                 │          ║
║            │  │ (53 morph)  │    │   (Product)  │                 │          ║
║            │  └─────────────┘    └──────┬───────┘                 │          ║
║            │                            │                          │          ║
║            └────────────────────────────┼──────────────────────────┘          ║
║                                         │                                     ║
║                         ┌───────────────┼───────────────┐                    ║
║                         │ π₁ (Model)    │ π₂ (Cmd)      │                    ║
║                         ▼               │               ▼                    ║
║            ┌─────────────────────┐      │      ┌─────────────────────┐       ║
║            │  standard_renderer  │      │      │      Runtime        │       ║
║            │       E: 55         │      │      │   (Cmd Executor)    │       ║
║            │   (View Functor)    │      │      │    → New Msg        │       ║
║            └──────────┬──────────┘      │      └──────────┬──────────┘       ║
║                       │                 │                 │                  ║
║                       │ 55 morphisms    │                 │ feedback         ║
║                       ▼                 │                 │                  ║
║            ┌─────────────────────┐      │                 │                  ║
║            │      screen.go      │      │                 │                  ║
║            │       E: 32         │      │                 │                  ║
║            │  (Screen Buffer)    │      │                 │                  ║
║            └──────────┬──────────┘      │                 │                  ║
║                       │                 │                 │                  ║
║                       │ render          │                 │                  ║
║                       ▼                 │                 │                  ║
║            ┌─────────────────────┐      │                 │                  ║
║            │     Terminal        │◄─────┴─────────────────┘                  ║
║            │   (η: IO monad)     │                                           ║
║            └─────────────────────┘                                           ║
║                                                                               ║
║  METRICS FROM ANALYSIS.JSON:                                                  ║
║  ═══════════════════════════                                                 ║
║                                                                               ║
║  Total Morphisms: 769                                                         ║
║  Core Hub (tea.go + renderer): 108 efferent (14% of total)                   ║
║  Supporting (screen + key + options): 71 efferent (9% of total)              ║
║  Examples: ~400 efferent (52% of total) - distributed across 50 apps         ║
║  Imports: 264 afferent (34% of total) - stable foundations                   ║
║                                                                               ║
╚═══════════════════════════════════════════════════════════════════════════════╝
```

---

## Interactive Animation Concept

For web-based visualization, animate the flow:

```
STATE 0: INIT
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ○ init                                                            │
│     ↓                                                              │
│  ● (Model₀, Cmd₀) ←── ACTIVE                                       │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

STATE 1: RECEIVING INPUT (Key 'a')
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ○ init                                                            │
│     ↓                                                              │
│  ○ (Model₀, Cmd₀)                                                   │
│     ↓                                                              │
│  ● KeyMsg('a') ←── ACTIVE                                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

STATE 2: UPDATE
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ○ KeyMsg('a')                                                      │
│     ↓                                                              │
│  ○ update(KeyMsg('a'), Model₀)                                      │
│     ↓                                                              │
│  ● (Model₁, Cmd.None) ←── ACTIVE                                   │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

STATE 3: VIEW
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ○ (Model₁, Cmd.None)                                               │
│     ↓                                                              │
│  ● view(Model₁) → VirtualDOM ←── ACTIVE                            │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘

STATE 4: RENDER
┌─────────────────────────────────────────────────────────────────────┐
│                                                                     │
│  ○ VirtualDOM                                                       │
│     ↓                                                              │
│  ● render → Terminal ←── ACTIVE                                    │
│     ↓                                                              │
│  ○ Wait for next input...                                          │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Summary

Enhanced categorical flow diagrams reveal:

1. **Product structure** (Model × Cmd) - paired state and effects
2. **Coproduct structure** (Msg) - sum type of events
3. **Functor nature** (view) - structure-preserving rendering
4. **Natural transformation** (render) - VirtualDOM → IO
5. **Fixed point** (program) - recursive loop structure
6. **Quantitative metrics** - morphism counts from analysis

These visualizations make the mathematical structure intuitive while preserving categorical precision.

---

## Next Steps

- `05-TOOLING-RECOMMENDATIONS.md`: Implementation tools
- `06-OIS-INTEGRATION.md`: Leveraging OIS dependency-observer
