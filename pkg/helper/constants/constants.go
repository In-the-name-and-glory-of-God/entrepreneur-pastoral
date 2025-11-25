package constants

const (
	DEVELOPMENT = "development"
	PRODUCTION  = "production"
)

const (
	ROLE_ADMIN int16 = iota + 1
	ROLE_MANAGER
	ROLE_ASSISTANT
	ROLE_ENTREPRENEUR
	ROLE_USER
)

const (
	QUEUE_NOTIFICATIONS = "notifications"
)

const (
	EMAIL_TEMPLATE_WELCOME        = "welcome.html"
	EMAIL_TEMPLATE_VERIFY_ACCOUNT = "verify_account.html"
	EMAIL_TEMPLATE_PASSWORD_RESET = "password_reset.html"
)
