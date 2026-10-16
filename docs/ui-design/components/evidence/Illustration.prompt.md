Frame for an illustration under a question: a drawing in a sandboxed frame, with an expand control that lays it over the whole window.
The control says "Expand illustration" on hover, on focus and to a screen reader.
Expanded, a close button, Escape and a press on the backdrop each call `onClose`, and the cursor returns to the expand control.

```jsx
<Illustration src="/sessions/abc/illustrations/c-1/q3.html" expanded={open} onExpand={() => setOpen(true)} onClose={() => setOpen(false)} />
```
