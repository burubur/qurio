To achieve 95% unit test coverage for the frontend (Vue.js), you must bridge the gap from the current 87.5% status by focusing on uncovered UI components and complex view-level orchestration. While core stores and feature components are already above 88%, the primary drag on the percentage is the components/ui directory (currently at ~78%).
The following strategic actions are required to reach the 95% target:
1. Exhaustive Testing of Shadcn UI Wrappers
The components/ui folder contains many base components (Badge, Button, Card, Select, Textarea) that are used throughout the application.
• Action: Implement tests for every variant and state defined in the cva (class-variance-authority) configs.
• Target Coverage: For components like Button and Badge, write tests that mount every combination of variant (default, destructive, outline, etc.) and size.
• Complex Components: Focus on the Select component suite, ensuring that the teleported SelectContent and SelectItem are correctly rendered and that interaction (opening/closing/selecting) is verified using attachTo: document.body.
2. Comprehensive Orchestration Testing for Views
Views like SourceDetailView.vue and DashboardView.vue handle significant orchestration logic that is often skipped in unit tests.
• Action: Test the interaction between vue-router and pinia stores within the views.
• Specific Gaps: In SourceDetailView.vue, you must test the polling interval logic that updates source status every 2 seconds and verify that the clearInterval is called on onUnmounted.
• View Mappings: Verify that DashboardView correctly triggers fetchStats() and fetchSources() upon mounting.
3. Deep-Dive into Store Error Paths and Polling
While stores like source.store.ts have high coverage, they contain "micro-logic" that must be 100% exercised.
• Action: Ensure every catch block in the store actions (e.g., fetchSources, addSource, uploadSource) has a corresponding test case that simulates a failed fetch call.
• Polling Logic: Specifically test the logic in startPolling which conditionally triggers a background fetch only if hasActiveSources is true.
4. Verify Custom Metadata and Formatting Logic
The frontend now handles more complex metadata fields like author, created_at, and language.
• Action: Update tests for SourceList.vue and SourceDetailView.vue to verify that they correctly parse and display these new top-level metadata fields.
• Utility Coverage: Ensure the cn (class merging) utility in lib/utils.ts has 100% coverage by testing various combinations of Tailwind classes and conditional logic.
5. Standardize Test Stubs
To prevent parsing errors and "false negatives" when components use many sub-components (like Lucide icons or complex UI wrappers), utilize global stubs in your Vitest configuration.
• Goal: This allows you to focus coverage on the specific logic of the component under test rather than its children, which should have their own isolated specs.

--------------------------------------------------------------------------------
Analogy for Frontend Coverage: Reaching 95% coverage is like calibrating a high-end monitor. Currently, the screen is bright and clear in the center (Core Features), but the edges and the fine settings menu (UI Components and View Orchestration) haven't been fully adjusted. To reach professional grade, you must ensure that every single pixel (component variant) and every adjustment slider (store action) works exactly as intended under every possible lighting condition (error state).