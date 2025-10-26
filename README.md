# entrepreneur-pastoral
Marketplace for catholic entrepreneurs where they can share their products and services.

## Goal
The main goal of this project is to bring together all catholic's entrepreneurs so they can share their products and services with the catholic community as well as offer jobs for other Catholics and find workers for their business.

## Architecture
- **Structure**: Modular monolith following DDD principles and the repository pattern.
- **Language**: Golang with Chi as router.
- **Database**: PostgreSQL.
- **Logging**: Uber Zap for logs.
- **DevOps**: Docker and Docker Compose for environment.
- **Cache**: Redis for cache and session store.
- **Queue**: RabbitMQ for notification queues.
- **UI**: Golang templates for web with HTMX and Tailwind.

<img width="1326" height="753" alt="Image" src="https://github.com/user-attachments/assets/d48514e2-f516-45ee-8605-d58ac7794e6f" />

## Modules
The system will have the following modules:
- **User**: Users data, Job profile, Roles and Notifications management.
- **Admin**: Configuration related stuff (users role, business config/approval, permissions/role management and pastoral info).
- **Entrepreneur**: Business data and contact info, products and services offered and jobs available.
- **Marketplace**: A meeting point for all users and business, where business can find suppliers and workers and users can find services/products or jobs.
- **Notification**: Handles all stuff related to notifications, based on user config.

