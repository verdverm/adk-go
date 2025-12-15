# ADK Internal LLM Agent Components (`internal/llminternal/`)

This package houses the entire, complex execution pipeline for `LLMAgent` implementations, defining the multi-turn conversational loop, state management, and all necessary request/response processing.

## Core Execution and State

| File/Sub-Package | Description |
| :--- | :--- |
| `agent.go` | Defines the internal `State` for an `LLMAgent`, including the `model.LLM`, lists of `tool.Tool` and `tool.Toolset`, instruction strings, schemas (`InputSchema`, `OutputSchema`), and configuration flags for agent transfers. |
| `base_flow.go` | **The Core Orchestrator.** Defines the `Flow` struct, which contains pipelines of `RequestProcessors` and `ResponseProcessors`, and execution callbacks. The `Flow.Run` method implements the primary iterative agent loop, continually calling `runOneStep` to process model calls, handle function calls (`handleFunctionCalls`), and manage flow control until a final response. |
| `stream_aggregator.go` | Contains the `streamingResponseAggregator` to efficiently combine fragmented text parts from a streaming LLM response into full, coherent chunks for yielding, ensuring data integrity during streaming. |
| `converters/` | Houses `converters.go`, which provides `Genai2LLMResponse` for converting raw `genai` responses into the ADK's standard `model.LLMResponse`. |
| `googlellm/` | Contains `variant.go`, which provides utilities to determine the specific Google LLM environment (e.g., Vertex AI vs. Gemini API) based on runtime configuration/environment variables. |

## Request Processors (Preprocessing before LLM call)

These functions modify the `model.LLMRequest` to prepare it for the LLM, running in a configurable pipeline defined in `Flow`.

| File | Functionality |
| :--- | :--- |
| `basic_processor.go` | Sets fundamental LLM generation configurations (`genai.GenerateContentConfig`) and applies the agent's `OutputSchema` to the request. |
| `instruction_processor.go` | Implements logic to load and inject session-scoped variables (e.g., `{app:key}`, `{artifact.file_name}`) into the agent's instructions and global instructions before they are sent to the LLM. |
| `contents_processor.go` | Constructs the conversation history (`req.Contents`). It is responsible for: 1) Filtering history by the current `Branch`, 2) Converting peer agent replies into contextual user-authored content, and 3) Rearranging function call/response event pairs in the history to maintain conversational integrity. |
| `agent_transfer.go` | Implements `AgentTransferRequestProcessor` which dynamically enables the `transfer_to_agent` tool based on the agent's position in the tree and configured transfer policies (`DisallowTransferToParent`, `DisallowTransferToPeers`). |
| `file_uploads_processor.go` | Removes the `DisplayName` from file parts for compatibility with specific LLM backends (e.g., Gemini API). |
| `other_processors.go` | Contains unexported placeholder functions for functionality like NL Planning (`nlPlanningRequestProcessor`) and Code Execution (`codeExecutionRequestProcessor`). |

## Deep Dive: LLM Agent Execution Flow (internal/llminternal/base_flow.go)

The file `internal/llminternal/base_flow.go` is the central orchestrator and execution engine for all Large Language Model (LLM) based agents (`LLMAgent`) within the ADK. It defines the iterative, multi-turn loop that drives the agent's reasoning, tool use, and response generation cycle.

---

### 1. Core Struct: Flow and Configuration

The `Flow` struct holds the configurable pipelines and callbacks that determine how an LLM agent executes a request. It is the runtime configuration layer over the LLM.

#### Flow Struct Definition

```go
type Flow struct {
	Model model.LLM // The model instance (implements model.LLM)

	// Pipeline of functions run BEFORE the LLM call
	RequestProcessors    []func(ctx agent.InvocationContext, req *model.LLMRequest) error
	
	// Pipeline of functions run AFTER the LLM call
	ResponseProcessors   []func(ctx agent.InvocationContext, req *model.LLMRequest, resp *model.LLMResponse) error
	
	// Callbacks for introspection/modification
	BeforeModelCallbacks []BeforeModelCallback
	AfterModelCallbacks  []AfterModelCallbacks
	BeforeToolCallbacks  []BeforeToolCallback
	AfterToolCallbacks   []AfterToolCallbacks
}
```

#### Relationship to Non-Internal Interfaces

The `Flow` manages interaction with the core public interfaces:

| Internal Component | Public Interface/Type | Role |
| :--- | :--- | :--- |
| `Model` field | `model.LLM` | Used directly to call `GenerateContent` for text generation. |
| `RequestProcessors` | `agent.InvocationContext` | All processors receive this context for full access to the running environment. |
| Callbacks | `agent.CallbackContext`, `tool.Context` | Custom logic is injected using these restricted-access contexts. |
| Execution Output | `session.Event` | The entire flow yields public `session.Event` structures as output. |

---

### 2. The Core Agent Execution Loop (Flow.Run)

The `Run` method implements the high-level iterative nature of the LLM agent, constantly calling `runOneStep` until a final state is reached. This enables the agent to handle multi-turn reasoning and tool-use cycles without external orchestration.

```go
func (f *Flow) Run(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
	return func(yield func(*session.Event, error) bool) {
		for {
			var lastEvent *session.Event
			for ev, err := range f.runOneStep(ctx) {
				// ... Yields events ...
				lastEvent = ev
			}
			
			// The loop terminates only if the last event signals a final response.
			if lastEvent == nil || lastEvent.IsFinalResponse() {
				return
			}
		}
	}
}
```

The loop continues as long as `lastEvent.IsFinalResponse()` returns `false`, which typically means the agent either called a function that requires a response (starting a new turn) or was interrupted.

---

### 3. Single Turn Execution (runOneStep)

The `runOneStep` function is the core logic that defines a single interaction turn, encompassing request preparation, model interaction, and post-processing.

#### Phase 1: Preprocessing (`f.preprocess`)

This phase ensures the `model.LLMRequest` is fully prepared before being sent to the LLM.

1.  **Request Processors**: Sequential functions run against the request (e.g., `contentsRequestProcessor` to assemble chat history).
2.  **Tool Preprocessing (`toolPreprocess`)**: Iterates over all active tools. Tools implementing the internal `toolinternal.RequestProcessor` interface are called to inject necessary components (like function declaration schemas, handled by the `toolinternal` package) into `req.Config.Tools`.

#### Phase 2: Calling the LLM (`f.callLLM`)

1.  **Before Model Callbacks**: Executes `BeforeModelCallbacks`. These can modify the request or halt execution by yielding a response directly.
2.  **Model Call**: Calls `f.Model.GenerateContent(...)`, retrieving an iterator of responses (supporting streaming).
3.  **After Model Callbacks**: Executes `AfterModelCallbacks` on each response chunk or error, allowing for post-call processing or error suppression.

#### Phase 3: Post-processing and Event Finalization

After the model generates content, this phase finalizes the event for the outside world.

1.  **Response Processors**: Executes `f.postprocess` using configured `ResponseProcessors` (e.g., for cleaning up internal thoughts added during preprocessing).
2.  **Event Finalization (`f.finalizeModelResponseEvent`)**:
    *   Generates a unique, internal client ID for any function call requested by the model (`utils.PopulateClientFunctionCallID`).
    *   Creates a `session.Event`, populating its `LLMResponse`, `Branch`, `Author`, and any state changes (`stateDelta`) collected during model callbacks.
    *   Identifies and records any `LongRunningToolIDs`.
3.  The completed `session.Event` is yielded.

---

### 4. Tool and Agent Transfer Handling

If the yielded event contains one or more `FunctionCall` requests from the LLM, the flow enters this crucial execution block.

#### `f.handleFunctionCalls`

1.  **Execution and Callbacks**: For each requested function call:
    *   `f.callTool` is executed, running `BeforeToolCallbacks`.
    *   The tool's `Run` method (from the internal `toolinternal.FunctionTool` interface) is executed.
    *   `AfterToolCallbacks` are run on the result.
2.  **Function Response Event**: The result of the tool execution is packaged into a new `session.Event` containing a `genai.FunctionResponse` part (Role: "user").
3.  **Merging**: Multiple concurrent tool responses are combined into a single event using `mergeParallelFunctionResponseEvents`.
4.  **Yield**: The response event is yielded. Because this event signals tool output that needs to be seen by the LLM, the outer `Flow.Run` loop restarts, and the new event becomes part of the history for the *next* LLM call.

#### Agent Transfer

After handling function calls, `runOneStep` checks the resultant event's actions:

```go
if ev.Actions.TransferToAgent != "" {
    // 1. Find the target agent in the agent tree.
    nextAgent := f.agentToRun(ctx, ev.Actions.TransferToAgent)
    // 2. Execute the new agent.
    for ev, err := range nextAgent.Run(ctx) {
        // ... Yields events from the new agent run ...
    }
}
```

If a `transfer_to_agent` action is set (which is triggered by the model calling a dynamic tool), control is immediately handed over to the target agent's `Run` method, effectively ending the current LLM agent's turn and initiating a run on a different agent in the tree.