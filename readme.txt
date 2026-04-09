notification-service/
├── cmd/
│   └── main.go                 # Entry point
├── internal/
│   ├── core/
│   │   ├── models/             # Data structs (Notification, User)
│   │   └── ports/              # Interfaces (Notifier, Scheduler, TemplateEngine)
│   ├── providers/              # Concrete implementations (Slack, Email, InApp)
│   ├── services/               # Business logic (Routing, Processing)
│   └── templates/              # HTML/Text templates logic
└── go.mod