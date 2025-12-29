# Airbnb-clone

> A simplified Airbnb platform connecting people who have rooms/apartments for short-term rental (Hosts) with tourists who need accommodation (Guests).

## 🛠 Tech Stack
* **Language:** Go (Golang)
* **Database:** PostgreSQL
* **Cache:** Redis
* **DevOps:** Docker, Docker Compose

## 🚀 Getting Started

### 1. Prerequisites
* Go 1.25.5+
* Docker & Docker Compose

### 2. Installation & Setup

Clone the repository:
```bash
git clone [https://github.com/katatrina/airbnb-clone.git](https://github.com/katatrina/airbnb-clone.git)
cd airbnb-clone
```

Setup configuration
```bash
cp .env.example .env
```

Start infrastructure (Postgres, Redis)
```bash
make docker-up
```

Run the application
```bash
make run-user
```
