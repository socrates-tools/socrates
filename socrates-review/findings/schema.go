package findings

const Schema = `{
  "type": "object",
  "additionalProperties": false,
  "required": ["findings", "risk_level", "risk_rationale"],
  "properties": {
    "findings": {
      "type": "array",
      "items": {
        "type": "object",
        "additionalProperties": false,
        "required": ["id", "severity", "file", "line", "description", "action"],
        "properties": {
          "id": {
            "type": "string",
            "minLength": 1,
            "description": "An identifier unique within this document, such as f1."
          },
          "severity": {
            "type": "string",
            "enum": ["error", "warning", "info"],
            "description": "Whether this blocks. error and warning block; info does not."
          },
          "file": {
            "type": "string",
            "minLength": 1,
            "description": "The path of the file, relative to the repository root."
          },
          "line": {
            "type": "integer",
            "minimum": 1,
            "description": "The line the finding is about, counted from one."
          },
          "description": {
            "type": "string",
            "minLength": 1,
            "description": "What is wrong, and how it fails."
          },
          "action": {
            "type": "string",
            "enum": ["no-op", "auto-fix", "ask-user"],
            "description": "Who resolves this. no-op is nothing to do; auto-fix is a mechanical change the author can make from this description alone; ask-user is anything needing a decision."
          }
        }
      }
    },
    "risk_level": {
      "type": "string",
      "enum": ["low", "medium", "high"],
      "description": "How risky this change is overall. Recorded, and it does not decide the verdict."
    },
    "risk_rationale": {
      "type": "string",
      "description": "One sentence saying why that risk level."
    }
  }
}`
