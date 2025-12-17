# Package `telemetry`

This package allows to set up custom telemetry processors that the ADK events will be emitted to.

## Functions

*   `RegisterSpanProcessor(processor sdktrace.SpanProcessor)`: Registers a span processor to the local trace provider instance.
