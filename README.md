# IndividualProjectP2
# Sport Venue Rental API

Sport Venue Rental API adalah REST API untuk aplikasi rental tempat/lapangan olahraga seperti futsal, badminton, basket, dan tennis. User dapat melakukan register, login, top up deposit menggunakan payment link Xendit, melihat daftar venue, melakukan rental venue, dan melihat rental history.

## Tech Stack

- Golang
- Echo Framework
- PostgreSQL
- JWT Authentication
- Xendit Payment Link
- Postman Documentation
- Railway Deployment

## Main Features

- Register user
- Login user
- Role authorization user/admin
- Admin create venue
- Admin update venue
- User view venues
- User top up deposit via Xendit
- Xendit callback simulation
- User rental venue using deposit
- User rental history

## Database Entities

### users

- id
- username
- email
- password
- deposit_amount
- role
- created_at
- updated_at

### venues

- id
- name
- category
- location
- stock_availability
- rental_cost
- created_at
- updated_at

### deposit_histories

- id
- user_id
- amount
- payment_link
- payment_status
- external_id
- created_at
- updated_at

### rentals

- id
- user_id
- venue_id
- rental_date
- start_time
- end_time
- total_cost
- status
- created_at
- updated_at

## ERD

```text
users 1 ---- many deposit_histories
users 1 ---- many rentals
venues 1 ---- many rentals