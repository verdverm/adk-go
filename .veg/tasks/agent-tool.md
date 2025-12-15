Evaluate, Understand, and Report about the `agentTool`, the issues I'm seeing, and the fixes or changes needed.

### References

- ./tool/agenttool/agent_tool.go (where `agentTool` is defined)
- ./tool/functiontool/function.go (where `functionTool` is defined)
- ./agent/llmagent/llmagent.go (where `llmAgent` is defined)
- /internal/llminternal/base_flow.go (where the main loop for `req -> loop(llm -> tool calls) -> resp`)

### Context

I'm using the `agentTool` type to have my main <agent> (`llmAgent`) call a <subagent> (`llmAgent`) as a <tool>.
It is not behaving the same as `functionTool`, the issues are described next.

### Issues

1. When the `agentTool` returns, it seems to return all the way to the user, not the main <agent> that called `agentTool`. How can we fix this so that it returns to the <agent>? The <agent> needs to continue working on the query, problem, or solution, but it is not... :[
2. The `agentTool` runs an agent, which itself supports chat, returns messages, and calls <tools> as well. The `agentTool` events are not propegating back / through the <agent>. My UI suffers because it looks like it's hanging, and I cannot validate with the `agentTool` is doing. How do we fix this?
3. The `functionTool` shares state correctly, (a) it sees the parent state (b) <tool> state changes are available in the calling <agent>. With `agentTool`, this is not so. It looks like the copy is only beforehand, not afterwards. How do we fix this?


### Reporting Guidance:

You MUST summarize your findings, relevant code snippets, analysis, and suggested code fixes in a final and comprehensive report.
You MUST the follow organization and formating guidance:

1. Use markdown for your report, wrap code snippets in ```<lang> for highlighting, prefer nested bullet lists to blocks of text. 1-3 sentences are fine.
2. Write a short introduction section, with a list each that covers the key points for: the user query, cause of the issues, and suggested fixes or changes.
2. Write sections for the core components, types, main loop, and how the request flow is processed by them.
3. Write sections for the <issues> and fixes. Include code snippets for both before and after. Explain why it is not working and why the change results in the correct behavior. Explain the before and after flow for early exit, event propagation, and state propagation. Divide this issue-by-issue or comprehensive before-after. Base your choice on your analysis of how the code works, the reasons for the issues, and the suggested changes.

IMPORTANT: You should write the report to `agentTool.md`. You should update this file as you iterate on understanding the problem and formulating a solution.
IMPORTANT: You should make several iterations of a evaluate-write-reflect loop.
