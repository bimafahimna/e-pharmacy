# E-Pharmacy Web Application Server-Side

## Overview

This repository hosts a Server-Side **E-Pharmacy Web Application** designed for seamless online purchasing of pharmaceuticals. Users can select their preferred pharmacy and enjoy the convenience of home delivery, all while benefiting from our optimized system that prioritizes speed, proximity, and cost-effectiveness.

### Key Components

- **Pharmacy Management System**: Enables pharmacists to efficiently manage inventory, orders, and customer interactions for their respective pharmacies.
- **Admin Management System**: Provides application administrators with tools to oversee user accounts, manage pharmacies, and analyze overall performance metrics.

## Features

- **User-Friendly Interface**: Intuitive design for easy navigation and shopping.
- **Pharmacy Selection**: Choose from various pharmacies to find the best options.
- **Optimized SQL Queries**: Fastest results for product and pharmacy searches, tailored to provide the cheapest options.
- **Go with Gin Framework**: Built using Go's Gin framework for robust performance and scalability.
- **Caching with Memcached**: Improves response times by caching frequent queries.
- **Pub/Sub with Redis**: Efficient real-time updates and notifications.
- **Docker Optimized**: Fully containerized for easy deployment and scalability.
- **Monitoring with Grafana**: Real-time metrics and logs to ensure application health.
- **Nginx Load Balancer**: Distributes incoming traffic for improved reliability and performance.
- **RDBMS Transactions**: Ensures atomicity and consistency for all database operations.
- **Goroutines**: Leverages Go's concurrency model for faster database interactions.

## Getting Started

### Prerequisite:
- golang migrate
- make
- docker
- postgresql

To get started with the E-Pharmacy Server, follow these steps:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/bimafahimna/E-Pharmacy-ServerSide.git
   ```
1. **Migrate the Database**:
   ```bash
   make migrateup
   ```
1. **Seed Data**:
   ```bash
   make data
   ```
1. **Create Docker Network**:
   ```bash
   make network
   ```
1. **Create Memcache Container**:
   ```bash
   make cacheup
   ```
1. **Start the Server**:
   ```bash
   make server
   ```

#### Notes:
start UI from [here](https://github.com/bimafahimna/E-Pharmacy-ClientSide) to have visualization