# Package `internal/llminternal`

This package houses the entire, complex execution pipeline for `LLMAgent` implementations, defining the multi-turn conversational loop, state management, and all necessary request/response processing.

## Core Architecture

The execution model revolves around the `Flow` struct, which orchestrates the interaction with the LLM.

### State Management (`agent.go`)

The `State` struct holds the internal configuration for an `LLMAgent`.

```go
type State struct {
	Model model.LLM

	Tools    []tool.Tool
	Toolsets []tool.Toolset

	IncludeContents string

	GenerateContentConfig *genai.GenerateContentConfig

	Instruction               string
	InstructionProvider       InstructionProvider
	GlobalInstruction         string
	GlobalInstructionProvider InstructionProvider

	DisallowTransferToParent bool
	DisallowTransferToPeers  bool

	InputSchema  *genai.Schema
	OutputSchema *genai.Schema

	OutputKey string
}
```

### The Execution Engine (`base_flow.go`)

The `Flow` struct defines the pipeline for processing requests and responses.

```go
type Flow struct {
	Model model.LLM

	RequestProcessors    []func(ctx agent.InvocationContext, req *model.LLMRequest) error
	ResponseProcessors   []func(ctx agent.InvocationContext, req *model.LLMRequest, resp *model.LLMResponse) error
	BeforeModelCallbacks []BeforeModelCallback
	AfterModelCallbacks  []AfterModelCallback
	BeforeToolCallbacks  []BeforeToolCallback
	AfterToolCallbacks   []AfterToolCallback
}
```

**Key Methods:**
*   `Run(ctx)`: The main entry point. It returns an iterator of events. It loops calling `runOneStep` until a final response is reached.
*   `runOneStep(ctx)`: Executes a single turn:
    1.  **Preprocess**: Runs `RequestProcessors` (e.g., building history, injecting instructions).
    2.  **Call LLM**: Invokes the model, handling streaming and callbacks.
    3.  **Postprocess**: Runs `ResponseProcessors`.
    4.  **Handle Tools**: If the model requests tool calls, executes them via `handleFunctionCalls` and yields the results.
    5.  **Agent Transfer**: Checks if control should be passed to another agent.

## Request Processors

These functions populate and modify the `model.LLMRequest` before it is sent to the LLM.

| Processor | File | Description |
| :--- | :--- | :--- |
| `basicRequestProcessor` | `basic_processor.go` | Sets up generation config and output schemas. |
| `instructionsRequestProcessor` | `instruction_processor.go` | Injects session state variables (e.g., `{app:key}`) into instructions. |
| `ContentsRequestProcessor` | `contents_processor.go` | Assembles conversation history from session events, filtering by branch and formatting cross-agent replies. |
| `AgentTransferRequestProcessor` | `agent_transfer.go` | Dynamically adds the `transfer_to_agent` tool if applicable. |
| `removeDisplayNameIfExists` | `file_uploads_processor.go` | Sanitizes file parts for compatibility with Gemini API. |

## Helper Components

### Stream Aggregation (`stream_aggregator.go`)
`streamingResponseAggregator` combines partial text chunks from streaming responses into coherent `LLMResponse` objects.

### Converters (`converters/`)
Provides utilities to map between `genai` types and ADK `model` types.
*   `Genai2LLMResponse`: Converts `genai.GenerateContentResponse` to `model.LLMResponse`.

### Google LLM Variants (`googlellm/`)
`variant.go` detects whether to use Vertex AI or Gemini API based on environment variables.

## Key Instructions Logic (`instruction_processor.go`)

Supports dynamic template substitution in instructions:
*   `InjectSessionState`: Replaces placeholders like `{variable}` with session state values or `{artifact.filename}` with file contents.

## Deep Dive: LLM Agent Execution Flow (`internal/llminternal/base_flow.go`)

The file `internal/llminternal/base_flow.go` is the central orchestrator and execution engine for all Large Language Model (LLM) based agents (`LLMAgent`) within the ADK. It defines the iterative, multi-turn loop that drives the agent's reasoning, tool use, and response generation cycle.

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
	AfterModelCallbacks  []AfterModelCallback
	BeforeToolCallbacks  []BeforeToolCallback
	AfterToolCallbacks   []AfterToolCallback
}
```

### 2. The Core Agent Execution Loop (`Flow.Run`)

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

### 3. Single Turn Execution (`runOneStep`)

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
3.  **Yield**: The completed `session.Event` is yielded.

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
